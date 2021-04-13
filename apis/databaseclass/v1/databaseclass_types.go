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
	"github.com/bedag/kubernetes-dbaas/pkg/database"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DatabaseClassSpec defines the desired state of DatabaseClass
type DatabaseClassSpec struct {
	Driver     string                         `json:"driver,omitempty"`
	Operations map[string]*database.Operation `json:"operations,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=databaseclasses,scope=Cluster,shortName=dbc
// DatabaseClass is the Schema for the databaseclasses API
type DatabaseClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DatabaseClassSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true
// DatabaseClassList contains a list of DatabaseClass
type DatabaseClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DatabaseClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DatabaseClass{}, &DatabaseClassList{})
}
