# Feature Specification: mariadb-shell-command

## Overview
Add a `phpier maria` command that connects directly to the MariaDB database container shell. This allows developers to run MariaDB commands, inspect databases, and perform database operations without needing to remember Docker commands or container names.

## Requirements
- Add `phpier maria` command that connects to the MariaDB container
- Use the MariaDB client (`mariadb` or `mysql`) to connect to the database
- Automatically detect MariaDB container name from docker-compose.yml
- Use database credentials from .phpier.yml configuration
- Connect to the correct database specified in project configuration
- Handle cases where MariaDB container is not running
- Provide clear error messages for connection failures
- Support both interactive shell and command execution modes
- Also support `phpier mariadb` as a full name alias

## Implementation Notes
- Add new Cobra command in `cmd/` directory
- Integrate with existing container detection logic
- Read database configuration from .phpier.yml
- Use `docker exec -it` to connect to MariaDB container
- Follow existing patterns from `sh` command implementation
- Container name format: `{project-name}_mariadb_1` or similar
- Default database connection: `mariadb -u {user} -p{password} {database}` or `mysql -u {user} -p{password} {database}`
- Handle interactive vs non-interactive modes
- Test with different MariaDB versions and configurations
- Ensure command works when MariaDB is the selected database type
- Add `mariadb` alias command for full name reference
- MariaDB containers typically have both `mariadb` and `mysql` clients available

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete