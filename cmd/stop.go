package cmd

import (
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop project containers and global services (shortcut for 'down --global')",
	Long: `Stop and remove project containers and services, then stop global services.

This command is a convenient shortcut for 'phpier down --global' and will:
- Stop and remove project-specific containers defined in docker-compose.yml
- Clean up project networks and volumes (non-persistent by default)
- Stop global services after stopping the project
- Preserve persistent data volumes by default

Examples:
  phpier stop    # Stop project and global services`,
	RunE: runStopAsDownGlobal,
}

func init() {
	rootCmd.AddCommand(stopCmd)
	// No flags - keep it simple
}

// runStopAsDownGlobal executes the stop command as a shortcut for down --global
func runStopAsDownGlobal(cmd *cobra.Command, args []string) error {
	// Set the globalFlag to true to ensure global services are stopped
	originalGlobalFlag := globalFlag
	globalFlag = true

	// Execute the down command logic
	err := runDown(cmd, args)

	// Restore original flag value
	globalFlag = originalGlobalFlag

	return err
}
