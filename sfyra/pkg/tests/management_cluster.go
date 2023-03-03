// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"testing"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/siderolabs/sidero/sfyra/pkg/capi"
	"github.com/siderolabs/sidero/sfyra/pkg/talos"
)

const (
	managementClusterName   = "management-cluster"
	managementClusterLBPort = 10000
)

// TestManagementCluster deploys the management cluster via CAPI.
func TestManagementCluster(ctx context.Context, metalClient client.Client, cluster talos.Cluster, capiManager *capi.Manager, talosRelease, kubernetesVersion string) TestFunc {
	return func(t *testing.T) {
		deployCluster(ctx, t, metalClient, cluster, capiManager, managementClusterName, serverClassName, managementClusterLBPort, 1, 1, talosRelease, kubernetesVersion)
	}
}
