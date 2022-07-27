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
	gvk "github.com/KnVerey/kustomize/pkg/gvk"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DummyDBSpec defines the desired state of DummyDB
type DummyDBSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas,omitempty"`

	// +kubebuilder:validation:Required
	Size resource.Quantity `json:"size,omitempty"`

	// +kubebuilder:validation:Required
	Gvk gvk.Gvk `json:"gvk,omitempty"`

	// +kubebuilder:validation:MinLength=10
	// +kubebuilder:validation:Required
	AdminPassword string `json:"adminPassword,omitempty"`
}

// DummyDBStatus defines the observed state of DummyDB
type DummyDBStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Minimum=0
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// +kubebuilder:validation:Minimum=0
	AvailabeSize int32 `json:"availabeSize,omitempty"`

	VolumeResizingInProgress bool `json:"volumeResizingInProgress,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DummyDB is the Schema for the dummydbs API
type DummyDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DummyDBSpec   `json:"spec,omitempty"`
	Status DummyDBStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DummyDBList contains a list of DummyDB
type DummyDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DummyDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DummyDB{}, &DummyDBList{})
}
