// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"net"
	"time"

	"github.com/siderolabs/go-retry/retry"
	"github.com/siderolabs/go-smbios/smbios"

	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/api"
)

func reconcileIPs(ctx context.Context, client api.AgentClient, s *smbios.SMBIOS, ips []net.IP) error {
	addresses := make([]*api.Address, len(ips))
	for i := range addresses {
		addresses[i] = &api.Address{
			Type:    "InternalIP",
			Address: ips[i].String(),
		}
	}

	return retry.Constant(5*time.Minute, retry.WithUnits(30*time.Second)).Retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		_, err := client.ReconcileServerAddresses(ctx, &api.ReconcileServerAddressesRequest{
			Uuid:    s.SystemInformation.UUID,
			Address: addresses,
		})
		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})
}

// IPAddrs finds and returns a list of non-loopback IP addresses of the
// current machine.
func IPAddrs() (ips []net.IP, err error) {
	ips = []net.IP{}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if ipnet.IP.IsGlobalUnicast() && !ipnet.IP.IsLinkLocalUnicast() {
				ips = append(ips, ipnet.IP)
			}
		}
	}

	return ips, nil
}
