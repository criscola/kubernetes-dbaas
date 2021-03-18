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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatabaseSpec defines the desired state of Database.
//
// Important: Run "make" to regenerate code after modifying this file. Json tags are required.
type DatabaseSpec struct {
	// Endpoint associates this resource with a particular endpoint (must be already configured on the operator side)
	Endpoint string `json:"endpoint,omitempty"`
	// Params is a map containing parameters to be mapped to the database instance
	Params map[string]string `json:"params,omitempty"`
}

// DatabaseStatus defines the observed state of Database.
type DatabaseStatus struct {
	// LastError if not nil, the resource in an error state
	LastError string `json:"lastError,omitempty"`
	// LastUpdate specifies the last time the Status field has been updated
	LastUpdate string `json:"lastUpdate,omitempty"`
	// LastErrorUpdateCount specifies how many times the LastError field has been updated
	LastErrorUpdateCount int `json:"lastErrorUpdateCount,omitempty"`
	// If Unrecoverable is set to true, the controller was unable to fix the issue by itself
	//
	// TODO: Do something like 'kubectl get pods', i.e. create a set of state and enable users to print column with the current state
	Unrecoverable bool `json:"unrecoverable,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Database is the Schema for the database API
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseSpec   `json:"spec,omitempty"`
	Status DatabaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DatabaseList contains a list of Database
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Database `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Database{}, &DatabaseList{})
}
