# Feature Specification: self-update

## Overview
Implement a `phpier self-update` command that allows users to automatically update phpier to the latest version from GitHub releases. This provides a convenient way for users to keep their phpier installation current without manual download and installation.

## Requirements

### Core Functionality
- `phpier self-update` command to update to the latest stable release
- `phpier self-update --version <version>` to update to a specific version
- `phpier self-update --check` to check for available updates without installing
- Automatic detection of current platform (OS/architecture)
- Download and replace the current binary atomically
- Backup current binary before replacement
- Rollback capability if update fails
- Progress indicator during download

### User Interface Requirements
- Clear output showing current version, latest version, and update status
- Confirmation prompt before proceeding with update (unless `--force` flag used)
- `--force` flag to skip confirmation prompts
- `--check-only` or `--check` flag to only check for updates
- Verbose output with `--verbose` flag showing detailed progress
- Error messages with helpful troubleshooting information

### Integration Requirements
- Use GitHub API to fetch latest release information
- Support GitHub release assets for different platforms (linux/darwin/windows, amd64/arm64)
- Respect existing CLI patterns and error handling from phpier
- Integrate with existing version display in `phpier version` command
- Handle network connectivity issues gracefully

### Security Requirements
- Verify downloaded binary integrity using checksums
- Use HTTPS for all network requests
- Validate GitHub API responses
- Ensure atomic binary replacement to prevent corruption
- Clean up temporary files after update

### Performance Requirements
- Show download progress for files >1MB
- Use resume capability for interrupted downloads
- Timeout handling for network requests
- Efficient binary replacement without service interruption

## Implementation Notes

### Technical Considerations
- Use Go's built-in HTTP client with progress tracking
- Leverage GitHub API v4 (REST) for release information
- Binary naming convention: `phpier-<os>-<arch>` (e.g., `phpier-darwin-amd64`)
- Atomic file operations using temp files and rename
- Cross-platform executable permissions handling

### Files to Create/Modify
- `cmd/self_update.go` - New command implementation
- `internal/updater/` - New package for update logic
  - `updater.go` - Core update functionality
  - `github.go` - GitHub API integration
  - `progress.go` - Download progress tracking
- `cmd/root.go` - Add self-update command registration
- Tests for all new components

### Dependencies to Add
- HTTP client for downloads (use stdlib)
- JSON parsing for GitHub API (use stdlib)
- File operations and permissions (use stdlib)
- Progress bar library (consider adding to go.mod if needed)

### Integration Points
- Use existing error handling patterns from `internal/errors`
- Follow existing command structure and flag patterns
- Integrate with existing version information in `cmd/root.go`
- Use existing logging configuration

### Testing Strategy
- Unit tests for updater logic with mocked HTTP responses
- Integration tests with test GitHub repository
- Cross-platform compatibility tests
- Error scenario testing (network failures, permission issues)
- Atomic update operation testing

### GitHub Release Requirements
- Releases must include pre-built binaries for supported platforms
- Release assets should follow naming convention
- Checksums should be provided for integrity verification
- Release notes should be available via GitHub API

## TODO
- [x] Design and plan implementation architecture
- [x] Create updater package structure
- [x] Implement GitHub API integration
- [x] Implement download and progress tracking
- [x] Implement updater functionality using install script approach
- [x] Add self-update command with all flags
- [x] Add comprehensive error handling
- [x] Write unit tests for updater logic
- [x] Test basic implementation functionality
- [x] Add documentation and help text
- [x] Review and refine implementation
- [x] Mark specification as complete