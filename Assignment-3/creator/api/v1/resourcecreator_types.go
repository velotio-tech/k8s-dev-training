/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ResourceSpec defines the desired resource to be managed
type ResourceSpec struct {
	// Group is the API group of the resource to manage
	Group string `json:"group"`
	// Version is the version of the resource to manage
	Version string `json:"version"`
	// Kind is the kind of the resource to manage
	Kind string `json:"kind"`
	// Name is the name of the resource to manage
	Name string `json:"name"`
	// Spec is the spec of the resource to manage
	Spec apiextensionsv1.JSON `json:"spec"`
}

// ResourceCreatorSpec defines the desired state of ResourceCreator.
type ResourceCreatorSpec struct {
	// Resources is a list of resources to be managed by this controller
	// +kubebuilder:validation:MinItems=1
	Resources []ResourceSpec `json:"resources"`
}

// ResourceCreatorStatus defines the observed state of ResourceCreator.
type ResourceCreatorStatus struct {
	// Resource contains the resource specification
	Resource ResourceSpec `json:"resource"`
	// Status indicates if the resource is created/updated/error
	Status string `json:"status"`
	// LastUpdateTime is the last time the resource was updated
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
	// Message contains additional information about the resource status
	// +optional
	Message string `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="CreatedAt",type="date",JSONPath=".metadata.creationTimestamp"

// ResourceCreator is the Schema for the resourcecreators API.
type ResourceCreator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceCreatorSpec   `json:"spec,omitempty"`
	Status ResourceCreatorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ResourceCreatorList contains a list of ResourceCreator.
type ResourceCreatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceCreator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceCreator{}, &ResourceCreatorList{})
}
