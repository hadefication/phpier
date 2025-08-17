package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"phpier/internal/config"
	"phpier/internal/docker"
	"phpier/internal/errors"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	shUser    string
	shCommand string
)

// shCmd represents the sh command
var shCmd = &cobra.Command{
	Use:   "sh",
	Short: "Open an interactive shell in the app container",
	Long: `Open an interactive shell session in the PHP app container.

This command provides direct access to the container environment, allowing you to:
- Explore the container filesystem 
- Run commands interactively
- Debug issues
- Perform manual operations

The shell starts in the /var/www/html directory as the www-data user by default.

Examples:
  phpier sh                          # Open interactive bash shell
  phpier sh -c "php -v"             # Execute single command
  phpier sh --user root             # Open shell as root user
  phpier sh -c "composer install"   # Run composer install`,
	RunE: runSh,
}

func runSh(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Load project configuration
	projectConfig, err := config.LoadProjectConfig()
	if err != nil {
		return errors.NewProjectNotInitializedError()
	}

	// Create Docker client
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}
	defer dockerClient.Close()

	logrus.Debugf("Looking for app container for project: %s", projectConfig.Name)

	// Get container ID for the app service
	containerID, err := dockerClient.GetContainerID(projectConfig.Name, "app")
	if err != nil {
		// Check if it's a container not found error vs other errors
		if strings.Contains(err.Error(), "Container not found") {
			return fmt.Errorf("app container is not running for project '%s'\n\nTry running 'phpier start' to start the services", projectConfig.Name)
		}
		return fmt.Errorf("failed to find app container for project '%s': %w\n\nMake sure Docker is running and try 'phpier start' to start the services", projectConfig.Name, err)
	}

	logrus.Debugf("Found container ID: %s", containerID)

	// Check if container is running
	isRunning, err := dockerClient.IsContainerRunningByID(ctx, containerID)
	if err != nil {
		return fmt.Errorf("failed to check container status: %w", err)
	}

	if !isRunning {
		return fmt.Errorf("app container is not running\n\nTry running 'phpier start' to start the services")
	}

	// Prepare shell command
	var shellCommand []string
	if shCommand != "" {
		// Execute single command
		shellCommand = []string{"/bin/bash", "-c", shCommand}
	} else {
		// Interactive shell - try bash first, fallback to sh
		shellCommand = []string{"/bin/bash"}
	}

	// Set up execution config
	execConfig := &docker.ExecConfig{
		Container:    containerID,
		Command:      shellCommand,
		WorkingDir:   "/var/www/html",
		User:         getEffectiveUser(),
		Tty:          shCommand == "", // TTY only for interactive mode
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  shCommand == "", // Stdin only for interactive mode
	}

	logrus.Debugf("Executing shell in container %s as user %s", containerID, execConfig.User)

	// Execute the shell
	exitCode, err := dockerClient.ExecInteractive(ctx, execConfig)
	if err != nil {
		// Try fallback to /bin/sh if bash failed and we're in interactive mode
		if shCommand == "" && contains(err.Error(), "bash") {
			logrus.Debug("Bash not available, falling back to /bin/sh")
			execConfig.Command = []string{"/bin/sh"}
			exitCode, err = dockerClient.ExecInteractive(ctx, execConfig)
			if err != nil {
				return fmt.Errorf("failed to start shell: %w", err)
			}
		} else {
			return fmt.Errorf("failed to execute shell command: %w", err)
		}
	}

	// Exit with the same code as the container command
	if exitCode != 0 && shCommand != "" {
		os.Exit(exitCode)
	}

	return nil
}

// getEffectiveUser returns the user that will be used for shell execution
func getEffectiveUser() string {
	if shUser != "" {
		return shUser
	}
	return "www-data"
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsAt(s, substr))))
}

// containsAt checks if substr exists anywhere in s
func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(shCmd)

	// Flags
	shCmd.Flags().StringVarP(&shCommand, "command", "c", "", "Execute a single command instead of opening interactive shell")
	shCmd.Flags().StringVar(&shUser, "user", "", "User to execute as (default: www-data)")
}
