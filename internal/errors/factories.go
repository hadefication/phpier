package errors

import "fmt"

// Docker-related error factories

// NewDockerNotFoundError creates a Docker not found error
func NewDockerNotFoundError() *PhpierError {
	return NewPhpierError(ErrorTypeDockerNotFound, "Docker is not installed or not found in PATH").
		WithSuggestion("Install Docker from https://docs.docker.com/get-docker/").
		WithSuggestion("Ensure Docker is added to your system PATH")
}

// NewDockerNotRunningError creates a Docker not running error
func NewDockerNotRunningError() *PhpierError {
	return NewPhpierError(ErrorTypeDockerNotRunning, "Docker daemon is not running").
		WithSuggestion("Start Docker Desktop or Docker daemon").
		WithSuggestion("Check Docker service status: 'docker version'")
}

// NewDockerComposeNotFoundError creates a Docker Compose not found error
func NewDockerComposeNotFoundError() *PhpierError {
	return NewPhpierError(ErrorTypeDockerComposeError, "Docker Compose is not available").
		WithSuggestion("Install Docker Compose or use Docker with compose plugin").
		WithSuggestion("For newer Docker installations, try 'docker compose' instead of 'docker-compose'")
}

// NewContainerNotFoundError creates a container not found error
func NewContainerNotFoundError(containerName string) *PhpierError {
	return NewPhpierError(ErrorTypeContainerNotFound, fmt.Sprintf("Container '%s' not found", containerName)).
		WithContext("container", containerName).
		WithSuggestion("Check if the container exists: 'docker ps -a'").
		WithSuggestion("Start the phpier environment: 'phpier up'")
}

// NewContainerNotRunningError creates a container not running error
func NewContainerNotRunningError(containerName string) *PhpierError {
	return NewPhpierError(ErrorTypeContainerNotRunning, fmt.Sprintf("Container '%s' is not running", containerName)).
		WithContext("container", containerName).
		WithSuggestion("Start the container: 'phpier up'").
		WithSuggestion("Check container status: 'docker ps'")
}

// NewBuildFailedError creates a build failed error
func NewBuildFailedError(service string, cause error) *PhpierError {
	return WrapError(ErrorTypeBuildFailed, fmt.Sprintf("Failed to build %s service", service), cause).
		WithContext("service", service).
		WithSuggestion("Check Docker build logs for detailed error information").
		WithSuggestion("Ensure all required system dependencies are available").
		WithSuggestion("Try rebuilding with: 'phpier build --no-cache'")
}

// Configuration-related error factories

// NewInvalidConfigError creates an invalid configuration error
func NewInvalidConfigError(field string, value interface{}) *PhpierError {
	return NewPhpierError(ErrorTypeInvalidConfig, fmt.Sprintf("Invalid configuration for field '%s': %v", field, value)).
		WithContext("field", field).
		WithContext("value", value).
		WithSuggestion("Check your .phpier.yml configuration file").
		WithSuggestion("Run 'phpier init' to regenerate configuration")
}

// NewConfigNotFoundError creates a configuration not found error
func NewConfigNotFoundError() *PhpierError {
	return NewPhpierError(ErrorTypeConfigNotFound, "Configuration file '.phpier.yml' not found").
		WithSuggestion("Run 'phpier init <php-version>' to initialize a new project").
		WithSuggestion("Ensure you're in the correct directory")
}

// NewProjectNotInitializedError creates a project not initialized error
func NewProjectNotInitializedError() *PhpierError {
	return NewPhpierError(ErrorTypeConfigNotFound, "Project not initialized - no .phpier.yml found").
		WithSuggestion("Run 'phpier init <php-version>' to initialize a new project").
		WithSuggestion("Ensure you're in the correct project directory")
}

// NewInvalidPHPVersionError creates an invalid PHP version error
func NewInvalidPHPVersionError(version string, supported []string) *PhpierError {
	return NewPhpierError(ErrorTypeInvalidPHPVersion, fmt.Sprintf("Unsupported PHP version: %s", version)).
		WithContext("version", version).
		WithContext("supported_versions", supported).
		WithSuggestion(fmt.Sprintf("Use one of the supported versions: %v", supported)).
		WithSuggestion("Check available PHP versions with 'phpier init --help'")
}

// NewInvalidDatabaseTypeError creates an invalid database type error
func NewInvalidDatabaseTypeError(dbType string, supported []string) *PhpierError {
	return NewPhpierError(ErrorTypeInvalidDatabaseType, fmt.Sprintf("Unsupported database type: %s", dbType)).
		WithContext("database_type", dbType).
		WithContext("supported_types", supported).
		WithSuggestion(fmt.Sprintf("Use one of the supported types: %v", supported))
}

// NewPortConflictError creates a port conflict error
func NewPortConflictError(port int, services []string) *PhpierError {
	return NewPhpierError(ErrorTypePortConflict, fmt.Sprintf("Port %d is used by multiple services: %v", port, services)).
		WithContext("port", port).
		WithContext("conflicting_services", services).
		WithSuggestion("Change port configuration in .phpier.yml").
		WithSuggestion("Use different ports for conflicting services")
}

