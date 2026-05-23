// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	metal "github.com/siderolabs/sidero/app/sidero-controller-manager/api/v1alpha2"
)

func TestServerInstallDiskValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		spec      metal.ServerSpec
		wantError bool
	}{
		{name: "default valid", spec: metal.ServerSpec{}},
		{name: "wipe only requires install disk", spec: metal.ServerSpec{WipeOnlyInstallDisk: true}, wantError: true},
		{name: "device name valid", spec: metal.ServerSpec{InstallDisk: &metal.InstallDisk{DeviceName: "/dev/disk/by-id/test"}}},
		{name: "selector valid", spec: metal.ServerSpec{InstallDisk: &metal.InstallDisk{Selector: &metal.InstallDiskSelector{Serial: "abc"}}}},
		{name: "device and selector invalid", spec: metal.ServerSpec{InstallDisk: &metal.InstallDisk{DeviceName: "/dev/sda", Selector: &metal.InstallDiskSelector{Serial: "abc"}}}, wantError: true},
		{name: "empty selector invalid", spec: metal.ServerSpec{InstallDisk: &metal.InstallDisk{Selector: &metal.InstallDiskSelector{}}}, wantError: true},
		{name: "invalid selector type", spec: metal.ServerSpec{InstallDisk: &metal.InstallDisk{Selector: &metal.InstallDiskSelector{Type: "flash"}}}, wantError: true},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			obj := &metal.Server{Spec: tt.spec}
			_, err := obj.ValidateCreate(context.Background(), obj)

			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestServerClassInstallDiskValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		spec      metal.ServerClassSpec
		wantError bool
	}{
		{name: "default valid", spec: metal.ServerClassSpec{}},
		{name: "wipe only requires install disk", spec: metal.ServerClassSpec{WipeOnlyInstallDisk: true}, wantError: true},
		{name: "device name valid", spec: metal.ServerClassSpec{InstallDisk: &metal.InstallDisk{DeviceName: "/dev/disk/by-id/test"}}},
		{name: "selector valid", spec: metal.ServerClassSpec{InstallDisk: &metal.InstallDisk{Selector: &metal.InstallDiskSelector{Serial: "abc"}}}},
		{name: "device and selector invalid", spec: metal.ServerClassSpec{InstallDisk: &metal.InstallDisk{DeviceName: "/dev/sda", Selector: &metal.InstallDiskSelector{Serial: "abc"}}}, wantError: true},
		{name: "empty selector invalid", spec: metal.ServerClassSpec{InstallDisk: &metal.InstallDisk{Selector: &metal.InstallDiskSelector{}}}, wantError: true},
		{name: "invalid selector type", spec: metal.ServerClassSpec{InstallDisk: &metal.InstallDisk{Selector: &metal.InstallDiskSelector{Type: "flash"}}}, wantError: true},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			obj := &metal.ServerClass{Spec: tt.spec}
			_, err := obj.ValidateCreate(context.Background(), obj)

			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
