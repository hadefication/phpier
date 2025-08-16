package cmd

import (
	"phpier/internal/config"
	"phpier/internal/docker"
	"phpier/internal/errors"
	"phpier/internal/generator"
	"phpier/internal/templates"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	startDetached bool
	startBuild    bool
	startForce    bool
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the global phpier services",
	Long: `Start the global phpier services including Traefik, databases, and other tools.

This command provides a convenient way to start the global infrastructure without 
using the full 'phpier global up' command. The services will run in the background 
and be available to all phpier projects.

This command will:
- Start Traefik reverse proxy for domain routing
- Start shared databases and caching services  
- Start monitoring and development tools
- Ensure proper service dependency order

Examples:
  phpier start                 # Start global services in background
  phpier start --build         # Rebuild and start global services
  phpier start --force         # Force restart services if already running`,
	RunE: runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Flags
	startCmd.Flags().BoolVarP(&startDetached, "detach", "d", true, "Run services in the background (default: true)")
	startCmd.Flags().BoolVar(&startBuild, "build", false, "Build images before starting services")
	startCmd.Flags().BoolVar(&startForce, "force", false, "Force restart services if already running")
}

func runStart(cmd *cobra.Command, args []string) error {
	logrus.Infof("🚀 Starting global services...")

	// Load global config
	globalCfg, err := config.LoadGlobalConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load global config", err)
	}

	// If force flag is set, try to stop services first
	if startForce {
		logrus.Infof("🔄 Force flag detected, stopping existing services first...")
		if err := stopGlobalServicesIfRunning(globalCfg); err != nil {
			logrus.Warnf("⚠️  Warning: Could not stop existing services: %v", err)
		}
	}

	// Create directories
	if err := generator.CreateGlobalDirectories(); err != nil {
		return errors.WrapError(errors.ErrorTypeFileSystemError, "Failed to create global directories", err)
	}

	// Create template engine
	engine := templates.NewEngine()

	// Generate global files
	if err := generator.GenerateGlobalFiles(engine, globalCfg); err != nil {
		return errors.WrapError(errors.ErrorTypeTemplateError, "Failed to generate global files", err)
	}

	// Create Docker Compose manager for the global stack
	composeManager, err := docker.NewGlobalComposeManager(globalCfg)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client for global stack", err)
	}

	// Build images if requested
	if startBuild {
		logrus.Infof("🔨 Building global service images...")
		if err := composeManager.Build(false); err != nil {
			return errors.WrapError(errors.ErrorTypeDockerError, "Failed to build global service images", err)
		}
	}

	// Start services
	if err := composeManager.Up(startDetached); err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to start global services", err)
	}

	logrus.Infof("✅ Global services started successfully!")

	if startDetached {
		logrus.Infof("📝 Services are running in the background. Use 'phpier stop' to stop them.")
	}

	return nil
}

// stopGlobalServicesIfRunning attempts to stop global services if they are running
func stopGlobalServicesIfRunning(globalCfg *config.GlobalConfig) error {
	composeManager, err := docker.NewGlobalComposeManager(globalCfg)
	if err != nil {
		return err
	}

	// Try to stop services (ignore errors if services aren't running)
	return composeManager.Down(false)
}
