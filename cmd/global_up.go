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

// globalUpCmd represents the 'global up' command
var globalUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the global shared services stack",
	Long: `Starts the persistent, global stack of services including Traefik, databases,
and other tools. These services will run in the background and be available
to all phpier projects.`,
	RunE: runGlobalUp,
}

func init() {
	globalCmd.AddCommand(globalUpCmd)
}

func runGlobalUp(cmd *cobra.Command, args []string) error {
	logrus.Infof("🚀 Starting global services...")

	// Load global config
	globalCfg, err := config.LoadGlobalConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load global config", err)
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

	// Start services
	if err := composeManager.Up(true); err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to start global services", err)
	}

	logrus.Infof("✅ Global services started successfully!")
	return nil
}
