package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"phpier/internal/config"

	"github.com/sirupsen/logrus"
)

// DownOptions represents options for the down operation
type DownOptions struct {
	RemoveVolumes bool
	Force         bool
}

// ComposeManager interface for Docker Compose operations.
type ComposeManager interface {
	Up(detached bool) error
	Down(removeVolumes bool) error
	DownWithOptions(options DownOptions) error
	Build(noCache bool, services ...string) error
}

// GlobalServiceChecker interface for checking global service status.
type GlobalServiceChecker interface {
	IsGlobalServiceRunning() (bool, error)
	StartGlobalServicesIfNeeded() error
}

// GlobalComposeManager handles Docker Compose operations for the global services stack.
type GlobalComposeManager struct {
	client     *Client
	globalCfg  *config.GlobalConfig
	composeCmd string
	workDir    string
}

// ProjectComposeManager handles Docker Compose operations for a specific project.
type ProjectComposeManager struct {
	client     *Client
	projectCfg *config.ProjectConfig
	globalCfg  *config.GlobalConfig
	composeCmd string
	workDir    string
}

// NewProjectComposeManager creates a new Docker Compose manager for a project.
func NewProjectComposeManager(projectCfg *config.ProjectConfig, globalCfg *config.GlobalConfig) (*ProjectComposeManager, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	return &ProjectComposeManager{
		client:     client,
		projectCfg: projectCfg,
		globalCfg:  globalCfg,
		composeCmd: client.GetDockerComposeCommand(),
		workDir:    ".phpier",
	}, nil
}

// Up starts the Docker Compose services for a project.
func (cm *ProjectComposeManager) Up(detached bool) error {
	if !cm.client.IsDockerRunning() {
		return fmt.Errorf("Docker daemon is not running. Please start Docker")
	}

	args := cm.buildComposeArgs("up")
	if detached {
		args = append(args, "-d")
	}

	return cm.runComposeCommand(args...)
}

// Down stops the Docker Compose services for a project.
func (cm *ProjectComposeManager) Down(removeVolumes bool) error {
	args := cm.buildComposeArgs("down")
	if removeVolumes {
		args = append(args, "-v")
	}

	return cm.runComposeCommand(args...)
}

// DownWithOptions stops the Docker Compose services for a project with additional options.
func (cm *ProjectComposeManager) DownWithOptions(options DownOptions) error {
	args := cm.buildComposeArgs("down")
	if options.RemoveVolumes {
		args = append(args, "-v")
	}
	if options.Force {
		args = append(args, "--remove-orphans")
	}

	return cm.runComposeCommand(args...)
}

// Build builds the Docker image for a project.
func (cm *ProjectComposeManager) Build(noCache bool, services ...string) error {
	args := cm.buildComposeArgs("build")
	if noCache {
		args = append(args, "--no-cache")
	}
	args = append(args, services...)

	return cm.runComposeCommand(args...)
}

// buildComposeArgs builds the base arguments for project docker-compose commands.
func (cm *ProjectComposeManager) buildComposeArgs(command string) []string {
	var args []string
	if strings.Contains(cm.composeCmd, " ") {
		args = []string{"compose"}
	}

	args = append(args, "-f", "docker-compose.yml")
	args = append(args, "-p", cm.projectCfg.Name)
	args = append(args, command)

	return args
}

// runComposeCommand runs a docker-compose command in the project's .phpier directory.
func (cm *ProjectComposeManager) runComposeCommand(args ...string) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(cm.workDir); err != nil {
		return fmt.Errorf("failed to change to project's .phpier directory %s: %w", cm.workDir, err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			logrus.Errorf("Failed to change back to original directory: %v", err)
		}
	}()

	if strings.Contains(cm.composeCmd, " ") {
		return cm.client.RunCommand("docker", args...)
	}
	return cm.client.RunCommand(cm.composeCmd, args...)
}

// NewGlobalComposeManager creates a new Docker Compose manager for the global services.
func NewGlobalComposeManager(globalCfg *config.GlobalConfig) (*GlobalComposeManager, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	return &GlobalComposeManager{
		client:     client,
		globalCfg:  globalCfg,
		composeCmd: client.GetDockerComposeCommand(),
		workDir:    filepath.Join(home, ".phpier"),
	}, nil
}

