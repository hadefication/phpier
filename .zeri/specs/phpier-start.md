# Feature Specification: phpier-start

## Overview
The `phpier start` command starts the global phpier services. This command provides a convenient way to start the global infrastructure (Traefik, shared databases, monitoring services) without needing to use the full `phpier global up` command.

## Requirements

### Core Functionality
- Start global phpier services (Traefik, shared databases, etc.)
- Check if global services are already running to avoid conflicts
- Provide clear feedback about service startup status
- Ensure proper dependency order when starting services

### Command Interface
- Primary command: `phpier start`
- Optional flag: `--detach` or `-d` to run services in background (default behavior)
- Optional flag: `--build` to rebuild global services before starting
- Optional flag: `--force` to force restart services if already running

### Service Management
- Start global compose stack using existing global compose manager
- Verify Docker daemon is running before attempting to start services
- Check if services are already running and handle appropriately
- Provide clear status messages during startup process

### Error Handling
- Handle cases where Docker daemon is not running
- Provide clear error messages for Docker connectivity issues
- Handle cases where global configuration is missing or invalid
- Handle port conflicts and other Docker-related errors

## Implementation Notes

### Technical Considerations
- Reuse existing global compose manager for consistency
- Follow same patterns as `phpier global up` command
- Use Docker Compose up command for global services
- Ensure proper error handling and user feedback

### Files to Modify
- `cmd/start.go` - Main start command implementation
- Use existing `internal/docker/compose.go` - GlobalComposeManager
- Use existing `internal/config/config.go` - Global configuration loading

### Integration Points
- Integrate with existing Docker Compose management
- Use same configuration loading as global up/down commands
- Follow existing error handling patterns
- Use same logging and output formatting

### Architectural Patterns
- Follow existing Cobra command structure
- Use Viper for flag and configuration management
- Implement proper error handling with Go error types
- Use existing global compose manager interface

### Testing Strategy
- Unit tests for command logic and flag parsing
- Integration tests with existing global compose functionality
- Test error scenarios (Docker unavailable, services already running)
- Test flag combinations and edge cases

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete