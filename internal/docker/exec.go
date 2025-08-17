package docker

import (
	"context"
	"fmt"
	"strings"

	"phpier/internal/config"
	"phpier/internal/errors"

	"github.com/sirupsen/logrus"
)

// ProxyCommand represents a command that can be proxied to the app container
type ProxyCommand struct {
	Name        string
	Command     string
	Description string
	Args        []string
	User        string
	WorkingDir  string
	Interactive bool // Whether the command needs interactive TTY
}

// ExecuteProxyCommand executes a command in the app container
func (c *Client) ExecuteProxyCommand(ctx context.Context, proxyCmd *ProxyCommand) (int, error) {
	// Load project configuration
	projectConfig, err := config.LoadProjectConfig()
	if err != nil {
		return 1, errors.NewProjectNotInitializedError()
	}

	logrus.Debugf("Looking for app container for project: %s", projectConfig.Name)

	// Get container ID for the app service
	containerID, err := c.GetContainerID(projectConfig.Name, "app")
	if err != nil {
		// Check if it's a container not found error vs other errors
		if strings.Contains(err.Error(), "Container not found") {
			return 1, fmt.Errorf("app container is not running for project '%s'\n\nTry running 'phpier start' to start the services", projectConfig.Name)
		}
		return 1, fmt.Errorf("failed to find app container for project '%s': %w\n\nMake sure Docker is running and try 'phpier start' to start the services", projectConfig.Name, err)
	}

	logrus.Debugf("Found container ID: %s", containerID)

	// Check if container is running
	isRunning, err := c.IsContainerRunningByID(ctx, containerID)
	if err != nil {
		return 1, fmt.Errorf("failed to check container status: %w", err)
	}

	if !isRunning {
		return 1, fmt.Errorf("app container is not running\n\nTry running 'phpier start' to start the services")
	}

	// Prepare command
	command := []string{proxyCmd.Command}
	command = append(command, proxyCmd.Args...)

	// Set default values
	user := proxyCmd.User
	if user == "" {
		user = "www-data"
	}

	workingDir := proxyCmd.WorkingDir
	if workingDir == "" {
		workingDir = "/var/www/html"
	}

	// Determine if command needs interactivity
	needsInteractive := proxyCmd.Interactive || isInteractiveCommand(proxyCmd.Command, proxyCmd.Args)
	
	// Set up execution config based on interactivity
	execConfig := &ExecConfig{
		Container:    containerID,
		Command:      command,
		WorkingDir:   workingDir,
		User:         user,
		Tty:          needsInteractive,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  needsInteractive,
	}

	logrus.Debugf("Executing '%s %s' in container %s as user %s", proxyCmd.Command, strings.Join(proxyCmd.Args, " "), containerID, execConfig.User)

	// Execute the command
	exitCode, err := c.ExecInteractive(ctx, execConfig)
	if err != nil {
		return 1, fmt.Errorf("failed to execute %s command: %w", proxyCmd.Name, err)
	}

	return exitCode, nil
}

// CheckToolAvailability checks if a tool is available in the app container
func (c *Client) CheckToolAvailability(ctx context.Context, tool string) error {
	// Load project configuration
	projectConfig, err := config.LoadProjectConfig()
	if err != nil {
		return errors.NewProjectNotInitializedError()
	}

	// Get container ID for the app service
	containerID, err := c.GetContainerID(projectConfig.Name, "app")
	if err != nil {
		return err
	}

	// Check if container is running
	isRunning, err := c.IsContainerRunningByID(ctx, containerID)
	if err != nil {
		return err
	}

	if !isRunning {
		return fmt.Errorf("app container is not running")
	}

	// Check if tool exists
	checkCmd := []string{"which", tool}
	_, err = c.ExecInContainerOutput(containerID, checkCmd)
	if err != nil {
		return fmt.Errorf("%s is not available in the container", tool)
	}

	return nil
}

// GetToolVersion gets the version of a tool in the app container
func (c *Client) GetToolVersion(ctx context.Context, tool string, versionFlag string) (string, error) {
	// Load project configuration
	projectConfig, err := config.LoadProjectConfig()
	if err != nil {
		return "", errors.NewProjectNotInitializedError()
	}

	// Get container ID for the app service
	containerID, err := c.GetContainerID(projectConfig.Name, "app")
	if err != nil {
		return "", err
	}

	// Check if container is running
	isRunning, err := c.IsContainerRunningByID(ctx, containerID)
	if err != nil {
		return "", err
	}

	if !isRunning {
		return "", fmt.Errorf("app container is not running")
	}

	// Get tool version
	versionCmd := []string{tool, versionFlag}
	output, err := c.ExecInContainerOutput(containerID, versionCmd)
	if err != nil {
		return "", fmt.Errorf("failed to get %s version: %w", tool, err)
	}

	return strings.TrimSpace(output), nil
}

// isInteractiveCommand determines if a command needs interactive TTY based on command and args
func isInteractiveCommand(command string, args []string) bool {
	// Commands that typically need interactive mode
	interactiveCommands := map[string][]string{
		"composer": {"require", "create-project", "install"}, // Some composer commands can be interactive
		"npm":      {"init", "install"},                      // Some npm commands can be interactive
		"php":      {"artisan", "tinker"},                    // PHP artisan tinker is interactive
		"node":     {},                                       // Node REPL is interactive when no args
	}
	
	// Check if command is potentially interactive
	if interactiveArgs, exists := interactiveCommands[command]; exists {
		// If no args and command can be interactive (like node REPL)
		if len(args) == 0 && command == "node" {
			return true
		}
		
		// Check for specific interactive arguments
		for _, arg := range args {
			for _, interactiveArg := range interactiveArgs {
				if arg == interactiveArg {
					return false // These can be interactive but not always, default to non-interactive
				}
			}
		}
	}
	
	// Default to non-interactive for most commands
	return false
}
