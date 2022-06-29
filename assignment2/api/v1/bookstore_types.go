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

// BookStoreSpec defines the desired state of BookStore
type BookStoreSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name of the book store
	// +kubebuilder:validation:MinLength=4
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Required
	// Location of the book store
	Location string `json:"location"`
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:Required
	// Owner of the book store
	Owner string `json:"owner"`
	// +kubebuilder:validation:MinLength=10
	// +kubebuilder:validation:MaxLength=10
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`
	// Date of establishment of the book store (YYYY-MM-DD)
	Established string `json:"established"`
}

// BookStoreStatus defines the observed state of BookStore
type BookStoreStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +kubebuilder:validation:Minimum=0
	CurrentBooksCount int32  `json:"currentBooksCount"`
	LastBookAdded     string `json:"lastBookAdded"`
	LastUpdateDate    string `json:"lastUpdateDate"`
	LastUpdateTime    string `json:"lastUpdateTime"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// BookStore is the Schema for the bookstores API
type BookStore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BookStoreSpec   `json:"spec,omitempty"`
	Status BookStoreStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BookStoreList contains a list of BookStore
type BookStoreList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BookStore `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BookStore{}, &BookStoreList{})
}
