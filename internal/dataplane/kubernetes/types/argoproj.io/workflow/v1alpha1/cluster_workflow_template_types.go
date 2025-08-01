package v1alpha1

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterWorkflowTemplate is the definition of a workflow template resource in cluster scope
// +genclient
// +genclient:noStatus
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterWorkflowTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Spec              WorkflowSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
}

type ClusterWorkflowTemplates []ClusterWorkflowTemplate

func (w ClusterWorkflowTemplates) Len() int {
	return len(w)
}

func (w ClusterWorkflowTemplates) Less(i, j int) bool {
	return strings.Compare(w[j].Name, w[i].Name) > 0
}

func (w ClusterWorkflowTemplates) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

// ClusterWorkflowTemplateList is list of ClusterWorkflowTemplate resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterWorkflowTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           ClusterWorkflowTemplates `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// GetTemplateByName retrieves a defined template by its name
func (cwftmpl *ClusterWorkflowTemplate) GetTemplateByName(name string) *Template {
	for _, t := range cwftmpl.Spec.Templates {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

// GetPodMetadata returns the PodMetadata of cluster workflow template.
func (cwftmpl *ClusterWorkflowTemplate) GetPodMetadata() *Metadata {
	return cwftmpl.Spec.PodMetadata
}

// GetWorkflowSpec returns the WorkflowSpec of cluster workflow template.
func (cwftmpl *ClusterWorkflowTemplate) GetWorkflowSpec() *WorkflowSpec {
	return &cwftmpl.Spec
}
