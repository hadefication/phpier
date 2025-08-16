# Feature Specification: replace-nvm-with-configurable-node

## Overview
Replace NVM (Node Version Manager) with direct Node.js installation based on a configurable version specified in .phpier.yml. This simplifies the Docker build process, reduces image size, and provides project-specific Node.js version control.

## Requirements
- Add `node: <version>` configuration key to .phpier.yml
- Support version formats: "lts", "16", "18", "20", "22", "none" 
- Replace NVM installation with direct Node.js installation in all Dockerfile templates
- Default to "lts" for new projects when node key is not specified
- Remove all NVM-related code from Dockerfile templates
- Maintain npm availability when Node.js is installed
- Support "none" option to skip Node.js installation entirely

## Implementation Notes
- **Config Parsing**: Update Go configuration structs to handle `node` field
- **Dockerfile Templates**: Replace NVM installation with NodeSource repository setup
- **Version Resolution**: Map "lts" to latest LTS version during template rendering
- **Template Files to Update**: 
  - `internal/templates/files/dockerfiles/php.Dockerfile.tpl`
  - `internal/templates/files/dockerfiles/php56-73.Dockerfile.tpl` 
  - `internal/templates/files/dockerfiles/php74-80.Dockerfile.tpl`
  - `internal/templates/files/dockerfiles/php81-84.Dockerfile.tpl`
- **Testing Strategy**: Test with different Node versions across multiple PHP versions
- **Backward Compatibility**: Existing projects without `node` key should default to "lts"

## TODO
- [x] Design and plan implementation
- [x] Examine current NVM implementation in Dockerfile templates
- [x] Update config parsing to handle 'node' key in .phpier.yml
- [x] Replace NVM with direct Node installation in all Dockerfile templates
- [x] Test with different Node versions (lts, specific version, none)
- [x] Update documentation and complete specification

## Implementation Summary

### Changes Made:

**Configuration:**
- Added `Node string` field to `ProjectConfig` struct in `internal/config/config.go`
- Updated `SaveProjectConfig` to include node configuration
- Modified `createProjectConfig` in `cmd/init.go` to default to "lts"

**Template Engine:**
- Added `resolveNodeVersion` function to resolve "lts", version numbers, and "none"
- Added `shouldInstallNode` function to determine if Node.js should be installed
- Added `split` helper function for template string manipulation

**Dockerfile Templates:**
- Updated all 4 Dockerfile templates to use conditional Node.js installation
- Replaced NVM installation with direct NodeSource repository installation
- Uses `apt-get install nodejs=<version>` for specific version control
- Templates now conditionally install based on `node` configuration

**Supported Node Versions:**
- `"lts"` → Latest LTS version (determined by NodeSource)
- `"16"` → Latest Node 16.x (determined by NodeSource)
- `"18"` → Latest Node 18.x (determined by NodeSource) 
- `"20"` → Latest Node 20.x (determined by NodeSource)
- `"22"` → Latest Node 22.x (determined by NodeSource)
- `"none"` → Skip Node.js installation entirely
- Full versions (e.g., "18.20.4") → Install specific version

**Reference:** 
- [Node.js Release Schedule](https://nodejs.org/en/about/releases/)
- [Node.js Downloads](https://nodejs.org/en/download/)
- [NodeSource Binary Distributions](https://github.com/nodesource/distributions)

### Benefits Achieved:
- **Faster Builds**: Direct Node installation vs NVM complexity
- **Smaller Images**: No NVM overhead and installation scripts
- **Version Control**: Project-specific Node versions via .phpier.yml
- **Simpler Dockerfiles**: Cleaner, more maintainable templates
- **Flexibility**: Can skip Node entirely for pure PHP projects
- **Dynamic Versioning**: Always gets latest patch versions from NodeSource repository
- **No Hard-coded Versions**: Automatically stays current with Node.js releases