// Up starts the Docker Compose services for the global stack.
func (gcm *GlobalComposeManager) Up(detached bool) error {
	if !gcm.client.IsDockerRunning() {
		return fmt.Errorf("Docker daemon is not running. Please start Docker")
	}

	args := gcm.buildComposeArgs("up")
	if detached {
		args = append(args, "-d")
	}

	return gcm.runComposeCommand(args...)
}

// Down stops the Docker Compose services for the global stack.
func (gcm *GlobalComposeManager) Down(removeVolumes bool) error {
	args := gcm.buildComposeArgs("down")
	if removeVolumes {
		args = append(args, "-v")
	}

	return gcm.runComposeCommand(args...)
}

// DownWithOptions stops the Docker Compose services for the global stack with additional options.
func (gcm *GlobalComposeManager) DownWithOptions(options DownOptions) error {
	args := gcm.buildComposeArgs("down")
	if options.RemoveVolumes {
		args = append(args, "-v")
	}
	if options.Force {
		args = append(args, "--remove-orphans")
	}

	return gcm.runComposeCommand(args...)
}

// Build builds the Docker images for the global stack.
func (gcm *GlobalComposeManager) Build(noCache bool, services ...string) error {
	args := gcm.buildComposeArgs("build")
	if noCache {
		args = append(args, "--no-cache")
	}
	args = append(args, services...)

	return gcm.runComposeCommand(args...)
}

// buildComposeArgs builds the base arguments for global docker-compose commands.
func (gcm *GlobalComposeManager) buildComposeArgs(command string) []string {
	var args []string
	if strings.Contains(gcm.composeCmd, " ") {
		args = []string{"compose"}
	}

	args = append(args, "-f", "docker-compose.yml")
	args = append(args, "-p", "phpier")
	args = append(args, command)

	return args
}

// runComposeCommand runs a docker-compose command in the global .phpier directory.
func (gcm *GlobalComposeManager) runComposeCommand(args ...string) error {
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(gcm.workDir); err != nil {
		return fmt.Errorf("failed to change to global .phpier directory %s: %w", gcm.workDir, err)
	}
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			logrus.Errorf("Failed to change back to original directory: %v", err)
		}
	}()

	if strings.Contains(gcm.composeCmd, " ") {
		return gcm.client.RunCommand("docker", args...)
	}
	return gcm.client.RunCommand(gcm.composeCmd, args...)
}

// IsGlobalServiceRunning checks if global services (particularly Traefik) are running.
func (gcm *GlobalComposeManager) IsGlobalServiceRunning() (bool, error) {
	if !gcm.client.IsDockerRunning() {
		return false, fmt.Errorf("Docker daemon is not running")
	}

	// Check if Traefik container is running by looking for phpier project containers
	isRunning, err := gcm.client.IsContainerRunning(gcm.client.ctx, "phpier-traefik-1")
	if err != nil {
		// Try alternative container name pattern
		isRunning, err = gcm.client.IsContainerRunning(gcm.client.ctx, "phpier_traefik_1")
		if err != nil {
			logrus.Debugf("Could not check Traefik container status: %v", err)
			return false, nil // Don't error, just assume not running
		}
	}

	return isRunning, nil
}

// StartGlobalServicesIfNeeded starts global services if they are not running.
func (gcm *GlobalComposeManager) StartGlobalServicesIfNeeded() error {
	isRunning, err := gcm.IsGlobalServiceRunning()
	if err != nil {
		return fmt.Errorf("failed to check global service status: %w", err)
	}

	if isRunning {
		logrus.Debugf("Global services are already running")
		return nil
	}

	logrus.Infof("🚀 Starting global services (Traefik)...")

	// Start global services in detached mode
	if err := gcm.Up(true); err != nil {
		return fmt.Errorf("failed to start global services: %w", err)
	}

	logrus.Infof("✅ Global services started successfully")
	return nil
}
