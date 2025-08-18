package cmd

import (
	"os"
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
	Use:   "down [app]",
	Short: "Stop and remove project containers and services",
	Long: `Stop and remove project containers and services. Optionally stop global services.

This command will:
- Stop and remove project-specific containers defined in docker-compose.yml
- Clean up project networks and volumes (non-persistent by default)
- Preserve persistent data volumes by default (use --remove-volumes to remove all)
- Optionally stop global services with --stop-global flag
- Optionally stop a specific project by name from any directory

Examples:
  phpier down                      # Stop project services only
  phpier down --global             # Stop project and global services
  phpier down --force              # Force remove containers without graceful shutdown
  phpier down myapp                # Stop 'myapp' project from any directory
  phpier down myapp --global       # Stop 'myapp' project and global services`,
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
	var projectPath string
	var projectCfg *config.ProjectConfig
	var err error

	// Check if app argument is provided
	if len(args) > 0 {
		projectName := args[0]
		logrus.Infof("🔍 Looking for project '%s'...", projectName)

		// Find project by name
		projectInfo, err := config.FindProjectByName(projectName)
		if err != nil {
			return err
		}

		projectPath = projectInfo.Path
		logrus.Infof("📂 Found project at: %s", projectPath)

		// Load project config from specific path
		projectCfg, err = config.LoadProjectConfigFromPath(projectPath)
		if err != nil {
			return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load project config", err)
		}
	} else {
		// Use current directory
		if !isProjectInitialized() {
			return errors.NewProjectNotInitializedError()
		}

		// Load configurations from current directory
		projectCfg, err = config.LoadProjectConfig()
		if err != nil {
			return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load project config", err)
		}

		// Get current working directory
		projectPath, err = os.Getwd()
		if err != nil {
			return errors.WrapError(errors.ErrorTypeFileSystemError, "Failed to get current directory", err)
		}
	}

	// Create Docker Compose manager for project
	var composeManager *docker.ProjectComposeManager
	if len(args) > 0 {
		// Use specific project path
		composeManager, err = docker.NewProjectComposeManagerWithPath(projectCfg, nil, projectPath)
		if err != nil {
			return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client", err)
		}
	} else {
		// Use current directory
		composeManager, err = docker.NewProjectComposeManager(projectCfg, nil)
		if err != nil {
			return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client", err)
		}
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
