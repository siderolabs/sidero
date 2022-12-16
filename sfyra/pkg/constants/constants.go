// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package constants provides default values for some parameters.
package constants

import "net/netip"

// Nameservers are defaults to use in bootstrap cluster.
var Nameservers = []netip.Addr{netip.MustParseAddr("8.8.8.8"), netip.MustParseAddr("1.1.1.1")}

// MTU default setting.
const MTU = 1440

// BootstrapControlPlane is a bootstrap cluster control-plane node name.
const BootstrapControlPlane = "bootstrap-control-plane"

// BootstrapWorker is a bootstrap cluster worker node name.
const BootstrapWorker = "bootstrap-worker"

// SideroAPIVersion is a string we need for creating Sidero resources.
const SideroAPIVersion = "metal.sidero.dev/v1alpha1"
