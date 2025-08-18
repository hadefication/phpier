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

// ExecuteGlobalServiceCommand executes a command in a global service container
func (c *Client) ExecuteGlobalServiceCommand(ctx context.Context, serviceName string, proxyCmd *ProxyCommand) (int, error) {
	// Global service containers are named phpier-<serviceName>
	containerName := fmt.Sprintf("phpier-%s", serviceName)

	logrus.Debugf("Looking for global service container: %s", containerName)

	// Check if container is running
	isRunning, err := c.IsContainerRunning(ctx, containerName)
	if err != nil {
		return 1, fmt.Errorf("failed to check %s container status: %w", serviceName, err)
	}

	if !isRunning {
		return 1, fmt.Errorf("%s service is not running\n\nTry running 'phpier global up' to start the global services", strings.Title(serviceName))
	}

	// Prepare command
	command := []string{proxyCmd.Command}
	command = append(command, proxyCmd.Args...)

	// Set default values for global services
	user := proxyCmd.User
	if user == "" {
		user = "" // Use container default user for global services
	}

	workingDir := proxyCmd.WorkingDir
	if workingDir == "" {
		workingDir = "" // Use container default working directory
	}

	// Determine if command needs interactivity
	needsInteractive := proxyCmd.Interactive || isInteractiveCommand(proxyCmd.Command, proxyCmd.Args)

	// Set up execution config based on interactivity
	execConfig := &ExecConfig{
		Container:    containerName,
		Command:      command,
		WorkingDir:   workingDir,
		User:         user,
		Tty:          needsInteractive,
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  needsInteractive,
	}

	logrus.Debugf("Executing '%s %s' in global container %s", proxyCmd.Command, strings.Join(proxyCmd.Args, " "), containerName)

	// Execute the command
	exitCode, err := c.ExecInteractive(ctx, execConfig)
	if err != nil {
		return 1, fmt.Errorf("failed to execute %s command: %w", proxyCmd.Name, err)
	}

	return exitCode, nil
}

// ExecuteGlobalProxyCommand executes a command in a specific app container by project name
func (c *Client) ExecuteGlobalProxyCommand(ctx context.Context, projectName string, proxyCmd *ProxyCommand) (int, error) {
	logrus.Debugf("Looking for app container for project: %s", projectName)

	// Get container ID for the app service of the specified project
	containerID, err := c.GetContainerID(projectName, "app")
	if err != nil {
		// Check if it's a container not found error vs other errors
		if strings.Contains(err.Error(), "Container not found") {
			// Try to provide helpful error message by checking if project exists
			projects, listErr := c.GetProjectServices(ctx, projectName)
			if listErr == nil && len(projects) == 0 {
				return 1, fmt.Errorf("project '%s' not found or has no running containers\n\nTry 'phpier services' to see available projects", projectName)
			}
			return 1, fmt.Errorf("app container is not running for project '%s'\n\nTry running 'phpier up' in the project directory to start the services", projectName)
		}
		return 1, fmt.Errorf("failed to find app container for project '%s': %w\n\nMake sure Docker is running and the project containers are started", projectName, err)
	}

	logrus.Debugf("Found container ID: %s", containerID)

	// Check if container is running
	isRunning, err := c.IsContainerRunningByID(ctx, containerID)
	if err != nil {
		return 1, fmt.Errorf("failed to check container status: %w", err)
	}

	if !isRunning {
		return 1, fmt.Errorf("app container is not running for project '%s'\n\nTry running 'phpier up' in the project directory to start the services", projectName)
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
		return 1, fmt.Errorf("failed to execute %s command in project '%s': %w", proxyCmd.Name, projectName, err)
	}

	return exitCode, nil
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
