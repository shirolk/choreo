// Copyright 2025 The OpenChoreo Authors
// SPDX-License-Identifier: Apache-2.0

package create

import (
	"strings"

	"github.com/spf13/cobra"

	openchoreov1alpha1 "github.com/openchoreo/openchoreo/api/v1alpha1"
	"github.com/openchoreo/openchoreo/pkg/cli/common/builder"
	"github.com/openchoreo/openchoreo/pkg/cli/common/constants"
	"github.com/openchoreo/openchoreo/pkg/cli/flags"
	"github.com/openchoreo/openchoreo/pkg/cli/types/api"
)

// Helper functions for common flag sets
func getBasicFlags() []flags.Flag {
	return []flags.Flag{
		flags.Name,
		flags.Interactive,
	}
}

func getOrgScopedFlags() []flags.Flag {
	return append(getBasicFlags(),
		flags.Organization,
	)
}

func getProjectLevelFlags() []flags.Flag {
	return append(getOrgScopedFlags(),
		flags.Project,
	)
}

func getComponentLevelFlags() []flags.Flag {
	return append(getProjectLevelFlags(),
		flags.Component,
	)
}

func getMetadataFlags() []flags.Flag { //nolint:unused // Used by temporarily disabled create commands
	return append(getBasicFlags(),
		flags.DisplayName,
		flags.Description,
	)
}

func NewCreateCmd(impl api.CommandImplementationInterface) *cobra.Command {
	createCmd := &cobra.Command{
		Use:   constants.Create.Use,
		Short: constants.Create.Short,
		Long:  constants.Create.Long,
	}

	createCmd.AddCommand(
		// newCreateOrganizationCmd(impl),
		// newCreateProjectCmd(impl),
		// newCreateComponentCmd(impl),
		// newCreateBuildCmd(impl),
		// newCreateDeploymentCmd(impl),
		// newCreateDataPlaneCmd(impl),
		// newCreateDeploymentTrackCmd(impl),
		// newCreateEnvironmentCmd(impl),
		// newCreateDeployableArtifactCmd(impl),
		// newCreateDeploymentPipelineCmd(impl),
		newCreateWorkloadCmd(impl),
	)

	return createCmd
}

//nolint:unused // Temporarily disabled
func newCreateOrganizationCmd(impl api.CommandImplementationInterface) *cobra.Command {
	return (&builder.CommandBuilder{
		Command: constants.CreateOrganization,
		Flags:   getMetadataFlags(),
		RunE: func(fg *builder.FlagGetter) error {
			return impl.CreateOrganization(api.CreateOrganizationParams{
				Name:        fg.GetString(flags.Name),
				DisplayName: fg.GetString(flags.DisplayName),
				Description: fg.GetString(flags.Description),
				Interactive: fg.GetBool(flags.Interactive),
			})
		},
	}).Build()
}

//nolint:unused // Temporarily disabled
func newCreateProjectCmd(impl api.CommandImplementationInterface) *cobra.Command {
	projectFlags := append(getOrgScopedFlags(),
		flags.DisplayName,
		flags.Description,
	)
	return (&builder.CommandBuilder{
		Command: constants.CreateProject,
		Flags:   projectFlags,
		RunE: func(fg *builder.FlagGetter) error {
			return impl.CreateProject(api.CreateProjectParams{
				Name:               fg.GetString(flags.Name),
				Organization:       fg.GetString(flags.Organization),
				DisplayName:        fg.GetString(flags.DisplayName),
				Description:        fg.GetString(flags.Description),
				DeploymentPipeline: fg.GetString(flags.DeploymentPipeline),
				Interactive:        fg.GetBool(flags.Interactive),
			})
		},
	}).Build()
}

//nolint:unused // Temporarily disabled
func newCreateComponentCmd(impl api.CommandImplementationInterface) *cobra.Command {
	componentFlags := append(getProjectLevelFlags(),
		flags.DisplayName,
		flags.GitRepositoryURL,
		flags.ComponentType,
		flags.DockerContext,
		flags.DockerfilePath,
		flags.Branch,
		flags.BuildpackName,
		flags.BuildpackVersion,
		flags.Path,
	)
	return (&builder.CommandBuilder{
		Command: constants.CreateComponent,
		Flags:   componentFlags,
		RunE: func(fg *builder.FlagGetter) error {
			return impl.CreateComponent(api.CreateComponentParams{
				Name:             fg.GetString(flags.Name),
				Organization:     fg.GetString(flags.Organization),
				Project:          fg.GetString(flags.Project),
				DisplayName:      fg.GetString(flags.DisplayName),
				GitRepositoryURL: fg.GetString(flags.GitRepositoryURL),
				Type:             openchoreov1alpha1.ComponentType(fg.GetString(flags.ComponentType)),
				Interactive:      fg.GetBool(flags.Interactive),
				Branch:           fg.GetString(flags.Branch),
				Path:             fg.GetString(flags.Path),
				DockerFile:       fg.GetString(flags.DockerfilePath),
				DockerContext:    fg.GetString(flags.DockerContext),
				BuildpackName:    fg.GetString(flags.BuildpackName),
				BuildpackVersion: fg.GetString(flags.BuildpackVersion),
			})
		},
	}).Build()
}

