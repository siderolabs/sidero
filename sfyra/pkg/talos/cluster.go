// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package talos provides interfaces for Talos objects.
package talos

import (
	"net"

	"github.com/talos-systems/talos/pkg/cluster"
)

// Cluster is an abstract interface for the Talos cluster.
//
// It might be provided by `provision` library created cluster, or by the CAPI built cluster.
type Cluster interface {
	// Name of the cluster.
	Name() string
	// IP of the bridge which controls the cluster.
	BridgeIP() net.IP
	// IP for the Sidero components (TFTP, iPXE, etc.).
	SideroComponentsIP() net.IP
	// K8s client source.
	KubernetesClient() cluster.K8sProvider
}
