// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metal "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
)

// TestMatchServersMetalMachines verifies that number of metal machines and servers match.
func TestMatchServersMetalMachines(ctx context.Context, metalClient client.Client) TestFunc {
	return func(t *testing.T) {
		var machines v1alpha3.MachineList

		require.NoError(t, metalClient.List(ctx, &machines))

		var servers metal.ServerList

		require.NoError(t, metalClient.List(ctx, &servers))

		inUseServers := 0

		for _, server := range servers.Items {
			if server.Status.InUse {
				inUseServers++
			}
		}

		assert.Equal(t, len(machines.Items), inUseServers)
	}
}
