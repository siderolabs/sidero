// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import (
	"fmt"
	"sort"

	"github.com/talos-systems/talos/pkg/machinery/kernel"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnvironmentDefault is an automatically created Environment.
const EnvironmentDefault = "default"

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

// EnvironmentSpec defines the desired state of Environment.
type EnvironmentSpec struct {
	Kernel Kernel `json:"kernel,omitempty"`
	Initrd Initrd `json:"initrd,omitempty"`
}

type AssetCondition struct {
	Asset  `json:",inline"`
	Status string `json:"status"`
	Type   string `json:"type"`
}

// EnvironmentStatus defines the observed state of Environment.
type EnvironmentStatus struct {
	Conditions []AssetCondition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Kernel",type="string",JSONPath=".spec.kernel.url",description="the kernel for the environment"
// +kubebuilder:printcolumn:name="Initrd",type="string",JSONPath=".spec.initrd.url",description="the initrd for the environment"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type==\"Ready\")].status",description="indicates the readiness of the environment"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="The age of this resource"
// +kubebuilder:storageversion

// Environment is the Schema for the environments API.
type Environment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EnvironmentSpec   `json:"spec,omitempty"`
	Status EnvironmentStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EnvironmentList contains a list of Environment.
type EnvironmentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Environment `json:"items"`
}

// EnvironmentDefaultSpec returns EnvironmentDefault's spec.
func EnvironmentDefaultSpec(talosRelease, apiEndpoint string, apiPort uint16) *EnvironmentSpec {
	args := make([]string, 0, len(kernel.DefaultArgs)+6)
	args = append(args, kernel.DefaultArgs...)
	args = append(args, "console=tty0", "console=ttyS0", "earlyprintk=ttyS0")
	args = append(args, "initrd=initramfs.xz", "talos.platform=metal")
	sort.Strings(args)

	return &EnvironmentSpec{
		Kernel: Kernel{
			Asset: Asset{
				URL: fmt.Sprintf("https://github.com/talos-systems/talos/releases/download/%s/vmlinuz-amd64", talosRelease),
			},
			Args: args,
		},
		Initrd: Initrd{
			Asset: Asset{
				URL: fmt.Sprintf("https://github.com/talos-systems/talos/releases/download/%s/initramfs-amd64.xz", talosRelease),
			},
		},
	}
}

// IsReady returns aggregated Environment readiness.
func (env *Environment) IsReady() bool {
	assetURLs := map[string]struct{}{}

	if env.Spec.Kernel.URL != "" {
		assetURLs[env.Spec.Kernel.URL] = struct{}{}
	}

	if env.Spec.Initrd.URL != "" {
		assetURLs[env.Spec.Initrd.URL] = struct{}{}
	}

	for _, cond := range env.Status.Conditions {
		if cond.Status == "True" && cond.Type == "Ready" {
			delete(assetURLs, cond.URL)
		}
	}

	return len(assetURLs) == 0
}

func init() {
	SchemeBuilder.Register(&Environment{}, &EnvironmentList{})
}
