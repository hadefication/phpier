# phpier Development Guide

This document provides comprehensive information for developers working on the phpier project.

## 🏗️ Project Architecture

### Technology Stack
- **Language**: Go 1.20+
- **CLI Framework**: Cobra + Viper
- **Template Engine**: Go text/template
- **Configuration**: YAML with Viper
- **Testing**: Go standard testing + testify
- **Containerization**: Docker & Docker Compose

### Project Structure
```
phpier/
├── main.go                     # CLI entry point
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── cmd/                        # Cobra commands
│   ├── root.go                 # Root command setup
│   ├── init.go                 # phpier init command
│   ├── up.go                   # phpier up command
│   ├── down.go                 # phpier down command
│   └── build.go                # phpier build command
├── internal/                   # Internal packages
│   ├── config/                 # Configuration management
│   │   ├── config.go           # Viper configuration setup
│   │   └── validation.go       # Config validation
│   ├── docker/                 # Docker operations
│   │   ├── client.go           # Docker client wrapper
│   │   └── compose.go          # Docker compose operations
│   ├── templates/              # Template management
│   │   ├── engine.go           # Go template engine
│   │   └── files/              # Template files (embedded)
│   │       ├── docker-compose/ # Docker compose templates
│   │       ├── dockerfiles/    # Dockerfile templates
│   │       └── configs/        # Configuration templates
│   └── services/               # Service management
├── configs/                    # Default configuration files
│   ├── defaults.yaml           # Default settings
│   └── php-versions.yaml       # PHP version definitions
├── tests/                      # Test files
│   ├── integration/            # Integration tests
│   └── unit/                   # Unit tests
├── .zeri/                      # Zeri project management
│   ├── project.md              # Project overview
│   ├── development.md          # Development practices
│   └── specs/                  # Feature specifications
├── README.md                   # User documentation
├── DEV.md                      # This file
└── LICENSE                     # License file
```

## 🚀 Development Setup

### Prerequisites
```bash
# Install Go 1.20+
brew install go              # macOS
apt install golang-go        # Ubuntu
# Or download from https://golang.org/dl/

# Install Docker
# Follow instructions at https://docs.docker.com/get-docker/

# Verify installations
go version                   # Should be 1.20+
docker --version
docker-compose --version
```

### Clone and Setup
```bash
# Clone the repository
git clone <repository-url>
cd phpier

# Install dependencies
go mod download

# Build the project
go build -o phpier main.go

# Run tests
go test ./...

# Run with verbose output
./phpier --verbose --help
```

### Development Workflow
```bash
# Run without building
go run main.go [command]

# Build and test locally
go build -o phpier main.go
./phpier init 8.3 --verbose

# Format code
go fmt ./...

# Run linter (install with: go install golang.org/x/lint/golint@latest)
golint ./...

# Run static analysis
go vet ./...
```

## 🧪 Testing

### Test Structure
```bash
tests/
├── unit/                       # Unit tests
│   ├── config_test.go         # Configuration tests
│   ├── docker_test.go         # Docker client tests
│   └── template_test.go       # Template engine tests
└── integration/                # Integration tests
    ├── init_test.go           # Init command tests
    ├── up_down_test.go        # Docker compose tests
    └── e2e_test.go            # End-to-end tests

Generated Project Structure:
your-project/
├── .phpier/                  # All generated files
│   ├── docker-compose.yml     # Service orchestration
│   ├── Dockerfile.php         # PHP container
│   └── docker/                # Configuration files
└── .phpier.yaml              # Project configuration
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific package
go test ./internal/config

# Run integration tests (requires Docker)
go test ./tests/integration

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Dependencies
```bash
# Install testify for assertions
go get github.com/stretchr/testify

# Install testcontainers for integration tests
go get github.com/testcontainers/testcontainers-go
```

### Writing Tests

#### Unit Test Example
```go
// internal/config/config_test.go
package config

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestValidatePHPVersion(t *testing.T) {
    tests := []struct {
        name     string
        version  string
        expected bool
    }{
        {"valid version", "8.3", true},
        {"invalid version", "9.0", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := IsValidPHPVersion(tt.version)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

#### Integration Test Example
```go
// tests/integration/init_test.go
package integration

import (
    "os"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestInitCommand(t *testing.T) {
    // Create temporary directory
    tmpDir := t.TempDir()
    os.Chdir(tmpDir)

    // Run init command
    err := runphpierCommand("init", "8.3")
    assert.NoError(t, err)

    // Verify files are created
    assert.FileExists(t, "docker-compose.yml")
    assert.FileExists(t, "Dockerfile.php")
}
```

## 🏗️ Architecture Patterns

### Command Pattern
Each CLI command is implemented as a separate Cobra command:

```go
// cmd/example.go
var exampleCmd = &cobra.Command{
    Use:   "example",
    Short: "Example command",
    RunE:  runExample,
}

func init() {
    rootCmd.AddCommand(exampleCmd)
    // Add flags here
}

func runExample(cmd *cobra.Command, args []string) error {
    // Implementation here
    return nil
}
```

### Configuration Management
Viper handles configuration with precedence:
1. Command line flags
2. Environment variables
3. Configuration files
4. Defaults

```go
// Set defaults
viper.SetDefault("php.version", "8.3")

// Bind flags
viper.BindPFlag("php.version", cmd.Flags().Lookup("php-version"))

// Read environment variables
viper.SetEnvPrefix("PHPIER")
viper.AutomaticEnv()
```

### Template Engine
Go templates with custom functions:

```go
// Template with custom function
tmpl := `FROM php:{{.Config.PHP.Version}}-fpm
{{- if serviceEnabled "redis" .Config }}
# Redis configuration
{{- end }}`

// Custom template functions
funcMap := template.FuncMap{
    "serviceEnabled": func(service string, config *Config) bool {
        // Implementation
    },
}
```

### Error Handling
Structured error handling with context:

```go
func validateConfig(config *Config) error {
    if config.PHP.Version == "" {
        return fmt.Errorf("PHP version is required")
    }
    
    if !IsValidPHPVersion(config.PHP.Version) {
        return fmt.Errorf("unsupported PHP version: %s", config.PHP.Version)
    }
    
    return nil
}
```

## 📝 Template System

### Template Structure
Templates are embedded in the binary using Go embed:

```go
//go:embed files
var templateFS embed.FS
```

### Template Syntax
Templates use Go's text/template syntax with custom functions:

```yaml
# docker-compose template
version: '3.8'
services:
  app:
    image: php:{{.Config.PHP.Version}}-fpm
    {{- if serviceEnabled "redis" .Config }}
    depends_on:
      - redis
    {{- end }}
```

### Custom Functions
Available in all templates:

| Function | Description | Example |
|----------|-------------|---------|
| `serviceEnabled` | Check if service is enabled | `{{serviceEnabled "redis" .Config}}` |
| `phpExtension` | Check if PHP extension is enabled | `{{phpExtension "gd" .Config}}` |
| `default` | Provide default value | `{{.Value \| default "fallback"}}` |
| `upper` | Convert to uppercase | `{{.Value \| upper}}` |
| `lower` | Convert to lowercase | `{{.Value \| lower}}` |

### Adding New Templates
1. Create template file in `internal/templates/files/`
2. Use `.tpl` extension
3. Follow existing naming conventions
4. Test with various configurations

## 🔧 Configuration System

### Configuration Hierarchy
1. Command line flags (highest priority)
2. Environment variables (`PHPIER_*`)
3. Project config file (`.phpier.yaml`)
4. Global config file (`~/.phpier.yaml`)
5. Built-in defaults (lowest priority)

### Configuration Schema
```go
type Config struct {
    Docker   DockerConfig   `mapstructure:"docker"`
    PHP      PHPConfig      `mapstructure:"php"`
    Services ServicesConfig `mapstructure:"services"`
    Traefik  TraefikConfig  `mapstructure:"traefik"`
}
```

### Adding New Configuration
1. Add to struct in `internal/config/config.go`
2. Add default in `setDefaults()`
3. Add validation in `validation.go`
4. Update templates as needed
5. Add tests

## 🐳 Docker Integration

### Docker Client Wrapper
The `internal/docker/client.go` provides a simplified interface:

```go
// Create client
client, err := docker.NewClient()

// Run command
err = client.RunCommand("docker", "ps")

// Get output
output, err := client.RunCommandOutput("docker", "ps", "-q")
```

### Docker Compose Manager
The `internal/docker/compose.go` handles compose operations:

```go
// Create manager
manager, err := docker.NewComposeManager(config)

// Start services
err = manager.Up(true) // detached mode

// Stop services
err = manager.Down(false) // preserve volumes
```

## 🚀 Adding New Commands

### 1. Create Command File
```go
// cmd/newcommand.go
package cmd

import (
    "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
    Use:   "new",
    Short: "Description",
    Long:  `Detailed description`,
    RunE:  runNew,
}

func init() {
    rootCmd.AddCommand(newCmd)
    
    // Add flags
    newCmd.Flags().StringVar(&flag, "flag", "default", "description")
}

func runNew(cmd *cobra.Command, args []string) error {
    // Implementation
    return nil
}
```

### 2. Add Tests
```go
// cmd/newcommand_test.go
package cmd

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
    // Test implementation
}
```

### 3. Update Documentation
- Add to README.md usage section
- Update help text
- Add examples

## 🏷️ Release Process

### Version Management
```bash
# Update version in main.go
var version = "v1.2.3"

# Tag release
git tag v1.2.3
git push origin v1.2.3
```

### Building Releases
```bash
# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o phpier-linux-amd64 main.go
GOOS=darwin GOARCH=amd64 go build -o phpier-darwin-amd64 main.go
GOOS=windows GOARCH=amd64 go build -o phpier-windows-amd64.exe main.go

# Or use build script
./scripts/build-release.sh v1.2.3
```

### Release Checklist
- [ ] Update version number
- [ ] Update CHANGELOG.md
- [ ] Run full test suite
- [ ] Build cross-platform binaries
- [ ] Create GitHub release
- [ ] Update documentation

## 🔍 Debugging

### Debug Mode
```bash
# Enable verbose logging
phpier --verbose init 8.3

# Debug template rendering
PHPIER_DEBUG_TEMPLATES=true phpier init 8.3

# Debug Docker commands
PHPIER_DEBUG_DOCKER=true phpier up
```

### Logging Levels
```go
// Set log level
logrus.SetLevel(logrus.DebugLevel)

// Log messages
logrus.Debug("Debug message")
logrus.Info("Info message")
logrus.Warn("Warning message")
logrus.Error("Error message")
```

### Common Debug Scenarios
```bash
# Template rendering issues
go run main.go init 8.3 --verbose 2>&1 | grep -i template

# Docker command failures
docker-compose config

# Configuration issues
phpier --verbose --config=/dev/null init 8.3
```

## 🤝 Contributing

### Code Style
- Follow Go conventions (`gofmt`, `golint`, `go vet`)
- Use meaningful variable names
- Write tests for new functionality
- Document public functions
- Keep functions small and focused

### Commit Messages
```
type(scope): description

- feat: new feature
- fix: bug fix
- docs: documentation
- style: formatting
- refactor: code restructuring
- test: adding tests
- chore: maintenance

Examples:
feat(init): add PostgreSQL support
fix(docker): handle missing docker-compose
docs(readme): update installation instructions
```

### Pull Request Process
1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Make changes and add tests
4. Ensure all tests pass (`go test ./...`)
5. Update documentation
6. Commit changes (`git commit -m 'feat: add amazing feature'`)
7. Push to branch (`git push origin feature/amazing-feature`)
8. Open Pull Request

### Code Review Guidelines
- Code should be self-documenting
- Tests should cover new functionality
- Breaking changes require major version bump
- Documentation should be updated
- Performance implications should be considered

## 📊 Project Metrics

### Code Quality
- Test coverage should be > 80%
- No critical security vulnerabilities
- All linting checks should pass
- Documentation should be up to date

### Performance Targets
- CLI startup time < 500ms
- Init command < 10s
- Template rendering < 1s
- Docker operations depend on Docker performance

## 🔗 Useful Resources

### Go Development
- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Cobra CLI](https://github.com/spf13/cobra)
- [Viper Configuration](https://github.com/spf13/viper)

### Docker
- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Docker SDK for Go](https://docs.docker.com/engine/api/sdk/)

### Testing
- [Go Testing](https://golang.org/pkg/testing/)
- [Testify](https://github.com/stretchr/testify)
- [Testcontainers](https://github.com/testcontainers/testcontainers-go)

---

**Happy contributing to phpier! 🛠️**