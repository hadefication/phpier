package cmd

import (
	"phpier/internal/config"
	"phpier/internal/docker"
	"phpier/internal/errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	stopGlobal    bool
	globalFlag    bool
	removeVolumes bool
	force         bool
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop and remove project containers and services",
	Long: `Stop and remove project containers and services. Optionally stop global services.

This command will:
- Stop and remove project-specific containers defined in docker-compose.yml
- Clean up project networks and volumes (non-persistent by default)
- Preserve persistent data volumes by default (use --remove-volumes to remove all)
- Optionally stop global services with --stop-global flag

Examples:
  phpier down                      # Stop project services only
  phpier down --global             # Stop project and global services
  phpier down --force              # Force remove containers without graceful shutdown`,
	RunE: runDown,
}

func init() {
	rootCmd.AddCommand(downCmd)

	// Flags
	downCmd.Flags().BoolVar(&globalFlag, "global", false, "Also stop global services after stopping project")
	downCmd.Flags().BoolVar(&stopGlobal, "stop-global", false, "Also stop global services after stopping project (legacy)")
	downCmd.Flags().BoolVar(&removeVolumes, "remove-volumes", false, "Remove all volumes including persistent data")
	downCmd.Flags().BoolVar(&force, "force", false, "Force remove containers without graceful shutdown")
}

func runDown(cmd *cobra.Command, args []string) error {
	if !isProjectInitialized() {
		return errors.NewProjectNotInitializedError()
	}

	// Load configurations
	projectCfg, err := config.LoadProjectConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load project config", err)
	}

	// Create Docker Compose manager for project
	composeManager, err := docker.NewProjectComposeManager(projectCfg, nil)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client", err)
	}

	// Stop project services
	logrus.Infof("🛑 Stopping project containers...")
	downOptions := docker.DownOptions{
		RemoveVolumes: removeVolumes,
		Force:         force,
	}

	if err := composeManager.DownWithOptions(downOptions); err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to stop project containers", err)
	}

	logrus.Infof("✅ Project containers stopped successfully!")

	// Handle global services if requested (support both --global and --stop-global flags)
	if globalFlag || stopGlobal {
		if err := handleGlobalServicesDown(); err != nil {
			return err
		}
	}

	return nil
}

func handleGlobalServicesDown() error {
	// Load global config
	globalCfg, err := config.LoadGlobalConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load global config", err)
	}

	// Create Docker client to check for running projects
	client, err := docker.NewClient()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client", err)
	}
	defer client.Close()

	// Check for other running phpier projects
	runningProjects, err := client.GetRunningPhpierProjects()
	if err != nil {
		logrus.Warnf("⚠️  Warning: Could not check for running projects: %v", err)
	} else if len(runningProjects) > 0 {
		logrus.Warnf("⚠️  Warning: Other phpier projects are still running: %v", runningProjects)
		logrus.Warnf("   Stopping global services will affect these projects.")

		// For now, we'll proceed anyway, but in the future we might want to ask for confirmation
		// or provide a flag to override this check
	}

	// Create global compose manager
	globalComposeManager, err := docker.NewGlobalComposeManager(globalCfg)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client for global services", err)
	}

	// Stop global services
	logrus.Infof("🛑 Stopping global services...")
	globalDownOptions := docker.DownOptions{
		RemoveVolumes: false, // Never remove global volumes by default
		Force:         force,
	}

	if err := globalComposeManager.DownWithOptions(globalDownOptions); err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to stop global services", err)
	}

	if len(runningProjects) > 0 {
		logrus.Infof("✅ Global services stopped. Note: %d other project(s) were also affected.", len(runningProjects))
	} else {
		logrus.Infof("✅ Global services stopped successfully!")
	}

	return nil
}
