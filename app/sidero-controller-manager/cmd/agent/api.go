// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/siderolabs/go-blockdevice/blockdevice/util/disk"
	"github.com/siderolabs/go-retry/retry"
	"github.com/siderolabs/go-smbios/smbios"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/api"
)

func create(ctx context.Context, client api.AgentClient, s *smbios.SMBIOS) (*api.CreateServerResponse, error) {
	disks, err := disk.List()
	if err != nil {
		log.Printf("encountered error fetching disks: %q", err)
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Printf("encountered error fetching network interfaces: %q", err)
	}

	req := &api.CreateServerRequest{
		Hardware: MapHardwareInformation(s, disks, interfaces),
		Hostname: "",
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("encountered error fetching hostname: %q", err)
	} else {
		req.Hostname = hostname
	}

	var resp *api.CreateServerResponse

	err = retry.Constant(5*time.Minute, retry.WithUnits(30*time.Second), retry.WithErrorLogging(true)).Retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		resp, err = client.CreateServer(ctx, req)
		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})

	return resp, err
}

func wipe(ctx context.Context, client api.AgentClient, s *smbios.SMBIOS) error {
	return retry.Constant(5*time.Minute, retry.WithUnits(30*time.Second), retry.WithErrorLogging(true)).Retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		_, err := client.MarkServerAsWiped(ctx, &api.MarkServerAsWipedRequest{Uuid: s.SystemInformation.UUID})
		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})
}

func connect(endpoint string) (*grpc.ClientConn, error) {
	return grpc.NewClient(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
