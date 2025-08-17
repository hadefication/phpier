# Feature Specification: dockerfile-template-standardization

## Overview
Standardize all Dockerfile templates across different PHP versions to ensure consistency, maintainability, and proper functionality. This addresses inconsistencies found between the main template and version-specific templates that could cause build failures or missing features.

## Requirements
- **Consistent Dependencies**: All templates must include the same system dependencies (especially `libmagickwand-dev` for imagick extension)
- **Uniform Nginx Configuration**: All templates must include the nginx symlink creation for proper site enablement
- **Conditional Composer Logic**: All templates must use the same conditional logic for Composer version selection based on PHP version
- **Standardized Structure**: All templates should follow the same ordering and commenting patterns
- **Multi-Version Compatibility**: Templates must work correctly across all supported PHP versions (5.6, 7.2, 7.3, 7.4, 8.0, 8.1, 8.2, 8.3, 8.4)

## Implementation Notes
- **Files Modified**: 
  - `/internal/templates/files/dockerfiles/php.Dockerfile.tpl` (main template)
  - `/internal/templates/files/dockerfiles/php56-73.Dockerfile.tpl` (older PHP versions)
  - `/internal/templates/files/dockerfiles/php74-80.Dockerfile.tpl` (transitional PHP versions)
  - `/internal/templates/files/dockerfiles/php81-84.Dockerfile.tpl` (modern PHP versions)
- **Key Changes**:
  - Added `libmagickwand-dev` dependency to all version-specific templates
  - Added nginx symlink creation to main template
  - Standardized Composer version selection logic across all templates
  - Fixed template variable references (`.Project.PHP` instead of `.Config.PHP.Version`)
- **Testing Strategy**: Test template generation and Docker builds across multiple PHP versions
- **Critical Dependencies**: Ensures imagick PHP extension can be compiled correctly

## TODO
- [x] Design and plan implementation
- [x] Audit existing templates for inconsistencies
- [x] Add missing `libmagickwand-dev` dependency to version-specific templates
- [x] Add nginx symlink creation to main template
- [x] Standardize Composer version selection logic across all templates
- [x] Fix template variable syntax errors
- [x] Test template generation with PHP 7.3
- [x] Test template generation with PHP 8.3
- [x] Verify generated Dockerfiles contain correct dependencies
- [x] Verify nginx symlink creation in generated files
- [x] Clean up test artifacts and Docker resources
- [x] Mark specification as complete