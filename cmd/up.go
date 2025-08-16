package cmd

import (
	"os"
	"phpier/internal/config"
	"phpier/internal/docker"
	"phpier/internal/errors"
	"phpier/internal/generator"
	"phpier/internal/templates"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	detached   bool
	build      bool
	skipGlobal bool
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the project's app container",
	Long: `Start the project's app container and connect it to the global services.

This command will:
- Start the PHP/Nginx container for the current project.
- Ensure the global services network is available.

Examples:
  phpier up                 # Start app container in the foreground
  phpier up -d              # Start app container in the background (detached)
  phpier up --build         # Rebuild the app container image before starting
  phpier up --skip-global   # Start only project services, skip global service check`,
	RunE: runUp,
}

func init() {
	rootCmd.AddCommand(upCmd)

	// Flags
	upCmd.Flags().BoolVarP(&detached, "detach", "d", false, "Run services in the background")
	upCmd.Flags().BoolVar(&build, "build", false, "Build images before starting services")
	upCmd.Flags().BoolVar(&skipGlobal, "skip-global", false, "Skip automatic global service startup check")
}

func runUp(cmd *cobra.Command, args []string) error {
	if !isProjectInitialized() {
		return errors.NewProjectNotInitializedError()
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
	if !skipGlobal {
		if err := ensureGlobalServicesRunning(globalCfg); err != nil {
			return errors.WrapError(errors.ErrorTypeDockerError, "Failed to ensure global services are running", err)
		}
	} else {
		logrus.Infof("⏭️  Skipping global service startup check (--skip-global flag used)")
	}

	// Regenerate files
	engine := templates.NewEngine()
	if err := generator.GenerateProjectFiles(engine, projectCfg, globalCfg); err != nil {
		return errors.WrapError(errors.ErrorTypeUnknown, "Failed to generate project files", err)
	}

	// Create Docker Compose manager
	composeManager, err := docker.NewProjectComposeManager(projectCfg, globalCfg)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client", err)
	}

	// Build image if requested
	if build {
		logrus.Infof("🔨 Building project image...")
		if err := composeManager.Build(false); err != nil {
			return errors.WrapError(errors.ErrorTypeDockerError, "Failed to build project image", err)
		}
	}

	// Start services
	logrus.Infof("🚀 Starting project container...")
	if err := composeManager.Up(detached); err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to start project container", err)
	}

	logrus.Infof("✅ Project services started successfully!")
	if detached {
		logrus.Infof("📝 Services are running in the background. Use 'phpier down' to stop them.")
	}

	return nil
}

// ensureGlobalServicesRunning checks if global services are running and starts them if needed
func ensureGlobalServicesRunning(globalCfg *config.GlobalConfig) error {
	// Create global compose manager
	globalManager, err := docker.NewGlobalComposeManager(globalCfg)
	if err != nil {
		return err
	}

	// Check if global services are running
	isRunning, err := globalManager.IsGlobalServiceRunning()
	if err != nil {
		return err
	}

	if !isRunning {
		logrus.Infof("🔍 Global services not detected, starting them...")

		// Create directories and generate files if needed
		if err := generator.CreateGlobalDirectories(); err != nil {
			return errors.WrapError(errors.ErrorTypeFileSystemError, "Failed to create global directories", err)
		}

		engine := templates.NewEngine()
		if err := generator.GenerateGlobalFiles(engine, globalCfg); err != nil {
			return errors.WrapError(errors.ErrorTypeTemplateError, "Failed to generate global files", err)
		}

		// Start global services
		if err := globalManager.StartGlobalServicesIfNeeded(); err != nil {
			return err
		}
	} else {
		logrus.Debugf("Global services are already running")
	}

	return nil
}

// isProjectInitialized checks if the current directory has been initialized as a phpier project
func isProjectInitialized() bool {
	if _, err := os.Stat(".phpier.yml"); os.IsNotExist(err) {
		return false
	}
	return true
}
