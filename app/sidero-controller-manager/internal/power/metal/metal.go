// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package metal provides interfaces to manage metal machines.
package metal

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/power/api"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/power/ipmi"
)

// ManagementClient control power and boot order of metal machine.
type ManagementClient interface {
	PowerOn() error
	PowerOff() error
	PowerCycle() error
	IsPoweredOn() (bool, error)
	SetPXE() error
	IsFake() bool
}

// NewManagementClient builds ManagementClient from the server spec.
func NewManagementClient(ctx context.Context, client client.Client, spec *v1alpha1.ServerSpec) (ManagementClient, error) {
	switch {
	case spec.BMC != nil:
		var err error

		bmcSpec := *spec.BMC

		if bmcSpec.User == "" {
			bmcSpec.User, err = bmcSpec.UserFrom.Resolve(ctx, client)
			if err != nil {
				return nil, err
			}
		}

		if bmcSpec.Pass == "" {
			bmcSpec.Pass, err = bmcSpec.PassFrom.Resolve(ctx, client)
			if err != nil {
				return nil, err
			}
		}

		if bmcSpec.Interface == "" {
			bmcSpec.Interface = "lanplus"
		}

		return ipmi.NewClient(bmcSpec)
	case spec.ManagementAPI != nil:
		return api.NewClient(*spec.ManagementAPI)
	default:
		return fakeClient{}, nil
	}
}
