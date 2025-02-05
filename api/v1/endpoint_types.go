/*
 * Copyright (c) 2025, WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
 *
 * WSO2 Inc. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EndpointServiceSpec defines the configuration of the upstream service
type EndpointServiceSpec struct {
	// URL of the upstream service
	URL string `json:"url,omitempty"`

	// Base path of the upstream service
	// +optional
	BasePath string `json:"basePath,omitempty"`

	// Port of the upstream service
	// +required
	Port int32 `json:"port"`
}

// EndpointSchemaSpec defines the schema configuration of the endpoint
type EndpointSchemaSpec struct {
	// File path of the schema relative to the component source code
	FilePath string `json:"filePath,omitempty"`

	// Inline content of the schema
	Content string `json:"content,omitempty"`
}

// EndpointAPISettingsSpec defines configuration parameters for managed endpoints
type EndpointAPISettingsSpec struct {
	SecuritySchemes                      []string                  `json:"securitySchemes,omitempty"`
	AuthorizationHeader                  string                    `json:"authorizationHeader,omitempty"`
	BackendJWT                           *BackendJWTConfig         `json:"backendJwt,omitempty"`
	OperationPolicies                    []OperationPolicy         `json:"operationPolicies,omitempty"`
	PublicVisibilityConfigurations       *VisibilityConfigurations `json:"publicVisibilityConfigurations,omitempty"`
	OrganizationVisibilityConfigurations *VisibilityConfigurations `json:"organizationVisibilityConfigurations,omitempty"`
}

// BackendJWTConfig defines JWT configuration for backend services
type BackendJWTConfig struct {
	Enabled       bool                    `json:"enabled"`
	Configuration BackendJWTConfigDetails `json:"configuration"`
}

// BackendJWTConfigDetails contains the detailed JWT configuration
type BackendJWTConfigDetails struct {
	Audiences []string `json:"audiences"`
}

// OperationPolicy defines authentication policy for an API operation
type OperationPolicy struct {
	Target             string `json:"target"`
	AuthenticationType string `json:"authenticationType"`
}

// VisibilityConfigurations defines configurations for different visibility levels
type VisibilityConfigurations struct {
	CORS      *CORSConfig      `json:"cors,omitempty"`
	RateLimit *RateLimitConfig `json:"rateLimit,omitempty"`
}

// CORSConfig defines Cross-Origin Resource Sharing configuration
type CORSConfig struct {
	Enabled       bool     `json:"enabled"`
	AllowOrigins  []string `json:"allowOrigins"`
	AllowMethods  []string `json:"allowMethods"`
	AllowHeaders  []string `json:"allowHeaders"`
	ExposeHeaders []string `json:"exposeHeaders"`
}

// RateLimitConfig defines rate limiting configuration
type RateLimitConfig struct {
	Tier string `json:"tier"`
}

// EndpointSpec defines the desired state of Endpoint
type EndpointSpec struct {
	// Type indicates the protocol of the endpoint
	// +kubebuilder:validation:Enum=HTTP;REST;gRPC;GraphQL;Websocket;TCP;UDP
	Type string `json:"type"`

	// Configuration of the upstream service
	// +required
	Service EndpointServiceSpec `json:"service"`

	// Schema of the endpoint if available
	// +optional
	Schema *EndpointSchemaSpec `json:"schema,omitempty"`

	// Network visibility levels that the endpoint is exposed
	// +optional
	NetworkVisibilities []string `json:"networkVisibilities,omitempty"`

	// Configuration parameters related to the managed endpoint
	// +optional
	APISettings *EndpointAPISettingsSpec `json:"apiSettings,omitempty"`

	// Configuration parameters related to the webapp gateway
	// +optional
	WebappGatewaySettings map[string]string `json:"webappGatewaySettings,omitempty"`
}

// EndpointStatus defines the observed state of Endpoint
type EndpointStatus struct {
	ObservedGeneration int64              `json:"observedGeneration,omitempty"`
	Conditions         []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// Endpoint is the Schema for the endpoints API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="url",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].message"
type Endpoint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EndpointSpec   `json:"spec,omitempty"`
	Status EndpointStatus `json:"status,omitempty"`
}

// EndpointList contains a list of Endpoint
// +kubebuilder:object:root=true
type EndpointList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Endpoint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Endpoint{}, &EndpointList{})
}
