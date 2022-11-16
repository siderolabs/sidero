// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package api provides metal machine management via API.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	metalv1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
	"github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/types"
)

// Client provides management over simple API.
type Client struct {
	endpoint string
}

// NewClient returns new API client to manage metal machine.
func NewClient(spec metalv1.ManagementAPI) (*Client, error) {
	return &Client{
		endpoint: spec.Endpoint,
	}, nil
}

// Close the client.
func (c *Client) Close() error {
	return nil
}

func (c *Client) postRequest(path string) error {
	failureMode := DefaultDice.Roll()

	switch failureMode { //nolint:exhaustive
	case ExplicitFailure:
		return fmt.Errorf("simulated failure from the power management")
	case SilentFailure:
		// don't do anything
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s%s", c.endpoint, path), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.Body != nil {
		defer func() {
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}()
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API error: %s", resp.Status)
	}

	return nil
}

// PowerOn will power on a given machine.
func (c *Client) PowerOn() error {
	return c.postRequest("/poweron")
}

// PowerOff will power off a given machine.
func (c *Client) PowerOff() error {
	return c.postRequest("/poweroff")
}

// PowerCycle will power cycle a given machine.
func (c *Client) PowerCycle() error {
	return c.postRequest("/reboot")
}

// SetPXE makes sure the node will pxe boot next time.
func (c *Client) SetPXE(mode types.PXEMode) error {
	// no way to enforce mode via QEMU API
	return c.postRequest("/pxeboot")
}

// IsPoweredOn checks current power state.
func (c *Client) IsPoweredOn() (bool, error) {
	failureMode := DefaultDice.Roll()

	switch failureMode { //nolint:exhaustive
	case ExplicitFailure:
		return false, fmt.Errorf("simulated failure from the power management")
	case SilentFailure:
		return time.Now().Second()%2 == 0, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/status", c.endpoint), nil)
	if err != nil {
		return false, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}

	if resp.Body != nil {
		defer func() {
			_, _ = io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
		}()
	}

	var status struct {
		PoweredOn bool
	}

	if err = json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return false, err
	}

	return status.PoweredOn, nil
}

// IsFake returns false.
func (c *Client) IsFake() bool {
	return false
}
