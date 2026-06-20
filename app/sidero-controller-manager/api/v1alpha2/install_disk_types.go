// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

// InstallDisk defines Talos install disk configuration.
// +kubebuilder:validation:XValidation:rule="has(self.deviceName) != has(self.selector)",message="exactly one of deviceName or selector must be set"
type InstallDisk struct {
	// DeviceName is an explicit block device path, e.g. "/dev/sda" or "/dev/disk/by-id/...".
	// +optional
	DeviceName string `json:"deviceName,omitempty"`

	// Selector selects a disk by stable hardware attributes.
	// +optional
	Selector *InstallDiskSelector `json:"selector,omitempty"`
}

// InstallDiskSelector selects a disk by stable attributes.
// +kubebuilder:validation:MinProperties=1
type InstallDiskSelector struct {
	Name   string `json:"name,omitempty"`
	Model  string `json:"model,omitempty"`
	Serial string `json:"serial,omitempty"`
	UUID   string `json:"uuid,omitempty"`
	WWID   string `json:"wwid,omitempty"`
	// +kubebuilder:validation:Enum=ssd;hdd;nvme;sd
	Type string `json:"type,omitempty"`
	Size uint64 `json:"size,omitempty"`
}
