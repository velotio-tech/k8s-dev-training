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

// CustomDeploymentSpec defines the desired state of CustomDeployment
type CustomDeploymentSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	
	// +kubebuilder:validation:MaxLength=64	
	DeploymentName string `json:"deploymentName"`


	// If not specified, one replica will be created.
	// +optional
	// +kubebuilder:validation:Minimum=0
	Replicas  int32  `json:"replicas,omitempty"`
}


// CustomDeploymentStatus defines the observed state of CustomDeployment
type CustomDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	
	// +optional
	// +kubebuilder:validation:Minimum=0
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CustomDeployment is the Schema for the customdeployments API
type CustomDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomDeploymentSpec   `json:"spec,omitempty"`
	Status CustomDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CustomDeploymentList contains a list of CustomDeployment
type CustomDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomDeployment{}, &CustomDeploymentList{})
}
