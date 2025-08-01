// Copyright 2025 The OpenChoreo Authors
// SPDX-License-Identifier: Apache-2.0

package deployableartifact

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	openchoreov1alpha1 "github.com/openchoreo/openchoreo/api/v1alpha1"
	"github.com/openchoreo/openchoreo/internal/controller"
	"github.com/openchoreo/openchoreo/internal/labels"
)

// DeployableArtifactCleanupFinalizer is the finalizer that is used to clean up deployable artifact resources.
const DeployableArtifactCleanupFinalizer = "openchoreo.dev/deployableartifact-cleanup"

// ensureFinalizer ensures that the finalizer is added to the deployable artifact.
// The first return value indicates whether the finalizer was added to the deployable artifact.
func (r *Reconciler) ensureFinalizer(ctx context.Context, deployableArtifact *openchoreov1alpha1.DeployableArtifact) (bool, error) {
	// If the deployable artifact is being deleted, no need to add the finalizer
	if !deployableArtifact.DeletionTimestamp.IsZero() {
		return false, nil
	}

	if controllerutil.AddFinalizer(deployableArtifact, DeployableArtifactCleanupFinalizer) {
		return true, r.Update(ctx, deployableArtifact)
	}

	return false, nil
}

func (r *Reconciler) finalize(ctx context.Context, old, deployableArtifact *openchoreov1alpha1.DeployableArtifact) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("deployableArtifact", deployableArtifact.Name)

	if !controllerutil.ContainsFinalizer(deployableArtifact, DeployableArtifactCleanupFinalizer) {
		// Nothing to do if the finalizer is not present
		return ctrl.Result{}, nil
	}

	// Mark the condition as finalizing and return so that the deployableArtifact will indicate that it is being finalized.
	// The actual finalization will be done in the next reconcile loop triggered by the status update.
	if meta.SetStatusCondition(&deployableArtifact.Status.Conditions, NewDeployableArtifactFinalizingCondition(deployableArtifact.Generation)) {
		if err := controller.UpdateStatusConditions(ctx, r.Client, old, deployableArtifact); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// Perform cleanup logic for referenced deployments
	artifactsDeleted, err := r.deleteDeploymentsAndWait(ctx, deployableArtifact)
	if err != nil {
		logger.Error(err, "Failed to delete deployments")
		return ctrl.Result{}, err
	}
	if !artifactsDeleted {
		logger.Info("Deployments are still being deleted", "name", deployableArtifact.Name)
		return ctrl.Result{}, nil
	}

	// Remove the finalizer once cleanup is done
	if controllerutil.RemoveFinalizer(deployableArtifact, DeployableArtifactCleanupFinalizer) {
		if err := r.Update(ctx, deployableArtifact); err != nil {
			return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
		}
	}

	logger.Info("Successfully finalized deployable artifact")
	return ctrl.Result{}, nil
}

// deleteDeploymentsAndWait deletes referenced deployments and waits for them to be fully deleted
func (r *Reconciler) deleteDeploymentsAndWait(ctx context.Context, deployableArtifact *openchoreov1alpha1.DeployableArtifact) (bool, error) {
	logger := log.FromContext(ctx).WithValues("deployableArtifact", deployableArtifact.Name)
	logger.Info("Cleaning up deployments")

	// Find all Deployments referred to by this DeployableArtifact
	deploymentList := &openchoreov1alpha1.DeploymentList{}
	listOpts := []client.ListOption{
		client.InNamespace(deployableArtifact.Namespace),
		client.MatchingLabels{
			labels.LabelKeyOrganizationName:    controller.GetOrganizationName(deployableArtifact),
			labels.LabelKeyProjectName:         controller.GetProjectName(deployableArtifact),
			labels.LabelKeyComponentName:       controller.GetComponentName(deployableArtifact),
			labels.LabelKeyDeploymentTrackName: controller.GetDeploymentTrackName(deployableArtifact),
		},
		client.MatchingFields{
			deployableArtifactRefIndexKey: deployableArtifact.Name,
		},
	}

	if err := r.List(ctx, deploymentList, listOpts...); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Deployments not found. Continuing with deletion.")
			return true, nil
		}
		return false, fmt.Errorf("failed to list deployments: %w", err)
	}

	pendingDeletion := false

	// Check if any deployments still exist
	if len(deploymentList.Items) > 0 {
		// Process each Deployment
		for i := range deploymentList.Items {
			deployment := &deploymentList.Items[i]

			// Check if the deployment is already being deleted
			if !deployment.DeletionTimestamp.IsZero() {
				// Still in the process of being deleted
				pendingDeletion = true
				logger.Info("Deployment is still being deleted", "name", deployment.Name)
				continue
			}

			// If not being deleted, trigger deletion
			logger.Info("Deleting deployment", "name", deployment.Name)
			if err := r.Delete(ctx, deployment); err != nil {
				if errors.IsNotFound(err) {
					logger.Info("Deployment already deleted", "name", deployment.Name)
					continue
				}
				return false, fmt.Errorf("failed to delete deployment %s: %w", deployment.Name, err)
			}

			// Mark as pending since we just triggered deletion
			pendingDeletion = true
		}

		// If there are still deployments being deleted, go to next iteration to check again later
		if pendingDeletion {
			return false, nil
		}
	}

	logger.Info("All deployments are deleted")
	return true, nil
}
