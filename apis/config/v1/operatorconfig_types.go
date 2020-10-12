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
	cfg "sigs.k8s.io/controller-runtime/pkg/config/v1alpha1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// OperatorConfig is the Schema for the operatorconfigs API
type OperatorConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// +kubebuilder:validation:Optional
	// ControllerManagerConfigurationSpec returns the configurations for controllers. See https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/config/v1alpha1#ControllerManagerConfigurationSpec
	cfg.ControllerManagerConfigurationSpec `json:",inline"`

	// rps configures the rate limiter to allow only a certain amount of operations per second per endpoint. If set to 0,
	// operations won't be rate-limited.
	int `json:"rps,omitempty"`

	// +kubebuilder:kubebuilder:validation:MinItems=1
	// DbmsList returns the configuration for the database endpoints.
	database.DbmsList `json:"dbms"`
}

// +kubebuilder:object:root=true
func (c OperatorConfig) Complete() (cfg.ControllerManagerConfigurationSpec, error) {
	return c.ControllerManagerConfigurationSpec, nil
}

func init() {
	SchemeBuilder.Register(&OperatorConfig{})
}
