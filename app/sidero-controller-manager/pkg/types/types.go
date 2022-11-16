// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package types

// BootFromDisk defines a way to boot from disk.
type BootFromDisk string

const (
	BootIPXEExit BootFromDisk = "ipxe-exit"    // Use iPXE script with `exit` command.
	Boot404      BootFromDisk = "http-404"     // Return HTTP 404 response to iPXE.
	BootSANDisk  BootFromDisk = "ipxe-sanboot" // Use iPXE script with `sanboot` command.
)

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
