package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MultiTierAppSpec defines the desired state of MultiTierApp
type MultiTierAppSpec struct {
	// Replicas for Level 2 Deployment
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// TargetNamespace restricts management to a single namespace
	TargetNamespace string `json:"targetNamespace"`

	// JobCommand is the script/command Level 3 Job will execute
	JobCommand string `json:"jobCommand"`
}

// MultiTierAppStatus defines the observed state of MultiTierApp
type MultiTierAppStatus struct {
	Phase          string `json:"phase,omitempty"`
	DeploymentName string `json:"deploymentName,omitempty"`
	JobName        string `json:"jobName,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

type MultiTierApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MultiTierAppSpec   `json:"spec,omitempty"`
	Status MultiTierAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type MultiTierAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MultiTierApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MultiTierApp{}, &MultiTierAppList{})
}
