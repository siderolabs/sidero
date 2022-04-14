// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/talos-systems/go-retry/retry"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "github.com/talos-systems/sidero/app/caps-controller-manager/api/v1alpha3"
	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
)

// TestServerReset verifies that all the servers got reset.
func TestServerReset(ctx context.Context, metalClient client.Client) TestFunc {
	return func(t *testing.T) {
		var machines infrav1.MetalMachineList

		labelSelector, err := labels.Parse("cluster.x-k8s.io/cluster-name=management-cluster,cluster.x-k8s.io/deployment-name=management-cluster-workers")
		require.NoError(t, err)

		err = metalClient.List(ctx, &machines, client.MatchingLabelsSelector{Selector: labelSelector})
		require.NoError(t, err)

		serverNamesToCheck := []string{}

		for i := range machines.Items {
			if machines.Items[i].Spec.ServerRef == nil {
				continue
			}

			serverNamesToCheck = append(serverNamesToCheck, machines.Items[i].Spec.ServerRef.Name)
		}

		err = scaleWorkers(ctx, metalClient, 0)
		require.NoError(t, err)

		err = retry.Constant(5*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			var servers metalv1.ServerList

			err = metalClient.List(ctx, &servers)
			if err != nil {
				return err
			}

			cleanedCount := 0

			for _, server := range servers.Items {
				for _, name := range serverNamesToCheck {
					if name != server.Name {
						continue
					}

					if !server.Status.IsClean {
						continue
					}

					cleanedCount++
				}
			}

			if cleanedCount != len(serverNamesToCheck) {
				return retry.ExpectedError(fmt.Errorf("expected %d servers to be clean, got %d", len(serverNamesToCheck), cleanedCount))
			}

			return nil
		})

		require.NoError(t, err)
	}
}
