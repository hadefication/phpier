package cmd

import (
	"context"
	"fmt"
	"strings"

	"phpier/internal/config"
	"phpier/internal/docker"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// mariaCmd represents the maria command
var mariaCmd = &cobra.Command{
	Use:   "maria",
	Short: "Connect to the MariaDB database shell",
	Long: `Connect directly to the MariaDB database container shell.

This command provides direct access to the MariaDB database, allowing you to:
- Run MariaDB commands and queries
- Inspect database structure and data
- Perform database operations
- Debug database issues

The command automatically connects using the configured database credentials.

Examples:
  phpier maria                      # Open interactive MariaDB shell
  phpier maria -e "SHOW TABLES"    # Execute single query`,
	RunE: runMaria,
}

// mariadbCmd represents the mariadb alias command
var mariadbCmd = &cobra.Command{
	Use:   "mariadb",
	Short: "Connect to the MariaDB database shell (alias for maria)",
	Long: `Connect directly to the MariaDB database container shell.

This is an alias for the 'maria' command.

Examples:
  phpier mariadb                    # Open interactive MariaDB shell
  phpier mariadb -e "SHOW TABLES"  # Execute single query`,
	RunE: runMaria,
}

func runMaria(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load global configuration to get database settings
	globalConfig, err := config.LoadGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	// Check if MariaDB is the configured database type
	if globalConfig.Services.Database.Type != "mariadb" {
		return fmt.Errorf("MariaDB is not configured as the database type (current: %s)\n\nUpdate your global configuration to use MariaDB", globalConfig.Services.Database.Type)
	}

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	logrus.Debugf("Looking for MariaDB container")

	// Get container ID for the MariaDB service
	containerID, err := dockerClient.GetContainerID("phpier", "mariadb")
	if err != nil {
		// Check if it's a container not found error vs other errors
		if strings.Contains(err.Error(), "Container not found") {
			return fmt.Errorf("MariaDB container is not running\n\nTry running 'phpier global up' to start the global services")
		}
		return fmt.Errorf("failed to find MariaDB container: %w\n\nMake sure Docker is running and try 'phpier global up' to start the services", err)
	}

	logrus.Debugf("Found MariaDB container ID: %s", containerID)

	// Check if container is running
	isRunning, err := dockerClient.IsContainerRunningByID(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to check MariaDB container status: %w", err)
	}

	if !isRunning {
		return fmt.Errorf("MariaDB container is not running\n\nTry running 'phpier global up' to start the global services")
	}

	// Prepare MariaDB command
	var mariaCommand []string
	if len(args) > 0 {
		// Execute SQL query from arguments
		query := strings.Join(args, " ")
		mariaCommand = []string{"mariadb", "-u", "phpier", "-pphpier", "phpier", "-e", query}
	} else {
		// Interactive MariaDB shell
		mariaCommand = []string{"mariadb", "-u", "phpier", "-pphpier", "phpier"}
	}

	// Set up execution config
	execConfig := &docker.ExecConfig{
		Container:    containerID,
		Command:      mariaCommand,
		WorkingDir:   "/",
		User:         "root",
		Tty:          len(args) == 0, // TTY only for interactive mode
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  len(args) == 0, // Stdin only for interactive mode
	}

	logrus.Debugf("Executing MariaDB command in container %s", containerID)

	// Execute the MariaDB command
	exitCode, err := dockerClient.ExecInteractive(ctx, execConfig)
	if err != nil {
		// Try fallback to mysql client if mariadb command failed
		if strings.Contains(err.Error(), "mariadb") {
			logrus.Debug("MariaDB client not available, falling back to mysql client")
			if len(args) > 0 {
				query := strings.Join(args, " ")
				execConfig.Command = []string{"mysql", "-u", "phpier", "-pphpier", "phpier", "-e", query}
			} else {
				execConfig.Command = []string{"mysql", "-u", "phpier", "-pphpier", "phpier"}
			}
			exitCode, err = dockerClient.ExecInteractive(ctx, execConfig)
			if err != nil {
				return fmt.Errorf("failed to execute MariaDB command: %w", err)
			}
		} else {
			return fmt.Errorf("failed to execute MariaDB command: %w", err)
		}
	}

	// Exit with the same code as the container command
	if exitCode != 0 && len(args) > 0 {
		return fmt.Errorf("MariaDB command failed with exit code %d", exitCode)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(mariaCmd)
	rootCmd.AddCommand(mariadbCmd)
}
