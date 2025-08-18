# PHPier Installation Guide

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) (20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (2.0+)
- [Go](https://golang.org/dl/) (1.20+) - for building from source

## Installation Methods

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

### Option 5: Development Mode

```bash
# Run directly with Go
cd phpier
go run main.go [command]
```

## Uninstallation

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

## Verification

After installation, verify phpier is working:

```bash
phpier version
phpier --help
```