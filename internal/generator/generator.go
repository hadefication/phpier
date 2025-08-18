package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"phpier/internal/config"
	"phpier/internal/templates"

	"github.com/sirupsen/logrus"
)

// GenerateProjectFiles generates all necessary files for a new project.
func GenerateProjectFiles(engine *templates.Engine, projectCfg *config.ProjectConfig, globalCfg *config.GlobalConfig) error {
	// Note: docker-compose.yml is now generated directly in project root by init command

	// Generate Dockerfile for the project
	dockerfile, err := engine.RenderPHPDockerfile(projectCfg)
	if err != nil {
		return fmt.Errorf("failed to render Dockerfile: %w", err)
	}
	if err := WriteFile(".phpier/Dockerfile.php", dockerfile); err != nil {
		return err
	}

	// Generate Nginx, Supervisor, and PHP configs
	// Note: These templates might need to be created or updated to work with the new config structure
	supervisorConf := `[supervisord]
nodaemon=true
user=root
# Move PID file to proper location
pidfile=/var/run/supervisor/supervisord.pid
# Move log files to proper location
logfile=/var/log/supervisor/supervisord.log
childlogdir=/var/log/supervisor
loglevel=info
silent=false

[unix_http_server]
file=/var/run/supervisor/supervisor.sock
chmod=0700
username=supervisor
password=supervisor

[supervisorctl]
serverurl=unix:///var/run/supervisor/supervisor.sock
username=supervisor
password=supervisor

[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface

# PHP-FPM program
[program:php-fpm]
command=/usr/local/sbin/php-fpm --nodaemonize --fpm-config /usr/local/etc/php-fpm.conf
autostart=true
autorestart=true
priority=5
stdout_logfile=/var/log/supervisor/php-fpm.log
stderr_logfile=/var/log/supervisor/php-fpm-error.log
user=root
killasgroup=true
stopasgroup=true

# Nginx program
[program:nginx]
command=/usr/sbin/nginx -g "daemon off;"
autostart=true
autorestart=true
priority=10
stdout_logfile=/var/log/supervisor/nginx.log
stderr_logfile=/var/log/supervisor/nginx-error.log
user=root
killasgroup=true
stopasgroup=true`
	if err := WriteFile(".phpier/docker/supervisor/supervisord.conf", supervisorConf); err != nil {
		return err
	}

	// Generate PHP configuration
	phpIni, err := engine.RenderPHPConfig()
	if err != nil {
		return fmt.Errorf("failed to render php.ini: %w", err)
	}
	if err := WriteFile(".phpier/docker/php/php.ini", phpIni); err != nil {
		return err
	}

	// Generate Nginx main configuration
	nginxConf, err := engine.RenderNginxConfig(projectCfg)
	if err != nil {
		return fmt.Errorf("failed to render nginx.conf: %w", err)
	}
	if err := WriteFile(".phpier/docker/nginx/nginx.conf", nginxConf); err != nil {
		return err
	}

	// Generate Nginx site configuration (default.conf)
	nginxSiteConf, err := engine.RenderNginxSiteConfig(projectCfg, globalCfg)
	if err != nil {
		return fmt.Errorf("failed to render nginx site config: %w", err)
	}
	if err := WriteFile(".phpier/docker/nginx/default.conf", nginxSiteConf); err != nil {
		return err
	}

	return nil
}

// GenerateGlobalFiles generates all necessary files for the global services stack.
func GenerateGlobalFiles(engine *templates.Engine, globalCfg *config.GlobalConfig) error {
	// Generate docker-compose.yml for the global stack
	dockerCompose, err := engine.RenderGlobalDockerCompose(globalCfg)
	if err != nil {
		return fmt.Errorf("failed to render global docker-compose.yml: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	globalPath := filepath.Join(home, ".phpier")
	if err := WriteFile(filepath.Join(globalPath, "docker-compose.yml"), dockerCompose); err != nil {
		return err
	}

	// Generate Traefik configuration
	traefikConfig, err := engine.RenderTraefikConfig(globalCfg)
	if err != nil {
		return fmt.Errorf("failed to render traefik config: %w", err)
	}
	if err := WriteFile(filepath.Join(globalPath, "traefik", "traefik.yml"), traefikConfig); err != nil {
		return err
	}

	// Generate Traefik dynamic configuration
	traefikDynamic, err := engine.RenderTraefikDynamicConfig(globalCfg)
	if err != nil {
		return fmt.Errorf("failed to render traefik dynamic config: %w", err)
	}
	if err := WriteFile(filepath.Join(globalPath, "traefik", "dynamic", "api.yml"), traefikDynamic); err != nil {
		return err
	}

	return nil
}

// CreateProjectDirectories creates the directory structure for a new project.
func CreateProjectDirectories() error {
	dirs := []string{
		".phpier",
		".phpier/docker",
		".phpier/docker/php",
		".phpier/docker/nginx",
		".phpier/docker/supervisor",
		".phpier/logs",
		".phpier/logs/nginx",
		".phpier/logs/php",
		".phpier/logs/supervisor",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		logrus.Debugf("Created directory: %s", dir)
	}

	// Create .gitignore in logs directory to prevent log files from being committed
	logsGitignore := `# Ignore all log files
*
# But keep this .gitignore file
!.gitignore`
	if err := WriteFile(".phpier/logs/.gitignore", logsGitignore); err != nil {
		return fmt.Errorf("failed to create logs .gitignore: %w", err)
	}

	// Create entrypoint script for permission mapping
	entrypointScript := `#!/usr/bin/env bash

# Exit on any error
set -e

echo "Starting entrypoint script..."

# Handle WWWUSER environment variable for permission mapping
if [ ! -z "$WWWUSER" ]; then
    echo "Setting up user mapping for WWWUSER=$WWWUSER"
    
    # Check if the UID is already in use by another user
    if id "$WWWUSER" >/dev/null 2>&1; then
        echo "UID $WWWUSER already exists, skipping usermod"
    else
        # Map the container user to the host user ID
        echo "Mapping phpier user to UID $WWWUSER"
        usermod -u $WWWUSER phpier 2>/dev/null || echo "Warning: usermod failed, continuing..."
    fi
fi

# Initialize Composer directory if it doesn't exist
if [ ! -d /.composer ]; then
    echo "Initializing Composer directory..."
    
    # Update PHP-FPM pool configuration with correct user
    for php_conf in /etc/php/*/fpm/pool.d/www.conf; do
        if [ -f "$php_conf" ]; then
            echo "Updating PHP-FPM config: $php_conf"
            sed -i "s/user\ \=.*/user\ \= ${WWWUSER:-phpier}/g" "$php_conf" 2>/dev/null || echo "Warning: Failed to update $php_conf"
        fi
    done

    # Create Composer directory
    mkdir -p /.composer
fi

# Set proper permissions for Composer directory
echo "Setting Composer directory permissions..."
chmod -R ugo+rw /.composer 2>/dev/null || echo "Warning: Failed to set Composer permissions"

echo "Entrypoint initialization complete."

# Execute command if provided, otherwise start services
if [ $# -gt 0 ]; then
    echo "Executing command: $@"
    # Use gosu to execute command as the mapped user
    exec gosu ${WWWUSER:-phpier} "$@"
else
    echo "Starting container services..."
    
    # Skip permission changes when WWWUSER is set (user mapping handles this)
    if [ -z "$WWWUSER" ]; then
        echo "Fixing permissions..."
        chown -R www-data:www-data /var/www/html
        chmod -R 755 /var/www/html
    else
        echo "WWWUSER set ($WWWUSER), skipping permission changes (user mapping active)"
    fi
    
    # Ensure supervisor directories exist with proper permissions
    mkdir -p /var/run/supervisor /var/log/supervisor
    chown -R root:root /var/run/supervisor /var/log/supervisor
    
    # Test nginx configuration
    echo "Testing Nginx configuration..."
    nginx -t
    
    echo "Starting supervisord..."
    # Use explicit config file and PID file location
    exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf \
        --pidfile=/var/run/supervisor/supervisord.pid \
        --logfile=/var/log/supervisor/supervisord.log
fi`
	if err := WriteFile(".phpier/docker/entrypoint.sh", entrypointScript); err != nil {
		return fmt.Errorf("failed to create entrypoint script: %w", err)
	}

	return nil
}

// CreateGlobalDirectories creates the directory structure for the global services.
func CreateGlobalDirectories() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	dirs := []string{
		filepath.Join(home, ".phpier"),
		filepath.Join(home, ".phpier", "traefik"),
		filepath.Join(home, ".phpier", "traefik", "dynamic"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		logrus.Debugf("Created global directory: %s", dir)
	}

	return nil
}

// WriteFile writes content to a file, creating directories as needed.
func WriteFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}

	logrus.Debugf("Generated file: %s", path)
	return nil
}
