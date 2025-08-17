# Feature Specification: postgresql-shell-command

## Overview
Add a `phpier psql` command that connects directly to the PostgreSQL database container shell. This allows developers to run PostgreSQL commands, inspect databases, and perform database operations without needing to remember Docker commands or container names.

## Requirements
- Add `phpier psql` command that connects to the PostgreSQL container
- Use the PostgreSQL client (`psql`) to connect to the database
- Automatically detect PostgreSQL container name from docker-compose.yml
- Use database credentials from .phpier.yml configuration
- Connect to the correct database specified in project configuration
- Handle cases where PostgreSQL container is not running
- Provide clear error messages for connection failures
- Support both interactive shell and command execution modes
- Also support `phpier postgres` as an alias for convenience

## Implementation Notes
- Add new Cobra command in `cmd/` directory
- Integrate with existing container detection logic
- Read database configuration from .phpier.yml
- Use `docker exec -it` to connect to PostgreSQL container
- Follow existing patterns from `sh` command implementation
- Container name format: `{project-name}_postgres_1` or similar
- Default database connection: `psql -U {user} -d {database}`
- Handle PGPASSWORD environment variable for password authentication
- Handle interactive vs non-interactive modes
- Test with different PostgreSQL versions and configurations
- Ensure command works when PostgreSQL is the selected database type
- Add `postgres` alias command for convenience

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete