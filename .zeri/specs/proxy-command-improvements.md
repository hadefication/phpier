# Feature Specification: proxy-command-improvements

## Overview
Replace the existing static proxy commands (composer, php, npm, node, npx, artisan) with a unified `phpier proxy` command that dynamically executes tools in app containers. The command will have context-aware behavior: when run inside a phpier project, it executes tools in the local project's app container (`phpier proxy <tool>`), and when run globally, it allows targeting specific apps (`phpier proxy <app> <tool>`).

This consolidates all proxy functionality into a single, flexible command while maintaining backward compatibility through improved user experience.

## Requirements

### Core Functionality
- **Project Context**: `phpier proxy <tool> [args...]` - Execute tool in current project's app container
- **Global Context**: `phpier proxy <app> <tool> [args...]` - Execute tool in specified app's container
- Support command arguments and flags passthrough to the target tool
- Maintain interactive terminal support for commands that require it
- Handle commands that don't exist in the container gracefully
- Automatic context detection (project vs global scope)

### Context Detection
- Automatically detect if running inside a phpier project directory (presence of `.phpier.yml`)
- Parse command arguments based on context (project vs global)
- Validate app names when running in global context
- List available apps when global context is used without specifying an app

### Command Execution
- Execute any tool available in the app container dynamically
- Support common PHP development tools: composer, php, npm, yarn, artisan, etc.
- Handle arbitrary commands and tools not pre-defined
- Provide command existence validation before execution

### User Experience
- Context-aware command parsing and execution
- Clear error messages for invalid app names or missing containers
- Seamless command execution as if running locally
- Preserve exit codes from container commands
- Support stdin/stdout/stderr streams properly
- Handle interactive prompts and TTY allocation

### Integration Requirements
- Work with existing phpier project structure and global app management
- Require active phpier environment (containers running)
- Support both project-local and global app container targeting
- Integrate with existing error handling patterns

## Implementation Notes

### Technical Considerations
- Context detection logic to determine project vs global scope
- Dynamic argument parsing based on detected context
- Use Docker exec to run commands in the target app container
- Implement proper TTY allocation for interactive commands
- Handle signal forwarding (Ctrl+C, etc.)
- Support working directory mapping between host and container

### Files to Modify/Create
- **Remove existing static proxy commands**: Delete `cmd/composer.go`, `cmd/php.go`, `cmd/npm.go`, `cmd/node.go`, `cmd/npx.go`, `cmd/artisan.go`
- **Create new proxy command**: Create `cmd/proxy.go` with unified proxy functionality
- **Consolidate existing functionality**: Leverage existing `docker.ProxyCommand` and `ExecuteProxyCommand` from `internal/docker/exec.go`
- **Implement context detection**: Add utilities to detect project vs global context (check for `.phpier.yml`)
- **Add dynamic argument parsing**: Parse arguments based on detected context
- **Enhance container discovery**: Support both project and global app container targeting
- **Extend error handling**: Add context-aware error messages and validation

### Container Integration
- Detect running app container for current project or specified app
- Map current working directory to container equivalent
- Ensure proper user permissions in container execution
- Handle cases where containers are not running
- Support app name resolution to container names

### Command Resolution Strategy
1. **Context Detection**: Check if `.phpier.yml` exists in current directory
2. **Argument Parsing**: 
   - Project context: `proxy <tool> [args...]`
   - Global context: `proxy <app> <tool> [args...]`
3. **Container Discovery**: Find target app container (current project or specified app)
4. **Validation**: Verify app container is running
5. **Execution**: Run command in container with proper TTY and signal handling
6. **Output**: Return appropriate exit codes and streams

### Testing Strategy
- Unit tests for context detection and argument parsing
- Integration tests with Docker containers for both project and global contexts
- Test common development workflows (composer, php, npm) through proxy command
- Test error scenarios (container not running, command not found, invalid app names)
- Test interactive command handling and TTY allocation
- Verify backward compatibility for existing workflows
- Test migration from static commands to proxy command

## TODO
- [x] Design and plan implementation approach
- [x] Remove existing static proxy command files (`composer.go`, `php.go`, `npm.go`, `node.go`, `npx.go`, `artisan.go`)
- [x] Create unified `cmd/proxy.go` with context-aware argument parsing
- [x] Implement project context detection (check for `.phpier.yml` file)
- [x] Add global app container discovery and validation
- [x] Enhance existing `docker.ExecuteProxyCommand` to support global app targeting
- [x] Implement dynamic tool execution without predefined command list
- [x] Add comprehensive error handling for both contexts
- [x] Create unit tests for context detection and argument parsing
- [x] Create integration tests for project and global proxy execution
- [x] Test common development workflows through new proxy command
- [x] Verify removal of old commands doesn't break existing functionality
- [x] Update any documentation or help text
- [x] Clean up Docker testing artifacts after implementation
- [x] Review and refine implementation
- [x] Mark specification as complete