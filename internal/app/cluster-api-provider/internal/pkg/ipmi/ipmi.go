package ipmi

import (
	goipmi "github.com/vmware/goipmi"

	metalv1alpha1 "github.com/talos-systems/sidero/internal/app/metal-controller-manager/api/v1alpha1"
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

// PowerOff will power on a given machine.
func (c *Client) PowerOff() error {
	return c.IPMIClient.Control(goipmi.ControlPowerDown)
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
	return c.IPMIClient.SetBootDevice(goipmi.BootDevicePxe)
}
