// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package constants provides default values for some parameters.
package constants

import "net"

// Nameservers are defaults to use in bootstrap cluster.
var Nameservers = []net.IP{net.ParseIP("8.8.8.8"), net.ParseIP("1.1.1.1")}

// MTU default setting.
const MTU = 1500

// BootstrapMaster is a bootstrap cluster master node name.
const BootstrapMaster = "bootstrap-master"

// SideroAPIVersion is a string we need for creating Sidero resources.
const SideroAPIVersion = "metal.sidero.dev/v1alpha2"
