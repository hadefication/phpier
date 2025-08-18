# phpier - Development Practices

## Code Standards & Quality

### Code Style
- Follow Go standard formatting (gofmt, go vet, golint)
- **⚠️ CRITICAL: Always run `go vet` and `go fmt` for every specification implementation or file changes during non-spec implementations**
- **⚠️ CRITICAL: Code MUST pass both `go fmt` and `go vet` before any commit or PR**
- Use Cobra command structure and naming conventions
- Implement proper error handling with Go error types
- Follow Go naming conventions (CamelCase for exported, camelCase for internal)
- Use Go modules for dependency management

### Naming Conventions
Consistent naming conventions based on language and framework

### File Organization
Organize by feature/domain

### Documentation Standards
Document all public APIs and complex business logic

### Security Guidelines
Sanitize all inputs, validate data, follow security best practices

### Performance Considerations
Optimize critical paths, cache expensive operations

---

## Architecture Decisions

### Decision Template
- **Date**: 
- **Decision**: 
- **Context**: 
- **Options Considered**: 
- **Chosen Option**: 
- **Rationale**: 
- **Consequences**: 

### Recent Decisions
- **2025-08-15**: Selected Go with Cobra + Viper framework for CLI development
- **2025-08-15**: Chosen for complex template management, single binary distribution, and mature ecosystem
- **2025-08-15**: Technology stack finalized - Go 1.21+ with industry-standard CLI patterns
- **2025-08-16**: Added service management commands (start, stop, down) for better UX
- **2025-08-16**: Implemented safety checks and user warnings for global service operations

### Key Architecture Decisions
Framework choice, database selection, deployment strategy

### Technology Choices
- **CLI Framework**: Cobra + Viper (Go ecosystem standard for CLI applications)
- **Language**: Go 1.21+ for performance, single binary distribution, and excellent tooling
- **Configuration Management**: Viper for config files, environment variables, and flags
- **Template Engine**: Go's built-in text/template for Dockerfile and config generation
- **Dependencies**: Docker, Docker Compose (single Go binary distribution)
- **Containerization**: Docker & Docker Compose for environment isolation
- **Reverse Proxy**: Traefik for automatic service discovery and routing
- **Target PHP Versions**: Support for PHP 5.6, 7.3, 7.4, 8.0, 8.1, 8.2, 8.3, 8.4
- **Database Options**: MySQL, PostgreSQL, MariaDB
- **Caching**: Valkey (Redis), Memcached

### Design Patterns
Appropriate patterns for the architecture and domain

---

## Code Patterns

### Standard Patterns
Consistent architectural patterns

### Component Patterns
Reusable components, consistent API structure

### Data Handling Patterns
Data models, validation, serialization

### Error Handling Patterns
Custom exceptions, error logging, user-friendly messages

### Testing Patterns
- Go standard testing framework with testify for assertions
- Unit tests for individual functions and command logic
- Integration tests using testcontainers for Docker operations
- Mock Docker commands using interfaces and dependency injection
- Table-driven tests for multiple scenarios
- Cobra command testing with custom test helpers

### Configuration Patterns
Environment-based config, feature flags

### Examples
- Cobra command structure for phpier operations
- Viper configuration management with precedence (flags > env > config files)
- Go template rendering for Dockerfile and docker-compose generation
- Docker SDK for Go for container management
- Structured logging with logrus or zap
- Configuration validation with Go structs and tags
- Service discovery patterns with Traefik

---

## Development Workflows

### Development Process
Feature branch workflow with code review

### Before Starting Development
Check latest main branch, create feature branch

### Implementation Steps
1. Write tests 2. Implement feature 3. Run tests 4. **⚠️ CRITICAL: Run `go fmt ./...` and `go vet ./...`** 5. Code review

### Testing Workflow
Unit tests, integration tests, manual testing

#### Testing Requirements
Write tests for all new features

### Code Review Process
Pull request review with at least one approval

#### Code Review Guidelines
All code must be reviewed before merge

### Deployment Steps
Deploy to staging, test, deploy to production

### Troubleshooting Common Issues
Check logs, reproduce issue, write failing test, fix, verify

### Docker Build Issues

#### PHP Extension Compilation Failures

**Issue**: PHP extensions failing to compile due to missing system dependencies
**Common Symptoms**:
- `Package 'libcurl', required by 'virtual:world', not found` for curl extension
- `configure: error: Package requirements not met` errors
- Missing development headers for libraries

**Solution Process**:
1. **Identify the missing dependency**: Check Docker build logs for specific package requirements
2. **Add system dependencies**: Update the Dockerfile template to include required dev packages
   ```dockerfile
   # Example: Adding curl development library
   RUN apt-get update && apt-get install -y \
       libcurl4-openssl-dev \
       # ... other dependencies
   ```
