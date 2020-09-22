// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import "strings"

type stringSlice []string

func (ss *stringSlice) String() string {
	return strings.Join(*ss, ",")
}

func (ss *stringSlice) Set(value string) error {
	*ss = append(*ss, strings.Split(value, ",")...)

	return nil
}
