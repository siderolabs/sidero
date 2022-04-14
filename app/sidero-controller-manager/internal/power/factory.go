// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package power provides common interface to manage power state.
package power

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/power/api"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/power/ipmi"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/power/metal"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/constants"
)

// NewManagementClient builds ManagementClient from the server spec.
func NewManagementClient(ctx context.Context, client client.Client, spec *metalv1.ServerSpec) (metal.ManagementClient, error) {
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

		if bmcSpec.User == "" || bmcSpec.Pass == "" {
			// no username and password, BMC information is not fully populated yet
			return fakeClient{}, nil
		}

		if bmcSpec.Interface == "" {
			bmcSpec.Interface = "lanplus"
		}

		if bmcSpec.Port == 0 {
			bmcSpec.Port = constants.DefaultBMCPort
		}

		return ipmi.NewClient(bmcSpec)
	case spec.ManagementAPI != nil:
		return api.NewClient(*spec.ManagementAPI)
	default:
		return fakeClient{}, nil
	}
}
