package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
	"phpier/internal/errors"
)

// ExecConfig represents configuration for executing commands in containers
type ExecConfig struct {
	Container    string
	Command      []string
	WorkingDir   string
	User         string
	Tty          bool
	AttachStdout bool
	AttachStderr bool
	AttachStdin  bool
	Environment  []string
}

// Client represents a Docker client wrapper
type Client struct {
	ctx context.Context
}

// NewClient creates a new Docker client
func NewClient() (*Client, error) {
	// Check if Docker is available
	if err := checkDockerAvailable(); err != nil {
		return nil, err
	}

	return &Client{
		ctx: context.Background(),
	}, nil
}

// Close closes the Docker client
func (c *Client) Close() {
	// No cleanup needed for exec-based client
}

// checkDockerAvailable checks if Docker is installed and running
func checkDockerAvailable() error {
	// Check if docker command exists
	if _, err := exec.LookPath("docker"); err != nil {
		return errors.NewDockerNotFoundError()
	}

	// Check if Docker daemon is running
	cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	if err := cmd.Run(); err != nil {
		return errors.NewDockerNotRunningError()
	}

	// Check if docker-compose is available
	if _, err := exec.LookPath("docker-compose"); err != nil {
		// Try compose plugin
		if _, err2 := exec.LookPath("docker"); err2 == nil {
			cmd := exec.Command("docker", "compose", "version")
			if err3 := cmd.Run(); err3 != nil {
				return errors.NewDockerComposeNotFoundError()
			}
		} else {
			return errors.NewDockerComposeNotFoundError()
		}
	}

	return nil
}

// IsDockerRunning checks if Docker daemon is running
func (c *Client) IsDockerRunning() bool {
	cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
	return cmd.Run() == nil
}

// GetDockerComposeCommand returns the appropriate docker-compose command
func (c *Client) GetDockerComposeCommand() string {
	// Check if docker-compose is available
	if _, err := exec.LookPath("docker-compose"); err == nil {
		return "docker-compose"
	}

	// Check if docker compose plugin is available
	cmd := exec.Command("docker", "compose", "version")
	if err := cmd.Run(); err == nil {
		return "docker compose"
	}

	return "docker-compose" // fallback
}

// RunCommand executes a Docker command
func (c *Client) RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = logrus.StandardLogger().Out
	cmd.Stderr = logrus.StandardLogger().Out

	logrus.Debugf("Executing: %s %s", command, strings.Join(args, " "))

	if err := cmd.Run(); err != nil {
		return errors.NewCommandFailedError(command, args, err)
	}

	return nil
}

