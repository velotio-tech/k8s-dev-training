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

// FlightTicketSpec defines the desired state of FlightTicket
type FlightTicketSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// From - Departure location
	From string `json:"from,omitempty"`
	// To - Arrival location
	To string `json:"to,omitempty"`

	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Maximum=9
	// Number - Number of tickets for booking
	Number int `json:"number,omitempty"`

	// +kubebuilder:validation:Required
	Gvk GroupVersionKind `json:"gvk"`
}

type GroupVersionKind struct {
	Group   string `json:"group,omitempty"`
	Version string `json:"version,omitempty"`
	Kind    string `json:"kind,omitempty"`

	// +kubebuilder:validation:Minimum=1
	Replicas int32 `json:"replicas,omitempty"`
}

// FlightTicketStatus defines the observed state of FlightTicket
type FlightTicketStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Enum=InProgress;Done;Failed
	BookingStatus string `json:"bookingStatus,omitempty"`
	Fare          int    `json:"fare,omitempty"`
	ReadyReplicas int    `json:"readyReplicas,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// FlightTicket is the Schema for the flighttickets API
type FlightTicket struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FlightTicketSpec   `json:"spec,omitempty"`
	Status FlightTicketStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FlightTicketList contains a list of FlightTicket
type FlightTicketList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FlightTicket `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FlightTicket{}, &FlightTicketList{})
}
