/*
Copyright 2024.

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
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
	"time"
)

type CertificateConditionType string

const (
	ConditionPending  CertificateConditionType = "Pending"
	ConditionIssued   CertificateConditionType = "Issued"
	ConditionRenewing CertificateConditionType = "Renewing"
	ConditionExpired  CertificateConditionType = "Expired"
	ConditionFailed   CertificateConditionType = "Failed"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CertificateSpec defines the desired state of Certificate
type CertificateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Domain     string `json:"domain"`
	ValidFor   string `json:"validFor"`
	SecretName string `json:"secretName"`
}

// CertificateStatus defines the observed state of Certificate
type CertificateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []CertificateCondition `json:"conditions"`
	ExpiryDate metav1.Time            `json:"expiryDate"`
	RenewedAt  metav1.Time            `json:"renewedAt,omitempty"`
}

type CertificateCondition struct {
	Type               CertificateConditionType `json:"type"`
	Status             metav1.ConditionStatus   `json:"status"` // True, False, or Unknown
	LastTransitionTime metav1.Time              `json:"lastTransitionTime,omitempty"`
	Reason             string                   `json:"reason,omitempty"`
	Message            string                   `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Certificate is the Schema for the certificates API
type Certificate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CertificateSpec   `json:"spec,omitempty"`
	Status CertificateStatus `json:"status,omitempty"`
}

func (c *Certificate) ParseValidFor() (time.Duration, error) {
	validFor := c.Spec.ValidFor
	if strings.HasSuffix(validFor, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(validFor, "d"))
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	if strings.HasSuffix(validFor, "y") {
		year, err := strconv.Atoi(strings.TrimSuffix(validFor, "y"))
		if err != nil {
			return 0, err
		}
		return time.Duration(year) * 365 * 24 * time.Hour, nil
	}
	return 0, fmt.Errorf("invalid value %s for validFor for field, should end with `d`(days) or `y`(years) e:g 1y, 20d ", validFor)
}

// +kubebuilder:object:root=true

// CertificateList contains a list of Certificate
type CertificateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Certificate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Certificate{}, &CertificateList{})
}
