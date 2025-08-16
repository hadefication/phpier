# Feature Specification: down-command-improvements

## Overview
Enhance both the `phpier down` and `phpier up` commands with additional flag options for better global service management:
1. Add `--global` flag to `phpier down` command for stopping global services
2. Add `--skip-global` flag to `phpier up` command to skip automatic global service startup

## Requirements

### Down Command Requirements
- Add `--global` flag to the `phpier down` command
- When `--global` flag is used, stop both project services and global services
- Maintain existing behavior when `--global` flag is not used (only stop project services)
- Provide clear feedback about which services are being stopped
- Include safety checks and user warnings similar to existing global service commands
- Support force flags to override safety checks when needed

### Up Command Requirements
- Add `--skip-global` flag to the `phpier up` command
- When `--skip-global` flag is used, skip the automatic global service startup check
- Maintain existing behavior when `--skip-global` flag is not used (auto-start global services)
- Provide clear feedback when global service startup is being skipped
- Allow users to start only project services without touching global infrastructure

## Implementation Notes

### Down Command Implementation
- Modify the down command in `cmd/down.go` to accept the `--global` flag
- Integrate with existing global service management functionality
- Reuse safety check patterns from the stop command implementation
- Update command help text to document the new flag
- Ensure consistent error handling and user messaging
- Follow established Cobra flag patterns used in other commands

### Up Command Implementation  
- Modify the up command in `cmd/up.go` to accept the `--skip-global` flag
- Add conditional logic to bypass global service startup when flag is present
- Update command help text to document the new flag
- Provide user feedback when global service startup is being skipped
- Maintain existing global service startup logic as default behavior
- Ensure flag works with existing up command functionality

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete