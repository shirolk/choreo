// Copyright 2025 The OpenChoreo Authors
// SPDX-License-Identifier: Apache-2.0

package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	openchoreov1alpha1 "github.com/openchoreo/openchoreo/api/v1alpha1"
	kubernetesClient "github.com/openchoreo/openchoreo/internal/clients/kubernetes"
	buildcontroller "github.com/openchoreo/openchoreo/internal/controller/build"
	argo "github.com/openchoreo/openchoreo/internal/dataplane/kubernetes/types/argoproj.io/workflow/v1alpha1"
	"github.com/openchoreo/openchoreo/internal/labels"
	"github.com/openchoreo/openchoreo/internal/openchoreo-api/models"
)

// BuildService handles build-related business logic
type BuildService struct {
	k8sClient         client.Client
	logger            *slog.Logger
	buildPlaneService *BuildPlaneService
	bpClientMgr       *kubernetesClient.KubeMultiClientManager
}

// NewBuildService creates a new build service
func NewBuildService(k8sClient client.Client, buildPlaneService *BuildPlaneService, bpClientMgr *kubernetesClient.KubeMultiClientManager, logger *slog.Logger) *BuildService {
	return &BuildService{
		k8sClient:         k8sClient,
		logger:            logger,
		buildPlaneService: buildPlaneService,
		bpClientMgr:       bpClientMgr,
	}
}

// ListBuildTemplates retrieves cluster workflow templates (argo) available for an organization in the buildplane
func (s *BuildService) ListBuildTemplates(ctx context.Context, orgName string) ([]models.BuildTemplateResponse, error) {
	s.logger.Debug("Listing build templates", "org", orgName)

	// Get the build plane Kubernetes client
	buildPlaneClient, err := s.buildPlaneService.GetBuildPlaneClient(ctx, orgName)
	if err != nil {
		return nil, fmt.Errorf("failed to get build plane client: %w", err)
	}

	// List ClusterWorkflowTemplates using the build plane client
	var clusterWorkflowTemplates argo.ClusterWorkflowTemplateList
	err = buildPlaneClient.List(ctx, &clusterWorkflowTemplates)
	if err != nil {
		s.logger.Error("Failed to list ClusterWorkflowTemplates", "error", err)
		return nil, fmt.Errorf("failed to list ClusterWorkflowTemplates: %w", err)
	}

	s.logger.Debug("Found build templates", "count", len(clusterWorkflowTemplates.Items), "org", orgName)

	templateResponses := make([]models.BuildTemplateResponse, 0, len(clusterWorkflowTemplates.Items))
	for _, template := range clusterWorkflowTemplates.Items {
		parameters := make([]models.BuildTemplateParameter, 0, len(template.Spec.Arguments.Parameters))
		if template.Spec.Arguments.Parameters != nil {
			for _, param := range template.Spec.Arguments.Parameters {
				templateParam := models.BuildTemplateParameter{
					Name: param.Name,
				}

				if param.Default != nil {
					templateParam.Default = string(*param.Default)
				}

				parameters = append(parameters, templateParam)
			}
		}

		templateResponse := models.BuildTemplateResponse{
			Name:       template.Name,
			Parameters: parameters,
			CreatedAt:  template.CreationTimestamp.Time,
		}

		templateResponses = append(templateResponses, templateResponse)
	}

	return templateResponses, nil
}

// TriggerBuild creates a new build from a component
func (s *BuildService) TriggerBuild(ctx context.Context, orgName, projectName, componentName, commit string) (*models.BuildResponse, error) {
	s.logger.Debug("Triggering build", "org", orgName, "project", projectName, "component", componentName, "commit", commit)

	// Retrieve component and use that to create the build
	var component openchoreov1alpha1.Component
	err := s.k8sClient.Get(ctx, client.ObjectKey{
		Name:      componentName,
		Namespace: orgName,
	}, &component)

	if err != nil {
		s.logger.Error("Failed to get component", "error", err, "org", orgName, "project", projectName, "component", componentName)
		return nil, fmt.Errorf("failed to get component: %w", err)
	}

	buildUUID := uuid.New().String()
	buildID := strings.ReplaceAll(buildUUID[:8], "-", "")

	buildName := fmt.Sprintf("%s-build-%s", componentName, buildID)

	build := &openchoreov1alpha1.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name:      buildName,
			Namespace: orgName,
			Labels: map[string]string{
				labels.LabelKeyOrganizationName: orgName,
				labels.LabelKeyProjectName:      projectName,
				labels.LabelKeyComponentName:    componentName,
			},
		},
		Spec: openchoreov1alpha1.BuildSpec{
			Owner: openchoreov1alpha1.BuildOwner{
				ProjectName:   projectName,
				ComponentName: componentName,
			},
			Repository: openchoreov1alpha1.Repository{
				URL: component.Spec.Build.Repository.URL,
				Revision: openchoreov1alpha1.Revision{
					Branch: component.Spec.Build.Repository.Revision.Branch,
					Commit: commit,
				},
				AppPath: component.Spec.Build.Repository.AppPath,
			},
			TemplateRef: component.Spec.Build.TemplateRef,
		},
	}

	err = s.k8sClient.Create(ctx, build)
	if err != nil {
		s.logger.Error("Failed to create build", "error", err)
		return nil, fmt.Errorf("failed to create build: %w", err)
	}

	s.logger.Info("Build created successfully", "build", buildName)

	if commit == "" {
		commit = "latest"
	}

	return &models.BuildResponse{
		Name:          buildName,
		UUID:          string(build.UID),
		ComponentName: componentName,
		ProjectName:   projectName,
		OrgName:       orgName,
		Commit:        commit,
		Status:        "Created",
		CreatedAt:     build.CreationTimestamp.Time,
	}, nil
}

// ListBuilds retrieves builds for a component using spec.owner fields instead of labels
func (s *BuildService) ListBuilds(ctx context.Context, orgName, projectName, componentName string) ([]models.BuildResponse, error) {
	s.logger.Debug("Listing builds", "org", orgName, "project", projectName, "component", componentName)

	var builds openchoreov1alpha1.BuildList
	err := s.k8sClient.List(ctx, &builds, client.InNamespace(orgName))
	if err != nil {
		s.logger.Error("Failed to list builds", "error", err)
		return nil, fmt.Errorf("failed to list builds: %w", err)
	}

	buildResponses := make([]models.BuildResponse, 0, len(builds.Items))
	for _, build := range builds.Items {
		// Filter by spec.owner fields instead of labels
		if build.Spec.Owner.ProjectName != projectName || build.Spec.Owner.ComponentName != componentName {
			continue
		}

		// This commit hash should always be there since the build is triggered with a commit
		// If not provided, we can default to "latest" for now.
		commit := build.Spec.Repository.Revision.Commit
		if commit == "" {
			commit = "latest"
		}

		buildResponses = append(buildResponses, models.BuildResponse{
			Name:          build.Name,
			UUID:          string(build.UID),
			ComponentName: componentName,
			ProjectName:   projectName,
			OrgName:       orgName,
			Commit:        commit,
			Status:        GetLatestBuildStatus(build.Status.Conditions),
			CreatedAt:     build.CreationTimestamp.Time,
			Image:         build.Status.ImageStatus.Image,
		})
	}

	return buildResponses, nil
}

func GetLatestBuildStatus(buildConditions []metav1.Condition) string {
	if len(buildConditions) == 0 {
		return statusUnknown
	}

	// Define the order of priority for build conditions (latest to earliest)
	// WorkloadUpdated > BuildCompleted > BuildTriggered > BuildInitiated
	conditionOrder := []string{
		string(buildcontroller.ConditionWorkloadUpdated),
		string(buildcontroller.ConditionBuildCompleted),
		string(buildcontroller.ConditionBuildTriggered),
		string(buildcontroller.ConditionBuildInitiated),
	}

	// Find the latest condition based on priority order
	for _, conditionType := range conditionOrder {
		for _, condition := range buildConditions {
			if condition.Type == conditionType {
				if condition.Type == string(buildcontroller.ConditionWorkloadUpdated) && condition.Status == metav1.ConditionTrue {
					return "Completed"
				}
				return condition.Reason
			}
		}
	}

	return statusUnknown
}
