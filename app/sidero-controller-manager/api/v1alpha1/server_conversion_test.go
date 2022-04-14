// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	metalv1alpha1 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha1"
	metalv1alpha2 "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
)

func TestServerConvertV1alpha1V1Alpha2(t *testing.T) {
	src := &metalv1alpha1.Server{
		Spec: metalv1alpha1.ServerSpec{
			Hostname: "example.com",
			SystemInformation: &metalv1alpha1.SystemInformation{
				Manufacturer: "Sidero",
				ProductName:  "Server",
				Version:      "v1.0",
			},
			CPU: &metalv1alpha1.CPUInformation{
				Manufacturer: "Sidero CPU",
				Version:      "v1",
			},
		},
	}
	dst := &metalv1alpha2.Server{}

	require.NoError(t, src.ConvertTo(dst))

	assert.Equal(t, "example.com", dst.Spec.Hostname)
	assert.Equal(t,
		&metalv1alpha2.SystemInformation{
			Manufacturer: "Sidero",
			ProductName:  "Server",
			Version:      "v1.0",
		},
		dst.Spec.Hardware.System,
	)
	assert.Equal(t,
		&metalv1alpha2.ComputeInformation{
			Processors: []*metalv1alpha2.Processor{
				{
					Manufacturer: "Sidero CPU",
					ProductName:  "v1",
				},
			},
		},
		dst.Spec.Hardware.Compute,
	)
}

func TestServerConvertV1alpha2V1Alpha1(t *testing.T) {
	src := &metalv1alpha2.Server{
		Spec: metalv1alpha2.ServerSpec{
			Hostname: "example.com",
			Hardware: &metalv1alpha2.HardwareInformation{
				System: &metalv1alpha2.SystemInformation{
					Manufacturer: "Sidero",
					ProductName:  "Server",
					Version:      "v1.0",
				},
				Compute: &metalv1alpha2.ComputeInformation{
					Processors: []*metalv1alpha2.Processor{
						{
							Manufacturer: "Sidero CPU",
							ProductName:  "v1",
						},
					},
				},
			},
		},
	}
	dst := &metalv1alpha1.Server{}

	require.NoError(t, dst.ConvertFrom(src))

	assert.Equal(t, "example.com", dst.Spec.Hostname)
	assert.Equal(t,
		&metalv1alpha1.SystemInformation{
			Manufacturer: "Sidero",
			ProductName:  "Server",
			Version:      "v1.0",
		},
		dst.Spec.SystemInformation,
	)
	assert.Equal(t,
		&metalv1alpha1.CPUInformation{
			Manufacturer: "Sidero CPU",
			Version:      "v1",
		},
		dst.Spec.CPU,
	)
}