// RunCommandOutput executes a Docker command and returns output
func (c *Client) RunCommandOutput(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	logrus.Debugf("Executing: %s %s", command, strings.Join(args, " "))

	output, err := cmd.Output()
	if err != nil {
		return "", errors.NewCommandFailedError(command, args, err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GetContainerID gets the container ID for a service
func (c *Client) GetContainerID(projectName, serviceName string) (string, error) {
	composeCmdStr := c.GetDockerComposeCommand()
	var args []string

	if strings.Contains(composeCmdStr, " ") {
		parts := strings.Split(composeCmdStr, " ")
		composeCmdStr = parts[0]
		args = append(parts[1:], "-p", projectName, "ps", "-q", serviceName)
	} else {
		args = []string{"-p", projectName, "ps", "-q", serviceName}
	}

	containerID, err := c.RunCommandOutput(composeCmdStr, args...)
	if err != nil {
		return "", err
	}

	if containerID == "" {
		return "", errors.NewContainerNotFoundError(serviceName)
	}

	return containerID, nil
}

// ExecInContainer executes a command in a container
func (c *Client) ExecInContainer(containerID string, command []string) error {
	args := append([]string{"exec", "-it", containerID}, command...)
	return c.RunCommand("docker", args...)
}

// ExecInContainerOutput executes a command in a container and returns output
func (c *Client) ExecInContainerOutput(containerID string, command []string) (string, error) {
	args := append([]string{"exec", "-i", containerID}, command...)
	return c.RunCommandOutput("docker", args...)
}

// IsContainerRunning checks if a container is running
func (c *Client) IsContainerRunning(ctx context.Context, containerName string) (bool, error) {
	output, err := c.RunCommandOutput("docker", "ps", "--filter", fmt.Sprintf("name=%s", containerName), "--filter", "status=running", "--format", "{{.Names}}")
	if err != nil {
		return false, err
	}

	// Check if the container name appears in the output
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == containerName {
			return true, nil
		}
	}

	return false, nil
}

// IsContainerRunningByID checks if a container is running by container ID
func (c *Client) IsContainerRunningByID(ctx context.Context, containerID string) (bool, error) {
	if containerID == "" {
		return false, nil
	}

	output, err := c.RunCommandOutput("docker", "inspect", "--format", "{{.State.Status}}", containerID)
	if err != nil {
		return false, err
	}

	status := strings.TrimSpace(output)
	return status == "running", nil
}

// ExecInteractive executes a command interactively in a container
func (c *Client) ExecInteractive(ctx context.Context, config *ExecConfig) (int, error) {
	args := []string{"exec"}

	// Add interactive flags
	if config.Tty && config.AttachStdin {
		args = append(args, "-it")
	} else if config.AttachStdin {
		args = append(args, "-i")
	} else if config.Tty {
		args = append(args, "-t")
	}

	// Add user if specified
	if config.User != "" {
		args = append(args, "--user", config.User)
	}

	// Add working directory if specified
	if config.WorkingDir != "" {
		args = append(args, "-w", config.WorkingDir)
	}

	// Add environment variables if specified
	for _, env := range config.Environment {
		args = append(args, "-e", env)
	}

	// Add container name
	args = append(args, config.Container)

	// Add command
	args = append(args, config.Command...)

	// Create command
	cmd := exec.CommandContext(ctx, "docker", args...)

	// Attach stdio
	if config.AttachStdin {
		cmd.Stdin = os.Stdin
	}
	if config.AttachStdout {
		cmd.Stdout = os.Stdout
	}
	if config.AttachStderr {
		cmd.Stderr = os.Stderr
	}

	logrus.Debugf("Executing: docker %s", strings.Join(args, " "))

	// Run the command
	err := cmd.Run()
	if err != nil {
		// Check if it's an exit error to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus(), nil
			}
		}
		return 1, fmt.Errorf("failed to execute command: %w", err)
	}

	return 0, nil
}

// ProjectInfo represents information about a discovered phpier project
type ProjectInfo struct {
	Name      string
	Status    string // "running", "stopped", "created"
	ImageName string
	Path      string // Working directory if available
}

