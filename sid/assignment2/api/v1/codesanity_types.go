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

type ResourceType string

const (
	Job     ResourceType = "job"
	CronJob ResourceType = "cronjob"
)

// CodeSanitySpec defines the desired state of CodeSanity
type CodeSanitySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Regex to select test files
	// +kubebuilder:validation:MinLength=4
	TestFilesRegexStr string `json:"testFilesRegexStr"`

	// Required coverage percentage
	// +kubebuilder:validation:Minimum=30
	// +optional
	RequiredCoverage *int32 `json:"requiredCoverage,omitempty"`

	// Pods to cover
	// +optional
	PodNames []string `json:"podNames,omitempty"`

	// Images to cover
	// +optional
	Images []string `json:"images,omitempty"`

	// Resource to be used to run the tests. This can either be a Job or a CronJob
	// +kubebuilder:validation:Enum=job;cronjob
	Resource ResourceType `json:"resource,omitempty"`
}

// CodeSanityStatus defines the observed state of CodeSanity

type CodeSanityStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Was the sanity successful
	// A map that stores the latest test coverage corresponding to an image
	//CoverageMap CoverageMap `json:"coverageMap,omitempty"`

	// Healthy pods
	// A map that stores the latest test coverage for a pod's image against it's image
	HealthyPods []string `json:"healthyPods,omitempty"`

	// Unhealthy pods
	// A map that stores the latest test coverage for a pod's image against it's image
	UnhealthyPods []string `json:"unhealthyPods,omitempty"`

	// Processed Pods
	// A list that holds pod names that have been processed
	ProcessedPods []string `json:"processedPods,omitempty"`

	// Last run
	LastRunAt metav1.Time `json:"lastRunAt,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// CodeSanity is the Schema for the codesanities API
type CodeSanity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CodeSanitySpec   `json:"spec,omitempty"`
	Status CodeSanityStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CodeSanityList contains a list of CodeSanity
type CodeSanityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CodeSanity `json:"items"`
}

//type CoverageMap struct {
//	Images map[string]*int32 `json:"Images,omitempty"`
//}

func init() {
	SchemeBuilder.Register(&CodeSanity{}, &CodeSanityList{})
}
