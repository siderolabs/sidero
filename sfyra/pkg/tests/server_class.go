// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/talos-systems/go-retry/retry"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/sfyra/pkg/constants"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

const serverClassName = "default"

// TestServerClassDefault verifies server class creation.
func TestServerClassDefault(ctx context.Context, metalClient client.Client, vmSet *vm.Set) TestFunc {
	return func(t *testing.T) {
		classSpec := v1alpha1.ServerClassSpec{
			Qualifiers: v1alpha1.Qualifiers{
				CPU: []v1alpha1.CPUInformation{
					{
						Manufacturer: "QEMU",
					},
				},
			},
		}

		serverClass, err := createServerClass(ctx, metalClient, serverClassName, classSpec)
		require.NoError(t, err)

		numNodes := len(vmSet.Nodes())

		// wait for the server class to gather all nodes (all nodes should match)
		require.NoError(t, retry.Constant(2*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			if err := metalClient.Get(ctx, types.NamespacedName{Name: serverClassName}, &serverClass); err != nil {
				return retry.UnexpectedError(err)
			}

			if len(serverClass.Status.ServersAvailable)+len(serverClass.Status.ServersInUse) != numNodes {
				return retry.ExpectedError(fmt.Errorf("%d + %d != %d", len(serverClass.Status.ServersAvailable), len(serverClass.Status.ServersInUse), numNodes))
			}

			return nil
		}))

		assert.Len(t, append(serverClass.Status.ServersAvailable, serverClass.Status.ServersInUse...), numNodes)

		nodes := vmSet.Nodes()
		expectedUUIDs := make([]string, len(nodes))

		for i := range nodes {
			expectedUUIDs[i] = nodes[i].UUID.String()
		}

		actualUUIDs := append(serverClass.Status.ServersAvailable, serverClass.Status.ServersInUse...)

		sort.Strings(expectedUUIDs)
		sort.Strings(actualUUIDs)

		assert.Equal(t, expectedUUIDs, actualUUIDs)
	}
}

func createServerClass(ctx context.Context, metalClient client.Client, name string, spec v1alpha1.ServerClassSpec) (v1alpha1.ServerClass, error) {
	var retClass v1alpha1.ServerClass

	if err := metalClient.Get(ctx, types.NamespacedName{Name: name}, &retClass); err != nil {
		if !apierrors.IsNotFound(err) {
			return retClass, nil
		}

		retClass.APIVersion = constants.SideroAPIVersion
		retClass.Name = name
		retClass.Spec = spec

		err = metalClient.Create(ctx, &retClass)
		if err != nil {
			return retClass, err
		}
	}

	return retClass, nil
}
