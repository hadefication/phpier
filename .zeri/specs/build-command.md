# Feature Specification: build-command

## Overview
Add a `phpier build` command that builds (or rebuilds) the app container for the current project. This provides a convenient way to rebuild the app container when Dockerfile changes are made or when a fresh build is needed.

## Requirements
- Add `build` command to the Cobra CLI structure
- Build only the app container (not all services)
- Support forcing a rebuild (--no-cache flag)
- Work within the current project directory (must have .phpier/ folder)
- Provide clear feedback on build progress and success/failure
- Handle Docker build errors gracefully

## Implementation Notes
- Add new Cobra command in the appropriate command file
- Use Docker SDK for Go to execute the build operation
- Build the app container using the project's Dockerfile.php
- Target only the app service from docker-compose.yml
- Include --no-cache option for clean rebuilds
- Validate project directory before attempting build
- Use consistent error handling patterns with other commands
- Follow existing command structure and naming conventions

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete