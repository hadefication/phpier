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

**Specific Version:**
```bash
curl -sSL https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh | bash -s -- -v v1.0.0
```

**Custom Directory:**
```bash
curl -sSL https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh | bash -s -- -d /usr/local/bin
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

### Service Management
```bash
phpier start                 # Start global services (Traefik, databases)
phpier stop                  # Stop global services
phpier up -d                 # Start project container
phpier down                  # Stop project container
```

### Project Setup
```bash
phpier init                  # Initialize with PHP 8.3 (default)
phpier init 7.4              # Initialize with specific PHP version
phpier init --project-name=myapp  # Custom project name
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
- 📖 **Documentation**: [Wiki](https://github.com/your-org/phpier/wiki)

---

**Happy coding with PHPier! 🚀**