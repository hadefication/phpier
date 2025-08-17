# Feature Specification: global-database-service-management

## Overview
Improve global database service management by allowing users to enable/disable specific database services (MySQL, PostgreSQL, MariaDB) rather than being limited to a single database type. This provides better flexibility for developers who need multiple database environments or want to reduce resource usage by only running required services.

## Requirements
- **Multiple Database Support**: Allow multiple database services to run simultaneously
- **Service Management Commands**: Add commands to enable/disable specific database services
  - `phpier db enable mysql|postgresql|mariadb`
  - `phpier db disable mysql|postgresql|mariadb`
  - `phpier db list` - Show enabled/disabled status of all databases
  - `phpier db status` - Show running status of enabled databases with ports
  - `phpier db credentials` - Show database credentials for enabled services
- **Default Configuration**: MySQL enabled by default, PostgreSQL and MariaDB disabled
- **Port Management**: Automatic port assignment to avoid conflicts
  - MySQL: 3306
  - PostgreSQL: 5432
  - MariaDB: 3307 (to avoid conflict with MySQL)
- **Backwards Compatibility**: Maintain existing single database type configuration for projects
- **Configuration Persistence**: Store enabled database services in global config
- **Resource Optimization**: Only start containers for enabled services
- **Dependency Management**: Handle service dependencies (e.g., PHPMyAdmin for MySQL/MariaDB)

## Implementation Notes
- **Configuration Schema Updates**:
  - Extend `GlobalConfig` struct to support multiple enabled databases
  - Add `DatabaseServices` field with individual service configurations
  - Maintain backwards compatibility with existing `services.database.type` field
- **Command Structure**:
  - Add new `phpier db` command group at the root level
  - Cleaner, more intuitive command structure for database management
- **Docker Compose Template Changes**:
  - Modify `global.yml.tpl` to conditionally include database services
  - Update template logic to handle multiple databases simultaneously
  - Adjust port mappings to prevent conflicts
- **Service Discovery Updates**:
  - Update database shell commands to detect which services are enabled
  - Provide clear error messages when attempting to connect to disabled services
- **Migration Strategy**:
  - Automatically migrate existing configurations to new format
  - Default to current behavior (single database) with migration path to multi-database
- **Files to Modify**:
  - `internal/config/config.go` - Configuration schema
  - `internal/templates/files/docker-compose/global.yml.tpl` - Template updates
  - `cmd/db.go` - Add database management commands
  - Database shell commands - Update service detection logic
- **Status Display Format**:
  - Show service name, enabled/disabled status, running status, and port
  - Example output:
    ```
    Database Services Status:
    ✓ MySQL      [enabled]  [running]   localhost:3306
    ✗ PostgreSQL [disabled] [stopped]   localhost:5432
    ✓ MariaDB    [enabled]  [stopped]   localhost:3307
    ```
- **Credentials Display Format**:
  - Show database connection information for enabled services only
  - Include host, port, username, and password
  - Example output:
    ```
    Database Credentials (enabled services only):
    
    MySQL:
      Host:     localhost:3306
      Username: phpier
      Password: phpier
      
    MariaDB:
      Host:     localhost:3307
      Username: phpier
      Password: phpier
    ```
- **Testing Requirements**:
  - Test enable/disable functionality
  - Test multiple database services running simultaneously
  - Test port conflict resolution
  - Test status command output with port information
  - Test credentials command output with connection details
  - Test backwards compatibility with existing projects
  - Test migration from single to multi-database configuration

## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [x] Add tests
- [x] Update documentation
- [x] Review and refine
- [x] Mark specification as complete