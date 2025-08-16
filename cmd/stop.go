package cmd

import (
	"phpier/internal/config"
	"phpier/internal/docker"
	"phpier/internal/errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	stopForce         bool
	stopRemoveVolumes bool
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the global phpier services",
	Long: `Stop the global phpier services including Traefik, databases, and other tools.

This command provides a convenient way to stop the global infrastructure without 
using the full 'phpier global down' command. This will stop all global services
and make projects inaccessible until global services are started again.

This command will:
- Stop Traefik reverse proxy
- Stop shared databases and caching services
- Stop monitoring and development tools
- Check for running projects and warn if they will be affected

Examples:
  phpier stop                      # Stop global services (with warnings)
  phpier stop --force             # Force stop even if projects are running
  phpier stop --remove-volumes    # Stop and remove global volumes (dangerous)`,
	RunE: runStop,
}

func init() {
	rootCmd.AddCommand(stopCmd)

	// Flags
	stopCmd.Flags().BoolVar(&stopForce, "force", false, "Force stop services even if projects are running")
	stopCmd.Flags().BoolVar(&stopRemoveVolumes, "remove-volumes", false, "Also remove global volumes (dangerous operation)")
}

func runStop(cmd *cobra.Command, args []string) error {
	// Load global config
	globalCfg, err := config.LoadGlobalConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load global config", err)
	}

	// Check for running projects unless force flag is used
	if !stopForce {
		if shouldAbort, err := checkRunningProjectsAndWarn(); err != nil {
			logrus.Warnf("⚠️  Warning: Could not check for running projects: %v", err)
		} else if shouldAbort {
			return errors.NewUserAbortedError("Operation aborted to protect running projects")
		}
	}

	logrus.Infof("🛑 Stopping global services...")

	// Create Docker Compose manager for the global stack
	composeManager, err := docker.NewGlobalComposeManager(globalCfg)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client for global stack", err)
	}

	// Stop services with appropriate options
	if stopRemoveVolumes {
		logrus.Warnf("⚠️  Removing global volumes - this will delete all global database data!")
		if err := composeManager.Down(true); err != nil {
			return errors.WrapError(errors.ErrorTypeDockerError, "Failed to stop global services and remove volumes", err)
		}
	} else {
		if err := composeManager.Down(false); err != nil {
			return errors.WrapError(errors.ErrorTypeDockerError, "Failed to stop global services", err)
		}
	}

	logrus.Infof("✅ Global services stopped successfully!")

	if stopRemoveVolumes {
		logrus.Infof("📝 Global volumes have been removed. Database data has been deleted.")
	} else {
		logrus.Infof("📝 Use 'phpier start' to start global services again.")
	}

	return nil
}

// checkRunningProjectsAndWarn checks for running projects and warns the user
// Returns true if the operation should be aborted, false if it should continue
func checkRunningProjectsAndWarn() (bool, error) {
	// Create Docker client to check for running projects
	client, err := docker.NewClient()
	if err != nil {
		return false, err
	}
	defer client.Close()

	// Check for running phpier projects
	runningProjects, err := client.GetRunningPhpierProjects()
	if err != nil {
		return false, err
	}

	if len(runningProjects) > 0 {
		logrus.Warnf("⚠️  Warning: %d phpier project(s) are currently running:", len(runningProjects))
		for _, project := range runningProjects {
			logrus.Warnf("   - %s", project)
		}
		logrus.Warnf("   Stopping global services will make these projects inaccessible.")
		logrus.Warnf("   Use 'phpier stop --force' to stop anyway, or stop individual projects first.")

		// For now, we'll abort the operation to protect running projects
		// In the future, we might want to prompt for user confirmation
		return true, nil
	}

	return false, nil
}
