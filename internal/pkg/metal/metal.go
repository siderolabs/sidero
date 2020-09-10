// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package metal provides interfaces to manage metal machines.
package metal

import (
	"github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/internal/pkg/api"
	"github.com/talos-systems/sidero/internal/pkg/ipmi"
)

// ManagementClient control power and boot order of metal machine.
type ManagementClient interface {
	PowerOn() error
	PowerOff() error
	PowerCycle() error
	IsPoweredOn() (bool, error)
	SetPXE() error
}

// NewManagementClient builds ManagementClient from the server spec.
func NewManagementClient(spec *v1alpha1.ServerSpec) (ManagementClient, error) {
	switch {
	case spec.BMC != nil:
		return ipmi.NewClient(*spec.BMC)
	case spec.ManagementAPI != nil:
		return api.NewClient(*spec.ManagementAPI)
	default:
		return fakeClient{}, nil
	}
}
