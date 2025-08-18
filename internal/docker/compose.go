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
	Timeout       int
}

// ReloadOptions represents options for the reload operation
type ReloadOptions struct {
	Detached bool
	Build    bool
	Force    bool
	Timeout  int
	Pull     bool
	NoCache  bool
}

// ComposeManager interface for Docker Compose operations.
type ComposeManager interface {
	Up(detached bool) error
	Down(removeVolumes bool) error
	DownWithOptions(options DownOptions) error
	Build(noCache bool, services ...string) error
	Reload(options ReloadOptions) error
	Logs(service string, follow bool, tail int, since string) error
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

	// Find the project root directory (where docker-compose.yml exists)
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("failed to find project root: %w", err)
	}

	// Use the project root as working directory (where docker-compose.yml is)
	workDir := projectRoot

	return &ProjectComposeManager{
		client:     client,
		projectCfg: projectCfg,
		globalCfg:  globalCfg,
		composeCmd: client.GetDockerComposeCommand(),
		workDir:    workDir,
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
	if options.Timeout > 0 {
		args = append(args, "--timeout", fmt.Sprintf("%d", options.Timeout))
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

// Reload restarts the Docker Compose services for a project with various options.
func (cm *ProjectComposeManager) Reload(options ReloadOptions) error {
	if !cm.client.IsDockerRunning() {
		return fmt.Errorf("Docker daemon is not running. Please start Docker")
	}

	// Step 1: Stop services
	logrus.Infof("🛑 Stopping project services...")
	downOptions := DownOptions{
		RemoveVolumes: false,
		Force:         options.Force,
		Timeout:       options.Timeout,
	}
	if err := cm.DownWithOptions(downOptions); err != nil {
		return fmt.Errorf("failed to stop services: %w", err)
	}

	// Step 2: Build if requested
	if options.Build {
		logrus.Infof("🔨 Building project image...")

		// Pull latest images if requested
		if options.Pull {
			logrus.Infof("📥 Pulling latest base images...")
			pullArgs := cm.buildComposeArgs("pull")
			if err := cm.runComposeCommand(pullArgs...); err != nil {
				return fmt.Errorf("failed to pull images: %w", err)
			}
		}

		// Build with options
		if err := cm.Build(options.NoCache, "app"); err != nil {
			return fmt.Errorf("failed to build project image: %w", err)
		}
	}

	// Step 3: Start services
	logrus.Infof("🚀 Starting project services...")
	if err := cm.Up(options.Detached); err != nil {
		return fmt.Errorf("failed to start services: %w", err)
	}

	return nil
}

// Logs displays logs from the Docker Compose services for a project.
func (cm *ProjectComposeManager) Logs(service string, follow bool, tail int, since string) error {
	if !cm.client.IsDockerRunning() {
		return fmt.Errorf("Docker daemon is not running. Please start Docker")
	}

	args := cm.buildComposeArgs("logs")

	// Add follow flag
	if follow {
		args = append(args, "-f")
	}

	// Add tail flag
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}

	// Add since flag
	if since != "" {
		args = append(args, "--since", since)
	}

	// Add service if specified
	if service != "" {
		args = append(args, service)
	}

	return cm.runComposeCommand(args...)
}

// buildComposeArgs builds the base arguments for project docker-compose commands.
func (cm *ProjectComposeManager) buildComposeArgs(command string) []string {
	var args []string
	if strings.Contains(cm.composeCmd, " ") {
		args = []string{"compose"}
	}

	args = append(args, "-f", ".phpier.yml")
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
	if options.Timeout > 0 {
		args = append(args, "--timeout", fmt.Sprintf("%d", options.Timeout))
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

// Reload is not supported for global services - use 'phpier global down' and 'phpier global up' instead.
func (gcm *GlobalComposeManager) Reload(options ReloadOptions) error {
	return fmt.Errorf("reload is not supported for global services - use 'phpier global down' and 'phpier global up' instead")
}

// Logs displays logs from the Docker Compose services for the global stack.
func (gcm *GlobalComposeManager) Logs(service string, follow bool, tail int, since string) error {
	if !gcm.client.IsDockerRunning() {
		return fmt.Errorf("Docker daemon is not running. Please start Docker")
	}

	args := gcm.buildComposeArgs("logs")

	// Add follow flag
	if follow {
		args = append(args, "-f")
	}

	// Add tail flag
	if tail > 0 {
		args = append(args, "--tail", fmt.Sprintf("%d", tail))
	}

	// Add since flag
	if since != "" {
		args = append(args, "--since", since)
	}

	// Add service if specified
	if service != "" {
		args = append(args, service)
	}

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

// findProjectRoot finds the project root directory by looking for .phpier.yml file
func findProjectRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	// Walk up the directory tree looking for .phpier.yml
	dir := currentDir
	for {
		configPath := filepath.Join(dir, ".phpier.yml")
		if _, err := os.Stat(configPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the root directory
			break
		}
		dir = parent
	}

	return "", fmt.Errorf(".phpier.yml not found in current directory or any parent directory")
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
