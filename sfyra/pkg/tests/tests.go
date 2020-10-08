// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package tests provides the Sidero tests.
package tests

import (
	"context"
	"log"
	"testing"

	"github.com/talos-systems/sidero/sfyra/pkg/capi"
	"github.com/talos-systems/sidero/sfyra/pkg/talos"
	"github.com/talos-systems/sidero/sfyra/pkg/vm"
)

// TestFunc is a testing function prototype.
type TestFunc func(t *testing.T)

// Options for the test.
type Options struct {
	KernelURL, InitrdURL string
	InstallerImage       string

	RegistryMirrors []string
}

// Run all the tests.
func Run(ctx context.Context, cluster talos.Cluster, vmSet *vm.Set, capiManager *capi.Manager, options Options) (ok bool) {
	metalClient, err := capiManager.GetMetalClient(ctx)
	if err != nil {
		log.Printf("error creating metalClient: %s", err)

		return false
	}

	return testing.MainStart(matchStringOnly(func(pat, str string) (bool, error) { return true, nil }), []testing.InternalTest{
		{
			"TestServerRegistration",
			TestServerRegistration(ctx, metalClient, vmSet),
		},
		{
			"TestServerMgmtAPI",
			TestServerMgmtAPI(ctx, metalClient, vmSet),
		},
		{
			"TestServerPatch",
			TestServerPatch(ctx, metalClient, options.InstallerImage, options.RegistryMirrors),
		},
		{
			"TestServerAcceptance",
			TestServerAcceptance(ctx, metalClient, vmSet),
		},
		{
			"TestServersReady",
			TestServersReady(ctx, metalClient),
		},
		{
			"TestEnvironmentDefault",
			TestEnvironmentDefault(ctx, metalClient, cluster, options.KernelURL, options.InitrdURL),
		},
		{
			"TestServerClassDefault",
			TestServerClassDefault(ctx, metalClient, vmSet),
		},
		{
			"TestManagementCluster",
			TestManagementCluster(ctx, metalClient, cluster, vmSet, capiManager),
		},
		{
			"TestMatchServersMetalMachines",
			TestMatchServersMetalMachines(ctx, metalClient),
		},
		{
			"TestServerReset",
			TestServerReset(ctx, metalClient, vmSet),
		},
	}, nil, nil).Run() == 0
}
