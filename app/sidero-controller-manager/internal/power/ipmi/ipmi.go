// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipmi

import (
	"fmt"

	goipmi "github.com/pensando/goipmi"

	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/internal/power/metal"
)

// Link to the IPMI spec: https://www.intel.com/content/dam/www/public/us/en/documents/product-briefs/ipmi-second-gen-interface-spec-v2-rev1-1.pdf
// Referenced in some of the comments below
// Future TODO(rsmitty): support "channels" other than #1

// Client is a holder for the IPMIClient.
type Client struct {
	IPMIClient *goipmi.Client
}

// NewClient creates an ipmi client to use.
func NewClient(bmcInfo metalv1.BMC) (*Client, error) {
	conn := &goipmi.Connection{
		Hostname:  bmcInfo.Endpoint,
		Port:      int(bmcInfo.Port),
		Username:  bmcInfo.User,
		Password:  bmcInfo.Pass,
		Interface: bmcInfo.Interface,
	}

	ipmiClient, err := goipmi.NewClient(conn)
	if err != nil {
		return nil, err
	}

	if err = ipmiClient.Open(); err != nil {
		return nil, fmt.Errorf("error opening client: %w", err)
	}

	return &Client{IPMIClient: ipmiClient}, nil
}

// Close the client.
func (c *Client) Close() error {
	return c.IPMIClient.Close()
}

// PowerOn will power on a given machine.
func (c *Client) PowerOn() error {
	return c.IPMIClient.Control(goipmi.ControlPowerUp)
}

// PowerOff will power off a given machine.
func (c *Client) PowerOff() error {
	return c.IPMIClient.Control(goipmi.ControlPowerDown)
}

// IsPoweredOn checks current power state.
func (c *Client) IsPoweredOn() (bool, error) {
	status, err := c.Status()
	if err != nil {
		return false, err
	}

	return status.IsSystemPowerOn(), nil
}

// PowerCycle will power cycle a given machine.
func (c *Client) PowerCycle() error {
	return c.IPMIClient.Control(goipmi.ControlPowerCycle)
}

// Status fetches the chassis status.
func (c *Client) Status() (*goipmi.ChassisStatusResponse, error) {
	req := &goipmi.Request{
		NetworkFunction: goipmi.NetworkFunctionChassis,
		Command:         goipmi.CommandChassisStatus,
		Data:            goipmi.ChassisStatusRequest{},
	}

	res := &goipmi.ChassisStatusResponse{}

	err := c.IPMIClient.Send(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// SetPXE makes sure the node will pxe boot next time.
func (c *Client) SetPXE(mode metal.PXEMode) error {
	switch mode {
	case metal.PXEModeBIOS:
		return c.IPMIClient.SetBootDevice(goipmi.BootDevicePxe)
	case metal.PXEModeUEFI:
		return c.IPMIClient.SetBootDeviceEFI(goipmi.BootDevicePxe)
	default:
		return fmt.Errorf("unsupported mode %q", mode)
	}
}

// GetLANConfig fetches a given param from the LAN Config. (see 23.2).
func (c *Client) GetLANConfig(param uint8) (*goipmi.LANConfigResponse, error) {
	req := &goipmi.Request{
		NetworkFunction: goipmi.NetworkFunctionTransport,
		Command:         goipmi.CommandGetLANConfig,
		Data: &goipmi.LANConfigRequest{
			ChannelNumber: 0x01,
			Param:         param,
		},
	}

	res := &goipmi.LANConfigResponse{}

	err := c.IPMIClient.Send(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

//
// User management functions
//

// GetUserSummary returns stats about user table, including max users allowed.
func (c *Client) GetUserSummary() (*goipmi.GetUserSummaryResponse, error) {
	req := &goipmi.Request{
		NetworkFunction: goipmi.NetworkFunctionApp,
		Command:         goipmi.CommandGetUserSummary,
		Data: &goipmi.GetUserSummaryRequest{
			ChannelNumber: 0x01,
			UserID:        0x01,
		},
	}

	res := &goipmi.GetUserSummaryResponse{}

	err := c.IPMIClient.Send(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetUserName fetches a un string given a uid. This is how we check if a user slot is available.
// nb: a "failure" here can actually mean that the slot is just open for use
//     or you can also have a user with "" as the name which won't
//     fail this check and is still open for use.
// (see 22.29).
func (c *Client) GetUserName(uid byte) (*goipmi.GetUserNameResponse, error) {
	req := &goipmi.Request{
		NetworkFunction: goipmi.NetworkFunctionApp,
		Command:         goipmi.CommandGetUserName,
		Data: &goipmi.GetUserNameRequest{
			UserID: uid,
		},
	}

	res := &goipmi.GetUserNameResponse{}

	err := c.IPMIClient.Send(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// SetUserName sets a string for the given uid (see 22.28).
func (c *Client) SetUserName(uid byte, name string) (*goipmi.SetUserNameResponse, error) {
	req := &goipmi.Request{
		NetworkFunction: goipmi.NetworkFunctionApp,
		Command:         goipmi.CommandSetUserName,
		Data: &goipmi.SetUserNameRequest{
			UserID:   uid,
			Username: name,
		},
	}

	res := &goipmi.SetUserNameResponse{}

	err := c.IPMIClient.Send(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// SetUserPass sets the password for a given uid (see 22.30).
// nb: This naively assumes you'll pass a 16 char or less pw string.
//     The goipmi function does not support longer right now.
func (c *Client) SetUserPass(uid byte, pass string) (*goipmi.SetUserPassResponse, error) {
	req := &goipmi.Request{
		NetworkFunction: goipmi.NetworkFunctionApp,
		Command:         goipmi.CommandSetUserPass,
		Data: &goipmi.SetUserPassRequest{
			UserID: uid,
			Pass:   []byte(pass),
		},
	}

	res := &goipmi.SetUserPassResponse{}

	err := c.IPMIClient.Send(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// SetUserAccess tweaks the privileges for a given uid (see 22.26).
func (c *Client) SetUserAccess(options, uid, limits, session byte) (*goipmi.SetUserAccessResponse, error) {
	req := &goipmi.Request{
		NetworkFunction: goipmi.NetworkFunctionApp,
		Command:         goipmi.CommandSetUserAccess,
		Data: &goipmi.SetUserAccessRequest{
			AccessOptions:    options,
			UserID:           uid,
			UserLimits:       limits,
			UserSessionLimit: session,
		},
	}

	res := &goipmi.SetUserAccessResponse{}

	err := c.IPMIClient.Send(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// EnableUser sets a user as enabled. Actually the same underlying command as SetUserPass (see 22.30).
func (c *Client) EnableUser(uid byte) (*goipmi.EnableUserResponse, error) {
	req := &goipmi.Request{
		NetworkFunction: goipmi.NetworkFunctionApp,
		Command:         goipmi.CommandEnableUser,
		Data: &goipmi.EnableUserRequest{
			UserID: uid,
		},
	}

	res := &goipmi.EnableUserResponse{}

	err := c.IPMIClient.Send(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// IsFake returns false.
func (c *Client) IsFake() bool {
	return false
}
