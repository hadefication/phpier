# Feature Specification: reload-command

## Overview
The `phpier reload` command provides a convenient way to restart project services without rebuilding containers. This is essential for applying configuration changes, clearing stuck processes, or getting a fresh start when services become unresponsive. Unlike `down` + `up`, reload is optimized for quick restarts with useful options for different scenarios.

## Requirements

### Core Functionality
- **Graceful restart**: Stop current project services cleanly, then start them back up
- **Preserve data**: Maintain volumes and persistent data by default
- **Configuration reload**: Pick up changes to `.phpier.yml` without regenerating supporting files
- **Error handling**: Handle cases where services are already stopped or failed to start
- **Status reporting**: Clear feedback on restart progress and any issues

### Command Interface
```bash
phpier reload [flags]
```

### Flags
- **`-d, --detach`**: Run services in background (detached mode) after restart
- **`--build`**: Rebuild container images before restarting (for code/dependency changes)
- **`--force`**: Force stop containers that don't respond to graceful shutdown
- **`--timeout=30`**: Timeout in seconds for stopping containers (default: 30)
- **`--skip-global`**: Skip checking/starting global services during reload
- **`--pull`**: Pull latest base images before rebuilding (requires --build)
- **`--no-cache`**: Don't use cache when rebuilding (requires --build)

### Usage Examples
```bash
phpier reload                    # Basic reload - stop and start services
phpier reload -d                 # Reload and run in background
phpier reload --build            # Rebuild images then reload
phpier reload --force --timeout=10  # Force stop with custom timeout
phpier reload --build --pull --no-cache  # Full rebuild with latest images
```

### User Experience Requirements
- **Fast execution**: Optimized for quick restarts (< 30 seconds for typical projects)
- **Clear feedback**: Progress indicators and helpful error messages
- **Consistent behavior**: Same patterns as other phpier commands (up, down, build)

### Integration Requirements
- **Project detection**: Work from any subdirectory within phpier project
- **Global services**: Option to ensure global services are running
- **Docker Compose**: Leverage existing compose manager and patterns
- **Configuration**: Respect current `.phpier.yml` settings

## Implementation Notes

### Technical Architecture
- **Reuse existing patterns**: Extend `ProjectComposeManager` with reload functionality
- **Command structure**: Follow Cobra command patterns used in `up.go`, `down.go`
- **Flag handling**: Use same flag patterns and validation as existing commands
- **Error handling**: Consistent error wrapping and user-friendly messages

### Files to Modify
1. **`cmd/reload.go`** - New command implementation
2. **`internal/docker/compose.go`** - Add `Reload()` method to `ProjectComposeManager`
3. **`cmd/root.go`** - Register reload command
4. **Documentation** - Update help text and examples

### Implementation Flow
1. **Pre-checks**: Verify project is initialized and Docker is running
2. **Load config**: Read current `.phpier.yml` and global configuration
3. **Stop services**: Gracefully stop containers (with timeout/force options)
4. **Rebuild (optional)**: If `--build` flag, rebuild images with specified options
5. **Start services**: Start containers with specified detach mode
6. **Global services**: Ensure global services if not skipped
7. **Status report**: Confirm successful reload or report errors

### Docker Operations Sequence
```go
// Pseudo-code flow
composeManager.DownWithOptions(DownOptions{
    RemoveVolumes: false,
    Force: options.Force,
    Timeout: options.Timeout,
})

if options.Build {
    composeManager.Build(options.NoCache, options.Pull)
}

composeManager.Up(options.Detached)
```

### Safety Considerations
- **Timeout handling**: Prevent indefinite hangs with configurable timeouts
- **State validation**: Check that services actually started after reload
- **Rollback**: Consider how to handle failed reloads

### Testing Strategy
- **Unit tests**: Test reload logic with mocked Docker operations
- **Integration tests**: Test actual Docker container stop/start cycles
- **Flag combinations**: Test various flag combinations and edge cases
- **Error scenarios**: Test behavior with failed stops, starts, builds
- **Performance**: Verify reload times meet UX requirements

### Error Handling Patterns
- **Service detection**: Handle cases where no services are running
- **Partial failures**: Continue operation if some containers fail
- **Resource conflicts**: Handle port conflicts or resource constraints
- **Configuration errors**: Validate `.phpier.yml` before operations

## TODO
- [x] Design and plan implementation approach
- [x] Create `cmd/reload.go` with Cobra command structure  
- [x] Extend `ProjectComposeManager` with `Reload()` method
- [x] Implement flag parsing and validation
- [x] Add safety checks for timeout and force operations
- [x] Create unit tests for reload logic
- [x] Create integration tests with actual Docker operations
- [x] Test various flag combinations and edge cases
- [x] Update documentation and help text
- [x] Review and refine implementation
- [x] Mark specification as complete