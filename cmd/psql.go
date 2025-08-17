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

// psqlCmd represents the psql command
var psqlCmd = &cobra.Command{
	Use:   "psql",
	Short: "Connect to the PostgreSQL database shell",
	Long: `Connect directly to the PostgreSQL database container shell.

This command provides direct access to the PostgreSQL database, allowing you to:
- Run PostgreSQL commands and queries
- Inspect database structure and data
- Perform database operations
- Debug database issues

The command automatically connects using the configured database credentials.

Examples:
  phpier psql                           # Open interactive psql shell
  phpier psql -c "SELECT version();"    # Execute single query`,
	RunE: runPSQL,
}

// postgresCmd represents the postgres alias command
var postgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Connect to the PostgreSQL database shell (alias for psql)",
	Long: `Connect directly to the PostgreSQL database container shell.

This is an alias for the 'psql' command.

Examples:
  phpier postgres                       # Open interactive psql shell
  phpier postgres -c "SELECT version();" # Execute single query`,
	RunE: runPSQL,
}

func runPSQL(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load global configuration to get database settings
	globalConfig, err := config.LoadGlobalConfig()
	if err != nil {
		return fmt.Errorf("failed to load global config: %w", err)
	}

	// Check if PostgreSQL is enabled
	if !globalConfig.IsDatabaseEnabled("postgresql") {
		return fmt.Errorf("PostgreSQL is not enabled\n\nRun 'phpier global db enable postgresql' to enable PostgreSQL")
	}

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	logrus.Debugf("Looking for PostgreSQL container")

	// Get container ID for the PostgreSQL service
	containerID, err := dockerClient.GetContainerID("phpier", "postgres")
	if err != nil {
		// Check if it's a container not found error vs other errors
		if strings.Contains(err.Error(), "Container not found") {
			return fmt.Errorf("PostgreSQL container is not running\n\nTry running 'phpier global up' to start the global services")
		}
		return fmt.Errorf("failed to find PostgreSQL container: %w\n\nMake sure Docker is running and try 'phpier global up' to start the services", err)
	}

	logrus.Debugf("Found PostgreSQL container ID: %s", containerID)

	// Check if container is running
	isRunning, err := dockerClient.IsContainerRunningByID(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to check PostgreSQL container status: %w", err)
	}

	if !isRunning {
		return fmt.Errorf("PostgreSQL container is not running\n\nTry running 'phpier global up' to start the global services")
	}

	// Get PostgreSQL configuration
	pgConfig := globalConfig.Services.Databases.PostgreSQL

	// Prepare PostgreSQL command
	var psqlCommand []string
	if len(args) > 0 {
		// Execute SQL query from arguments
		psqlCommand = append([]string{"psql", "-U", pgConfig.Username, "-d", pgConfig.Database}, args...)
	} else {
		// Interactive psql shell
		psqlCommand = []string{"psql", "-U", pgConfig.Username, "-d", pgConfig.Database}
	}

	// Set up execution config with PGPASSWORD environment variable
	execConfig := &docker.ExecConfig{
		Container:    containerID,
		Command:      psqlCommand,
		WorkingDir:   "/",
		User:         "postgres",
		Tty:          len(args) == 0, // TTY only for interactive mode
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  len(args) == 0, // Stdin only for interactive mode
		Environment:  []string{fmt.Sprintf("PGPASSWORD=%s", pgConfig.Password)},
	}

	logrus.Debugf("Executing PostgreSQL command in container %s", containerID)

	// Execute the psql command
	exitCode, err := dockerClient.ExecInteractive(ctx, execConfig)
	if err != nil {
		return fmt.Errorf("failed to execute PostgreSQL command: %w", err)
	}

	// Exit with the same code as the container command
	if exitCode != 0 && len(args) > 0 {
		return fmt.Errorf("PostgreSQL command failed with exit code %d", exitCode)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(psqlCmd)
	rootCmd.AddCommand(postgresCmd)
}
