# Feature Specification: services-command

## Overview
Add a `phpier services` command that displays comprehensive information about all running phpier services. This command provides developers with a quick overview of their phpier environment status, including container health, ports, volumes, and service dependencies.

## Requirements

### Core Functionality
- Display all running phpier containers with their status (running/stopped/exited)
- Show container names, images, and creation timestamps
- Display exposed ports and their mappings (container:host)
- Show mounted volumes and their bind paths
- Indicate service health status and uptime
- Display resource usage (CPU, memory) when available
- Show database connection status and service URLs

### User Interface Requirements
- Clean, tabular output with clear headers
- Color-coded status indicators (green=running, red=stopped, yellow=starting)
- Option for JSON output (`--json` flag) for scripting
- Option for verbose output (`--verbose` flag) with additional details
- Filter by service type (`--type app|db|cache|proxy`)
- Support for `--project` flag to show specific project services

### Service Information Display
- **Global Services**: Traefik, phpier-network status
- **Project Services**: App container, database, cache, development tools
- **Service URLs**: Show accessible URLs (traefik routes, adminer, mailpit)
- **Database Info**: Connection details, database names, user info
- **Cache Services**: Redis/Valkey, Memcached status and connections
- **Development Tools**: PHPMyAdmin, Mailpit, status and URLs

### Integration Requirements
- Compatible with existing Docker Compose configurations
- Work with both global and project-specific services
- Support multi-project environments
- Integrate with Traefik service discovery

## Implementation Notes

### Technical Considerations
- Use Docker SDK for Go to query container information
- Parse docker-compose.yml files to understand service relationships
- Query Traefik API for routing information and service discovery
- Handle cases where services are partially running or in error states

### Files to Modify/Create
- `cmd/services.go` - New Cobra command implementation
- `internal/docker/services.go` - Service discovery and status functions
- `internal/config/project.go` - Project service mapping and configuration
- `internal/display/table.go` - Table formatting utilities
- Add services command to root command in `cmd/root.go`

### Service Discovery Strategy
- Scan for containers with phpier labels or naming conventions
- Parse `.phpier.yml` configuration for service definitions
- Query Docker daemon for container status and metadata
- Cross-reference with Traefik service discovery

### Output Format Design
```
PHPIER SERVICES STATUS

Global Services:
NAME        STATUS   UPTIME    PORTS               IMAGE
traefik     running  2h 15m    80:80, 443:443     traefik:v2.10
network     active   2h 15m    -                  -

Project Services (project-name):
NAME            STATUS   UPTIME    PORTS      URLS                           IMAGE
app             running  1h 30m    -          http://project.localhost       phpier/php:8.3
mysql           running  1h 30m    3306       -                              mysql:8.0
redis           running  1h 30m    6379       -                              redis:7
phpmyadmin      running  1h 30m    -          http://pma.project.localhost   phpmyadmin:latest
mailpit         running  1h 30m    -          http://mail.project.localhost  axllent/mailpit
```

### Error Handling
- Handle Docker daemon connection failures gracefully
- Provide helpful messages when no phpier services are running
- Handle permission issues accessing Docker socket
- Warn about orphaned containers or configuration mismatches

### Testing Strategy
- Unit tests for service discovery functions
- Integration tests with running Docker containers
- Test output formatting with mock data
- Test filtering and flag combinations
- Test JSON output format validation

## TODO
- [x] Design and plan command structure and flags
- [x] Implement Docker service discovery functionality
- [x] Create service status and information gathering
- [x] Implement table formatting and display logic
- [x] Add JSON output format support
- [x] Implement filtering and project-specific views
- [x] Create Cobra command integration
- [x] Add comprehensive unit tests
- [ ] Add integration tests with Docker environment
- [x] Update documentation and help text
- [x] Review and refine implementation
- [x] Mark specification as complete