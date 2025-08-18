package main

import (
	"phpier/cmd"
)

// Build information (set via ldflags during build)
var (
	version = "v1.0.5"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	// Set version information for the CLI
	cmd.SetVersionInfo(version, commit, date)

	// Execute command - error handling is done in cmd.Execute()
	cmd.Execute()
}
