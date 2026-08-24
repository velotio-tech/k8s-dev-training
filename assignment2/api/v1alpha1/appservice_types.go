package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AppServiceSpec defines the desired state of AppService
type AppServiceSpec struct {
	// Replicas is the desired number of instances (1 to 10).
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// Environment specifies deployment tier (dev, staging, prod).
	// +kubebuilder:validation:Enum=dev;staging;prod
	Environment string `json:"environment"`

	// Image is the container image to run.
	// +kubebuilder:validation:MinLength=5
	Image string `json:"image"`

	// Ports is a list of internal container ports (1 to 5 items).
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=5
	// +optional
	Ports []int32 `json:"ports,omitempty"`
}

// AppServiceStatus defines the observed state of AppService
type AppServiceStatus struct {
	// ReadyReplicas represents current active instances.
	// +kubebuilder:validation:Minimum=0
	ReadyReplicas int32 `json:"readyReplicas"`

	// Phase tracks lifecycle state (Pending, Running, Failed).
	// +kubebuilder:validation:Enum=Pending;Running;Failed
	// +kubebuilder:default=Pending
	Phase string `json:"phase"`

	// ActiveEndpoints lists reachable IP addresses.
	// +optional
	ActiveEndpoints []string `json:"activeEndpoints,omitempty"`

	// LastSyncTime records the last controller reconciliation timestamp.
	// +optional
	LastSyncTime *metav1.Time `json:"lastSyncTime,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// AppService is the Schema for the appservices API
type AppService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppServiceSpec   `json:"spec,omitempty"`
	Status AppServiceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AppServiceList contains a list of AppService
type AppServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AppService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AppService{}, &AppServiceList{})
}