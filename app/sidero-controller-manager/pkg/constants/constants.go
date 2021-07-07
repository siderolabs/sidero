// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package constants

import "time"

const (
	DataDirectory    = "/var/lib/sidero"
	AgentEndpointArg = "sidero.endpoint"

	KernelAsset = "vmlinuz"
	InitrdAsset = "initramfs.xz"

	DefaultRequeueAfter = time.Second * 20

	DefaultServerRebootTimeout = time.Minute * 20

	DefaultBMCPort = uint32(623)
)
