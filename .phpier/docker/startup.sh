#!/bin/bash
set -e

echo "Starting container initialization..."

# Fix permissions for mounted files
echo "Fixing permissions..."
chown -R www-data:www-data /var/www/html
chmod -R 755 /var/www/html

# Create index.php if it doesn't exist
if [ ! -f /var/www/html/index.php ]; then
    echo "Creating default index.php..."
    echo "<?php phpinfo(); ?>" > /var/www/html/index.php
    chown www-data:www-data /var/www/html/index.php
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