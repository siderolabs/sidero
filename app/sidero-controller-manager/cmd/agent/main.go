// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/siderolabs/go-blockdevice/blockdevice"
	"github.com/siderolabs/go-blockdevice/blockdevice/util/disk"
	"github.com/siderolabs/go-debug"
	"github.com/siderolabs/go-procfs/procfs"
	"github.com/siderolabs/go-smbios/smbios"
	"golang.org/x/sync/errgroup"

	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/api"
	"github.com/siderolabs/sidero/app/sidero-controller-manager/pkg/constants"
)

const (
	debugAddr = ":9991"
)

func mainFunc() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		debugLogFunc := func(msg string) {
			log.Print(msg)
		}
		if err := debug.ListenAndServe(ctx, debugAddr, debugLogFunc); err != nil {
			log.Fatalf("failed to start debug server: %s", err)
		}
	}()

	if err := setup(); err != nil {
		return err
	}

	if err := setupNetworking(); err != nil {
		return err
	}

	var endpoint string
	if found := procfs.ProcCmdline().Get(constants.AgentEndpointArg).First(); found != nil {
		endpoint = *found
	} else {
		return fmt.Errorf("no endpoint found")
	}

	log.Printf("Using %q as API endpoint", endpoint)

	conn, err := connect(ctx, endpoint)
	if err != nil {
		return err
	}

	defer conn.Close()

	client := api.NewAgentClient(conn)

	log.Println("Reading SMBIOS")

	s, err := smbios.New()
	if err != nil {
		return err
	}

	createResp, err := create(ctx, client, s)
	if err != nil {
		return err
	}

	log.Println("Registration complete")

	ips, err := IPAddrs()
	if err != nil {
		log.Println("failed to discover IPs")
	} else {
		if err = reconcileIPs(ctx, client, s, ips); err != nil {
			return err
		}

		log.Printf("Reconciled IPs")
	}

	if createResp.GetWipe() {
		disks, err := disk.List()
		if err != nil {
			return err
		}

		var (
			eg errgroup.Group
			wg sync.WaitGroup
		)

		heartbeatCtx, stopHeartbeat := context.WithCancel(ctx)

		heartbeatInterval := (time.Duration(createResp.RebootTimeout) * time.Second) / 3

		ticker := time.NewTicker(heartbeatInterval)

		wg.Add(1)

		go func() {
			defer wg.Done()

			for {
				callCtx, cancel := context.WithTimeout(ctx, heartbeatInterval)

				if _, err := client.Heartbeat(callCtx, &api.HeartbeatRequest{Uuid: s.SystemInformation.UUID}); err != nil {
					log.Printf("Failed to send wipe heartbeat %s", err)
				}

				cancel()

				select {
				case <-ticker.C:
				case <-heartbeatCtx.Done():
					return
				}
			}
		}()

		defer func() {
			ticker.Stop()
			stopHeartbeat()
			wg.Wait()
		}()

		for _, d := range disks {
			func(disk *disk.Disk) {
				eg.Go(func() error {
					path := disk.DeviceName

					if disk.ReadOnly {
						log.Printf("Skipping read-only disk %s", path)

						return nil
					}

					log.Printf("Resetting %s", path)

					bd, err := blockdevice.Open(path)
					if err != nil {
						log.Printf("Skipping %s: %s", path, err)

						return nil
					}

					if createResp.GetInsecureWipe() {
						if err = bd.FastWipe(); err != nil {
							return fmt.Errorf("failed wiping %q: %w", path, err)
						}

						log.Printf("Fast wiped %s", path)
					} else {
						method, err := bd.Wipe()
						if err != nil {
							return fmt.Errorf("failed wiping %q: %w", path, err)
						}

						log.Printf("Wiped %s with %s", path, method)
					}

					return bd.Close()
				})
			}(d)
		}

		if err := eg.Wait(); err != nil {
			return err
		}

		if err := wipe(ctx, client, s); err != nil {
			return err
		}

		log.Println("Wipe complete")
	}

	if createResp.GetSetupBmc() {
		log.Println("Attempting to automatically discover and configure BMC")

		// Attempt to discover the BMC IP
		// nb: we don't consider failure to get BMC info a hard failure
		//     users can always patch the bmc info to the server themselves.
		err := attemptBMCIP(ctx, client, s)
		if err != nil {
			log.Printf("encountered error discovering BMC IP. skipping setup: %q", err.Error())
		} else {
			// Attempt to add sidero user to BMC only if we discovered the IP
			// nb: we don't consider failure to get BMC info a hard failure.
			//     users can always patch the bmc info to the server themselves.
			err = attemptBMCUserSetup(ctx, client, s)
			if err != nil {
				log.Printf("encountered error setting up BMC user. skipping setup: %q", err.Error())
			}
		}
	}

	return nil
}

func main() {
	shutdown(mainFunc())
}
