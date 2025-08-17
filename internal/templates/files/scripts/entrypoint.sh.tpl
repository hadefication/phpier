#!/usr/bin/env bash

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
    {{- if or (eq .PHPVersion "5.6") (eq .PHPVersion "7.2") (eq .PHPVersion "7.3") }}
    if [ -f "/etc/php/{{ .PHPVersion }}/fpm/pool.d/www.conf" ]; then
        echo "Updating PHP-FPM config for PHP {{ .PHPVersion }}"
        sed -i "s/user\ \=.*/user\ \= ${WWWUSER:-phpier}/g" /etc/php/{{ .PHPVersion }}/fpm/pool.d/www.conf 2>/dev/null || echo "Warning: Failed to update PHP config"
    fi
    {{- else if or (eq .PHPVersion "7.4") (eq .PHPVersion "8.0") }}
    if [ -f "/etc/php/{{ .PHPVersion }}/fpm/pool.d/www.conf" ]; then
        echo "Updating PHP-FPM config for PHP {{ .PHPVersion }}"
        sed -i "s/user\ \=.*/user\ \= ${WWWUSER:-phpier}/g" /etc/php/{{ .PHPVersion }}/fpm/pool.d/www.conf 2>/dev/null || echo "Warning: Failed to update PHP config"
    fi
    {{- else }}
    if [ -f "/etc/php/{{ .PHPVersion }}/fpm/pool.d/www.conf" ]; then
        echo "Updating PHP-FPM config for PHP {{ .PHPVersion }}"
        sed -i "s/user\ \=.*/user\ \= ${WWWUSER:-phpier}/g" /etc/php/{{ .PHPVersion }}/fpm/pool.d/www.conf 2>/dev/null || echo "Warning: Failed to update PHP config"
    fi
    {{- end }}

    # Create Composer directory
    mkdir -p /.composer
fi

# Set proper permissions for Composer directory
echo "Setting Composer directory permissions..."
chmod -R ugo+rw /.composer 2>/dev/null || echo "Warning: Failed to set Composer permissions"

echo "Entrypoint initialization complete."

# Execute command if provided, otherwise run startup script
if [ $# -gt 0 ]; then
    echo "Executing command: $@"
    # Use gosu to execute command as the mapped user
    exec gosu ${WWWUSER:-phpier} "$@"
else
    echo "Starting container services..."
    # Run the startup script which handles initialization and starts supervisord
    exec /usr/local/bin/startup.sh
fi