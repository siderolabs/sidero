// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package metal provides interfaces to manage metal machines.
package metal

import "github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/types"

// ManagementClient control power and boot order of metal machine.
type ManagementClient interface {
	PowerOn() error
	PowerOff() error
	PowerCycle() error
	IsPoweredOn() (bool, error)
	SetPXE(mode types.PXEMode) error
	IsFake() bool
	Close() error
}
