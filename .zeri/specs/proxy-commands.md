# Feature Specification: proxy-commands

## Overview
Implement proxy commands that allow running app container tools directly through the phpier CLI. This enables developers to execute commands like `phpier composer install`, `phpier node --version`, `phpier npm install`, etc. without manually accessing the Docker container.

## Requirements

### Core Functionality
- Support proxy commands for common development tools:
  - `phpier php <args>` - Execute PHP commands
  - `phpier composer <args>` - Execute Composer commands
  - `phpier node <args>` - Execute Node.js commands
  - `phpier npm <args>` - Execute NPM commands
  - `phpier npx <args>` - Execute NPX commands
  - `phpier yarn <args>` - Execute Yarn commands (if available)
  - `phpier artisan <args>` - Execute Laravel Artisan commands
- Forward all arguments and flags to the target command
- Preserve output formatting and colors
- Handle interactive commands (e.g., `composer require` with prompts)
- Return proper exit codes from the proxied commands

### User Interface Requirements
- Commands should feel native and transparent
- Support `--help` flag for each proxy command
- Provide clear error messages when tools are not available
- Show helpful suggestions when commands fail

### Integration Requirements
- Work with existing phpier project structure
- Detect if phpier environment is running before executing commands
- Auto-start containers if they're stopped (with user confirmation)
- Support both global and project-specific contexts

### Performance Requirements
- Fast command execution with minimal overhead
- Efficient container communication
- Support for long-running commands (e.g., `npm run watch`)

## Implementation Notes

### Technical Considerations
- Use `docker exec` to run commands in the app container
- Implement as Cobra subcommands with dynamic registration
- Handle TTY allocation for interactive commands
- Forward environment variables appropriately
- Support working directory context within container

### Files to Modify/Create
- `cmd/proxy.go` - Main proxy command implementation
- `cmd/php.go` - PHP-specific proxy command
- `cmd/composer.go` - Composer proxy command
- `cmd/node.go` - Node.js proxy command
- `cmd/npm.go` - NPM proxy command
- `cmd/npx.go` - NPX proxy command
- `cmd/artisan.go` - Laravel Artisan proxy command
- `internal/docker/exec.go` - Docker exec utilities
- Update `cmd/root.go` to register proxy commands

### Integration Points
- Leverage existing Docker client in `internal/docker/client.go`
- Use existing project detection logic
- Integrate with current error handling patterns
- Follow established Cobra command structure

### Architectural Decisions
- Each tool gets its own Cobra command for better help and validation
- Shared proxy execution logic in internal package
- Support for tool availability detection
- Consistent error handling across all proxy commands

### Testing Strategy
- Unit tests for command parsing and validation
- Integration tests with actual Docker containers
- Test interactive command handling
- Verify exit code forwarding
- Test with various PHP versions and tool availability

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete