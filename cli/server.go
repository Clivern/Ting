// Copyright 2026 Ting. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/clivern/ting/core"

	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Ting proxy server",
	Run: func(_ *cobra.Command, _ []string) {
		err := core.Load(config)
		if err != nil {
			panic(err.Error())
		}

		err = core.SetupLogging()
		if err != nil {
			panic(err.Error())
		}

		err = core.RunServer()
		if err != nil {
			panic(fmt.Sprintf("Server error: %s", err.Error()))
		}
	},
}

// init registers the server subcommand and flags.
func init() {
	serverCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.dist.yml",
		"Absolute path to config file (required)",
	)
	if err := serverCmd.MarkFlagRequired("config"); err != nil {
		panic(err.Error())
	}
	rootCmd.AddCommand(serverCmd)
}
