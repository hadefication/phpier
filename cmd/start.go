package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start services (global only if not in project, global + project if in project)",
	Long: `Start services based on current directory context.

When run outside a phpier project directory:
- Start global services only (Traefik, databases, etc.)

When run inside a phpier project directory:
- Start global services and project container in detached mode
- Connect the project to the global services network

Examples:
  phpier start    # Context-aware service startup`,
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
	// No flags - keep it simple
}

// runStart implements context-aware startup logic
func runStart(cmd *cobra.Command, args []string) error {
	if isPhpierProject() {
		// In a phpier project directory - start global services and project
		logrus.Infof("🔍 Detected phpier project - starting global services and project...")

		// Set detached to true for start command
		detached = true
		return runUp(cmd, args)
	} else {
		// Not in a phpier project directory - start global services only
		logrus.Infof("🌐 No phpier project detected - starting global services only...")
		return runGlobalUp(cmd, args)
	}
}
