package cmd

import (
	"phpier/internal/config"
	"phpier/internal/docker"
	"phpier/internal/errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// globalDownCmd represents the 'global down' command
var globalDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the global shared services stack",
	Long: `Stops the persistent, global stack of services. This will stop Traefik,
all shared databases, and other tools. Projects will not be accessible
until the global stack is started again.`,
	RunE: runGlobalDown,
}

func init() {
	globalCmd.AddCommand(globalDownCmd)
}

func runGlobalDown(cmd *cobra.Command, args []string) error {
	logrus.Infof("🛑 Stopping global services...")

	// Load global config
	globalCfg, err := config.LoadGlobalConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load global config", err)
	}

	// Create Docker Compose manager for the global stack
	composeManager, err := docker.NewGlobalComposeManager(globalCfg)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client for global stack", err)
	}

	// Stop services
	if err := composeManager.Down(false); err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to stop global services", err)
	}

	logrus.Infof("✅ Global services stopped successfully!")
	return nil
}
