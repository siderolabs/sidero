// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import (
	"fmt"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// AcceptedServerFilter matches Servers that have Spec.Accepted set to true.
func AcceptedServerFilter(s Server) (bool, error) {
	return s.Spec.Accepted, nil
}

// NotCordonedServerFilter matches Servers that have Spec.Paused set to false.
func NotCordonedServerFilter(s Server) (bool, error) {
	return !s.Spec.Cordoned, nil
}

// SelectorFilter returns a ServerFilter that matches servers against the
// serverclass's selector field.
func (sc *ServerClass) SelectorFilter() func(Server) (bool, error) {
	return func(server Server) (bool, error) {
		s, err := metav1.LabelSelectorAsSelector(&sc.Spec.Selector)
		if err != nil {
			return false, fmt.Errorf("failed to get selector from labelselector: %v", err)
		}

		return s.Matches(labels.Set(server.GetLabels())), nil
	}
}

// QualifiersFilter returns a ServerFilter that matches servers against the
// serverclass's qualifiers field.
func (sc *ServerClass) QualifiersFilter() func(Server) (bool, error) {
	return func(server Server) (bool, error) {
		q := sc.Spec.Qualifiers

		// check hardware qualifiers if they are present
		if filters := q.Hardware; len(filters) > 0 {
			var match bool

			for _, filter := range filters {
				if info := server.Spec.Hardware; info != nil && filter.PartialEqual(info) {
					match = true
					break
				}
			}

			if !match {
				return false, nil
			}
		}

		if filters := q.LabelSelectors; len(filters) > 0 {
			var match bool

			for _, filter := range filters {
				for labelKey, labelVal := range filter {
					if val, ok := server.ObjectMeta.Labels[labelKey]; ok && labelVal == val {
						match = true
						break
					}
				}
			}

			if !match {
				return false, nil
			}
		}

		return true, nil
	}
}

// FilterServers returns the subset of servers that pass all provided filters.
// In case of error the returned slice will be nil.
func FilterServers(servers []Server, filters ...func(Server) (bool, error)) ([]Server, error) {
	matches := make([]Server, 0, len(servers))

	for _, server := range servers {
		match := true

		for _, filter := range filters {
			var err error

			match, err = filter(server)
			if err != nil {
				return nil, fmt.Errorf("failed to filter server: %v", err)
			}

			if !match {
				break
			}
		}

		if match {
			matches = append(matches, server)
		}
	}

	sort.Slice(matches, func(i, j int) bool { return matches[i].Name < matches[j].Name })

	return matches, nil
}
