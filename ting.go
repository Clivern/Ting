// Copyright 2026 Ting. All rights reserved.
// License can be found in the LICENSE file.

// Package main is the Ting entry point.
package main

import "github.com/clivern/ting/cli"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

// main is the application entry point.
func main() {
	cli.Version = version
	cli.Commit = commit
	cli.Date = date
	cli.BuiltBy = builtBy

	cli.Execute()
}
