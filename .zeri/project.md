# phpier - Project Context

## Overview
A CLI tool to manage PHP development using Docker. Supports multiple PHP versions (5.6, 7.2, 7.3, 7.4, 8.0, 8.1, 8.2, 8.3, 8.4) with Traefik for folder-based domain routing (directory.localhost).

## Tech Stack
- **CLI Framework**: Go with Cobra + Viper (industry standard for CLI tools)
- **Language**: Go 1.21+
- **Configuration**: Viper for multi-source config management
- **Templating**: Go's built-in text/template and html/template
- **Containerization**: Docker & Docker Compose
- **Reverse Proxy**: Traefik
- **Web Server**: Nginx (in app container)
- **PHP Runtime**: PHP-FPM with multiple version support
- **Package Management**: Composer, NVM

## Architecture
Containerized development environment with:
- App container (PHP + tools)
- Database containers (MySQL/PostgreSQL/MariaDB)
- Caching services (Valkey/Redis, Memcached)
- Development tools (PHPMyAdmin, Mailpit)

## Key Components

### CLI Commands
- `phpier init <version>` - Initialize phpier environment with PHP version
- `phpier up` - Start docker-compose services
- `phpier down` - Stop docker-compose services
- `phpier build` - Build/rebuild services
- `phpier php` - Proxy PHP commands to app container
- `phpier <tool>` - Proxy app container tools

### Containers
- **App Container**: Pre-installed with PHP, PHP-FPM, extensions, Composer, NVM, Nginx
- **Database**: Configurable (MySQL, PostgreSQL, MariaDB)
- **Valkey/Redis**: Caching and session storage
- **Memcached**: Memory caching
- **Mailpit**: Email testing

### Features
- Multiple PHP version support
- Traefik-based routing with `.localhost` domains
- PHPMyAdmin integration (when MySQL selected)
- Development tool proxying through CLI

## Current Focus
CLI tool development and Docker environment setup

## Environment Setup
Docker-based development with automatic domain routing via Traefik

## Important Notes
- Supports PHP 5.6, 7.2, 7.3, 7.4, 8.0, 8.1, 8.2, 8.3, 8.4
- Uses folder-based domain naming convention
- All development tools accessible through CLI proxy commands

