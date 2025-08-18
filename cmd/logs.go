package cmd

import (
	"phpier/internal/config"
	"phpier/internal/docker"
	"phpier/internal/errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	follow bool
	tail   int
	since  string
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "View logs from project containers",
	Long: `View logs from Docker containers for the current project's services.

This command will:
- Display logs from all project containers when no service is specified
- Show logs from a specific service when service name is provided
- Support real-time log following with --follow flag
- Allow limiting output with --tail and --since flags

Available services depend on your project configuration but typically include:
- app (PHP/Nginx container)
- database (MySQL, PostgreSQL, or MariaDB)
- valkey (Redis-compatible cache)
- memcached (if enabled)

Examples:
  phpier logs                    # Show logs from all services
  phpier logs app                # Show logs from app container only
  phpier logs database           # Show logs from database container
  phpier logs -f                 # Follow/tail logs in real-time
  phpier logs --tail 100         # Show last 100 lines
  phpier logs --since "2023-01-01T00:00:00Z"  # Show logs since timestamp`,
	RunE: runLogs,
}

func init() {
	rootCmd.AddCommand(logsCmd)

	// Flags
	logsCmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output in real-time")
	logsCmd.Flags().IntVar(&tail, "tail", 0, "Number of lines to show from the end of the logs")
	logsCmd.Flags().StringVar(&since, "since", "", "Show logs since timestamp (e.g. 2013-01-02T13:23:37Z)")
}

func runLogs(cmd *cobra.Command, args []string) error {
	if !isProjectInitialized() {
		return errors.NewProjectNotInitializedError()
	}

	// Load project configuration
	projectCfg, err := config.LoadProjectConfig()
	if err != nil {
		return errors.WrapError(errors.ErrorTypeConfigNotFound, "Failed to load project config", err)
	}

	// Create Docker Compose manager
	composeManager, err := docker.NewProjectComposeManager(projectCfg, nil)
	if err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to create Docker client", err)
	}

	// Determine service to show logs for
	var service string
	if len(args) > 0 {
		service = args[0]
	}

	// Show logs
	logrus.Infof("📝 Showing logs for project '%s'...", projectCfg.Name)
	if service != "" {
		logrus.Infof("🔍 Filtering logs for service: %s", service)
	}

	if err := composeManager.Logs(service, follow, tail, since); err != nil {
		return errors.WrapError(errors.ErrorTypeDockerError, "Failed to retrieve logs", err)
	}

	return nil
}
