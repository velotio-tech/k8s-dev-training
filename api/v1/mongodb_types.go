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
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MongoDBSpec defines the desired state of MongoDB
type MongoDBSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// GVK denotes GroupVersionKind for the nested resource.
	//+kubebuilder:validation:XEmbeddedResource
	GVK GVK `json:"gvk"`
	// InitUser is a username field used to initialize the database
	InitUser string `json:"init_user,omitempty"`
	// InitPassword is a password field used to initialize the database
	InitPassword string `json:"init_password,omitempty"`
	//+kubebuilder:validation:default:=1
	//+kubebuilder:validation:Optional
	MaxUsers int `json:"max_users"`
	//+kubebuilder:validation:Maximum:=10000
	//+kubebuilder:validation:Minimum:=100
	MaxConcurrentConnections int `json:"max_concurrent_connections"`
}

// MongoDBStatus defines the observed state of MongoDB
type MongoDBStatus struct {
	//+kubebuilder:validation:Enum:=healthy;unhealthy
	Condition string `json:"condition,omitempty"`
	//+kubebuilder:validation:Enum:=running;terminating;pending
	Phase string `json:"phase,omitempty"`
}

// GVK unambiguously identifies a kind.  It doesn't anonymously include GroupVersion
// to avoid automatic coercion.  It doesn't use a GroupVersion to avoid custom marshalling
type GVK struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

func (gvk *GVK) GroupVersion() schema.GroupVersion {
	result, _ := schema.ParseGroupVersion(gvk.APIVersion)
	return result
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MongoDB is the Schema for the mongodbs API
type MongoDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MongoDBSpec   `json:"spec,omitempty"`
	Status MongoDBStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MongoDBList contains a list of MongoDB
type MongoDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoDB{}, &MongoDBList{})
}
