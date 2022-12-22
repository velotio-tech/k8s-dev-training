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

// AppdeployerSpec defines the desired state of Appdeployer
type AppdeployerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Fields of Appdeployer
	Name        string `json:"name"`
	Replicas    int    `json:"replicas"`
	Image       string `json:"image"`
	ServiceType string `json:"service-type"`
	Port        int    `json:"port,omitempty"`

	// GroupVersionKind - for the nested resource.
	//+kubebuilder:validation:XEmbeededResource
	GVK GroupVersionKind `json:"gvk"`
}

// AppdeployerStatus defines the observed state of Appdeployer
type AppdeployerStatus struct {
	// use of markers
	// markers more info - https://book.kubebuilder.io/reference/markers.html
	//+kubebuilder:validation:Enum:=created;deleted
	AppProgress string `json:"app-progress,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Appdeployer is the Schema for the appdeployers API
type Appdeployer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppdeployerSpec   `json:"spec,omitempty"`
	Status AppdeployerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AppdeployerList contains a list of Appdeployer
type AppdeployerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Appdeployer `json:"items"`
}

type GroupVersionKind struct {
	Version string `json:"version"`
	Kind    string `json:"kind"`
}

func init() {
	SchemeBuilder.Register(&Appdeployer{}, &AppdeployerList{})
}