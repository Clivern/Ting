// Copyright 2026 Ting. All rights reserved.
// License can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version buildinfo item
	Version = "dev"
	// Commit buildinfo item
	Commit = "none"
	// Date buildinfo item
	Date = "unknown"
	// BuiltBy buildinfo item
	BuiltBy = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf(
			"Current Ting version %v commit %v, built @%v by %v.\n",
			Version,
			Commit,
			Date,
			BuiltBy,
		)
	},
}

// init registers the version subcommand.
func init() {
	rootCmd.AddCommand(versionCmd)
}
