/*
Copyright 2022.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Container struct {

	// This field specify the selecltor for image template.
	Selector string `json:"selector,omitempty"`

	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:MinLength=1

	// Field name specify the name of container
	Name string `json:"name,omitempty"`

	// Image fields pulll the image from specify registery for deployment
	Image string `json:"image,omitempty"`

	// Specify the port on container that expose to the outside
	Port int `json:"port,omitempty"`
}

// DeploymentSpec defines the desired state of Deployment
type DeploymentSpec struct {
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10

	// Total number of replicas to be create. default (1)
	Replicas int32 `json:"replicas,omitempty"`

	// Selector to choose the replicas pods
	Selector string `json:"selector,omitempty"`

	Container Container `json:"container,omitempty"`

	// +kubebuilder:validation:Required

	// Specify the group version kind of the CRD
	GroupVersionKind GroupVersionKind `json:"gvk,omitempty"`
}

type GroupVersionKind struct {
	// +kubebuilder:validation:Required
	Group string `json:"group"`

	// +kubebuilder:validation:Required
	Version string `json:"version"`

	// +kubebuilder:validation:Required
	Kind string `json:"kind"`
}

// DeploymentStatus defines the observed state of Deployment
type DeploymentStatus struct {
	Replicas      int32 `json:"replicas,omitempty"`
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Deployment is the Schema for the deployments API
type Deployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              DeploymentSpec   `json:"spec,omitempty"`
	Status            DeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DeploymentList contains a list of Deployment
type DeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Deployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Deployment{}, &DeploymentList{})
}
