package cmd

import (
	"fmt"

	"phpier/internal/errors"
	"phpier/internal/updater"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// selfUpdateCmd represents the self-update command
var selfUpdateCmd = &cobra.Command{
	Use:   "self-update",
	Short: "Update phpier to the latest version",
	Long: `Update phpier to the latest version or a specific version.

This command downloads and installs the latest stable release of phpier,
or a specific version if specified with the --version flag.

Examples:
  # Update to the latest version
  phpier self-update
  
  # Update to a specific version
  phpier self-update --version v1.2.0
  
  # Check for updates without installing
  phpier self-update --check
  
  # Force update without confirmation
  phpier self-update --force`,
	RunE: runSelfUpdate,
}

var (
	updateVersion string
	checkOnly     bool
	forceUpdate   bool
)

func init() {
	rootCmd.AddCommand(selfUpdateCmd)

	// Add flags
	selfUpdateCmd.Flags().StringVar(&updateVersion, "version", "", "specific version to update to (e.g., v1.2.0)")
	selfUpdateCmd.Flags().BoolVarP(&checkOnly, "check", "c", false, "check for updates without installing")
	selfUpdateCmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "force update without confirmation prompts")

	// Bind flags to viper for consistency
	viper.BindPFlag("update.version", selfUpdateCmd.Flags().Lookup("version"))
	viper.BindPFlag("update.check", selfUpdateCmd.Flags().Lookup("check"))
	viper.BindPFlag("update.force", selfUpdateCmd.Flags().Lookup("force"))
}

func runSelfUpdate(cmd *cobra.Command, args []string) error {
	// Get verbose setting from global flag
	verbose := viper.GetBool("verbose")

	// Create updater instance
	u := updater.NewUpdater(buildVersion, verbose)

	// Prepare update options
	options := updater.UpdateOptions{
		Version:   updateVersion,
		CheckOnly: checkOnly,
		Force:     forceUpdate,
		Verbose:   verbose,
	}

	// Validate version format if specified
	if updateVersion != "" && !updater.IsValidVersionFormat(updateVersion) {
		return errors.NewUserError(
			fmt.Sprintf("Invalid version format: %s", updateVersion),
			"Version must be in format 'v1.2.3' or '1.2.3'",
		)
	}

	// Run update
	if err := u.Update(options); err != nil {
		return errors.NewUpdateError(
			"Failed to update phpier",
			err.Error(),
		)
	}

	return nil
}
