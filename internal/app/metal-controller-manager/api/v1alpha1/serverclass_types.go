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

type Qualifiers struct {
	CPU               []CPUInformation    `json:"cpu,omitempty"`
	SystemInformation []SystemInformation `json:"systemInformation,omitempty"`
	LabelSelectors    []map[string]string `json:"labelSelectors,omitempty"`
}

// ServerClassSpec defines the desired state of ServerClass
type ServerClassSpec struct {
	Qualifiers Qualifiers `json:"qualifiers"`
}

// ServerClassStatus defines the observed state of ServerClass
type ServerClassStatus struct {
	ServersAvailable []string `json:"serversAvailable"`
	ServersInUse     []string `json:"serversInUse"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Available",type="string",JSONPath=".status.serversAvailable",description="the number of available servers"
// +kubebuilder:printcolumn:name="In Use",type="string",JSONPath=".status.serversInUse",description="the number of servers in use"

// ServerClass is the Schema for the serverclasses API
type ServerClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerClassSpec   `json:"spec,omitempty"`
	Status ServerClassStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ServerClassList contains a list of ServerClass
type ServerClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServerClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServerClass{}, &ServerClassList{})
}
