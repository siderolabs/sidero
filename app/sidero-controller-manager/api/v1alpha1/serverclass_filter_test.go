// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	metalv1alpha1 "github.com/siderolabs/sidero/app/sidero-controller-manager/api/v1alpha1"
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
				"zone":         "west",
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
		s        metav1.LabelSelector
		q        metalv1alpha1.Qualifiers
		expected []metalv1alpha1.Server
	}{
		"empty selector - empty qualifier": {
			// Matches all servers
			expected: []metalv1alpha1.Server{atom, ryzen},
		},
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
		// This should probably only return atom. Leaving it as-is to
		// avoid breaking changes before we remove LabelSelectors in
		// favor of Selector.
		"with multiple labels - single selector": {
			q: metalv1alpha1.Qualifiers{
				LabelSelectors: []map[string]string{
					{
						"common-label": "true",
						"zone":         "central",
					},
				},
			},
			expected: []metalv1alpha1.Server{atom, ryzen},
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
		"with same label key different label value": {
			q: metalv1alpha1.Qualifiers{
				LabelSelectors: []map[string]string{
					{
						"zone": "central",
					},
				},
			},
			expected: []metalv1alpha1.Server{atom},
		},
		"selector - single MatchLabels single result": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"zone": "central",
				},
			},
			expected: []metalv1alpha1.Server{atom},
		},
		"selector - single MatchLabels multiple results": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"common-label": "true",
				},
			},
			expected: []metalv1alpha1.Server{atom, ryzen},
		},
		"selector - multiple MatchLabels": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"zone":         "central",
					"common-label": "true",
				},
			},
			expected: []metalv1alpha1.Server{atom},
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
			expected: []metalv1alpha1.Server{atom, ryzen},
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
			expected: []metalv1alpha1.Server{ryzen},
		},
		"selector and qualifiers both match": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"common-label": "true",
				},
			},
			q: metalv1alpha1.Qualifiers{
				SystemInformation: []metalv1alpha1.SystemInformation{
					{
						Manufacturer: "QEMU",
					},
				},
			},
			expected: []metalv1alpha1.Server{ryzen},
		},
		"selector and qualifiers with disqualifying selector": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"common-label": "no-match",
				},
			},
			q: metalv1alpha1.Qualifiers{
				SystemInformation: []metalv1alpha1.SystemInformation{
					{
						Manufacturer: "QEMU",
					},
				},
			},
			expected: []metalv1alpha1.Server{},
		},
		"selector and qualifiers with disqualifying qualifier": {
			s: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"common-label": "true",
				},
			},
			q: metalv1alpha1.Qualifiers{
				SystemInformation: []metalv1alpha1.SystemInformation{
					{
						Manufacturer: "Gateway",
					},
				},
			},
			expected: []metalv1alpha1.Server{},
		},
		metalv1alpha1.ServerClassAny: {
			expected: []metalv1alpha1.Server{atom, ryzen},
		},
	}

	for name, td := range testdata {
		name, td := name, td
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			sc := &metalv1alpha1.ServerClass{
				Spec: metalv1alpha1.ServerClassSpec{
					Selector:   td.s,
					Qualifiers: td.q,
				},
			}
			actual, err := metalv1alpha1.FilterServers(servers,
				metalv1alpha1.AcceptedServerFilter,
				sc.SelectorFilter(),
				sc.QualifiersFilter(),
			)
			assert.NoError(t, err)
			assert.Equal(t, td.expected, actual)
		})
	}
}
