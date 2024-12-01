/*


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

// JobType defines the type to create job or cronjob
// +kubebuilder:validation:Enum=job;cronjob
type JobType string

// +kubebuilder:validation:type=string
// Status specifies the status of WorkloadJob operating on
type Status string

// MyJobSpec defines the desired state of MyJob
type MyJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Type is the type of backup in the sequence of backups of an Application.
	ResourceType JobType `json:"resourceType"`

	// JobName is the name of the Job resource that the
	// controller should create.
	// This field must be specified.
	// +kubebuilder:validation:MaxLength=64
	JobName string `json:"jobName"`
}

// MyJobStatus defines the observed state of MyJob
type MyJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Status is the status of the condition.
	// +nullable:true
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=InProgress;Error;Completed;Failed
	Status Status `json:"status,omitempty"`

	// Timestamp is the time a condition occurred.
	// +nullable:true
	// +kubebuilder:validation:Optional
	// +kubebulder:validation:Format="date-time"
	Timestamp *metav1.Time `json:"timestamp,omitempty"`

	// A brief message indicating details about why the component is in this condition.
	// +nullable:true
	// +kubebuilder:validation:Optional
	Reason string `json:"reason,omitempty"`
}

// +kubebuilder:object:root=true

// MyJob is the Schema for the myjobs API
type MyJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyJobSpec   `json:"spec,omitempty"`
	Status MyJobStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MyJobList contains a list of MyJob
type MyJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MyJob{}, &MyJobList{})
}
