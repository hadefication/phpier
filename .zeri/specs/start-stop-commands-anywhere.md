# Feature Specification: start-stop-commands-anywhere

## Overview
Enable `phpier start` and `phpier stop` commands to be runnable from any directory. When executed outside of a phpier project directory, these commands should manage global services only. This improves user experience by allowing global service management without needing to navigate to a specific project directory.

## Requirements

### Core Functionality
- `phpier start` and `phpier stop` commands should work from any directory
- When run outside a phpier project directory:
  - `phpier start` should start global services only
  - `phpier stop` should stop global services only
- When run inside a phpier project directory:
  - `phpier start` should start global services AND bring up the project (equivalent to `phpier global up` + `phpier up`)
  - `phpier stop` should stop the project AND global services (equivalent to `phpier down` + `phpier global down`)

### Directory Detection Requirements
- Implement detection logic to determine if current directory contains a phpier project
- Check for presence of `.phpier.yml` configuration file
- Check for `.phpier/` directory with project configuration
- Handle edge cases like nested directories and symlinks

### User Interface Requirements
- Provide clear feedback about what services are being started/stopped
- Distinguish between global-only operations and project operations in output messages
- Maintain existing safety checks and user warnings
- Show different help text or examples based on context

### Integration and Compatibility Requirements
- Maintain backward compatibility with existing command behavior
- Ensure existing project-based functionality remains unchanged
- Work with existing Docker Compose configurations
- Support all current command flags and options

### Performance and Security Requirements
- Directory detection should be fast and not impact command startup time
- Avoid unnecessary file system operations
- Maintain existing security practices for Docker operations
- Proper error handling to prevent command hangs

## Implementation Notes

### Technical Considerations and Dependencies
- Existing global service management functionality in CLI
- Current project detection mechanisms
- Docker Compose service definitions for global services
- Cobra command framework structure

### Files or Components That Need Modification
- Command handlers for `start` and `stop` commands
- Service management utility functions
- Configuration loading logic
- Command help text and usage examples

### Integration Points with Existing Systems
- Cobra command definitions and argument parsing
- Docker Compose service management
- Configuration file loading (Viper)
- Error handling and user feedback systems

### Architectural Decisions and Patterns to Follow
- Add directory detection utility function following existing patterns
- Modify command handlers to determine operation scope before execution
- Update service management calls with scope parameter
- Enhance user feedback messages with scope information
- Follow existing error handling patterns

### Testing Strategy and Requirements
- Unit tests for directory detection functionality
- Integration tests for global-only command execution
- Test commands from various directory contexts (project, non-project, nested)
- Test error conditions (missing Docker, permission issues)
- Test backward compatibility with existing project-based workflows

### Command Behavior Matrix
| Command | In Project Dir | Outside Project Dir |
|---------|---------------|-------------------|
| `phpier start` | Start global services + bring up project | Start global services only |
| `phpier stop` | Stop project + stop global services | Stop global services only |

## TODO
- [x] Design directory detection utility function
- [x] Implement phpier project detection logic (`isPhpierProject()`)
- [x] Modify `start` command to work from any directory
- [x] Modify `stop` command to work from any directory
- [x] Update command output messages to indicate operation scope
- [x] Update command help text to reflect new global functionality
- [x] Add unit tests for directory detection functionality
- [x] Add integration tests for global-only command execution
- [x] Test commands from various directory contexts
- [x] Test error handling for missing Docker or permissions
- [x] Test backward compatibility with existing workflows
- [x] Review and refine implementation
- [x] Mark specification as complete