3. **Rebuild and test**: Regenerate project files and rebuild Docker containers

**Common Missing Dependencies**:
- `libcurl4-openssl-dev` - Required for curl PHP extension
- `libmagickwand-dev` - Required for imagick PHP extension  
- `libpq-dev` - Required for pgsql PHP extension
- `libicu-dev` - Required for intl PHP extension
- `libzip-dev` - Required for zip PHP extension

**Testing Docker Build**:
```bash
# Navigate to test project
cd /private/tmp/phpier-test2

# Test Docker build manually
docker build -f .phpier/Dockerfile.php -t test-php-build . 2>&1

# Check if build succeeds and extensions compile
docker run --rm test-php-build php -m | grep -E "(curl|gd|mysqli)"
```

**Template Location**: `/Users/glenbangkila/AI/phpier/internal/templates/files/dockerfiles/php.Dockerfile.tpl`

### Docker Development Workflow
1. Test changes in isolated containers
2. Verify multi-PHP version compatibility
3. Test service interconnectivity
4. Validate Traefik routing configuration
5. Test CLI command proxying

**⚠️ CRITICAL: Dockerfile Template Updates**
- **ALWAYS update ALL Dockerfile templates** when making changes that affect container behavior
- There are multiple templates for different PHP versions: `php56-73.Dockerfile.tpl`, `php74-80.Dockerfile.tpl`, `php81-84.Dockerfile.tpl`, and `php.Dockerfile.tpl`
- Changes to one template must be applied to all relevant templates to maintain consistency
- Test with multiple PHP versions to ensure changes work across all supported versions

### ⚠️ CRITICAL: Testing Environment Rules

**NEVER EVER add test instances in the main codebase directory:**

- **FORBIDDEN**: Running `phpier init` in the main codebase directory
- **FORBIDDEN**: Creating `.phpier/`, `.phpier.yml`, or `public/` directories in the main codebase
- **FORBIDDEN**: Any test containers, configurations, or initialization files in the project root

**Correct Testing Locations:**
- Use `/tmp/` or `/private/tmp/` for temporary test projects
- Create dedicated test directories outside the main codebase
- Use temporary directories that can be safely deleted after testing

**Why This Matters:**
- Keeps the main codebase clean and focused on source code
- Prevents accidental commits of test artifacts
- Avoids conflicts with the actual phpier binary being developed
- Maintains professional repository structure

**Example Correct Testing:**
```bash
# ✅ CORRECT: Test in temporary directory with phpier-test- prefix
cd /tmp && mkdir phpier-test-start-command && cd phpier-test-start-command
/path/to/phpier init 8.3 --project-name=phpier-test-start

# ❌ WRONG: Never test in main codebase
cd /Users/.../phpier-codebase
./phpier init 8.3  # DON'T DO THIS!
```

**Test Naming Convention:**
- **MANDATORY**: Use `phpier-test-<test_name>` prefix for all test projects
- **Directory names**: `/tmp/phpier-test-<feature>` or `/private/tmp/phpier-test-<feature>`
- **Project names**: `phpier-test-<feature>` when using `--project-name` flag
- **Why**: Enables safe, targeted cleanup without affecting other Docker resources

### ⚠️ CRITICAL: Docker Testing Cleanup Rules

**MANDATORY CLEANUP PROCEDURES:**

Docker testing generates significant artifacts that consume disk space. Follow these cleanup procedures:

**Track All Generated Artifacts:**
- Log all test project directories created during testing
- Note Docker images built during feature implementation
- Track volumes, networks, and containers created
- Document temporary directories and their contents

**Mandatory Cleanup After Specification Implementation:**
```bash
# Clean up phpier test directories (use specific prefix)
rm -rf /tmp/phpier-test-*
rm -rf /private/tmp/phpier-test-*

# Clean up ONLY phpier test Docker resources (safer approach)
# Containers with phpier-test- prefix
docker ps -a --filter "name=phpier-test-" -q | xargs docker rm -f

# Images with phpier-test- prefix
docker images --filter "reference=*phpier-test-*" -q | xargs docker rmi -f

# Volumes with phpier-test- prefix
docker volume ls --filter "name=phpier-test-" -q | xargs docker volume rm -f

# Networks with phpier-test- prefix (if any)
docker network ls --filter "name=phpier-test-" -q | xargs docker network rm
```

**Why This Matters:**
- Docker images can consume several GB per test
- Multiple PHP version builds compound storage usage
- Accumulated artifacts slow down development environment
- Professional development practices require resource management

