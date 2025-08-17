# Feature Specification: php-7.2-support

## Overview
Add support for PHP 7.2 to the phpier CLI tool. Currently phpier supports PHP versions 5.6, 7.3, 7.4, 8.0, 8.1, 8.2, 8.3, and 8.4. This enhancement will add PHP 7.2 to fill the gap between 5.6 and 7.3, providing complete coverage for legacy applications that require this specific version.

## Requirements
- Add PHP 7.2 as a valid version option for `phpier init` command
- Create appropriate Dockerfile template for PHP 7.2 environment
- Ensure all standard PHP extensions are available and properly configured
- Include development tools (Composer, NVM) in PHP 7.2 container
- Maintain compatibility with existing CLI commands and Docker infrastructure
- Support all database options (MySQL, PostgreSQL, MariaDB) with PHP 7.2
- Include PHPMyAdmin compatibility when MySQL is selected
- Ensure Traefik routing works correctly with PHP 7.2 containers

## Implementation Notes
- PHP 7.2 should use the same Dockerfile template structure as other PHP 7.x versions
- Check which template file PHP 7.2 should use based on existing patterns:
  - `php56-73.Dockerfile.tpl` (if 7.2 follows 5.6-7.3 pattern)
  - `php74-80.Dockerfile.tpl` (if 7.2 follows 7.4+ pattern)
  - May need to adjust template groupings or create new template
- Update version validation logic in CLI commands
- Add PHP 7.2 to supported versions list in project documentation
- Test with multiple database configurations
- Verify extension compatibility (some extensions may have different requirements)
- Test Composer and NVM functionality within PHP 7.2 container
- Ensure backwards compatibility with existing projects

## TODO
- [ ] Analyze existing PHP version template structure and determine correct template for 7.2
- [ ] Update CLI version validation to include PHP 7.2
- [ ] Add or modify Dockerfile template for PHP 7.2 support
- [ ] Test initialization with PHP 7.2 across all database options
- [ ] Verify container startup and service connectivity
- [ ] Test CLI command proxying with PHP 7.2 container
- [ ] Update project documentation to reflect PHP 7.2 support
- [ ] Clean up test artifacts and Docker resources
- [ ] Mark specification as complete