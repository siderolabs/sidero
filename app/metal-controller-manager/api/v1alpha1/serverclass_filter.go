// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha1

import (
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// FilterAcceptedServers returns a new slice of Servers that are accepted and qualify.
//
// Returned Servers are always sorted by name for stable results.
func FilterAcceptedServers(servers []Server, q Qualifiers) []Server {
	res := make([]Server, 0, len(servers))

	for _, server := range servers {
		// skip non-accepted servers
		if !server.Spec.Accepted {
			continue
		}

		// check CPU qualifiers if they are present
		if filters := q.CPU; len(filters) > 0 {
			var match bool

			for _, filter := range filters {
				if cpu := server.Spec.CPU; cpu != nil && filter.PartialEqual(cpu) {
					match = true
					break
				}
			}

			if !match {
				continue
			}
		}

		if filters := q.SystemInformation; len(filters) > 0 {
			var match bool

			for _, filter := range filters {
				if sysInfo := server.Spec.SystemInformation; sysInfo != nil && filter.PartialEqual(sysInfo) {
					match = true
					break
				}
			}

			if !match {
				continue
			}
		}

		if filters := q.LabelSelectors; len(filters) > 0 {
			var match bool

			for _, filter := range filters {
				ls := metav1.SetAsLabelSelector(labels.Set(filter))
				s, err := metav1.LabelSelectorAsSelector(ls)
				if err != nil {
					// TODO: At least log an error.
					continue
				}

				if !s.Matches(labels.Set(server.GetLabels())) {
					continue
				}

				match = true
				break
			}

			if !match {
				continue
			}
		}

		res = append(res, server)
	}

	sort.Slice(res, func(i, j int) bool { return res[i].Name < res[j].Name })

	return res
}
