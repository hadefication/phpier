package cmd

import (
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start global services and project container (shortcut for 'up -d')",
	Long: `Start the global services and project container in detached mode.

This command is a convenient shortcut for 'phpier up -d' and will:
- Start the global services (Traefik, databases, etc.) if not running
- Start the PHP/Nginx container for the current project in detached mode
- Connect the project to the global services network

Examples:
  phpier start    # Start global services and project container in detached mode`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Set detached to true for start command
		detached = true
		return runUp(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	// No flags - keep it simple
}