func newCreateWorkloadCmd(impl api.CommandImplementationInterface) *cobra.Command {
	workloadFlags := append(getComponentLevelFlags(),
		flags.Image,
		flags.Output,
		flags.WorkloadDescriptor,
	)

	return (&builder.CommandBuilder{
		Command: constants.CreateWorkload,
		Flags:   workloadFlags,
		RunE: func(fg *builder.FlagGetter) error {
			return impl.CreateWorkload(api.CreateWorkloadParams{
				FilePath:         fg.GetString(flags.WorkloadDescriptor),
				OrganizationName: fg.GetString(flags.Organization),
				ProjectName:      fg.GetString(flags.Project),
				ComponentName:    fg.GetString(flags.Component),
				ImageURL:         fg.GetString(flags.Image),
				OutputPath:       fg.GetString(flags.Output),
				Interactive:      fg.GetBool(flags.Interactive),
			})
		},
	}).Build()
}

//nolint:unused // Temporarily disabled
func newCreateBuildCmd(impl api.CommandImplementationInterface) *cobra.Command {
	buildFlags := append(getComponentLevelFlags(),
		flags.DockerContext,
		flags.DockerfilePath,
		flags.BuildpackName,
		flags.BuildpackVersion,
		flags.Branch,
		flags.Path,
		flags.Revision,
		flags.AutoBuild,
		flags.DeploymentTrack,
	)

	return (&builder.CommandBuilder{
		Command: constants.CreateBuild,
		Flags:   buildFlags,
		RunE: func(fg *builder.FlagGetter) error {
			return impl.CreateBuild(api.CreateBuildParams{
				Name:         fg.GetString(flags.Name),
				Organization: fg.GetString(flags.Organization),
				Project:      fg.GetString(flags.Project),
				Component:    fg.GetString(flags.Component),
				Branch:       fg.GetString(flags.Branch),
				Path:         fg.GetString(flags.Path),
				Revision:     fg.GetString(flags.Revision),
				AutoBuild:    fg.GetBool(flags.AutoBuild),
				Interactive:  fg.GetBool(flags.Interactive),
				Docker: &openchoreov1alpha1.DockerConfiguration{
					Context:        fg.GetString(flags.DockerContext),
					DockerfilePath: fg.GetString(flags.DockerfilePath),
				},
				Buildpack: &openchoreov1alpha1.BuildpackConfiguration{
					Name:    openchoreov1alpha1.BuildpackName(fg.GetString(flags.BuildpackName)),
					Version: fg.GetString(flags.BuildpackVersion),
				},
				DeploymentTrack: fg.GetString(flags.DeploymentTrack),
			})
		},
	}).Build()
}

//nolint:unused // Temporarily disabled
func newCreateDeploymentCmd(impl api.CommandImplementationInterface) *cobra.Command {
	deployFlags := append(getComponentLevelFlags(),
		flags.Environment,
		flags.DeploymentTrack,
		flags.DeployableArtifact,
	)
	return (&builder.CommandBuilder{
		Command: constants.CreateDeployment,
		Flags:   deployFlags,
		RunE: func(fg *builder.FlagGetter) error {
			return impl.CreateDeployment(api.CreateDeploymentParams{
				Name:               fg.GetString(flags.Name),
				Organization:       fg.GetString(flags.Organization),
				Project:            fg.GetString(flags.Project),
				Component:          fg.GetString(flags.Component),
				Environment:        fg.GetString(flags.Environment),
				DeploymentTrack:    fg.GetString(flags.DeploymentTrack),
				DeployableArtifact: fg.GetString(flags.DeployableArtifact),
				Interactive:        fg.GetBool(flags.Interactive),
			})
		},
	}).Build()
}

//nolint:unused // Temporarily disabled
func newCreateDataPlaneCmd(impl api.CommandImplementationInterface) *cobra.Command {
	dpFlags := append(getMetadataFlags(),
		flags.KubernetesClusterName,
		flags.EnableCilium,
		flags.EnableScaleToZero,
		flags.GatewayType,
		flags.PublicVirtualHost,
		flags.OrgVirtualHost,
		flags.Organization,
		flags.APIServerURL,
		flags.CaCert,
		flags.ClientCert,
		flags.ClientKey,
	)
	return (&builder.CommandBuilder{
		Command: constants.CreateDataPlane,
		Flags:   dpFlags,
		RunE: func(fg *builder.FlagGetter) error {
			return impl.CreateDataPlane(api.CreateDataPlaneParams{
				Name:                    fg.GetString(flags.Name),
				Organization:            fg.GetString(flags.Organization),
				KubernetesClusterName:   fg.GetString(flags.KubernetesClusterName),
				EnableCilium:            fg.GetBool(flags.EnableCilium),
				EnableScaleToZero:       fg.GetBool(flags.EnableScaleToZero),
				GatewayType:             fg.GetString(flags.GatewayType),
				PublicVirtualHost:       fg.GetString(flags.PublicVirtualHost),
				OrganizationVirtualHost: fg.GetString(flags.OrgVirtualHost),
				Interactive:             fg.GetBool(flags.Interactive),
				APIServerURL:            fg.GetString(flags.APIServerURL),
				CACert:                  fg.GetString(flags.CaCert),
				ClientCert:              fg.GetString(flags.ClientCert),
				ClientKey:               fg.GetString(flags.ClientKey),
			})
		},
	}).Build()
}

