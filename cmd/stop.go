package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop services (global only if not in project, project + global if in project)",
	Long: `Stop services based on current directory context.

When run outside a phpier project directory:
- Stop global services only

When run inside a phpier project directory:
- Stop project containers and services first
- Then stop global services
- Clean up project networks and volumes (non-persistent by default)
- Preserve persistent data volumes by default

Examples:
  phpier stop    # Context-aware service stopping`,
	RunE: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)
	// No flags - keep it simple
}

// runStop implements context-aware stopping logic
func runStop(cmd *cobra.Command, args []string) error {
	if isPhpierProject() {
		// In a phpier project directory - stop project and global services
		logrus.Infof("🔍 Detected phpier project - stopping project and global services...")

		// Set the globalFlag to true to ensure global services are stopped
		originalGlobalFlag := globalFlag
		globalFlag = true

		// Execute the down command logic
		err := runDown(cmd, args)

		// Restore original flag value
		globalFlag = originalGlobalFlag

		return err
	} else {
		// Not in a phpier project directory - stop global services only
		logrus.Infof("🌐 No phpier project detected - stopping global services only...")
		return runGlobalDown(cmd, args)
	}
}
