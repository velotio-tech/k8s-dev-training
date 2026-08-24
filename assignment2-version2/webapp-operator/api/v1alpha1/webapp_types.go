/*
Drop this in place of the generated api/v1alpha1/webapp_types.go
after `kubebuilder create api --group apps --version v1alpha1 --kind WebApp`.

Every +kubebuilder:... comment directly above a field or type is a MARKER.
controller-gen reads these comments (yes, actual Go comments) and turns
them into the OpenAPI validation schema inside the generated CRD YAML.
This is the single most important thing to understand about kubebuilder:
your comments are not documentation, they are the source of truth for
validation.
*/
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// WebAppSpec defines the desired state of WebApp.
// Everything in here is written by the user (or a GitOps pipeline) and
// read by the controller. Never written by the controller.
type WebAppSpec struct {
	// Image is the container image to deploy.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[a-zA-Z0-9._/-]+:[a-zA-Z0-9._-]+$`
	Image string `json:"image"`

	// Replicas is the desired pod count.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// +kubebuilder:default=1
	Replicas int32 `json:"replicas,omitempty"`

	// Port is the container port the app listens on.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +kubebuilder:default=8080
	Port int32 `json:"port,omitempty"`

	// Environment restricts which deployment tier this WebApp targets.
	// +kubebuilder:validation:Enum=Development;Staging;Production
	// +kubebuilder:default=Development
	Environment string `json:"environment,omitempty"`
}

// WebAppStatus defines the observed state of WebApp.
// Everything in here is written by the controller, read by users/tools.
// Never written directly by a user through normal apply.
type WebAppStatus struct {
	// AvailableReplicas is how many replicas are actually Ready right now.
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`

	// Phase is a short human-readable state machine value.
	// +kubebuilder:validation:Enum=Pending;Running;Failed
	Phase string `json:"phase,omitempty"`

	// LastUpdated records when the controller last touched this status.
	LastUpdated string `json:"lastUpdated,omitempty"`

	// Message carries a free-text explanation, useful for `kubectl describe`.
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
// +kubebuilder:printcolumn:name="Replicas",type=integer,JSONPath=`.spec.replicas`
// +kubebuilder:printcolumn:name="Environment",type=string,JSONPath=`.spec.environment`

// WebApp is the Schema for the webapps API.
type WebApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebAppSpec   `json:"spec,omitempty"`
	Status WebAppStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WebAppList contains a list of WebApp.
type WebAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WebApp{}, &WebAppList{})
}