//nolint:unused // Temporarily disabled
func newCreateDeploymentTrackCmd(impl api.CommandImplementationInterface) *cobra.Command {
	trackFlags := append(getComponentLevelFlags(),
		flags.APIVersion,
		flags.AutoDeploy,
	)
	return (&builder.CommandBuilder{
		Command: constants.CreateDeploymentTrack,
		Flags:   trackFlags,
		RunE: func(fg *builder.FlagGetter) error {
			return impl.CreateDeploymentTrack(api.CreateDeploymentTrackParams{
				Name:         fg.GetString(flags.Name),
				Organization: fg.GetString(flags.Organization),
				Project:      fg.GetString(flags.Project),
				Component:    fg.GetString(flags.Component),
				APIVersion:   fg.GetString(flags.APIVersion),
				AutoDeploy:   fg.GetBool(flags.AutoDeploy),
				Interactive:  fg.GetBool(flags.Interactive),
			})
		},
	}).Build()
}

//nolint:unused // Temporarily disabled
func newCreateEnvironmentCmd(impl api.CommandImplementationInterface) *cobra.Command {
	envFlags := append(getOrgScopedFlags(),
		flags.DisplayName,
		flags.Description,
		flags.IsProduction,
		flags.DNSPrefix,
		flags.DataPlaneRef,
	)
	return (&builder.CommandBuilder{
		Command: constants.CreateEnvironment,
		Flags:   envFlags,
		RunE: func(fg *builder.FlagGetter) error {
			return impl.CreateEnvironment(api.CreateEnvironmentParams{
				Name:         fg.GetString(flags.Name),
				Organization: fg.GetString(flags.Organization),
				DisplayName:  fg.GetString(flags.DisplayName),
				Description:  fg.GetString(flags.Description),
				Interactive:  fg.GetBool(flags.Interactive),
				DataPlaneRef: fg.GetString(flags.DataPlaneRef),
				IsProduction: fg.GetBool(flags.IsProduction),
				DNSPrefix:    fg.GetString(flags.DNSPrefix),
			})
		},
	}).Build()
}

//nolint:unused // Temporarily disabled
func newCreateDeployableArtifactCmd(impl api.CommandImplementationInterface) *cobra.Command {
	artifactFlags := append(getComponentLevelFlags(),
		flags.DeploymentTrack,
		flags.Build,
	)
	return (&builder.CommandBuilder{
		Command: constants.CreateDeployableArtifact,
		Flags:   artifactFlags,
		RunE: func(fg *builder.FlagGetter) error {
			return impl.CreateDeployableArtifact(api.CreateDeployableArtifactParams{
				Name:            fg.GetString(flags.Name),
				Organization:    fg.GetString(flags.Organization),
				Project:         fg.GetString(flags.Project),
				Component:       fg.GetString(flags.Component),
				DeploymentTrack: fg.GetString(flags.DeploymentTrack),
				Interactive:     fg.GetBool(flags.Interactive),
			})
		},
	}).Build()
}

//nolint:unused // Temporarily disabled
func newCreateDeploymentPipelineCmd(impl api.CommandImplementationInterface) *cobra.Command {
	dpFlags := []flags.Flag{
		flags.Organization,
		flags.Name,
		flags.EnvironmentOrder,
	}

	return (&builder.CommandBuilder{
		Command: constants.CreateDeploymentPipeline,
		Flags:   dpFlags,
		RunE: func(fg *builder.FlagGetter) error {
			// Get environment order from flag
			envOrderStr := fg.GetString(flags.EnvironmentOrder)
			var environmentOrder []string
			if envOrderStr != "" {
				environmentOrder = strings.Split(envOrderStr, ",")
				for i := range environmentOrder {
					environmentOrder[i] = strings.TrimSpace(environmentOrder[i])
				}
			}

			return impl.CreateDeploymentPipeline(api.CreateDeploymentPipelineParams{
				Name:             fg.GetString(flags.Name),
				Organization:     fg.GetString(flags.Organization),
				EnvironmentOrder: environmentOrder,
			})
		},
	}).Build()
}
