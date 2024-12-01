/*
Copyright 2021.

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

package v2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MyResourceSpec defines the desired state of MyResource

type MyResourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of MyResource. Edit myresource_types.go to remove/update
	//Foo string `json:"foo,omitempty"`
	// +kubebuilder:validation:MinLength=3
	JobName string `json:"jobName"`
	// +kubebuilder:validation:MinLength=1
	Command string `json:"command"`
	// +kubebuilder:validation:Pattern="[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])T(2[0-3]|[01][0-9]):[0-5][0-9]:[0-9][0-9]Z"
	Schedule string `json:"schedule"`

}

// MyResourceStatus defines the observed state of MyResource
type MyResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	JobState JobState `json:"jobState"`
}
type JobState string
const (
	Pending JobState = "Pending"
	Running JobState = "Running"
	Finished JobState = "Finished"
)
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MyResource is the Schema for the myresources API
type MyResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyResourceSpec   `json:"spec,omitempty"`
	Status MyResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MyResourceList contains a list of MyResource
type MyResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyResource{}, &MyResourceList{})
}
