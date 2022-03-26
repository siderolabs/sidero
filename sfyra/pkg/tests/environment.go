// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/talos-systems/go-procfs/procfs"
	"github.com/talos-systems/go-retry/retry"
	"github.com/talos-systems/talos/pkg/machinery/kernel"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/sfyra/pkg/constants"
	"github.com/talos-systems/sidero/sfyra/pkg/talos"
)

const environmentName = "sfyra"

// TestEnvironmentDefault verifies environment "default".
func TestEnvironmentDefault(ctx context.Context, metalClient client.Client, cluster talos.Cluster, kernelURL, initrdURL string) TestFunc {
	return func(t *testing.T) {
		var environment metalv1.Environment
		err := metalClient.Get(ctx, types.NamespacedName{Name: metalv1.EnvironmentDefault}, &environment)
		require.NoError(t, err)
		assert.True(t, environment.IsReady())

		// delete environment to see it being recreated
		err = metalClient.Delete(ctx, &environment)
		require.NoError(t, err)

		environment = metalv1.Environment{}
		err = retry.Constant(60 * time.Second).Retry(func() error {
			if err := metalClient.Get(ctx, types.NamespacedName{Name: metalv1.EnvironmentDefault}, &environment); err != nil {
				if apierrors.IsNotFound(err) {
					return retry.ExpectedError(err)
				}
				return err
			}

			if !environment.IsReady() {
				return retry.ExpectedErrorf("some assets are not ready")
			}

			return nil
		})
		require.NoError(t, err)
		assert.True(t, environment.IsReady())
	}
}

// TestEnvironmentCreate verifies environment creation.
func TestEnvironmentCreate(ctx context.Context, metalClient client.Client, cluster talos.Cluster, kernelURL, initrdURL string) TestFunc {
	return func(t *testing.T) {
		var environment metalv1.Environment

		if err := metalClient.Get(ctx, types.NamespacedName{Name: environmentName}, &environment); err != nil {
			if !apierrors.IsNotFound(err) {
				require.NoError(t, err)
			}

			cmdline := procfs.NewCmdline("")
			cmdline.SetAll(kernel.DefaultArgs)

			cmdline.Append("console", "ttyS0")
			cmdline.Append("reboot", "k")
			cmdline.Append("panic", "1")
			cmdline.Append("talos.platform", "metal")
			cmdline.Append("talos.shutdown", "halt")
			cmdline.Append("initrd", "initramfs.xz")

			environment.APIVersion = constants.SideroAPIVersion
			environment.Name = environmentName
			environment.Spec.Kernel.URL = kernelURL
			environment.Spec.Kernel.SHA512 = "" // TODO: add a test
			environment.Spec.Kernel.Args = cmdline.Strings()
			environment.Spec.Initrd.URL = initrdURL
			environment.Spec.Initrd.SHA512 = "" // TODO: add a test

			require.NoError(t, metalClient.Create(ctx, &environment))
		}

		// wait for the environment to report ready
		require.NoError(t, retry.Constant(5*time.Minute, retry.WithUnits(10*time.Second)).Retry(func() error {
			if err := metalClient.Get(ctx, types.NamespacedName{Name: environmentName}, &environment); err != nil {
				return err
			}

			if !environment.IsReady() {
				return retry.ExpectedErrorf("some assets are not ready")
			}

			return nil
		}))
	}
}
