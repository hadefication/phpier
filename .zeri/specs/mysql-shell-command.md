# Feature Specification: mysql-shell-command

## Overview
Add a `phpier mysql` command that connects directly to the MySQL database container shell. This allows developers to run MySQL commands, inspect databases, and perform database operations without needing to remember Docker commands or container names.

## Requirements
- Add `phpier mysql` command that connects to the MySQL container
- Use the MySQL client (`mysql`) to connect to the database
- Automatically detect MySQL container name from docker-compose.yml
- Use database credentials from .phpier.yml configuration
- Connect to the correct database specified in project configuration
- Handle cases where MySQL container is not running
- Provide clear error messages for connection failures
- Support both interactive shell and command execution modes

## Implementation Notes
- Add new Cobra command in `cmd/` directory
- Integrate with existing container detection logic
- Read database configuration from .phpier.yml
- Use `docker exec -it` to connect to MySQL container
- Follow existing patterns from `sh` command implementation
- Container name format: `{project-name}_mysql_1` or similar
- Default database connection: `mysql -u {user} -p{password} {database}`
- Handle interactive vs non-interactive modes
- Test with different MySQL versions and configurations
- Ensure command works when MySQL is the selected database type

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete