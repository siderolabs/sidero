// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"log"
	"math/big"
	"net"
	"time"

	"github.com/siderolabs/go-retry/retry"
	"github.com/siderolabs/go-smbios/smbios"

	metalv1 "github.com/siderolabs/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/api"
	"github.com/siderolabs/sidero/app/sidero-controller-manager/internal/power/ipmi"
)

func attemptBMCIP(ctx context.Context, client api.AgentClient, s *smbios.SMBIOS) error {
	bmcInfo := &api.BMCInfo{}

	// Create "open" client
	bmcSpec := metalv1.BMC{
		Interface: "open",
	}

	ipmiClient, err := ipmi.NewClient(bmcSpec)
	if err != nil {
		return err
	}

	defer ipmiClient.Close() //nolint:errcheck

	// Fetch BMC IP (param 3 in LAN config)
	ipResp, err := ipmiClient.GetLANConfig(0x03)
	if err != nil {
		return err
	}

	bmcIP := net.IP(ipResp.Data)
	bmcInfo.Ip = bmcIP.String()

	// Fetch BMC Port (param 8 in LAN config)
	portResp, err := ipmiClient.GetLANConfig(0x08)
	if err != nil {
		return err
	}

	// Port is only a 16bit piece of data,
	// but the smallest protobuf supports is 32bit, so we have this little conversion.
	bmcInfo.Port = uint32(binary.LittleEndian.Uint16(portResp.Data))

	// Attempt to update server object
	err = retry.Constant(5*time.Minute, retry.WithUnits(30*time.Second), retry.WithErrorLogging(true)).Retry(func() error {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		_, err = client.UpdateBMCInfo(
			ctx,
			&api.UpdateBMCInfoRequest{
				Uuid:    s.SystemInformation.UUID,
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
	bmcInfo := &api.BMCInfo{}

	// Create "open" client
	bmcSpec := metalv1.BMC{
		Interface: "open",
	}

	ipmiClient, err := ipmi.NewClient(bmcSpec)
	if err != nil {
		return err
	}

	defer ipmiClient.Close() //nolint:errcheck

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
				Uuid:    s.SystemInformation.UUID,
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
