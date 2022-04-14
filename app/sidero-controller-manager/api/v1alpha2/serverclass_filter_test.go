// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	metal "github.com/talos-systems/sidero/app/sidero-controller-manager/api/v1alpha2"
)

func TestFilterAcceptedServers(t *testing.T) {
	t.Parallel()

	atom := metal.Server{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"common-label": "true",
				"zone":         "central",
			},
		},
		Spec: metal.ServerSpec{
			Accepted: true,
			Hardware: &metal.HardwareInformation{
				Compute: &metal.ComputeInformation{
					TotalCoreCount:   4,
					TotalThreadCount: 4,
					ProcessorCount:   1,
					Processors: []*metal.Processor{
						{
							Manufacturer: "Intel(R) Corporation",
							ProductName:  "Intel(R) Atom(TM) CPU C3558 @ 2.20GHz",
							SerialNumber: "",
							Speed:        2200,
							CoreCount:    4,
							ThreadCount:  4,
						},
					},
				},
			},
		},
	}
	dualXeon := metal.Server{
		Spec: metal.ServerSpec{
			Accepted: true,
			Hardware: &metal.HardwareInformation{
				Compute: &metal.ComputeInformation{
					TotalCoreCount:   16,
					TotalThreadCount: 32,
					ProcessorCount:   2,
					Processors: []*metal.Processor{
						{
							Manufacturer: "Intel(R) Corporation",
							ProductName:  "Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz",
							SerialNumber: "",
							Speed:        2400,
							CoreCount:    8,
							ThreadCount:  16,
						},
						{
							Manufacturer: "Intel(R) Corporation",
							ProductName:  "Intel(R) Xeon(R) CPU E5-2630 v3 @ 2.40GHz",
							SerialNumber: "",
							Speed:        2400,
							CoreCount:    8,
							ThreadCount:  16,
						},
					},
				},
			},
		},
	}
	ryzen := metal.Server{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"common-label": "true",
				"zone":         "east",
			},
		},
		Spec: metal.ServerSpec{
			Accepted: true,
			Hardware: &metal.HardwareInformation{
				System: &metal.SystemInformation{
					Manufacturer: "QEMU",
				},
				Compute: &metal.ComputeInformation{
					TotalCoreCount:   8,
					TotalThreadCount: 16,
					ProcessorCount:   1,
					Processors: []*metal.Processor{
						{
							Manufacturer: "Advanced Micro Devices, Inc.",
							ProductName:  "AMD Ryzen 7 2700X Eight-Core Processor",
							SerialNumber: "",
							Speed:        3700,
							CoreCount:    8,
							ThreadCount:  16,
						},
					},
				},
			},
		},
	}
	notAccepted := metal.Server{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"common-label": "true",
				"zone":         "west",
			},
		},
		Spec: metal.ServerSpec{
			Accepted: false,
			Hardware: &metal.HardwareInformation{
				System: &metal.SystemInformation{
					Manufacturer: "QEMU",
				},
				Compute: &metal.ComputeInformation{
					TotalCoreCount:   8,
					TotalThreadCount: 16,
					ProcessorCount:   1,
					Processors: []*metal.Processor{
						{
							Manufacturer: "Advanced Micro Devices, Inc.",
							ProductName:  "AMD Ryzen 7 2700X Eight-Core Processor",
							SerialNumber: "",
							Speed:        3700,
							CoreCount:    8,
							ThreadCount:  16,
						},
					},
				},
			},
		},
	}

	servers := []metal.Server{atom, dualXeon, ryzen, notAccepted}

	testdata := map[string]struct {
		s        metav1.LabelSelector
		q        metal.Qualifiers
		expected []metal.Server
	}{
		"empty selector - empty qualifier": {
			// Matches all servers
			expected: []metal.Server{atom, dualXeon, ryzen},
		},
		"Intel only": {
			q: metal.Qualifiers{
				Hardware: []metal.HardwareInformation{
					{
						Compute: &metal.ComputeInformation{
							Processors: []*metal.Processor{
								{
									Manufacturer: "Intel(R) Corporation",
								},
							},
						},
					},
				},
			},
			expected: []metal.Server{atom, dualXeon},
		},
		"Intel and AMD": {
			q: metal.Qualifiers{
				Hardware: []metal.HardwareInformation{
					{
						Compute: &metal.ComputeInformation{
							Processors: []*metal.Processor{
								{
									Manufacturer: "Intel(R) Corporation",
								},
							},
						},
					},
					{
						Compute: &metal.ComputeInformation{
							Processors: []*metal.Processor{
								{
									Manufacturer: "Advanced Micro Devices, Inc.",
								},
							},
						},
					},
				},
			},
			expected: []metal.Server{atom, dualXeon, ryzen},
		},
		"QEMU only": {
			q: metal.Qualifiers{
				Hardware: []metal.HardwareInformation{
					{
						System: &metal.SystemInformation{
							Manufacturer: "QEMU",
						},
					},
				},
			},
			expected: []metal.Server{ryzen},
		},
		"with label": {
			q: metal.Qualifiers{
				LabelSelectors: []map[string]string{
					{
						"common-label": "true",
					},
				},
			},
			expected: []metal.Server{atom, ryzen},
		},
		// This should probably only return atom. Leaving it as-is to
		// avoid breaking changes before we remove LabelSelectors in
		// favor of Selector.
		"with multiple labels - single selector": {
			q: metal.Qualifiers{
				LabelSelectors: []map[string]string{
					{
						"common-label": "true",
						"zone":         "central",
					},
				},
			},
			expected: []metal.Server{atom, ryzen},
		},
		"with multiple labels - multiple selectors": {
			q: metal.Qualifiers{
				LabelSelectors: []map[string]string{
					{
						"common-label": "true",
					},
					{
						"zone": "central",
					},
				},
			},
			expected: []metal.Server{atom, ryzen},
		},
		"with same label key different label value": {
			q: metal.Qualifiers{
				LabelSelectors: []map[string]string{
					{
						"zone": "central",
					},
				},
			},
			expected: []metal.Server{atom},
		},
		"selector - single MatchLabels single result": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"zone": "central",
				},
			},
			expected: []metal.Server{atom},
		},
		"selector - single MatchLabels multiple results": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"common-label": "true",
				},
			},
			expected: []metal.Server{atom, ryzen},
		},
		"selector - multiple MatchLabels": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"zone":         "central",
					"common-label": "true",
				},
			},
			expected: []metal.Server{atom},
		},
		"selector - MatchExpressions common label key": {
			s: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "common-label",
						Operator: "Exists",
					},
				},
			},
			expected: []metal.Server{atom, ryzen},
		},
		"selector - MatchExpressions multiple values": {
			s: metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "zone",
						Operator: "In",
						Values: []string{
							"east",
							"west",
						},
					},
				},
			},
			expected: []metal.Server{ryzen},
		},
		"selector and qualifiers both match": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"common-label": "true",
				},
			},
			q: metal.Qualifiers{
				Hardware: []metal.HardwareInformation{
					{
						System: &metal.SystemInformation{
							Manufacturer: "QEMU",
						},
					},
				},
			},
			expected: []metal.Server{ryzen},
		},
		"selector and qualifiers with disqualifying selector": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"common-label": "no-match",
				},
			},
			q: metal.Qualifiers{
				Hardware: []metal.HardwareInformation{
					{
						System: &metal.SystemInformation{
							Manufacturer: "QEMU",
						},
					},
				},
			},
			expected: []metal.Server{},
		},
		"selector and qualifiers with disqualifying qualifier": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"common-label": "true",
				},
			},
			q: metal.Qualifiers{
				Hardware: []metal.HardwareInformation{
					{
						System: &metal.SystemInformation{
							Manufacturer: "Gateway",
						},
					},
				},
			},
			expected: []metal.Server{},
		},
		metal.ServerClassAny: {
			expected: []metal.Server{atom, dualXeon, ryzen},
		},
	}

	for name, td := range testdata {
		name, td := name, td
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sc := &metal.ServerClass{
				Spec: metal.ServerClassSpec{
					Selector:   td.s,
					Qualifiers: td.q,
				},
			}
			actual, err := metal.FilterServers(servers,
				metal.AcceptedServerFilter,
				sc.SelectorFilter(),
				sc.QualifiersFilter(),
			)
			assert.NoError(t, err)
			assert.Equal(t, td.expected, actual)
		})
	}
}
