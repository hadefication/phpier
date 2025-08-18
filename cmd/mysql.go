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

// mysqlCmd represents the mysql command
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Connect to the MySQL database shell",
	Long: `Connect directly to the MySQL database container shell.

This command provides direct access to the MySQL database, allowing you to:
- Run MySQL commands and queries
- Inspect database structure and data
- Perform database operations
- Debug database issues

The command automatically connects using the configured database credentials.

Examples:
  phpier mysql                    # Open interactive MySQL shell
  phpier mysql -e "SHOW TABLES"  # Execute single query`,
	RunE: runMySQL,
}

func runMySQL(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load global configuration to get database settings
	globalConfig, err := config.LoadGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	// Check if MySQL is enabled
	if !globalConfig.IsDatabaseEnabled("mysql") {
		return fmt.Errorf("MySQL is not enabled\n\nRun 'phpier global db enable mysql' to enable MySQL")
	}

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	logrus.Debugf("Looking for MySQL container")

	// Get container ID for the MySQL service
	containerID, err := dockerClient.GetContainerID("phpier", "mysql")
	if err != nil {
		// Check if it's a container not found error vs other errors
		if strings.Contains(err.Error(), "Container not found") {
			return fmt.Errorf("MySQL container is not running\n\nTry running 'phpier global up' to start the global services")
		}
		return fmt.Errorf("failed to find MySQL container: %w\n\nMake sure Docker is running and try 'phpier global up' to start the services", err)
	}

	logrus.Debugf("Found MySQL container ID: %s", containerID)

	// Check if container is running
	isRunning, err := dockerClient.IsContainerRunningByID(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to check MySQL container status: %w", err)
	}

	if !isRunning {
		return fmt.Errorf("MySQL container is not running\n\nTry running 'phpier global up' to start the global services")
	}

	// Get MySQL configuration
	mysqlConfig := globalConfig.Services.Databases.MySQL

	// Prepare MySQL command
	var mysqlCommand []string
	if len(args) > 0 {
		// Execute SQL query from arguments
		query := strings.Join(args, " ")
		mysqlCommand = []string{"mysql", "-u", "root", fmt.Sprintf("-p%s", mysqlConfig.Password), "-e", query}
	} else {
		// Interactive MySQL shell
		mysqlCommand = []string{"mysql", "-u", "root", fmt.Sprintf("-p%s", mysqlConfig.Password)}
	}

	// Set up execution config
	execConfig := &docker.ExecConfig{
		Container:    containerID,
		Command:      mysqlCommand,
		WorkingDir:   "/",
		User:         "root",
		Tty:          len(args) == 0, // TTY only for interactive mode
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  len(args) == 0, // Stdin only for interactive mode
	}

	logrus.Debugf("Executing MySQL command in container %s", containerID)

	// Execute the MySQL command
	exitCode, err := dockerClient.ExecInteractive(ctx, execConfig)
	if err != nil {
		return fmt.Errorf("failed to execute MySQL command: %w", err)
	}

	// Exit with the same code as the container command
	if exitCode != 0 && len(args) > 0 {
		return fmt.Errorf("MySQL command failed with exit code %d", exitCode)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(mysqlCmd)
}
