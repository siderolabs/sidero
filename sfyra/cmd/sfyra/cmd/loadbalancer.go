// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"github.com/spf13/cobra"
)

var loadbalancerCmd = &cobra.Command{
	Use:   "loadbalancer",
	Short: "Manage load balancer for the control plane nodes.",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(loadbalancerCmd)
}
