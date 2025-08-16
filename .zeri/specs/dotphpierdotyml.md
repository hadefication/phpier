# Feature Specification: .phpier.yml Configuration Improvements

## Overview
Transform .phpier.yml into an ultra-minimal, clean source of truth for essential project settings only. The file contains just the project name, PHP version, and app container essentials (volumes, environment), while all complex configurations (PHP extensions, settings, infrastructure) are handled by generated template files. This provides maximum simplicity with direct file editing for advanced customization.

## Requirements

### Core Functionality Requirements
- **Ultra-Minimal Configuration**: .phpier.yml should contain only essential project settings
  - Project name (top-level, simple string - also used for domain generation)
  - PHP version (simple string, determines Dockerfile template choice)
  - App container essentials (volumes, environment variables only)
- **Hardcoded Sensible Defaults**: Common settings have smart defaults (port 80, webroot `/var/www/html`, PHP memory limits)
- **Extension Management**: PHP extensions should be managed in Dockerfile templates as default installations
- **Direct PHP Customization**: Users edit generated `.phpier/docker/php/php.ini` file directly
- **Domain Generation**: Traefik domains auto-generated from project name (e.g., `myproject.localhost`)
- **Ultra-Clean Structure**: Flat structure, minimal nesting, only essential configuration
- **File Extension Consistency**: Use `.yml` extension for all YAML files (not `.yaml`)
- **Version Switching**: Easy PHP version changes regenerate appropriate Dockerfile template
- **Configuration Validation**: Validate all configuration fields during init and up commands

### User Interface Requirements
- **Simplified Init Process**: `phpier init <version>` should create .phpier.yml with sensible defaults
- **Live Regeneration**: `phpier up` regenerates .phpier/ files based on current .phpier.yml
- **Clear Configuration Structure**: Generated .phpier.yml should be self-documenting
- **Error Messages**: Provide helpful error messages for invalid configurations

### Integration Requirements
- **Docker Integration**: .phpier.yml values must integrate with Docker Compose generation
- **Template System**: Configuration must work with existing template engine
- **CLI Commands**: All phpier commands must recognize the simplified configuration
- **Global Config Compatibility**: Project config must work alongside global configuration

### Performance and Security Requirements
- **Fast Parsing**: Configuration loading should be optimized for CLI responsiveness
- **Secure Defaults**: Default configurations should follow security best practices
- **Validation**: Input validation to prevent configuration injection attacks

## Implementation Notes

### Technical Considerations
- **Current Structure**: ProjectConfig in `internal/config/config.go:11-15` contains Docker and PHP configuration
- **Dependencies**: Viper for YAML parsing, existing template engine for file generation
- **File Locations**: 
  - Project config: `.phpier.yml` (currently `.phpier.yaml`)
  - Global config: `~/.phpier/config.yaml`

### Files Requiring Modification
1. **internal/config/config.go**:
   - Remove Extensions and Settings fields from ProjectConfig struct
   - Update LoadProjectConfig and SaveProjectConfig functions
   - Change all references from `.phpier.yaml` to `.phpier.yml`

2. **cmd/init.go**:
   - Remove extension and settings configuration generation
   - Move getDefaultExtensions logic to Dockerfile template data
   - Update file references to use `.yml` extension

3. **internal/templates/files/dockerfiles/*.Dockerfile.tpl**:
   - Hardcode PHP extensions as default installations in all templates
   - Remove template variables for user-configured extensions

4. **internal/templates/files/configs/php.ini.tpl**:
   - Move default PHP settings from user config to php.ini template
   - Include all common PHP.ini configurations as defaults

5. **internal/generator/generator.go**:
   - Remove extension and settings data from template rendering
   - Update file extension references to `.yml`

6. **Error messages and CLI help** (multiple files):
   - Update all references from `.phpier.yaml` to `.phpier.yml`
   - Remove extension-related error messages and suggestions

### Integration Points
- **Template Engine**: Pass simplified config data to templates
- **Docker Compose Generation**: Use project name for container naming and network configuration
- **CLI Commands**: Update validation in up, down, and other commands
- **Error Handling**: Update error messages in `internal/errors/factories.go`

### Architectural Decisions
- **Separation of Concerns**: User configuration vs. build-time configuration
- **Template-Based Extensions**: Extensions defined in Dockerfiles, not user config
- **Ultra-Minimal Approach**: Only essential settings in user configuration

### Testing Strategy
- **Unit Tests**: Test configuration loading, saving, and validation
- **Integration Tests**: Test init command with ultra-minimal configuration
- **Template Tests**: Verify Dockerfile generation with hardcoded extensions
- **Version Switching Tests**: Test PHP version changes regenerate correct templates

## Proposed .phpier.yml Structure

### Current Structure (Complex)
```yaml
docker:
  project_name: "myproject"
php:
  version: "8.3"
  extensions:
    - bcmath
    - curl
    - gd
    - mysqli
    # ... many more extensions
  settings:
    memory_limit: "512M"
```

### Final Structure (Ultra-Minimal & Practical)
```yaml
name: "myproject"               # Project name (simple and clean)
php: "8.3"                      # PHP version (simple string)
  
app:
  volumes:
    - "./:/var/www/html"         # Host:Container volume mappings
  environment:
    - "APP_ENV=local"            # Environment variables
    - "APP_DEBUG=true"
```

### Configuration Philosophy

#### What .phpier.yml Controls (Source of Truth)
- **Project Identity**: `name` - Simple project identifier (also used for domain generation)
- **Runtime Version**: `php` - PHP version string (determines Dockerfile template)
- **App Container**: Essential container configuration (only what users typically customize)
  - Volume mounts and file mappings
  - Environment variables

#### What Templates Handle (Generated Files)  
- **PHP Extensions**: Hardcoded in Dockerfile templates based on PHP version
- **PHP Settings**: Static php.ini template with sensible defaults (users edit directly)
- **Infrastructure**: Docker Compose, Nginx, Supervisor configurations  
- **Domain Routing**: Traefik rules auto-generated from project `name` (e.g., `myproject.localhost`)
- **Sensible Defaults**: Port 80, webroot `/var/www/html` (hardcoded, users can edit generated files)

#### User Workflow
1. **Init**: `phpier init 8.3` creates minimal .phpier.yml + generates .phpier/ files
2. **Project Settings**: Edit .phpier.yml for container behavior (volumes, environment, PHP version)
3. **PHP Settings**: Edit `.phpier/docker/php/php.ini` directly for memory_limit, upload sizes, etc.
4. **Version Switch**: Change `php: "8.1"` in .phpier.yml → `phpier up` regenerates Dockerfile only
5. **Run**: `phpier up` regenerates files and starts containers

#### File Generation Strategy
- **Smart Regeneration**: Files in .phpier/ are regenerated on every `phpier up`
- **Config-Driven**: Generated files reflect current .phpier.yml settings  
- **Static Templates**: PHP settings come from static template (no user config merging)
- **Direct Editing**: Users edit generated files directly for advanced customization
- **Version Flexibility**: Easy PHP version switching regenerates appropriate Dockerfile


## TODO
- [x] Remove Extensions and Settings fields from ProjectConfig struct
- [x] Update init command to generate ultra-minimal .phpier.yml (name + php version + app essentials)
- [x] Move PHP extensions to Dockerfile templates as default installations
- [x] Move PHP settings to static php.ini template file  
- [x] Update all file references from .phpier.yaml to .phpier.yml
- [x] Update error messages and CLI help text
- [x] Remove extension-related template variables and logic
- [x] Hardcode sensible defaults (port 80, webroot) in templates
- [x] Implement auto-domain generation from project name
- [x] Test ultra-minimal configuration approach
- [x] Verify implementation matches specification
- [x] Mark specification as complete

## ✅ Implementation Complete

The ultra-minimal .phpier.yml configuration has been successfully implemented and tested. The final structure provides the perfect balance of simplicity and control, with essential project settings in the config file and advanced customization through direct file editing.