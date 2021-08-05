// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"sync"
	"time"

	"github.com/talos-systems/go-blockdevice/blockdevice"
	"github.com/talos-systems/go-blockdevice/blockdevice/util/disk"
	debug "github.com/talos-systems/go-debug"
	kmsg "github.com/talos-systems/go-kmsg"
	"github.com/talos-systems/go-procfs/procfs"
	"github.com/talos-systems/go-retry/retry"
	"github.com/talos-systems/go-smbios/smbios"
	talosnet "github.com/talos-systems/net"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
	"google.golang.org/grpc"

	"github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/api"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/power/ipmi"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/constants"
)

const (
	debugAddr = ":9991"
)

func setup() error {
	if err := os.MkdirAll("/etc", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/dev", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/proc", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/sys", 0o777); err != nil {
		return err
	}

	if err := os.MkdirAll("/tmp", 0o777); err != nil {
		return err
	}

	if err := unix.Mount("devtmpfs", "/dev", "devtmpfs", unix.MS_NOSUID, "mode=0755"); err != nil {
		return err
	}

	if err := unix.Mount("proc", "/proc", "proc", unix.MS_NOSUID|unix.MS_NOEXEC|unix.MS_NODEV, ""); err != nil {
		return err
	}

	if err := unix.Mount("sysfs", "/sys", "sysfs", 0, ""); err != nil {
		return err
	}

	if err := unix.Mount("tmpfs", "/tmp", "tmpfs", 0, ""); err != nil {
		return err
	}

	if err := unix.Symlink("/proc/net/pnp", "/etc/resolv.conf"); err != nil {
		return err
	}

	if err := kmsg.SetupLogger(nil, "[sidero]", nil); err != nil {
		return err
	}

	// Set the PATH env var.
	if err := os.Setenv("PATH", "/sbin:/bin:/usr/sbin:/usr/bin:/usr/local/sbin:/usr/local/bin"); err != nil {
		return errors.New("error setting PATH")
	}

	return nil
}

func create(ctx context.Context, client api.AgentClient, s *smbios.SMBIOS) (*api.CreateServerResponse, error) {
	uuid, err := s.SystemInformation().UUID()
	if err != nil {
		return nil, err
	}

	req := &api.CreateServerRequest{
		SystemInformation: &api.SystemInformation{
			Uuid:         uuid.String(),
			Manufacturer: s.SystemInformation().Manufacturer(),
			ProductName:  s.SystemInformation().ProductName(),
			Version:      s.SystemInformation().Version(),
			SerialNumber: s.SystemInformation().SerialNumber(),
			SkuNumber:    s.SystemInformation().SKUNumber(),
			Family:       s.SystemInformation().Family(),
		},
		Cpu: &api.CPU{
			Manufacturer: s.ProcessorInformation().ProcessorManufacturer(),
			Version:      s.ProcessorInformation().ProcessorVersion(),
		},
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
	uuid, err := s.SystemInformation().UUID()
	if err != nil {
		return err
	}

	return retry.Constant(5*time.Minute, retry.WithUnits(30*time.Second), retry.WithErrorLogging(true)).Retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		_, err = client.MarkServerAsWiped(ctx, &api.MarkServerAsWipedRequest{Uuid: uuid.String()})
		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})
}

func reconcileIPs(ctx context.Context, client api.AgentClient, s *smbios.SMBIOS, ips []net.IP) error {
	uuid, err := s.SystemInformation().UUID()
	if err != nil {
		return err
	}

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

		_, err = client.ReconcileServerAddresses(ctx, &api.ReconcileServerAddressesRequest{
			Uuid:    uuid.String(),
			Address: addresses,
		})
		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})
}

func shutdown(err error) {
	if err != nil {
		log.Println(err)
	}

	for i := 10; i >= 0; i-- {
		log.Printf("rebooting in %d seconds\n", i)
		time.Sleep(1 * time.Second)
	}

	if unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART) == nil {
		select {}
	}

	os.Exit(1)
}

