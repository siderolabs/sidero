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
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BMC defines data about how to talk to the node via ipmitool.
type BMC struct {
	Endpoint string `json:"endpoint"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
}

type SystemInformation struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	ProductName  string `json:"productName,omitempty"`
	Version      string `json:"version,omitempty"`
	SerialNumber string `json:"serialNumber,omitempty"`
	SKUNumber    string `json:"skuNumber,omitempty"`
	Family       string `json:"family,omitempty"`
}

type CPUInformation struct {
	Manufacturer string `json:"manufacturer,omitempty"`
	Version      string `json:"version,omitempty"`
}

// nb: we use apiextensions.JSON for the value below b/c we can't use interface{} with controller-gen.
// found this workaround here: https://github.com/kubernetes-sigs/controller-tools/pull/126#issuecomment-630769075

type ConfigPatches struct {
	Op    string             `json:"op"`
	Path  string             `json:"path"`
	Value apiextensions.JSON `json:"value,omitempty"`
}

// ServerSpec defines the desired state of Server.
type ServerSpec struct {
	SystemInformation *SystemInformation `json:"system,omitempty"`
	CPU               *CPUInformation    `json:"cpu,omitempty"`
	BMC               *BMC               `json:"bmc,omitempty"`
	ConfigPatches     []ConfigPatches    `json:"configPatches,omitempty"`
}

// ServerStatus defines the observed state of Server.
type ServerStatus struct {
	Ready bool `json:"ready,omitempty"`
	InUse bool `json:"inUse,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// Server is the Schema for the servers API.
type Server struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerSpec   `json:"spec,omitempty"`
	Status ServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ServerList contains a list of Server.
type ServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Server `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Server{}, &ServerList{})
}