// NewRequiredFieldMissingError creates a required field missing error
func NewRequiredFieldMissingError(field string) *PhpierError {
	return NewPhpierError(ErrorTypeRequiredFieldMissing, fmt.Sprintf("Required field '%s' is missing", field)).
		WithContext("field", field).
		WithSuggestion(fmt.Sprintf("Add the '%s' field to your configuration", field)).
		WithSuggestion("Check your .phpier.yml file")
}

// File system-related error factories

// NewFileNotFoundError creates a file not found error
func NewFileNotFoundError(filename string) *PhpierError {
	return NewPhpierError(ErrorTypeFileNotFound, fmt.Sprintf("File not found: %s", filename)).
		WithContext("file", filename).
		WithSuggestion("Check if the file exists and path is correct").
		WithSuggestion("Ensure you have read permissions for the file")
}

// NewFilePermissionError creates a file permission error
func NewFilePermissionError(filename string, operation string) *PhpierError {
	return NewPhpierError(ErrorTypeFilePermission, fmt.Sprintf("Permission denied for %s operation on file: %s", operation, filename)).
		WithContext("file", filename).
		WithContext("operation", operation).
		WithSuggestion("Check file permissions and ownership").
		WithSuggestion("Ensure you have necessary permissions for the operation")
}

// NewDirectoryExistsError creates a directory exists error
func NewDirectoryExistsError(directory string) *PhpierError {
	return NewPhpierError(ErrorTypeDirectoryExists, fmt.Sprintf("Directory already exists: %s", directory)).
		WithContext("directory", directory).
		WithSuggestion("Use a different directory name").
		WithSuggestion("Remove the existing directory if it's no longer needed").
		WithSuggestion("Use --force flag to overwrite existing directory")
}

// NewTemplateError creates a template processing error
func NewTemplateError(template string, cause error) *PhpierError {
	return WrapError(ErrorTypeTemplateError, fmt.Sprintf("Failed to process template: %s", template), cause).
		WithContext("template", template).
		WithSuggestion("Check template syntax and variables").
		WithSuggestion("Ensure all required template variables are provided")
}

// Command execution-related error factories

// NewCommandFailedError creates a command failed error
func NewCommandFailedError(command string, args []string, cause error) *PhpierError {
	return WrapError(ErrorTypeCommandFailed, fmt.Sprintf("Command failed: %s", command), cause).
		WithContext("command", command).
		WithContext("arguments", args).
		WithSuggestion("Check command syntax and arguments").
		WithSuggestion("Ensure all required dependencies are installed")
}

// NewCommandNotFoundError creates a command not found error
func NewCommandNotFoundError(command string) *PhpierError {
	return NewPhpierError(ErrorTypeCommandNotFound, fmt.Sprintf("Command not found: %s", command)).
		WithContext("command", command).
		WithSuggestion(fmt.Sprintf("Install %s or ensure it's in your PATH", command)).
		WithSuggestion("Check if the required software is properly installed")
}

// NewInvalidArgumentsError creates an invalid arguments error
func NewInvalidArgumentsError(message string) *PhpierError {
	return NewPhpierError(ErrorTypeInvalidArguments, message).
		WithSuggestion("Check command usage with --help flag").
		WithSuggestion("Review command documentation")
}

// User interaction error factories

// NewUserAbortedError creates a user aborted error
func NewUserAbortedError(message string) *PhpierError {
	return NewPhpierError(ErrorTypeUserAborted, message).
		WithSuggestion("Use --force flag to override safety checks if needed").
		WithSuggestion("Resolve the underlying issue and try again")
}

// Internal error factories

// NewInternalError creates an internal error
func NewInternalError(message string, cause error) *PhpierError {
	return WrapError(ErrorTypeInternal, fmt.Sprintf("Internal error: %s", message), cause).
		WithSuggestion("This appears to be an internal error").
		WithSuggestion("Please report this issue with the full error message")
}

// Project management error factories

// NewProjectNotFoundError creates a project not found error
func NewProjectNotFoundError(projectName string) *PhpierError {
	return NewPhpierError(ErrorTypeProjectNotFound, fmt.Sprintf("Project '%s' not found", projectName)).
		WithContext("project_name", projectName).
		WithSuggestion("Check if the project exists and is initialized with 'phpier init'").
		WithSuggestion("Ensure the project name is correct").
		WithSuggestion("Use 'phpier list' to see available projects")
}

// NewMultipleProjectsFoundError creates a multiple projects found error
func NewMultipleProjectsFoundError(projectName string, paths []string) *PhpierError {
	return NewPhpierError(ErrorTypeMultipleProjects, fmt.Sprintf("Multiple projects named '%s' found", projectName)).
		WithContext("project_name", projectName).
		WithContext("paths", paths).
		WithSuggestion("Specify the full path to the project you want to use").
		WithSuggestion("Rename one of the projects to avoid conflicts").
		WithSuggestion("Navigate to the specific project directory and use 'phpier up' without arguments")
}

// NewProjectDiscoveryFailedError creates a project discovery failed error
func NewProjectDiscoveryFailedError(cause error) *PhpierError {
	return WrapError(ErrorTypeProjectDiscoveryFailed, "Failed to discover phpier projects", cause).
		WithSuggestion("Check file system permissions in your project directories").
		WithSuggestion("Ensure you have access to the directories being searched")
}