func connect(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	return grpc.DialContext(ctx, endpoint, grpc.WithInsecure())
}

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

	ips, err := talosnet.IPAddrs()
	if err != nil {
		log.Println("failed to discover IPs")
	} else {
		if err = reconcileIPs(ctx, client, s, ips); err != nil {
			shutdown(err)
		}

		log.Printf("Reconciled IPs")
	}

	if createResp.GetWipe() {
		disks, err := disk.List()
		if err != nil {
			shutdown(err)
		}

		uuid, err := s.SystemInformation().UUID()
		if err != nil {
			shutdown(err)
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

				if _, err := client.Heartbeat(callCtx, &api.HeartbeatRequest{Uuid: uuid.String()}); err != nil {
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

		for _, disk := range disks {
			func(path string) {
				eg.Go(func() error {
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
			}(disk.DeviceName)
		}

		if err := eg.Wait(); err != nil {
			shutdown(err)
		}

		if err := wipe(ctx, client, s); err != nil {
			shutdown(err)
		}

		log.Println("Wipe complete")
	}

	return nil
}

func main() {
	shutdown(mainFunc())
}

func attemptBMCIP(ctx context.Context, client api.AgentClient, s *smbios.SMBIOS) error {
	uuid, err := s.SystemInformation().UUID()
	if err != nil {
		return err
	}

	bmcInfo := &api.BMCInfo{}

	// Create "open" client
	bmcSpec := v1alpha1.BMC{
		Interface: "open",
	}

	ipmiClient, err := ipmi.NewClient(bmcSpec)
	if err != nil {
		return err
	}

	// Fetch BMC IP
	ipResp, err := ipmiClient.GetBMCIP()
	if err != nil {
		return err
	}

	bmcIP := net.IP(ipResp.Data)
	bmcInfo.Ip = bmcIP.String()

	// Attempt to update server object
	err = retry.Constant(5*time.Minute, retry.WithUnits(30*time.Second), retry.WithErrorLogging(true)).Retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		_, err = client.UpdateBMCInfo(
			ctx,
			&api.UpdateBMCInfoRequest{
				Uuid:    uuid.String(),
				BmcInfo: bmcInfo,
			},
		)

		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})

	return nil
}

func attemptBMCUserSetup(ctx context.Context, client api.AgentClient, s *smbios.SMBIOS) error {
	uuid, err := s.SystemInformation().UUID()
	if err != nil {
		return err
	}

	bmcInfo := &api.BMCInfo{}

	// Create "open" client
	bmcSpec := v1alpha1.BMC{
		Interface: "open",
	}

	ipmiClient, err := ipmi.NewClient(bmcSpec)
	if err != nil {
		return err
	}

	// Get user summary to see how many user slots
	summResp, err := ipmiClient.GetUserSummary()
	if err != nil {
		return err
	}

	maxUsers := summResp.MaxUsers & 0x1F // Only bits [0:5] provide this number

	// Check if sidero user already exists by combing through all userIDs
	// nb: we start looking at user id 2, because 1 should always be an unamed admin user and
	//     we don't want to confuse that unnamed admin with an open slot we can take over.
	exists := false
	sideroUserID := uint8(0)

	for i := uint8(2); i <= maxUsers; i++ {
		userRes, err := ipmiClient.GetUserName(i)
		if err != nil {
			// nb: A failure here actually seems to mean that the user slot is unused,
			// even though you can also have a slot with empty user as well. *scratches head*
			// We'll take note of this spot if we haven't already found another empty one.
			if sideroUserID == 0 {
				sideroUserID = i
			}

			continue
		}

		// Found pre-existing sidero user
		if userRes.Username == "sidero" {
			exists = true
			sideroUserID = i
			log.Printf("Sidero user already present in slot %d. We'll claim it as our very own.\n", i)

			break
		} else if userRes.Username == "" && sideroUserID == 0 {
			// If this is the first empty user that's not the UID 1 (which we skip),
			// we'll take this spot for sidero user
			log.Printf("Found empty user slot %d. Noting as a possible place for sidero user.\n", i)
			sideroUserID = i
		}
	}

	// User didn't pre-exist and there's no room
	// Return without sidero user :(
	if sideroUserID == 0 {
		return errors.New("no slot available for sidero user")
	}

	// Not already present and there's an empty slot so we add sidero user
	if !exists {
		log.Printf("Adding sidero user to slot %d\n", sideroUserID)

		_, err := ipmiClient.SetUserName(sideroUserID, "sidero")
		if err != nil {
			return err
		}
	}

	// Reset pass for sidero user
	// nb: we _always_ reset the user pass because we can't ever get
	//     it back out when we find an existing sidero user.
	pass, err := genPass16()
	if err != nil {
		return err
	}

	_, err = ipmiClient.SetUserPass(sideroUserID, pass)
	if err != nil {
		return err
	}

	// Make sidero an admin
	// Options: 0x91 == Callin true, Link false, IPMI Msg true, Channel 1
	// Limits: 0x03 == Administrator
	// Session: 0x00 No session limit
	_, err = ipmiClient.SetUserAccess(0x91, sideroUserID, 0x04, 0x00)
	if err != nil {
		return err
	}

	// Enable the sidero user
	_, err = ipmiClient.EnableUser(sideroUserID)
	if err != nil {
		return err
	}

	// Finally fill in info for update request
	bmcInfo.User = "sidero"
	bmcInfo.Pass = pass

	// Attempt to update server object
	err = retry.Constant(5*time.Minute, retry.WithUnits(30*time.Second), retry.WithErrorLogging(true)).Retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		_, err = client.UpdateBMCInfo(
			ctx,
			&api.UpdateBMCInfoRequest{
				Uuid:    uuid.String(),
				BmcInfo: bmcInfo,
			},
		)

		if err != nil {
			return retry.ExpectedError(err)
		}

		return nil
	})

	return nil
}

// Returns a random pass string of len 16.
func genPass16() (string, error) {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, 16)
	for i := range b {
		rando, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(letterRunes))),
		)
		if err != nil {
			return "", err
		}

		b[i] = letterRunes[rando.Int64()]
	}

	return string(b), nil
}
