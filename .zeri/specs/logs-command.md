# Feature Specification: logs-command

## Overview
Add a `phpier logs` command to view Docker container logs for the current project's services. This command should provide easy access to logs from all running containers with options for filtering, following, and formatting.

## Requirements
- [x] Implement `phpier logs` command that displays logs from all project containers
- [x] Support `phpier logs <service>` to view logs from a specific service (app, database, valkey, etc.)
- [x] Add `-f` or `--follow` flag to follow/tail logs in real-time
- [x] Add `--tail <n>` flag to show only the last n lines of logs
- [x] Add `--since <timestamp>` flag to show logs since a specific time
- [x] Support color output for better readability (inherited from Docker Compose)
- [x] Handle cases where no containers are running gracefully
- [x] Validate that the command is run from a phpier project directory
- [x] Support both short and long flag formats for all options
- [x] Add help text and usage examples for the command

## Implementation Notes
- Use Docker Compose logs functionality under the hood
- Leverage existing project detection logic from other commands
- Follow established Cobra command patterns used in the codebase
- Ensure the command works with all supported service types (MySQL, PostgreSQL, MariaDB, Valkey, Memcached, etc.)
- Consider container naming conventions used by phpier
- Handle Docker Compose v1 vs v2 compatibility if needed
- Add proper error handling for Docker daemon connectivity issues

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete