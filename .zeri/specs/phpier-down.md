# Feature Specification: phpier-down

## Overview
The `phpier down` command stops and removes project containers and services. It provides an optional `--stop-global` flag to also stop global phpier services (Traefik, shared databases, etc.) when shutting down a project.

## Requirements

### Core Functionality
- Stop and remove project-specific containers defined in docker-compose.yml
- Clean up project networks and volumes (non-persistent)
- Preserve persistent data volumes by default
- Graceful shutdown of services in proper order

### Command Interface
- Primary command: `phpier down`
- Optional flag: `--stop-global` to also stop global services
- Optional flag: `--remove-volumes` to remove all volumes including persistent data
- Optional flag: `--force` to force remove containers without graceful shutdown

### Global Services Management
- When `--stop-global` is used, stop global services after project shutdown
- Global services include: Traefik, shared databases, monitoring services
- Check if other projects are running before stopping global services
- Warn user if stopping global services will affect other running projects

### Error Handling
- Handle cases where containers are already stopped
- Provide clear error messages for Docker connectivity issues
- Continue with remaining cleanup if some containers fail to stop

## Implementation Notes

### Technical Considerations
- Use Docker Compose down command for project services
- Implement global service detection and management
- Check for running containers before stopping global services
- Follow Docker Compose service dependency order for shutdown

### Files to Modify
- `cmd/down.go` - Main down command implementation
- `internal/docker/compose.go` - Docker Compose operations
- `internal/config/config.go` - Configuration for global services
- Add global service management functions

### Integration Points
- Integrate with existing Docker Compose management
- Use same configuration loading as other commands
- Coordinate with global service management (similar to global up/down)

### Architectural Patterns
- Follow existing Cobra command structure
- Use Viper for flag and configuration management
- Implement proper error handling with Go error types
- Use dependency injection for Docker operations

### Testing Strategy
- Unit tests for command logic and flag parsing
- Integration tests with testcontainers for Docker operations
- Test global service detection and management
- Test error scenarios (containers already stopped, Docker unavailable)

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete