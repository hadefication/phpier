# phpier Scripts

This directory contains utility scripts for building, installing, and maintaining the phpier CLI.

## 🔨 Local Installation

### Install Script

Use `local-install.sh` to build and install phpier locally:

```bash
# Run from the project root directory
./scripts/local-install.sh
```

**What it does:**
- Builds the phpier binary with version information from git
- Sets executable permissions on the binary
- Removes any existing phpier installation in `/usr/local/bin`
- Moves the new binary to `/usr/local/bin/phpier`
- Verifies the installation and tests the command
- Provides usage examples

**Features:**
- ✅ **Version embedding**: Includes git version, commit hash, and build date
- ✅ **Existing installation handling**: Automatically removes old versions
- ✅ **Permission management**: Handles sudo requirements automatically
- ✅ **Path verification**: Checks if install location is in PATH
- ✅ **Installation testing**: Verifies the installation works
- ✅ **Colored output**: User-friendly colored terminal output

**Example output:**
```
🔨 Building phpier CLI...
📦 Version: v1.0.0
📦 Commit: abc12345
📦 Date: 2025-08-15T15:30:00Z
⚙️  Compiling binary...
🔐 Setting executable permissions...
⚠️  Existing phpier installation found
   Current version: phpier dev
🗑️  Removing existing installation...
✅ Existing installation removed
📦 Installing to /usr/local/bin...
✅ Successfully installed phpier to /usr/local/bin/phpier
✅ /usr/local/bin is in your PATH
🧪 Testing installation...
✅ Installation test successful
📋 Installed version: phpier v1.0.0

🎉 Installation complete!
📚 Usage examples:
   phpier init 8.3
   phpier up
   phpier php -v
   phpier --help
```

### Uninstall Script

Use `local-uninstall.sh` to remove phpier from your system:

```bash
./scripts/local-uninstall.sh
```

**What it does:**
- Locates the phpier binary (in `/usr/local/bin` or elsewhere in PATH)
- Shows current version before removal
- Prompts for confirmation
- Removes the binary with appropriate permissions
- Verifies complete removal

## 🚀 Usage

### Quick Install
```bash
# Clone the repository
git clone <repository-url>
cd phpier

# Install phpier locally
./scripts/local-install.sh
```

### Upgrade Existing Installation
```bash
# Pull latest changes
git pull

# Reinstall (automatically removes old version)
./scripts/local-install.sh
```

### Complete Removal
```bash
# Remove phpier from system
./scripts/local-uninstall.sh
```

## 📋 Requirements

- **Go 1.20+**: Required for building the binary
- **Git**: Used for version information (optional)
- **sudo privileges**: May be required for installing to `/usr/local/bin`

## 🔧 Customization

You can customize the installation by modifying variables in the scripts:

```bash
# In local-install.sh
BINARY_NAME="phpier"        # Name of the binary
INSTALL_PATH="/usr/local/bin" # Installation directory
BUILD_PATH="./phpier"       # Temporary build location
```

## 🐛 Troubleshooting

### Permission Denied
If you get permission errors:
```bash
# Make scripts executable
chmod +x scripts/*.sh

# Or run with explicit bash
bash scripts/local-install.sh
```

### Installation Path Not in PATH
If `/usr/local/bin` is not in your PATH:
```bash
# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export PATH="/usr/local/bin:$PATH"

# Reload your shell
source ~/.bashrc  # or ~/.zshrc
```

### Multiple Installations
If you have phpier installed in multiple locations:
```bash
# Find all installations
which -a phpier

# Remove specific installation
sudo rm /path/to/phpier

# Or use the uninstall script
./scripts/local-uninstall.sh
```

## 📝 Development Notes

The install script automatically:
- Detects git repository information for versioning
- Handles different shell environments
- Manages sudo requirements intelligently
- Provides comprehensive error handling and user feedback

The build process includes:
- **Optimized binary**: Built with `-s -w` flags for smaller size
- **Version embedding**: Git version, commit, and build date
- **Cross-platform support**: Builds for the current platform

## 🔗 Related Commands

After installation, you can use these phpier commands:
```bash
phpier version          # Show version information
phpier init 8.3        # Initialize new project
phpier up              # Start development environment
phpier down            # Stop development environment
phpier php -v          # Run PHP commands
phpier --help          # Show help information
```