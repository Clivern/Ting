// Copyright 2026 Ting. All rights reserved.
// License can be found in the LICENSE file.

// Package cli provides the Ting command line.
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var config string

var rootCmd = &cobra.Command{
	Use: "ting",
	Short: `Ting is a reverse proxy for the OpenRouter API


If you have any suggestions, bug reports, or annoyances please report
them to our issue tracker at <https://github.com/clivern/ting/issues>`,
}

// Execute runs the command line tool.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
