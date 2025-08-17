# Feature Specification: sh-command

## Overview
Add a `phpier sh` command that provides direct shell access to the app container. This allows developers to interactively explore the container environment, run commands, debug issues, and perform manual operations within the PHP development environment.

## Requirements

### Core Functionality
- `phpier sh` command opens an interactive shell in the app container
- Use `/bin/bash` as the default shell (fallback to `/bin/sh` if bash unavailable)
- Execute as the `www-data` user to match the web server environment
- Start in the `/var/www/html` working directory (project root)
- Support both interactive mode and single command execution
- Handle cases where app container is not running with clear error messages

### Command Interface
- `phpier sh` - Opens interactive shell session
- `phpier sh -c "command"` - Execute single command and return output
- `phpier sh --user <user>` - Optional: specify different user (default: www-data)
- Exit shell with standard exit commands (exit, ctrl+d)
- Preserve exit codes from commands executed in the container

### Error Handling
- Check if app container exists and is running before attempting shell access
- Provide clear error messages for common failure scenarios
- Suggest `phpier start` or `phpier up` if services are not running
- Handle Docker connectivity issues gracefully

### Integration Requirements
- Follow existing phpier command patterns and error handling
- Use the same Docker client configuration as other commands
- Integrate with existing project detection logic
- Maintain consistency with other container interaction commands

## Implementation Notes

### Technical Considerations
- Use Docker SDK for Go to execute `docker exec` equivalent operations
- Implement interactive terminal handling for proper shell experience
- Handle TTY allocation for interactive sessions
- Support both Windows and Unix-like systems for terminal interaction

### Files to Modify/Create
- `cmd/sh.go` - New Cobra command for shell access
- `cmd/root.go` - Register the new sh command
- `internal/docker/container.go` - Add shell execution functionality
- `internal/project/project.go` - Container status checking utilities

### Docker Integration
- Use `docker exec -it <container> /bin/bash` equivalent
- Container name follows existing naming patterns (project-name_app_1)
- Handle container selection when multiple containers exist
- Respect Docker context and connection settings

### Error Handling Patterns
- Follow existing phpier error types and patterns
- Provide actionable error messages with suggested solutions
- Log appropriate debug information for troubleshooting

### Testing Strategy
- Unit tests for command parsing and validation
- Integration tests using testcontainers
- Test interactive and non-interactive modes
- Test error scenarios (container not running, Docker unavailable)
- Cross-platform terminal handling tests

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete