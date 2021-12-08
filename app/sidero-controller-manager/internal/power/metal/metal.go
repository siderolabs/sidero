// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package metal provides interfaces to manage metal machines.
package metal

// ManagementClient control power and boot order of metal machine.
type ManagementClient interface {
	PowerOn() error
	PowerOff() error
	PowerCycle() error
	IsPoweredOn() (bool, error)
	SetPXE(mode PXEMode) error
	IsFake() bool
	Close() error
}

// PXEMode specifies PXE boot mode.
type PXEMode string

const (
	PXEModeBIOS = "bios"
	PXEModeUEFI = "uefi"
)

func (mode PXEMode) IsValid() bool {
	switch mode {
	case PXEModeBIOS:
		return true
	case PXEModeUEFI:
		return true
	default:
		return false
	}
}
