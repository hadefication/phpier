# Feature Specification: core-features

## Overview
Core CLI commands and Docker environment management for phpier - a PHP development environment management tool. This specification covers the essential commands needed to initialize, manage, and interact with containerized PHP development environments.

## Requirements

### CLI Commands
- **phpier init <version>**: Initialize a phpier environment with specified PHP version
  - Support PHP versions: 5.6, 7.3, 7.4, 8.0, 8.1, 8.2, 8.3, 8.4
  - Generate docker-compose.yml with selected services
  - Create Traefik configuration for folder-based routing
  - Set up directory structure and configuration files

- **phpier up**: Start the Docker environment
  - Run docker-compose up with proper configuration
  - Start Traefik reverse proxy
  - Display connection information and URLs

- **phpier down**: Stop the Docker environment
  - Run docker-compose down
  - Clean up containers and networks
  - Preserve volumes and data

- **phpier build**: Build or rebuild services
  - Force rebuild of Docker images
  - Update container configurations
  - Handle dependency updates

- **phpier php**: Proxy PHP commands to app container
  - Execute PHP commands inside the app container
  - Pass through all arguments and options
  - Maintain proper exit codes

- **phpier <tool>**: Proxy app container tools
  - Support common tools: composer, npm, node, yarn, artisan
  - Execute commands within the app container context
  - Preserve working directory and file permissions

### Container Services
- **App Container**: PHP runtime with development tools
  - PHP-FPM with selected version
  - Nginx web server
  - Composer, NVM pre-installed
  - Common PHP extensions
  - Development utilities

- **Database Options**: Configurable database services
  - MySQL with PHPMyAdmin
  - PostgreSQL with pgAdmin
  - MariaDB with PHPMyAdmin

- **Caching Services**: Optional caching layers
  - Valkey (Redis) for caching and sessions
  - Memcached for object caching

- **Development Tools**
  - Mailpit for email testing
  - Traefik for reverse proxy and SSL

### Domain Routing
- Folder-based domain routing: `<directory>.localhost`
- Automatic SSL certificate generation
- Service discovery through Traefik

## Implementation Notes

### Technical Dependencies
- Docker and Docker Compose must be installed
- Go 1.21+ for development and compilation
- Cobra + Viper frameworks for CLI functionality
- Go's text/template package for file generation

### Technology Stack: Go + Cobra + Viper
**Selected Solution Benefits:**
- **Cobra**: Industry standard for Go CLI applications (used by kubectl, GitHub CLI, Hugo)
- **Viper**: Powerful configuration management (config files, env vars, flags with precedence)
- **Go Templates**: Built-in template engine for complex file generation
- **Single Binary**: Cross-platform compilation to single executable
- **Docker SDK**: Native Go Docker client for container management
- **Rich Ecosystem**: Mature testing, logging, and validation libraries

**Go Libraries:**
- `github.com/spf13/cobra` - CLI command structure
- `github.com/spf13/viper` - Configuration management
- `github.com/docker/docker` - Docker SDK for Go
- `github.com/stretchr/testify` - Enhanced testing capabilities
- `github.com/sirupsen/logrus` - Structured logging

### Go Project Structure
```
phpier/                    # Go project root
├── main.go                  # Main CLI entry point
├── go.mod                   # Go module definition
├── go.sum                   # Go module checksums
├── cmd/                     # Cobra commands
│   ├── root.go              # Root command setup
│   ├── init.go              # phpier init command
│   ├── up.go                # phpier up command
│   ├── down.go              # phpier down command
│   ├── build.go             # phpier build command
│   └── proxy.go             # phpier php/tool proxy commands
├── internal/                # Internal packages
│   ├── config/              # Configuration management
│   │   ├── config.go        # Viper configuration setup
│   │   └── validation.go    # Config validation
│   ├── docker/              # Docker operations
│   │   ├── client.go        # Docker client wrapper
│   │   ├── compose.go       # Docker compose operations
│   │   └── proxy.go         # Container command proxying
│   ├── templates/           # Template management
│   │   ├── engine.go        # Go template engine
│   │   ├── renderer.go      # Template rendering logic
│   │   └── validator.go     # Template validation
│   └── services/            # Service management
│       ├── php.go           # PHP version management
│       ├── database.go      # Database service setup
│       └── traefik.go       # Traefik configuration
├── templates/               # Template files (embedded)
│   ├── docker-compose/      # Docker compose templates
│   ├── dockerfiles/         # Dockerfile templates
│   └── configs/             # Configuration templates
├── configs/                 # Default configuration files
│   ├── php-versions.yaml    # PHP version definitions
│   ├── services.yaml        # Service configurations
│   └── defaults.yaml        # Default settings
└── tests/                   # Test files
    ├── integration/         # Integration tests
    └── unit/                # Unit tests
```

### Go Command Architecture
- **Cobra Command Structure**: Each command as separate Go file in cmd/
- **Dependency Injection**: Interfaces for Docker, templates, and config
- **Viper Configuration**: Multi-source config with precedence (flags > env > files)
- **Template Embedding**: Go embed directive for bundling templates in binary
- **Error Handling**: Structured error types with user-friendly messages
- **Logging**: Structured logging with configurable levels

### Error Handling
- Validate Docker installation and availability
- Check for port conflicts before starting services
- Provide clear error messages for common issues
- Graceful handling of missing dependencies

### Testing Strategy
- **Unit Tests**: Go standard testing with testify assertions
- **Integration Tests**: testcontainers for real Docker testing
- **Command Tests**: Cobra command testing with mock dependencies
- **Table-Driven Tests**: Multiple scenarios and edge cases
- **Cross-Platform**: Go build constraints for platform-specific code
- **CI/CD**: GitHub Actions with matrix builds for multiple platforms

## TODO
- [x] Initialize Go module and project structure
- [x] Set up Cobra CLI with root command and basic structure
- [x] Implement Viper configuration management system
- [x] Create Go template engine for Dockerfile and docker-compose generation
- [x] Implement `phpier init` command with PHP version selection
- [x] Implement `phpier up/down` commands with Docker SDK integration
- [x] Implement `phpier build` command for container management
- [x] Add proper Composer and Node.js version compatibility across PHP versions
- [ ] Implement `phpier php` command proxy functionality using Docker exec
- [ ] Implement generic tool proxy for container commands (composer, npm, etc.)
- [ ] Create Traefik configuration templates and setup logic
- [ ] Add container service templates (database, caching, tools) with Go templates
- [ ] Implement comprehensive error handling with structured error types
- [ ] Add unit tests with testify and integration tests with testcontainers
- [ ] Set up cross-platform builds and GitHub Actions CI/CD
- [ ] Create installation scripts and documentation
- [ ] Review and refine implementation
- [ ] Mark specification as complete