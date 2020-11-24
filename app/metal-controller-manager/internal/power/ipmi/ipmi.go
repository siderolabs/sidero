// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package ipmi

import (
	goipmi "github.com/pensando/goipmi"

	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
)

// Note (rsmitty): This pkg is pretty sparse right now, but I wanted to go ahead and create
// it in case we want to do something more complex w/ IPMI in the future.

// Client is a holder for the IPMIClient.
type Client struct {
	IPMIClient *goipmi.Client
}

// NewClient creates an ipmi client to use.
func NewClient(bmcInfo metalv1alpha1.BMC) (*Client, error) {
	conn := &goipmi.Connection{
		Hostname:  bmcInfo.Endpoint,
		Username:  bmcInfo.User,
		Password:  bmcInfo.Pass,
		Interface: "lanplus",
	}

	ipmiClient, err := goipmi.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return &Client{IPMIClient: ipmiClient}, nil
}

// Note (rsmitty): I think checking this system power isn't really necessary, but we may want
// to make more complex power decisions later on.

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
func (c *Client) SetPXE() error {
	return c.IPMIClient.SetBootDeviceEFI(goipmi.BootDevicePxe)
}

// IsFake returns false.
func (c *Client) IsFake() bool {
	return false
}
