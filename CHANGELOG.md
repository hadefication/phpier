# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial phpier CLI implementation
- Multi-PHP version support (5.6, 7.3, 7.4, 8.0, 8.1, 8.2, 8.3, 8.4)
- Docker Compose orchestration
- Traefik reverse proxy integration
- Database support (MySQL, PostgreSQL, MariaDB)
- Caching services (Redis, Memcached)
- Development tools (PHPMyAdmin, pgAdmin, Mailpit)
- Template-based file generation
- Configuration management with Viper
- Version compatibility matrix for Composer and Node.js
- Clean `.phpier` directory structure
- Cross-platform build system

### Features
- `phpier init` - Initialize PHP development environment
- `phpier up` - Start Docker services
- `phpier down` - Stop Docker services  
- `phpier build` - Build/rebuild Docker images
- `phpier version` - Show version information
- Automatic domain routing with `<project>.localhost`
- Customizable Docker and configuration files
- Comprehensive CLI help and validation

### Technical
- Go 1.20+ with Cobra + Viper frameworks
- Embedded template system with Go templates
- Docker SDK integration for container management
- Multi-platform release builds (Linux, macOS, Windows)
- Comprehensive testing framework setup
- Development tooling (Makefile, build scripts)

## [1.0.0] - TBD

### Added
- Initial release of phpier
- Core CLI functionality
- Multi-platform binaries
- Complete documentation

<!-- Template for future releases:

## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security improvements

-->