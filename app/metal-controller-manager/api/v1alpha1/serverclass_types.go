// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ServerClassAny is an automatically created ServerClass that includes all Servers.
const ServerClassAny = "any"

type Qualifiers struct {
	CPU               []CPUInformation    `json:"cpu,omitempty"`
	SystemInformation []SystemInformation `json:"systemInformation,omitempty"`
	LabelSelectors    []map[string]string `json:"labelSelectors,omitempty"`
}

// ServerClassSpec defines the desired state of ServerClass.
type ServerClassSpec struct {
	EnvironmentRef *corev1.ObjectReference `json:"environmentRef,omitempty"`
	Qualifiers     Qualifiers              `json:"qualifiers"`
	ConfigPatches  []ConfigPatches         `json:"configPatches,omitempty"`
}

// ServerClassStatus defines the observed state of ServerClass.
type ServerClassStatus struct {
	ServersAvailable []string `json:"serversAvailable"`
	ServersInUse     []string `json:"serversInUse"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Available",type="string",JSONPath=".status.serversAvailable",description="the number of available servers"
// +kubebuilder:printcolumn:name="In Use",type="string",JSONPath=".status.serversInUse",description="the number of servers in use"

// ServerClass is the Schema for the serverclasses API.
type ServerClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerClassSpec   `json:"spec,omitempty"`
	Status ServerClassStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ServerClassList contains a list of ServerClass.
type ServerClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServerClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ServerClass{}, &ServerClassList{})
}