// GetRunningPhpierProjects returns a list of running phpier projects
func (c *Client) GetRunningPhpierProjects() ([]string, error) {
	output, err := c.RunCommandOutput("docker", "ps", "--filter", "label=com.docker.compose.project", "--format", "{{.Label \"com.docker.compose.project\"}}")
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(output) == "" {
		return []string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	projects := make(map[string]bool)
	var result []string

	for _, line := range lines {
		project := strings.TrimSpace(line)
		if project != "" && project != "phpier" && !projects[project] {
			projects[project] = true
			result = append(result, project)
		}
	}

	return result, nil
}

// GetAllPhpierProjects returns all phpier projects (running and stopped)
func (c *Client) GetAllPhpierProjects() ([]ProjectInfo, error) {
	// Get all containers (running and stopped) with phpier labels
	output, err := c.RunCommandOutput("docker", "ps", "-a", 
		"--filter", "label=com.docker.compose.project",
		"--filter", "label=phpier.managed=true",
		"--format", "{{.Label \"com.docker.compose.project\"}}\t{{.Status}}\t{{.Image}}\t{{.Label \"com.docker.compose.working-dir\"}}")
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(output) == "" {
		return []ProjectInfo{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	projectMap := make(map[string]*ProjectInfo)

	for _, line := range lines {
		parts := strings.Split(strings.TrimSpace(line), "\t")
		if len(parts) < 3 {
			continue
		}

		projectName := strings.TrimSpace(parts[0])
		status := strings.TrimSpace(parts[1])
		imageName := strings.TrimSpace(parts[2])
		workingDir := ""
		if len(parts) > 3 {
			workingDir = strings.TrimSpace(parts[3])
		}

		// Skip global phpier services
		if projectName == "" || projectName == "phpier" {
			continue
		}

		// Determine project status from container status
		projectStatus := "stopped"
		if strings.Contains(status, "Up") {
			projectStatus = "running"
		} else if strings.Contains(status, "Created") {
			projectStatus = "created"
		}

		// Use the first container found for each project or update with running status
		if existing, exists := projectMap[projectName]; !exists || (existing.Status != "running" && projectStatus == "running") {
			projectMap[projectName] = &ProjectInfo{
				Name:      projectName,
				Status:    projectStatus,
				ImageName: imageName,
				Path:      workingDir,
			}
		}
	}

	// Convert map to slice
	var result []ProjectInfo
	for _, project := range projectMap {
		result = append(result, *project)
	}

	return result, nil
}

// GetPhpierProjectByName finds a specific phpier project by name using Docker
func (c *Client) GetPhpierProjectByName(projectName string) (*ProjectInfo, error) {
	projects, err := c.GetAllPhpierProjects()
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == projectName {
			return &project, nil
		}
	}

	return nil, fmt.Errorf("project '%s' not found in Docker containers", projectName)
}

// GetPhpierImages returns all phpier-built images (prefixed with phpier-)
func (c *Client) GetPhpierImages() ([]string, error) {
	output, err := c.RunCommandOutput("docker", "images", 
		"--filter", "reference=phpier-*",
		"--format", "{{.Repository}}:{{.Tag}}")
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(output) == "" {
		return []string{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	var result []string

	for _, line := range lines {
		imageName := strings.TrimSpace(line)
		if imageName != "" && strings.HasPrefix(imageName, "phpier-") {
			result = append(result, imageName)
		}
	}

	return result, nil
}

// GetPhpierProjectsFromImages discovers projects by scanning phpier- prefixed Docker images
func (c *Client) GetPhpierProjectsFromImages() ([]ProjectInfo, error) {
	// Get all phpier- prefixed images
	output, err := c.RunCommandOutput("docker", "images", 
		"--filter", "reference=phpier-*",
		"--format", "{{.Repository}}:{{.Tag}}")
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(output) == "" {
		return []ProjectInfo{}, nil
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	projectMap := make(map[string]ProjectInfo)

	for _, line := range lines {
		imageName := strings.TrimSpace(line)
		if imageName == "" || !strings.HasPrefix(imageName, "phpier-") {
			continue
		}

		// Extract project name from image name (remove phpier- prefix and :tag)
		projectName := strings.TrimPrefix(imageName, "phpier-")
		if colonIndex := strings.Index(projectName, ":"); colonIndex != -1 {
			projectName = projectName[:colonIndex]
		}

		if projectName == "" {
			continue
		}

		// Check if this project has running containers
		status := "stopped"
		workingDir := ""
		
		// Try to get container info for this project
		containerOutput, err := c.RunCommandOutput("docker", "ps", "-a",
			"--filter", fmt.Sprintf("ancestor=%s", imageName),
			"--format", "{{.Status}}\t{{.Label \"com.docker.compose.working-dir\"}}")
		
		if err == nil && strings.TrimSpace(containerOutput) != "" {
			parts := strings.Split(strings.TrimSpace(containerOutput), "\t")
			if len(parts) > 0 {
				containerStatus := strings.TrimSpace(parts[0])
				if strings.Contains(containerStatus, "Up") {
					status = "running"
				} else if strings.Contains(containerStatus, "Created") {
					status = "created"
				}
				
				if len(parts) > 1 {
					workingDir = strings.TrimSpace(parts[1])
				}
			}
		}

		projectMap[projectName] = ProjectInfo{
			Name:      projectName,
			Status:    status,
			ImageName: imageName,
			Path:      workingDir,
		}
	}

	// Convert map to slice
	var result []ProjectInfo
	for _, project := range projectMap {
		result = append(result, project)
	}

	return result, nil
}

// GetProjectWorkingDirectory attempts to find the working directory for a project
func (c *Client) GetProjectWorkingDirectory(projectName string) (string, error) {
	// Try to get working directory from running container
	output, err := c.RunCommandOutput("docker", "ps", 
		"--filter", fmt.Sprintf("label=com.docker.compose.project=%s", projectName),
		"--format", "{{.Label \"com.docker.compose.working-dir\"}}")
	if err != nil {
		return "", err
	}

	if workDir := strings.TrimSpace(output); workDir != "" {
		return workDir, nil
	}

	// If not running, try to get from stopped containers
	output, err = c.RunCommandOutput("docker", "ps", "-a",
		"--filter", fmt.Sprintf("label=com.docker.compose.project=%s", projectName),
		"--format", "{{.Label \"com.docker.compose.working-dir\"}}")
	if err != nil {
		return "", err
	}

	workDir := strings.TrimSpace(output)
	if workDir == "" {
		return "", fmt.Errorf("working directory not found for project '%s'", projectName)
	}

	return workDir, nil
}
