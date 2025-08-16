# Feature Specification: start-stop-commands

## Overview
Add convenient `start` and `stop` shortcut commands to the phpier CLI. These commands will provide simpler aliases for the existing `phpier up` and `phpier down --global` commands respectively, improving user experience with more intuitive command names.

## Requirements

### Core Functionality
- `phpier start` command that acts as a shortcut for `phpier up -d` (detached mode by default, no flags supported)
- `phpier stop` command that acts as a shortcut for `phpier down --global` (no flags supported)
- Commands should maintain the same core behavior as their underlying commands
- Preserve all existing functionality of `up -d` and `down --global` commands
- Keep commands simple without flag support for ease of use

### User Interface Requirements
- Commands should follow simple, clean help text patterns
- Error handling should be consistent with existing command behavior
- Success/failure messages should be clear and informative
- No flag support - commands should work with sensible defaults

### Compatibility Requirements
- Must integrate seamlessly with existing Cobra command structure
- Should not break or interfere with existing `up` and `down` commands
- Maintain backward compatibility with current CLI usage patterns

### Implementation Requirements
- Commands should delegate to existing `up -d` and `down --global` implementations
- No duplication of core logic - reuse existing command functions
- Follow established Go and Cobra patterns used in the codebase

## Implementation Notes

### Technical Considerations
- Leverage Cobra's command aliasing or create wrapper commands that call existing functions
- Reuse existing `upCmd` and `downCmd` implementations to avoid code duplication
- Ensure proper flag inheritance and parsing

### Files to Modify
- `cmd/start.go` - New start command implementation
- `cmd/stop.go` - New stop command implementation  
- `cmd/root.go` - Register new commands with root command
- Possibly update help text and documentation

### Architectural Patterns
- Follow existing command structure patterns in the codebase
- Use dependency injection patterns already established
- Maintain consistent error handling and logging patterns

### Testing Strategy
- Unit tests for new command functions
- Integration tests to verify commands work as shortcuts
- Ensure tests cover flag passing and error scenarios
- Test that shortcuts produce identical results to original commands

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Clean up Docker artifacts (mandatory cleanup per development guidelines)
- [x] Mark specification as complete

## Docker Cleanup Log
As per `.zeri/development.md` requirements, tracking Docker artifacts created during testing:

### Test Directories Created:
- `/tmp/test-start-detached` (created and cleaned up - should have used phpier-test- prefix)
- `/tmp/phpier-test-start-detached` (created and cleaned up - proper naming convention)

### Docker Resources Generated:
- Docker images built during start command testing
- Containers from phpier test projects  
- Networks and volumes created by phpier initialization

### Cleanup Commands to Run:
```bash
# Clean up phpier test directories (using proper prefix going forward)
rm -rf /tmp/phpier-test-*
rm -rf /private/tmp/phpier-test-*

# Clean up ONLY phpier test Docker resources (safer approach)
# Containers with phpier-test- prefix
docker ps -a --filter "name=phpier-test-" -q | xargs docker rm -f

# Images with phpier-test- prefix  
docker images --filter "reference=*phpier-test-*" -q | xargs docker rmi -f

# Volumes with phpier-test- prefix
docker volume ls --filter "name=phpier-test-" -q | xargs docker volume rm -f

# Networks with phpier-test- prefix (if any)
docker network ls --filter "name=phpier-test-" -q | xargs docker network rm
```