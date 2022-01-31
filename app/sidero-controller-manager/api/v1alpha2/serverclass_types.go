// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	siderotypes "github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/types"
)

// ServerClassAny is an automatically created ServerClass that includes all Servers.
const ServerClassAny = "any"

type Qualifiers struct {
	Hardware       []HardwareInformation `json:"hardware,omitempty"`
	LabelSelectors []map[string]string   `json:"labelSelectors,omitempty"`
}

// ServerClassSpec defines the desired state of ServerClass.
type ServerClassSpec struct {
	// Reference to the environment which should be used to provision the servers via this server class.
	// +optional
	EnvironmentRef *corev1.ObjectReference `json:"environmentRef,omitempty"`
	// Qualifiers to match on the server spec.
	//
	// If qualifiers are empty, they match all servers.
	// Server should match both qualifiers and selector conditions to be included into the server class.
	// +optional
	Qualifiers Qualifiers `json:"qualifiers"`
	// Label selector to filter the matching servers based on labels.
	// A label selector is a label query over a set of resources. The result of matchLabels and
	// matchExpressions are ANDed. An empty label selector matches all objects. A null label selector matches no objects.
	// +optional
	Selector metav1.LabelSelector `json:"selector"`
	// Set of config patches to apply to the machine configuration to the servers provisioned via this server class.
	// +optional
	ConfigPatches []ConfigPatches `json:"configPatches,omitempty"`
	// BootFromDiskMethod specifies the method to exit iPXE to force boot from disk.
	//
	// If not set, controller default is used.
	// Valid values: ipxe-exit, http-404, ipxe-sanboot.
	//
	// +optional
	BootFromDiskMethod siderotypes.BootFromDisk `json:"bootFromDiskMethod,omitempty"`
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
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of this resource"
// +kubebuilder:storageversion

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
