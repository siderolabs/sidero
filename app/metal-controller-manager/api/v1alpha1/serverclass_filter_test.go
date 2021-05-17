// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	metalv1alpha1 "github.com/talos-systems/sidero/app/metal-controller-manager/api/v1alpha1"
)

func TestFilterAcceptedServers(t *testing.T) {
	t.Parallel()

	atom := metalv1alpha1.Server{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"common-label": "true",
				"zone":         "central",
			},
		},
		Spec: metalv1alpha1.ServerSpec{
			Accepted: true,
			CPU: &metalv1alpha1.CPUInformation{
				Manufacturer: "Intel(R) Corporation",
				Version:      "Intel(R) Atom(TM) CPU C3558 @ 2.20GHz",
			},
		},
	}
	ryzen := metalv1alpha1.Server{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"common-label": "true",
				"zone":         "east",
			},
		},
		Spec: metalv1alpha1.ServerSpec{
			Accepted: true,
			CPU: &metalv1alpha1.CPUInformation{
				Manufacturer: "Advanced Micro Devices, Inc.",
				Version:      "AMD Ryzen 7 2700X Eight-Core Processor",
			},
			SystemInformation: &metalv1alpha1.SystemInformation{
				Manufacturer: "QEMU",
			},
		},
	}
	notAccepted := metalv1alpha1.Server{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"common-label": "true",
			},
		},
		Spec: metalv1alpha1.ServerSpec{
			Accepted: false,
			CPU: &metalv1alpha1.CPUInformation{
				Manufacturer: "Advanced Micro Devices, Inc.",
				Version:      "AMD Ryzen 7 2700X Eight-Core Processor",
			},
			SystemInformation: &metalv1alpha1.SystemInformation{
				Manufacturer: "QEMU",
			},
		},
	}

	servers := []metalv1alpha1.Server{atom, ryzen, notAccepted}

	testdata := map[string]struct {
		q        metalv1alpha1.Qualifiers
		expected []metalv1alpha1.Server
	}{
		"Intel only": {
			q: metalv1alpha1.Qualifiers{
				CPU: []metalv1alpha1.CPUInformation{
					{
						Manufacturer: "Intel(R) Corporation",
					},
				},
			},
			expected: []metalv1alpha1.Server{atom},
		},
		"QEMU only": {
			q: metalv1alpha1.Qualifiers{
				SystemInformation: []metalv1alpha1.SystemInformation{
					{
						Manufacturer: "QEMU",
					},
				},
			},
			expected: []metalv1alpha1.Server{ryzen},
		},
		"with label": {
			q: metalv1alpha1.Qualifiers{
				LabelSelectors: []map[string]string{
					{
						"common-label": "true",
					},
				},
			},
			expected: []metalv1alpha1.Server{atom, ryzen},
		},
		"with multiple labels - single selector": {
			q: metalv1alpha1.Qualifiers{
				LabelSelectors: []map[string]string{
					{
						"common-label": "true",
						"zone":         "central",
					},
				},
			},
			expected: []metalv1alpha1.Server{atom},
		},
		"with multiple labels - multiple selectors": {
			q: metalv1alpha1.Qualifiers{
				LabelSelectors: []map[string]string{
					{
						"common-label": "true",
					},
					{
						"zone": "central",
					},
				},
			},
			expected: []metalv1alpha1.Server{atom, ryzen},
		},
		metalv1alpha1.ServerClassAny: {
			expected: []metalv1alpha1.Server{atom, ryzen},
		},
	}

	for name, td := range testdata {
		name, td := name, td
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := metalv1alpha1.FilterAcceptedServers(servers, td.q)
			assert.Equal(t, actual, td.expected)
		})
	}
}
