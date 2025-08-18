package cmd

import (
	"phpier/internal/config"
	"phpier/internal/docker"
	"phpier/internal/errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	reloadDetached   bool
	reloadBuild      bool
	reloadForce      bool
	reloadTimeout    int
	reloadSkipGlobal bool
	reloadPull       bool
	reloadNoCache    bool
)

// reloadCmd represents the reload command
var reloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "Restart project services with optional rebuild",
	Long: `Restart the project's app container and services. This command gracefully stops
and then starts the services, optionally rebuilding images or pulling updates.

This command will:
- Stop current project services gracefully (or forcefully with --force)
- Optionally rebuild container images (with --build)
- Start services back up (detached with -d)
- Ensure global services are running (unless --skip-global)

Examples:
  phpier reload                          # Basic restart
  phpier reload -d                       # Restart in background
  phpier reload --build                  # Rebuild images then restart
  phpier reload --force --timeout=10     # Force stop with custom timeout
  phpier reload --build --pull --no-cache # Complete refresh`,
	RunE: runReload,
}

func init() {
	rootCmd.AddCommand(reloadCmd)

	// Flags
	reloadCmd.Flags().BoolVarP(&reloadDetached, "detach", "d", false, "Run services in the background after restart")
	reloadCmd.Flags().BoolVar(&reloadBuild, "build", false, "Rebuild container images before restarting")
	reloadCmd.Flags().BoolVar(&reloadForce, "force", false, "Force stop containers that don't respond to graceful shutdown")
	reloadCmd.Flags().IntVar(&reloadTimeout, "timeout", 30, "Timeout in seconds for stopping containers")
	reloadCmd.Flags().BoolVar(&reloadSkipGlobal, "skip-global", false, "Skip checking/starting global services during reload")
	reloadCmd.Flags().BoolVar(&reloadPull, "pull", false, "Pull latest base images before rebuilding (requires --build)")
	reloadCmd.Flags().BoolVar(&reloadNoCache, "no-cache", false, "Don't use cache when rebuilding (requires --build)")
}

func runReload(cmd *cobra.Command, args []string) error {
	if !isProjectInitialized() {
		return errors.NewProjectNotInitializedError()
	}

	// Set WWWUSER to current user ID if not already set
	if err := docker.SetWWWUser(); err != nil {
		logrus.Warnf("Failed to set WWWUSER: %v", err)
	}

	// Validate flag combinations
	if reloadPull && !reloadBuild {
		return errors.NewInvalidArgumentsError("--pull flag requires --build flag")
	}
	if reloadNoCache && !reloadBuild {
		return errors.NewInvalidArgumentsError("--no-cache flag requires --build flag")
	}
	if reloadTimeout <= 0 {
		return errors.NewInvalidArgumentsError("--timeout must be greater than 0")
	}

	// Load configurations
	projectCfg, err := config.LoadProjectConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load project config", err)
	}
	globalCfg, err := config.LoadGlobalConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load global config", err)
	}

	// Check and start global services if needed (unless --skip-global flag is used)
	if !reloadSkipGlobal {
		if err := ensureGlobalServicesRunning(globalCfg); err != nil {
			return errors.WrapError(errors.ErrorTypeDockerError, "Failed to ensure global services are running", err)
		}
	} else {
		logrus.Infof("⏭️  Skipping global service startup check (--skip-global flag used)")
	}

	// Create Docker Compose manager
	composeManager, err := docker.NewProjectComposeManager(projectCfg, globalCfg)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client", err)
	}

	// Perform reload operation
	reloadOptions := docker.ReloadOptions{
		Detached: reloadDetached,
		Build:    reloadBuild,
		Force:    reloadForce,
		Timeout:  reloadTimeout,
		Pull:     reloadPull,
		NoCache:  reloadNoCache,
	}

	logrus.Infof("🔄 Reloading project services...")
	if err := composeManager.Reload(reloadOptions); err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to reload project services", err)
	}

	logrus.Infof("✅ Project services reloaded successfully!")
	if reloadDetached {
		logrus.Infof("📝 Services are running in the background. Use 'phpier down' to stop them.")
	}

	return nil
}
