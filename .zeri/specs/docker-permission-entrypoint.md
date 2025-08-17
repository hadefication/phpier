# Feature Specification: docker-permission-entrypoint

## Overview
Implement a Docker entrypoint script to solve mounted file permission issues in PHP containers. This script will handle user ID mapping between host and container to prevent permission conflicts with mounted volumes.

## Requirements
- Create an entrypoint script that handles WWWUSER environment variable
- Map container user to host user ID when WWWUSER is provided
- Initialize Composer directory with proper permissions
- Support both command execution and default supervisord startup
- Use gosu for proper privilege dropping
- Update PHP-FPM pool configuration with correct user

## Implementation Notes
- Add entrypoint.sh script to Docker templates
- Install gosu package in Dockerfile templates
- Set script as ENTRYPOINT in Dockerfile
- Support all PHP versions (5.6-8.4)
- Handle Composer directory creation and permissions
- Update PHP-FPM pool.d/www.conf configuration
- Use supervisord as default process when no command provided

## TODO
- [x] Design and plan implementation
- [x] Create entrypoint.sh template script
- [x] Update Dockerfile templates to install gosu
- [x] Add entrypoint script to Dockerfile templates
- [x] Test with different PHP versions
- [x] Verify permission handling works correctly
- [x] Mark specification as complete