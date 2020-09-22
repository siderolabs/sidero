// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package constants provides default values for some parameters.
package constants

import "net"

// Default network parameters.
var (
	Nameservers = []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("1.1.1.1")}
	CNIBinPath  = []string{"/opt/cni/bin"}
)

// Default CNI paths.
const (
	CNIConfDir  = "/etc/cni/conf.d"
	CNICacheDir = "/var/lib/cni"
)

// MTU default setting.
const MTU = 1500

// BootstrapMaster is a bootstrap cluster master node name.
const BootstrapMaster = "bootstrap-master"
