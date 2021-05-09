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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
type ResourceType string
// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
const (
	Deployment ResourceType = "deployment"
	Pod  ResourceType ="pod"
)
// ReplicaManagerSpec defines the desired state of ReplicaManager
type ReplicaManagerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//+kubebuilder:validation:Required
	Image string `json:"image,omitempty"`

	//+kubebuilder:validation:MinLength=0
	//+kubebuilder:validation:MaxLength=20
	ContainerName string `json:"containerName,omitempty"`

	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Required
	Count *int32 `json:"count,omitempty"`

    //+kubebuilder:validation:Enum=deployment,pod
    //+kubebuilder:default=pod
	Type ResourceType `json:"type"`
}

// ReplicaManagerStatus defines the observed state of ReplicaManager
type ReplicaManagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//+optional
	ObjRef []corev1.ObjectReference `json:"objRef"`

	//+optional
	Healthy bool `json:"healthy"`

	Phase string `json:"phase"`

}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ReplicaManager is the Schema for the replicamanagers API
type ReplicaManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReplicaManagerSpec   `json:"spec,omitempty"`
	Status ReplicaManagerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ReplicaManagerList contains a list of ReplicaManager
type ReplicaManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReplicaManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ReplicaManager{}, &ReplicaManagerList{})
}
