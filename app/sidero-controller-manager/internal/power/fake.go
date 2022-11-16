// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package power

import "github.com/talos-systems/sidero/app/sidero-controller-manager/pkg/types"

type fakeClient struct{}

func (fakeClient) PowerOn() error {
	return nil
}

func (fakeClient) PowerOff() error {
	return nil
}

func (fakeClient) PowerCycle() error {
	return nil
}

func (fakeClient) SetPXE(mode types.PXEMode) error {
	return nil
}

func (fakeClient) IsPoweredOn() (bool, error) {
	return true, nil
}

func (fakeClient) IsFake() bool {
	return true
}

func (fakeClient) Close() error {
	return nil
}
