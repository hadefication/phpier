# Feature Specification: Database Authentication Fixes

## Overview
Fix authentication issues for MariaDB and PostgreSQL similar to the MySQL fix, and update documentation with proper connection credentials and instructions for external database clients.

## Requirements
- Fix MariaDB authentication plugin to use mysql_native_password for client compatibility
- Ensure PostgreSQL authentication works correctly with external clients
- Update both project (base.yml.tpl) and global (global.yml.tpl) Docker templates
- Update documentation with correct credentials for all databases (MySQL, MariaDB, PostgreSQL)
- Test external client connections (Sequel Ace, TablePlus, pgAdmin, etc.)
- Document web interface access (Adminer, PHPMyAdmin, pgAdmin)

## Implementation Notes
- MariaDB should use mysql_native_password like MySQL for client compatibility
- PostgreSQL authentication is typically correct by default but verify trust settings
- Update both base.yml.tpl and global.yml.tpl templates consistently
- Document default credentials for each database type
- Include connection instructions for popular database clients
- Test with fresh database volumes to ensure changes take effect

## TODO
- [x] Fix MariaDB authentication in both Docker templates
- [x] Verify PostgreSQL authentication configuration
- [x] Test database connections with external clients
- [x] Update CONFIGURATION.md with database credentials and access methods
- [x] Test web interface access for all database types
- [x] Clean up any test artifacts created during testing
- [x] Mark specification as complete