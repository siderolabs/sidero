// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// nolint: scopelint
package v1alpha1_test

import (
	"testing"

	"github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
)

func Test_PartialEqual(t *testing.T) {
	type args struct {
		a interface{}
		b interface{}
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "defaults are partially equal",
			args: args{
				a: &v1alpha1.CPUInformation{},
				b: &v1alpha1.CPUInformation{},
			},
			want: true,
		},
		{
			name: "is partially equal",
			args: args{
				a: &v1alpha1.CPUInformation{
					Manufacturer: "QEMU",
					// Skip Version to indicate that we don't want to compare it.
				},
				b: &v1alpha1.CPUInformation{
					Manufacturer: "QEMU",
					Version:      "1.2.0",
				},
			},
			want: true,
		},
		{
			name: "is not partially equal",
			args: args{
				a: &v1alpha1.CPUInformation{
					Manufacturer: "QEMU",
					Version:      "1.0.0",
				},
				b: &v1alpha1.CPUInformation{
					Manufacturer: "QEMU",
					Version:      "1.2.0",
				},
			},
			want: false,
		},
		{
			name: "partially equal value",
			args: args{
				a: v1alpha1.CPUInformation{
					Manufacturer: "QEMU",
				},
				b: v1alpha1.CPUInformation{
					Manufacturer: "QEMU",
					Version:      "1.2.0",
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v1alpha1.PartialEqual(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("PartialEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
