// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	debug "github.com/siderolabs/go-debug"
	"github.com/spf13/cobra"
)

const (
	debugAddr = ":9995"
)

var rootCmd = &cobra.Command{
	Use:   "sfyra",
	Short: "Sfyra is a tool to deploy Sidero and run integration tests against it.",
	Long:  ``,
}

// Execute root command.
func Execute() {
	go func() {
		debugLogFunc := func(msg string) {
			log.Print(msg)
		}
		if err := debug.ListenAndServe(context.TODO(), debugAddr, debugLogFunc); err != nil {
			log.Fatalf("failed to start debug server: %s", err)
		}
	}()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var options = DefaultOptions()
