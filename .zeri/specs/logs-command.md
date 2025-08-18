# Feature Specification: logs-command

## Overview
Add a `phpier logs` command to view Docker container logs for the current project's services. This command should provide easy access to logs from all running containers with options for filtering, following, and formatting.

## Requirements
- [ ] Implement `phpier logs` command that displays logs from all project containers
- [ ] Support `phpier logs <service>` to view logs from a specific service (app, database, valkey, etc.)
- [ ] Add `-f` or `--follow` flag to follow/tail logs in real-time
- [ ] Add `--tail <n>` flag to show only the last n lines of logs
- [ ] Add `--since <timestamp>` flag to show logs since a specific time
- [ ] Support color output for better readability
- [ ] Handle cases where no containers are running gracefully
- [ ] Validate that the command is run from a phpier project directory
- [ ] Support both short and long flag formats for all options
- [ ] Add help text and usage examples for the command

## Implementation Notes
- Use Docker Compose logs functionality under the hood
- Leverage existing project detection logic from other commands
- Follow established Cobra command patterns used in the codebase
- Ensure the command works with all supported service types (MySQL, PostgreSQL, MariaDB, Valkey, Memcached, etc.)
- Consider container naming conventions used by phpier
- Handle Docker Compose v1 vs v2 compatibility if needed
- Add proper error handling for Docker daemon connectivity issues

## TODO
- [ ] Design and plan implementation
- [ ] Implement core functionality
- [ ] Add tests
- [ ] Update documentation
- [ ] Review and refine
- [ ] Mark specification as complete