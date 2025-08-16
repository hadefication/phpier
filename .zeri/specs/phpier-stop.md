# Feature Specification: phpier-stop

## Overview
The `phpier stop` command stops the global phpier services. This command provides a convenient way to stop the global infrastructure (Traefik, shared databases, monitoring services) without needing to use the full `phpier global down` command.

## Requirements

### Core Functionality
- Stop global phpier services (Traefik, shared databases, etc.)
- Check if any projects are still running and warn user
- Provide clear feedback about service shutdown status
- Graceful shutdown of services in proper order

### Command Interface
- Primary command: `phpier stop`
- Optional flag: `--force` to force stop services even if projects are running
- Optional flag: `--remove-volumes` to also remove global volumes (dangerous operation)
- Provide warnings when stopping services will affect running projects

### Service Management
- Stop global compose stack using existing global compose manager
- Check for running phpier projects before stopping global services
- Warn user if stopping global services will affect running projects
- Provide clear status messages during shutdown process

### Error Handling
- Handle cases where services are already stopped
- Provide clear error messages for Docker connectivity issues
- Continue with shutdown even if some services fail to stop gracefully
- Handle cases where global configuration is missing

## Implementation Notes

### Technical Considerations
- Reuse existing global compose manager for consistency
- Follow same patterns as `phpier global down` command
- Use Docker Compose down command for global services
- Check for running projects before stopping global services

### Files to Modify
- `cmd/stop.go` - Main stop command implementation
- Use existing `internal/docker/compose.go` - GlobalComposeManager
- Use existing `internal/docker/client.go` - GetRunningPhpierProjects
- Use existing `internal/config/config.go` - Global configuration loading

### Integration Points
- Integrate with existing Docker Compose management
- Use same configuration loading as global up/down commands
- Use existing project detection functionality
- Follow existing error handling patterns

### Architectural Patterns
- Follow existing Cobra command structure
- Use Viper for flag and configuration management
- Implement proper error handling with Go error types
- Use existing global compose manager interface

### Testing Strategy
- Unit tests for command logic and flag parsing
- Integration tests with existing global compose functionality
- Test project detection and warning logic
- Test error scenarios (Docker unavailable, services already stopped)
- Test flag combinations and edge cases

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete