# PHPier Configuration Guide

## Configuration System

PHPier uses a dual configuration system:

- **Global Config** (`~/.phpier.yaml`): Default settings for all projects
- **Project Config** (`.phpier.yaml`): Per-project overrides and settings

The `init` command creates a project-specific `.phpier.yaml` with your chosen PHP version and project name.

## Generated File Structure

```
your-project/
├── .phpier/                  # All phpier-generated files
│   ├── docker-compose.yml     # Service orchestration
│   ├── Dockerfile.php         # PHP container definition
│   └── docker/
│       ├── php/
│       │   └── php.ini        # PHP settings
│       └── supervisor/
│           └── supervisord.conf # Process management
└── .phpier.yaml              # PHPier configuration
```

**Key Files:**
- **`.phpier.yaml`**: Project configuration (PHP version, project name)
- **`.phpier/Dockerfile.php`**: PHP container with your chosen version and extensions
- **`.phpier/docker-compose.yml`**: Orchestrates your project container with global services
- **`.phpier/docker/php/php.ini`**: PHP configuration settings
- **`.phpier/docker/supervisor/supervisord.conf`**: Process management (PHP-FPM, etc.)

## Project Configuration Example

```yaml
docker:
  composefile: "docker-compose.yml"
  projectname: "my-project"
  network: "my-project"

php:
  version: "8.3"
  extensions:
    - bcmath
    - curl
    - gd
    - mysqli
    - pdo
    - pdo_mysql
    - redis
    - zip
  settings:
    memory_limit: "256M"
    upload_max_filesize: "64M"
    max_execution_time: "300"

services:
  database:
    type: "mysql"
    version: "8.0"
    database: "my-project"
    username: "my-project"
    password: "my-project"
    port: 3306
  cache:
    redis:
      enabled: true
      port: 6379
    memcached:
      enabled: false
      port: 11211
  tools:
    phpmyadmin: true
    mailpit: true
    pgadmin: false
  webserver:
    port: 80
    ssl_port: 443
    doc_root: "/var/www/html"
    index_doc: "index.php"

traefik:
  enabled: true
  domain: "localhost"
  port: 80
  ssl_port: 443
  ssl: false
```

## Customization Examples

All generated files are fully editable for advanced customization.

### Add New Service

Edit `.phpier/docker-compose.yml`:
```yaml
elasticsearch:
  image: elasticsearch:8.8.0
  environment:
    - discovery.type=single-node
  ports:
    - "9200:9200"
  networks:
    - your-project
```

### Install PHP Extension

Edit `.phpier/Dockerfile.php`:
```dockerfile
# Find the RUN docker-php-ext-install line and add your extensions
RUN docker-php-ext-install bcmath gd imagick pdo_pgsql
```

### Modify PHP Settings

Edit `.phpier/docker/php/php.ini`:
```ini
memory_limit = 512M
upload_max_filesize = 100M
max_execution_time = 600
post_max_size = 100M
```

### Customize Process Management

Edit `.phpier/docker/supervisor/supervisord.conf`:
```ini
[program:php-fpm]
command=/usr/local/sbin/php-fpm --nodaemonize --fpm-config /usr/local/etc/php-fpm.conf
autorestart=true

[program:custom-worker]
command=php /var/www/html/worker.php
autorestart=true
```

## When to Rebuild vs Restart

**Rebuild Required** (`phpier up --build -d`):
- Changes to `.phpier/` directory contents (Dockerfile, nginx configs, PHP configs, supervisor configs)
- Changes to `.phpier.yml` Docker Compose configuration
- Adding new PHP extensions or system packages
- Modifying container build process

**Restart Only** (`phpier up -d`):
- Changes to application code (files in your project directory)
- Data or content files that are volume-mounted
- No container configuration changes

**Quick Reference:**
```bash
# After config changes - rebuild required
phpier down && phpier up --build -d

# After code changes - restart only
phpier down && phpier up -d

# Alternative reload command (restarts services without rebuilding)
phpier reload
```

**Why?** Files in `.phpier/` are copied into the Docker image during build time, while your application files are mounted as volumes and update in real-time.

## Domain Access

### With Traefik (Default)
- **Application**: `http://<project-name>.localhost`
- **HTTPS**: `https://<project-name>.localhost` (self-signed cert)
- **Traefik Dashboard**: `http://localhost:8080`
- **Adminer (Database)**: `http://phpier-mysql.localhost`
- **Mailpit (Email)**: `http://phpier-mailpit.localhost`

### Without Traefik
- **Application**: `http://localhost:80`
- **Direct port access based on configuration**

## Database Access

PHPier provides multiple database options with both external client access and web-based administration interfaces.

### MySQL
**External Client Access (Sequel Ace, TablePlus, etc.)**
- **Host**: `localhost`
- **Port**: `3306`
- **Username**: `root`
- **Password**: `phpier`
- **Database**: `phpier` (default)

**Alternative User Access**
- **Username**: `phpier`
- **Password**: `phpier`
- **Database**: `phpier`

**Web Interface**
- **PHPMyAdmin**: `http://phpmyadmin.localhost` (with Traefik) or `http://localhost:8080` (without Traefik)
- **Adminer**: `http://phpier-adminer.localhost` (global service)

### MariaDB
**External Client Access (Sequel Ace, TablePlus, etc.)**
- **Host**: `localhost`
- **Port**: `3307` (note: different port to avoid MySQL conflicts)
- **Username**: `root`
- **Password**: `phpier`
- **Database**: `phpier` (default)

**Alternative User Access**
- **Username**: `phpier`
- **Password**: `phpier`
- **Database**: `phpier`

**Web Interface**
- **Adminer**: `http://phpier-adminer.localhost` (supports both MySQL and MariaDB)

### PostgreSQL
**External Client Access (pgAdmin, TablePlus, etc.)**
- **Host**: `localhost`
- **Port**: `5432`
- **Username**: `phpier`
- **Password**: `phpier`
- **Database**: `phpier`

**Web Interface**
- **pgAdmin**: `http://pgadmin.localhost` (when enabled in global config)
- **Login**: `admin@example.com` / `admin`

### Database Authentication Notes

**MySQL 8.0+ Authentication Fix**
PHPier automatically configures MySQL and MariaDB to use `mysql_native_password` authentication plugin for better client compatibility. This resolves common connection issues with desktop database clients.

**Connection Troubleshooting**
- Ensure phpier global services are running: `phpier global up`
- For authentication errors, try stopping and restarting services: `phpier global down && phpier global up`
- If using custom credentials, check your global configuration: `~/.phpier/config.yaml`

**Creating New Databases**
Use the `root` user credentials to create additional databases through external clients or web interfaces. The `phpier` user has limited privileges by default.

## Email Testing

Mailpit is included for email testing:
- **Web Interface**: `http://localhost:8025` (direct port) or `http://phpier-mailpit.localhost` (Traefik domain)
- **SMTP Server**: `localhost:1025`

Configure your application:
```php
// PHP mail configuration
MAIL_HOST=mailpit
MAIL_PORT=1025
MAIL_ENCRYPTION=null
```