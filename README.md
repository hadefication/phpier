# PHPier

A powerful CLI tool for managing PHP development environments using Docker. Like a pier that connects the land to the water, PHPier connects your PHP projects to containerized development infrastructure. Supports multiple PHP versions with automatic service orchestration and folder-based domain routing.

## 🚀 Features

- **Multi-PHP Support**: PHP 5.6, 7.3, 7.4, 8.0, 8.1, 8.2, 8.3, 8.4
- **Global Services**: Shared Traefik, databases, and tools across all projects
- **Database Options**: MySQL, PostgreSQL, MariaDB with admin tools
- **Caching Services**: Redis and Memcached support
- **Development Tools**: Mailpit, PHPMyAdmin, pgAdmin
- **Traefik Integration**: Automatic reverse proxy with `<project>.localhost` domains
- **Dual Architecture**: Global services + individual project containers
- **Flexible Commands**: Simple service management with safety checks

## 📋 Prerequisites

- [Docker](https://docs.docker.com/get-docker/) (20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (2.0+)
- [Go](https://golang.org/dl/) (1.20+) - for building from source

## 🔧 Installation

### Option 1: Local Install Script (Recommended)

The easiest way to install phpier is using the automated install script:

```bash
# Clone the repository
git clone <repository-url>
cd phpier

# Install locally (handles everything automatically)
./scripts/local-install.sh

# Or use Make
make install
```

**What the install script does:**
- 🗑️ **Uninstalls** any existing phpier installation
- 🔨 **Builds** the binary with version information from git  
- 🔐 **Sets** executable permissions
- 📦 **Installs** to `/usr/local/bin/phpier`
- ✅ **Verifies** the installation works

### Option 2: Manual Build

```bash
# Clone and build manually
git clone <repository-url>
cd phpier

# Build with version information
go build -ldflags="-s -w -X main.version=dev" -o phpier

# Install globally (optional)
chmod +x phpier
sudo mv phpier /usr/local/bin/

# Or run locally
./phpier --help
```

### Option 3: Using Make

```bash
# Build only
make build

# Build and install
make install

# Clean build artifacts
make clean
```

### Option 4: Download Binary (Coming Soon)

Pre-built binaries will be available once the first release is published.

**Linux x64:**
```bash
curl -L https://github.com/your-org/phpier/releases/latest/download/phpier-linux-amd64 -o phpier
chmod +x phpier
sudo mv phpier /usr/local/bin/
```

**macOS x64 (Intel):**
```bash
curl -L https://github.com/your-org/phpier/releases/latest/download/phpier-darwin-amd64 -o phpier
chmod +x phpier
sudo mv phpier /usr/local/bin/
```

**macOS ARM64 (Apple Silicon):**
```bash
curl -L https://github.com/your-org/phpier/releases/latest/download/phpier-darwin-arm64 -o phpier
chmod +x phpier
sudo mv phpier /usr/local/bin/
```

**Auto-detect (Linux/macOS only):**
```bash
# This command auto-detects your platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then ARCH="amd64"; fi
if [ "$ARCH" = "aarch64" ]; then ARCH="arm64"; fi
curl -L "https://github.com/your-org/phpier/releases/latest/download/phpier-${OS}-${ARCH}" -o phpier
chmod +x phpier
sudo mv phpier /usr/local/bin/
```

### Option 3: Development Mode

```bash
# Run directly with Go
cd phpier
go run main.go [command]
```

## 🗑️ Uninstallation

To remove phpier from your system:

```bash
# Using the uninstall script (recommended)
./scripts/local-uninstall.sh

# Or using Make
make uninstall

# Or manually
sudo rm /usr/local/bin/phpier
```

The uninstall script will:
- Show the current version before removal
- Prompt for confirmation
- Remove the binary with appropriate permissions
- Verify complete removal

## 🏗️ How It Works

PHPier uses a **dual architecture** with global services and individual project containers:

### Global Services
- **Shared across all projects**: Traefik, databases, Redis, Mailpit, etc.
- **Always running**: Started once with `phpier start` or `phpier global up`
- **Persistent**: Data and configurations persist between project sessions

### Project Containers
- **One per project**: Each project gets its own PHP/Nginx container
- **Project-specific**: Uses your chosen PHP version and extensions
- **Connected**: Automatically connects to global services network

### Workflow
1. **Start global services** once: `phpier start`
2. **Initialize projects** as needed: `phpier init`
3. **Run projects** individually: `phpier up`
4. **Access** via `http://project-name.localhost`

This means you can run multiple PHP projects with different versions simultaneously, all sharing the same databases and tools.

## 🎯 Quick Start

### 1. Initialize a New Project

```bash
# Create your project directory
mkdir my-laravel-app
cd my-laravel-app

# Initialize with PHP 8.3 (default)
phpier init

# Or specify version and project name
phpier init 7.4 --project-name=myapp
phpier init 8.1 --php-version=8.1
```

### 2. Start the Development Environment

```bash
# Start global services first (Traefik, databases, etc.)
phpier start

# Then start your project container
phpier up -d

# Or start with image rebuilding
phpier up --build

# Start in foreground (see logs)
phpier up
```

### 3. Access Your Application

After starting, you'll see output like:
```
🌐 Your development environment is ready!
🔗 Application: http://my-laravel-app.localhost
🔧 Traefik Dashboard: http://localhost:8080
🗄️  PHPMyAdmin: http://localhost:8080
📧 Mailpit: http://localhost:8025
💾 Database: mysql (port 3306)
🔴 Redis: localhost:6379
```

### 4. Stop the Environment

```bash
# Stop project services only
phpier down

# Stop project and global services
phpier down --stop-global

# Stop project and remove data volumes
phpier down --remove-volumes

# Or stop global services separately
phpier stop

# Force stop global services (ignore running projects warning)
phpier stop --force
```

## 📖 Detailed Usage

### Available Commands

```bash
phpier --help                    # Show all commands
phpier init --help               # Initialize environment help
phpier start --help              # Start global services help
phpier stop --help               # Stop global services help
phpier up --help                 # Start project container help
phpier down --help               # Stop project container help
phpier global --help             # Manage global services help
phpier version                   # Show version information
```

### Service Management Commands

#### Global Service Commands
```bash
# Start global services (Traefik, databases, etc.)
phpier start                     # Start in background (default)
phpier start --build             # Rebuild and start
phpier start --force             # Force restart if already running

# Stop global services
phpier stop                      # Stop with safety checks
phpier stop --force             # Force stop (ignore warnings)
phpier stop --remove-volumes    # Stop and remove global volumes (dangerous)
```

#### Project Container Commands
```bash
# Start project container
phpier up                        # Start in foreground
phpier up -d                     # Start in background (detached)
phpier up --build               # Rebuild and start

# Stop project container
phpier down                      # Stop project container only
phpier down --stop-global       # Stop project and global services
phpier down --remove-volumes    # Stop and remove all volumes
phpier down --force             # Force remove containers

# Alternative global service management
phpier global up                 # Start global services (alternative to 'start')
phpier global down               # Stop global services (alternative to 'stop')
```

### Init Command Options

```bash
# Basic initialization
phpier init                      # PHP 8.3 with default settings

# Specify PHP version
phpier init 7.4                  # PHP 7.4
phpier init 8.1                  # PHP 8.1

# Using flags
phpier init --php-version=8.2    # Explicit PHP version flag
phpier init --project-name=myapp # Custom project name (defaults to directory name)

# Combined options
phpier init 7.4 --project-name=legacy-app
```

### Configuration System

PHPier uses a dual configuration system:

- **Global Config** (`~/.phpier.yaml`): Default settings for all projects
- **Project Config** (`.phpier.yaml`): Per-project overrides and settings

The `init` command creates a project-specific `.phpier.yaml` with your chosen PHP version and project name.

### Examples

#### Laravel Project
```bash
mkdir my-laravel-app
cd my-laravel-app
phpier init 8.3
phpier start              # Start global services
phpier up -d              # Start project container
# Visit http://my-laravel-app.localhost
```

#### Legacy PHP Project
```bash
mkdir legacy-project
cd legacy-project
phpier init 7.4 --project-name=legacy
phpier start              # Start global services
phpier up -d              # Start project container
# Visit http://legacy.localhost
```

#### Different PHP Versions
```bash
# PHP 8.2 Project
mkdir modern-api
cd modern-api
phpier init 8.2
phpier start
phpier up -d

# PHP 5.6 Legacy Project
mkdir very-old-project
cd very-old-project
phpier init 5.6 --project-name=oldapp
phpier up -d              # Global services already running
# Visit http://oldapp.localhost
```

## 🛠️ Customization

All generated files are fully editable:

### Generated File Structure
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

### Common Customizations

#### Add New Service
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

#### Install PHP Extension
Edit `.phpier/Dockerfile.php`:
```dockerfile
# Find the RUN docker-php-ext-install line and add your extensions
RUN docker-php-ext-install bcmath gd imagick pdo_pgsql
```

#### Modify PHP Settings
Edit `.phpier/docker/php/php.ini`:
```ini
memory_limit = 512M
upload_max_filesize = 100M
max_execution_time = 600
post_max_size = 100M
```

#### Customize Process Management
Edit `.phpier/docker/supervisor/supervisord.conf`:
```ini
[program:php-fpm]
command=/usr/local/sbin/php-fpm --nodaemonize --fpm-config /usr/local/etc/php-fpm.conf
autorestart=true

[program:custom-worker]
command=php /var/www/html/worker.php
autorestart=true
```

After making changes:
```bash
phpier down      # Stop current container
phpier up --build  # Rebuild and restart with changes
```

## 🌐 Domain Access

### With Traefik (Default)
- **Application**: `http://<project-name>.localhost`
- **HTTPS**: `https://<project-name>.localhost` (self-signed cert)
- **Traefik Dashboard**: `http://localhost:8080`

### Without Traefik
- **Application**: `http://localhost:80`
- **Direct port access based on configuration**

## 🗄️ Database Access

### MySQL/MariaDB
- **Host**: `localhost`
- **Port**: `3306`
- **Database**: `<project-name>`
- **Username**: `<project-name>`
- **Password**: `<project-name>`
- **Admin**: PHPMyAdmin at `http://localhost:8080`

### PostgreSQL
- **Host**: `localhost`
- **Port**: `5432`
- **Database**: `<project-name>`
- **Username**: `<project-name>`
- **Password**: `<project-name>`
- **Admin**: pgAdmin at `http://localhost:8081`

## 📧 Email Testing

Mailpit is included for email testing:
- **Web Interface**: `http://localhost:8025`
- **SMTP Server**: `localhost:1025`

Configure your application:
```php
// PHP mail configuration
MAIL_HOST=mailpit
MAIL_PORT=1025
MAIL_ENCRYPTION=null
```

## 🔄 Supported PHP Versions

phpier supports all major PHP versions with appropriate tooling:

| PHP Version | Status | Extensions Included |
|-------------|--------|-------------------|
| 5.6 | Legacy Support | Basic web development extensions |
| 7.3 | Legacy Support | Full extension set |
| 7.4 | Active Support | Full extension set + modern tools |
| 8.0+ | Active Support | Latest extensions + performance optimizations |

Each container includes Composer, NVM (Node.js), and appropriate tooling for the PHP version.

## 🐛 Troubleshooting

### Docker Issues
```bash
# Check Docker is running
docker --version
docker-compose --version

# Check if ports are available
docker ps
netstat -tulpn | grep :80
```

### Permission Issues
```bash
# Fix file permissions
sudo chown -R $USER:$USER .
chmod -R 755 .
```

### Container Issues
```bash
# View project logs
cd your-project/.phpier
docker-compose logs

# View specific service logs
docker-compose logs app

# Rebuild project from scratch
phpier down --remove-volumes
phpier up --build

# Reset global services
phpier stop
phpier start --force
```

### Domain Access Issues
```bash
# Test Traefik routing
curl -H "Host: your-project.localhost" http://localhost

# Check DNS resolution
ping your-project.localhost
```

## 🔧 Configuration

### Configuration Files

- **Global Config**: `~/.phpier.yaml` - Default settings for all projects
- **Project Config**: `.phpier.yaml` - Project-specific settings (created by `init`)

### Project Configuration Example

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

## 🤝 Contributing

See [DEV.md](DEV.md) for development information and contribution guidelines.

## 📄 License

[MIT License](LICENSE)

## 🙋 Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/your-org/phpier/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/your-org/phpier/discussions)
- 📖 **Documentation**: [Wiki](https://github.com/your-org/phpier/wiki)

---

**Happy coding with PHPier! 🚀**