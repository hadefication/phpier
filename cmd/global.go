package cmd

import (
	"github.com/spf13/cobra"
)

// globalCmd represents the base command for managing global services
var globalCmd = &cobra.Command{
	Use:   "global",
	Short: "Manage the global shared services stack (Traefik, databases, etc.)",
	Long: `The 'global' command provides a set of sub-commands to manage the lifecycle
of the shared services that are used across all phpier projects.

This includes services like Traefik, MySQL, PostgreSQL, Redis, and Mailpit.
These services run persistently in the background and are not tied to any
single project.`,
}

func init() {
	rootCmd.AddCommand(globalCmd)
}
