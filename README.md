# PHPier

A powerful CLI tool for managing PHP development environments using Docker. Like a pier that connects the land to the water, PHPier connects your PHP projects to containerized development infrastructure. Supports multiple PHP versions with automatic service orchestration and folder-based domain routing.

## 🚀 Features

- **Multi-PHP Support**: PHP 5.6, 7.2, 7.3, 7.4, 8.0, 8.1, 8.2, 8.3, 8.4
- **Global Services**: Shared Traefik, databases, and tools across all projects
- **Database Options**: MySQL, PostgreSQL, MariaDB with admin tools
- **Caching Services**: Redis and Memcached support
- **Development Tools**: Mailpit, PHPMyAdmin, pgAdmin
- **Traefik Integration**: Automatic reverse proxy with `<project>.localhost` domains
- **Dual Architecture**: Global services + individual project containers
- **Flexible Commands**: Simple service management with safety checks

## 📋 Installation

### One-Line Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh | bash
```

### Alternative Installation Methods

**Install Specific Version:**
```bash
curl -sSL https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh | bash -s -- -v v1.0.0
```

**Install to Custom Directory:**
```bash
curl -sSL https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh | bash -s -- -d /usr/local/bin
```

**Install with Both Options:**
```bash
curl -sSL https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh | bash -s -- -v v1.0.0 -d /usr/local/bin
```

**Force Installation (Skip Existing Check):**
```bash
curl -sSL https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh | bash -s -- --force
```

**Build from Source:**
```bash
git clone https://github.com/hadefication/phpier.git
cd phpier
./scripts/local-install.sh
```

### Platform Support
- ✅ **Linux** (AMD64, ARM64) - Ubuntu, CentOS, Debian, etc.
- ✅ **macOS** (Intel, Apple Silicon) - macOS 10.15+
- ✅ **Windows WSL** (Ubuntu, Debian) - WSL 2 recommended

> **Note:** Native Windows is not supported. Use WSL (Windows Subsystem for Linux) for Windows development.

See [INSTALLATION.md](INSTALLATION.md) for detailed installation instructions and troubleshooting.

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

```bash
# 1. Initialize a new project
mkdir my-project && cd my-project
phpier init 8.3

# 2. Start services
phpier start              # Start global services
phpier up -d              # Start project container

# 3. Access your app
# Visit: http://my-project.localhost
```

## 📚 Documentation

- **[Installation Guide](INSTALLATION.md)**: Detailed installation methods and options
- **[Configuration](CONFIGURATION.md)**: Project setup, customization, and settings
- **[Troubleshooting](TROUBLESHOOTING.md)**: Common issues and debugging tips

## 🛠️ Commands

PHPier provides a comprehensive set of commands for managing your PHP development environment.

### Service Management

#### Global Services
```bash
phpier start                 # Start global services (Traefik, databases, etc.)
phpier stop                  # Stop global services with safety checks
phpier global up             # Start global shared services stack
phpier global down           # Stop global shared services stack
```

#### Project Services
```bash
phpier up [-d]               # Start project container (detached mode optional)
phpier down                  # Stop project containers
phpier build                 # Build/rebuild project's app container
phpier reload                # Restart project services with optional rebuild
```

### Project Management

#### Project Setup
```bash
phpier init [version]        # Initialize project (default: PHP 8.3)
phpier init 8.1              # Initialize with specific PHP version
phpier init --project-name=myapp  # Custom project name
phpier init 7.4 --project-name=legacy  # Version + custom name
```

**File Structure After Init:**

```
my-project/
├── .phpier.yml                    # Project configuration (PHP version, settings)
├── .phpier/                       # PHPier-generated files directory
│   ├── docker-compose.yml         # Project container orchestration
│   ├── Dockerfile.php             # Custom PHP container with your version
│   └── docker/                    # Container configuration files
│       ├── nginx/                 # Nginx web server config
│       │   └── default.conf       # Site configuration
│       ├── php/                   # PHP runtime configuration
│       │   └── php.ini            # PHP settings (memory, upload limits, etc.)
│       └── supervisor/            # Process management
│           └── supervisord.conf   # PHP-FPM and service supervision
└── public/                        # Web root directory (auto-created)
    └── index.php                  # Default PHP info page
```

**Key Files:**
- **`.phpier.yml`**: Project configuration including PHP version, database settings, and container options
- **`.phpier/docker-compose.yml`**: Defines your project's app container and connects to global services
- **`.phpier/Dockerfile.php`**: Custom PHP container built with your chosen version and extensions  
- **`.phpier/docker/`**: All container configuration files (Nginx, PHP, Supervisor)
- **`public/index.php`**: Default landing page showing PHP info and environment details

#### Project Information
```bash
phpier list                  # List all discovered phpier projects
phpier services              # Show status of all phpier services
phpier services --project myapp  # Show services for specific project
phpier services --type app   # Filter by service type (app, db, cache, proxy, tools)
phpier services --status running  # Filter by status
phpier services --json       # Output in JSON format
phpier logs                  # View logs from project containers
```

### Database Access

#### Direct Database Connections
```bash
phpier mysql                 # Connect to MySQL database shell
phpier postgres              # Connect to PostgreSQL database shell
phpier psql                  # Connect to PostgreSQL database shell (same as postgres)
phpier maria                 # Connect to MariaDB database shell
phpier mariadb               # Connect to MariaDB database shell (same as maria)
```

#### Database Commands with Arguments
```bash
phpier mysql -e "SHOW TABLES"           # Execute MySQL query
phpier postgres -c "SELECT version();"  # Execute PostgreSQL query
phpier maria -e "SHOW DATABASES"        # Execute MariaDB query
```

### Caching Services
```bash
phpier redis                 # Execute Redis CLI commands
phpier memcached             # Connect to Memcached via telnet
```

### Database Management
```bash
phpier db                    # Manage database services
```

### Container Access

#### Shell Access
```bash
phpier sh                    # Open interactive shell in app container
```

#### Tool Proxying
```bash
# Context-aware tool execution:
phpier proxy <tool> [args...]           # In project directory
phpier proxy <app> <tool> [args...]     # From anywhere

# Examples:
phpier proxy composer install --no-dev  # Run Composer with flags
phpier proxy php -v                     # Show PHP version
phpier proxy npm run dev -- --watch     # Run npm with arguments  
phpier proxy php artisan migrate        # Laravel Artisan commands
phpier proxy myapp composer require phpunit/phpunit  # Global context
```

### Utility Commands
```bash
phpier version               # Show version information
phpier help                  # Show help for phpier
phpier [command] --help      # Show help for specific command
```

### Global Flags
```bash
--config string              # Config file (default: $HOME/.phpier.yml)
-v, --verbose               # Verbose output
-h, --help                  # Help for any command
```

### Examples
```bash
# Laravel project
mkdir my-laravel-app && cd my-laravel-app
phpier init 8.3 && phpier start && phpier up -d
# Visit: http://my-laravel-app.localhost

# Legacy project  
mkdir legacy-app && cd legacy-app
phpier init 7.4 --project-name=legacy && phpier up -d
# Visit: http://legacy.localhost
```

## 🤝 Contributing

See [DEV.md](DEV.md) for development information and contribution guidelines.

## 📄 License

[MIT License](LICENSE)

## 🙋 Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/your-org/phpier/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/your-org/phpier/discussions)

---

**Happy coding with PHPier! 🚀**