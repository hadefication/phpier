package cmd

import (
	"phpier/internal/errors"
	"phpier/internal/updater"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// selfUpdateCmd represents the self-update command
var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update phpier to the latest version",
	Long: `Update phpier to the latest version.

This command downloads and installs the latest stable release of phpier
using the official installation script. If you're already running the
latest version, no changes will be made.

Example:
  phpier self-update`,
	RunE: runSelfUpdate,
}

func init() {
	rootCmd.AddCommand(selfUpdateCmd)
}

func runSelfUpdate(cmd *cobra.Command, args []string) error {
	// Get verbose setting from global flag
	verbose := viper.GetBool("verbose")

	// Create updater instance
	u := updater.NewUpdater(buildVersion, verbose)

	// Simple update - no options needed
	if err := u.Update(); err != nil {
		return errors.NewUpdateError(
			"Failed to update phpier",
			err.Error(),
		)
	}

	return nil
}
