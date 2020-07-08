/*
Copyright 2020 Talos Systems, Inc.
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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Asset struct {
	URL    string `json:"url,omitempty"`
	SHA512 string `json:"sha512,omitempty"`
}

type Kernel struct {
	Asset `json:",inline"`

	Args []string `json:"args,omitempty"`
}

type Initrd struct {
	Asset `json:",inline"`
}

// EnvironmentSpec defines the desired state of Environment
type EnvironmentSpec struct {
	Kernel Kernel `json:"kernel,omitempty"`
	Initrd Initrd `json:"initrd,omitempty"`
}

type AssetCondition struct {
	Asset  `json:",inline"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// EnvironmentStatus defines the observed state of Environment
type EnvironmentStatus struct {
	Conditions []AssetCondition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Kernel",type="string",JSONPath=".spec.kernel.url",description="the kernel for the environment"
// +kubebuilder:printcolumn:name="Initrd",type="string",JSONPath=".spec.initrd.url",description="the initrd for the environment"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description="indicates the readiness of the environment"

// Environment is the Schema for the environments API
type Environment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EnvironmentSpec   `json:"spec,omitempty"`
	Status EnvironmentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EnvironmentList contains a list of Environment
type EnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Environment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Environment{}, &EnvironmentList{})
}
