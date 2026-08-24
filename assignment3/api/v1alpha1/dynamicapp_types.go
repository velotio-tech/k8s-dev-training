// step 3.2: Define the Schema
// defines target resources (GVKs and names) for dynamic management
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceTarget defines a dynamic GVK target resource to manage
type ResourceTarget struct {
	// Group is the API group (e.g., "", "apps")
	Group string `json:"group"`
	// Version is the API version (e.g., "v1")
	Version string `json:"version"`
	// Kind is the resource kind (e.g., "ConfigMap", "Service")
	Kind string `json:"kind"`
	// Name is the name of the child resource
	Name string `json:"name"`
}

// DynamicAppSpec defines the desired state of DynamicApp
type DynamicAppSpec struct {
	// Replicas for demo workload
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas"`

	// TargetResources holds GVK targets to create and manage
	// +optional
	TargetResources []ResourceTarget `json:"targetResources,omitempty"`
}

// DynamicAppStatus defines the observed state of DynamicApp
type DynamicAppStatus struct {
	// Phase of controller management
	Phase string `json:"phase,omitempty"`
	// ManagedResources lists formatted strings of created child resources
	ManagedResources []string `json:"managedResources,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// DynamicApp is the Schema for the dynamicapps API
type DynamicApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DynamicAppSpec   `json:"spec,omitempty"`
	Status DynamicAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DynamicAppList contains a list of DynamicApp
type DynamicAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DynamicApp `json:"items"`
}

// addKnownTypes adds the list of types defined in this package to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(GroupVersion,
		&DynamicApp{},
		&DynamicAppList{},
	)
	metav1.AddToGroupVersion(scheme, GroupVersion)
	return nil
}

func init() {
	SchemeBuilder.Register(addKnownTypes)
}
