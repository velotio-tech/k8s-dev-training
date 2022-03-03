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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MultiLevelSpec defines the desired state of MultiLevel
type MultiLevelSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:MinLength=0
	// The name of the resource
	Name string `json:"name"`

	// +kubebuilder:validation:
	// The name of the image to be deployed
	Image string `json:"image"`
	// Command string `json:"Command"`

	// kubebuilder:validation:Minimum=1
	// The number of replicas of the image
	Replicas int `json:"replicas"`
}

// MultiLevelStatus defines the observed state of MultiLevel
type MultiLevelStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	State string `json:"state"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MultiLevel is the Schema for the multilevels API
type MultiLevel struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MultiLevelSpec   `json:"spec,omitempty"`
	Status MultiLevelStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MultiLevelList contains a list of MultiLevel
type MultiLevelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MultiLevel `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MultiLevel{}, &MultiLevelList{})
}