**Implementation Requirements:**
- **MANDATORY**: Clean up all Docker artifacts after completing each specification
- **MANDATORY**: Log all test directories and Docker resources created during development
- **MANDATORY**: Verify cleanup completion before marking specifications as complete
- Include cleanup verification in specification TODO checklists

### PHP Environment Testing
- Test across all supported PHP versions (5.6-8.4)
- Verify extension availability in each version
- Test database connectivity for each option
- Validate tool availability (Composer, NVM)

### Go CLI Development
- Use `phpier <command>` Cobra command structure
- Test Go functions with standard testing framework and testify
- Mock Docker operations using interfaces and dependency injection
- Test Cobra command parsing and validation
- Validate error handling with proper Go error types
- Cross-compile for multiple platforms (Linux/macOS/Windows)
- Use Go modules for reproducible builds

### CLI Command Implementation Patterns

#### Service Management Commands
- **phpier start**: Start global services (alternative to `phpier global up`)
- **phpier stop**: Stop global services with project safety checks
- **phpier down**: Stop project services with optional global service stopping

#### Command Safety Features
- Project detection before stopping global services
- User warnings and abort protection for dangerous operations
- Force flags to override safety checks when needed
- Clear status messages and next-step guidance

#### Error Handling Extensions
- Added `ErrorTypeUserAborted` for user-cancelled operations
- Comprehensive safety checks with meaningful error messages
- Graceful degradation when services are already stopped/started

---

## Feature Planning

### Planning Process
Requirements gathering, technical design, estimation

### Requirements Gathering
Stakeholder interviews, user stories, acceptance criteria

### Technical Analysis
Architecture review, dependency analysis, risk assessment

### Design Considerations
User experience, performance, security, maintainability

### Implementation Planning
Break down into tasks, estimate effort, plan sprints

### Risk Assessment
Identify technical risks, mitigation strategies

### Timeline Estimation
Story points, velocity tracking, buffer for unknowns

---

## Debugging & Maintenance

### Debugging Process
Reproduce, isolate, identify root cause, fix, verify

### Common Issues
- Docker container startup failures
- PHP version conflicts
- Database connection issues
- Traefik routing configuration
- Port conflicts
- File permission issues in containers
- PHP extension missing errors

### Debugging Tools
Debugger, logging, profiler, monitoring tools

### Log Analysis
Check application logs, error logs, system logs

### Performance Debugging
Profiling, query analysis, resource monitoring

### Error Tracking
Use error tracking service, categorize errors, prioritize fixes

### Resolution Documentation
Document solution, update runbooks, share learnings

---

## Specification Implementation

### Creating Specifications

Use `zeri add-spec <name>` to create new feature specifications:

```bash
# Create a new specification
zeri add-spec "feature-name"

# This creates .zeri/specs/feature-name.md with the standard template
```

**Specification Structure:**
- **Overview**: Brief description of the feature or enhancement
- **Requirements**: Detailed list of functional requirements
- **Implementation Notes**: Technical considerations and dependencies
- **TODO**: Checklist for tracking implementation progress

### Specification Workflow

1. **Create Specification**: Use `zeri add-spec` command to create structured requirements
2. **Plan Implementation**: Break down requirements into actionable tasks
3. **Implement Features**: Follow the TODO checklist step by step
4. **Mark Progress**: Update TODOs in real-time during development
5. **⚠️ CRITICAL: Code Quality Check**: Run `go fmt ./...` and `go vet ./...` - MUST pass before completion
6. **Review and Complete**: Ensure all requirements are met

### Best Practices

**Specification Content:**
- Write clear, actionable requirements
- Include technical considerations and dependencies
- Reference existing patterns and conventions
- Consider testing and documentation needs

**Implementation Process:**
- Always start with a specification for non-trivial features
- Break complex features into smaller, manageable tasks
- Follow established coding patterns and conventions
- Write tests alongside implementation

### TODO Marking

Mark TODO items as complete when implementing specifications:

- Mark checkboxes as `- [x]` when completing each implementation step
- This helps track progress and manage development workflow
- Update TODOs in real-time during implementation

**Example:**
```markdown
## TODO
- [x] Design and plan implementation
- [x] Implement core functionality
- [ ] Add tests
- [ ] Update documentation
- [ ] Review and refine
- [ ] Mark specification as complete
```

### Specification Directory Structure

```
.zeri/
├── specs/                    # Feature specifications
│   ├── feature-name.md      # Individual specification files
│   └── another-feature.md   # Each spec is self-contained
└── templates/
    └── spec.md              # Template for new specifications
```