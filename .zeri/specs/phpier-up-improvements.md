# Feature Specification: phpier-up-improvements

## Overview
Enhance the `phpier up` command to intelligently manage global services. The command should automatically check and start global services (Traefik) before starting project services.

## Requirements

### Functional Requirements
- **Global Service Detection**: Check if global services (Traefik) are running before starting project services
- **Automatic Global Start**: If global services are stopped, start them automatically before proceeding
- **Status Feedback**: Provide clear user feedback about global service status and actions taken
- **Error Handling**: Handle cases where global services fail to start or projects aren't initialized

### User Interface Requirements
- **Clear Output**: Display status messages for global service checks and actions
- **Consistent Messaging**: Follow existing phpier CLI messaging patterns

### Integration Requirements
- **Docker Compose**: Integrate with existing docker-compose up functionality
- **Global Services**: Leverage existing global service management (start/stop commands)
- **Project Detection**: Use existing project validation logic
- **Error Types**: Use established error handling patterns from the codebase

### Performance Requirements
- **Fast Status Check**: Global service status check should be lightweight and fast
- **Parallel Operations**: Where possible, run operations in parallel to minimize startup time
- **Cached Status**: Avoid redundant Docker API calls for service status

## Implementation Notes

### Technical Considerations
- **Service Detection Logic**: Use Docker API to check if Traefik container is running
- **Command Structure**: Extend existing Cobra command with enhanced logic
- **State Management**: Track global service startup to avoid duplicate operations

### Files to Modify
- `cmd/up.go` - Main up command implementation
- `internal/docker/` - Docker service management utilities
- `internal/services/global.go` - Global service management functions
- `cmd/up_test.go` - Add comprehensive tests for new functionality

### Integration Points
- **Global Service Commands**: Reuse logic from `phpier start` command
- **Project Validation**: Use existing project detection from current up command
- **Docker Compose**: Extend current docker-compose integration
- **Error Handling**: Integrate with existing error types and user messaging

### Architectural Patterns
- **Command Pattern**: Follow Cobra command structure with RunE function
- **Service Layer**: Use service layer for Docker operations and global service management
- **Configuration**: Leverage Viper for any new configuration options
- **Dependency Injection**: Use interfaces for Docker operations to enable testing

### Testing Strategy
- **Unit Tests**: Test global service detection logic in isolation
- **Integration Tests**: Test full up workflow with mock Docker containers
- **Command Tests**: Use Cobra command testing helpers for CLI interactions
- **Error Scenarios**: Test failure cases (Docker down, invalid projects, etc.)

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete