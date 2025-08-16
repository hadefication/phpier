# Feature Specification: Init Command File Placement Fixes

## Overview
Fix file placement issues in the `phpier init` command to ensure proper directory structure and avoid conflicts with existing projects. Address issues with index.php creation, supervisor file placement, and log file exposure.

## Requirements

### File Placement Requirements
- **No public/index.php Creation**: Never create public/index.php during init - this conflicts with existing projects
- **Supervisor Files in .phpier**: All supervisor configuration files must be placed in .phpier directory, not in public
- **Log Files in .phpier**: Expose and organize log files within the .phpier directory structure
- **Clean Project Structure**: Maintain separation between project files and phpier configuration

### Directory Structure Requirements
- **Configuration Files**: All phpier-related files must be in .phpier/ directory
- **No Public Directory Pollution**: Never create files in public/ or project root (except .phpier.yml)
- **Log Organization**: Structured log directory within .phpier for easy access
- **Template Consistency**: Ensure all templates respect the .phpier directory structure

### Existing Project Compatibility
- **Safe Initialization**: init command should be safe to run in existing projects
- **No File Conflicts**: Avoid overwriting or creating conflicting project files
- **Minimal Footprint**: Only create .phpier.yml and .phpier/ directory structure

## Implementation Notes

### Technical Considerations
- **Current Issue**: init command may be creating public/index.php unnecessarily
- **Supervisor Files**: Currently placed in public/, should be in .phpier/supervisor/
- **Log Files**: Need proper exposure and organization in .phpier/logs/
- **Template Updates**: Multiple template files may need directory path updates

### Files Requiring Investigation and Modification
1. **cmd/init.go**: Review file creation logic and remove public/index.php generation
2. **Template Files**: Update supervisor and logging template paths
3. **Docker Compose Templates**: Ensure volume mappings point to correct .phpier paths
4. **Dockerfile Templates**: Update supervisor configuration paths

### Directory Structure Goal
```
project-root/
├── .phpier.yml                    # Project configuration
├── .phpier/                       # All phpier files
│   ├── docker-compose.yml         # Generated compose file
│   ├── Dockerfile.php             # Generated Dockerfile
│   ├── docker/
│   │   ├── nginx/                 # Nginx configurations
│   │   ├── php/                   # PHP configurations
│   │   └── supervisor/            # Supervisor files (moved from public)
│   ├── logs/                      # Exposed log files
│   │   ├── nginx/                 # Nginx logs
│   │   ├── php/                   # PHP-FPM logs
│   │   └── supervisor/            # Supervisor logs
│   └── traefik/                   # Traefik configurations
├── (existing project files)       # User's actual project
└── (no public/index.php created)  # Never create this
```

### Integration Points
- **Docker Volume Mappings**: Update to expose logs correctly
- **Supervisor Configuration**: Ensure log paths point to .phpier/logs/
- **Template Engine**: Update all path references in templates
- **CLI Commands**: Ensure other commands work with new structure

## TODO
- [x] Investigate current init command file creation behavior
- [x] Identify where public/index.php is being created and remove it
- [x] Move supervisor file templates to .phpier/supervisor/ directory
- [x] Create .phpier/logs/ directory structure with proper subdirectories
- [x] Update Docker volume mappings to expose logs from containers
- [x] Update template paths for supervisor configurations
- [x] Test init command in existing project to ensure no conflicts
- [x] Verify log files are properly accessible in .phpier/logs/
- [x] Add .gitignore to .phpier/logs to prevent log files from being committed
- [x] Mark specification as complete

## ✅ Implementation Complete

The init command file placement issues have been successfully fixed:

### Fixed Issues
1. **Removed public/index.php creation**: The init command no longer creates public/index.php files that could conflict with existing projects
2. **Fixed volume mapping**: Changed from `../public:/var/www/html` to `./:/var/www/html` to mount the entire project directory
3. **Removed public directory creation**: No longer creates a public directory during init
4. **Enhanced log structure**: Created proper .phpier/logs/ directory structure with subdirectories for nginx, php, and supervisor logs
5. **Added log volume mappings**: Docker containers now properly expose logs to .phpier/logs/ directories
6. **Improved supervisor configuration**: Added proper log file paths and rotation settings
7. **Added logs .gitignore**: Prevents log files from being committed to git while preserving directory structure

### Directory Structure
The final structure ensures clean separation between project files and phpier configuration:
```
project-root/
├── .phpier.yml                    # Project configuration
├── .phpier/                       # All phpier files
│   ├── docker-compose.yml         # Generated compose file
│   ├── Dockerfile.php             # Generated Dockerfile  
│   ├── docker/
│   │   ├── nginx/                 # Nginx configurations
│   │   ├── php/                   # PHP configurations
│   │   └── supervisor/            # Supervisor configurations
│   └── logs/                      # Log files (with .gitignore)
│       ├── .gitignore             # Ignores all log files
│       ├── nginx/                 # Nginx logs
│       ├── php/                   # PHP-FPM logs
│       └── supervisor/            # Supervisor logs
├── (existing project files)       # User's actual project files
└── (no conflicts created)         # Safe for existing projects
```

### Testing Verified
- Init command tested in directory with existing index.php and README.md
- No conflicts or overwrites occurred
- Log directories properly created and mapped
- Volume mappings correctly configured for current directory mount