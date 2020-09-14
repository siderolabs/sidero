// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package metal

type fakeClient struct{}

func (fakeClient) PowerOn() error {
	return nil
}

func (fakeClient) PowerOff() error {
	return nil
}

func (fakeClient) SetPXE() error {
	return nil
}

func (fakeClient) IsPoweredOn() (bool, error) {
	return true, nil
}
