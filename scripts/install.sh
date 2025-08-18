#!/bin/bash

# phpier Installation Script
# This script downloads and installs the latest version of phpier
# Usage: curl -sSL https://raw.githubusercontent.com/your-org/phpier/main/scripts/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="phpier"
REPO="hadefication/phpier"
INSTALL_DIR="$HOME/.local/bin"
TEMP_DIR="/tmp/phpier-install"

# Functions
log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

banner() {
    echo -e "${PURPLE}"
    cat << 'EOF'
    ____  __  ______     _           
   / __ \/ / / / __ \   (_)__  _____
  / /_/ / /_/ / /_/ /  / / _ \/ ___/
 / ____/ __  / ____/  / /  __/ /    
/_/   /_/ /_/_/      /_/\___/_/     
                                   
EOF
    echo -e "${NC}"
    echo -e "${BLUE}PHPier - PHP Development Environment Manager${NC}"
    echo
}

usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Install phpier - A CLI tool for managing PHP development environments.
Supports Linux, macOS, and WSL (Windows Subsystem for Linux).

Options:
    -v, --version VERSION    Install specific version (e.g., v1.0.0)
    -d, --dir DIRECTORY      Install directory (default: ~/.local/bin)
    -h, --help              Show this help message
    --no-verify             Skip checksum verification
    --force                 Force installation even if already installed

Examples:
    # Install latest version
    curl -sSL https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh | bash

    # Install specific version
    curl -sSL https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh | bash -s -- -v v1.0.0

    # Install to custom directory
    curl -sSL https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh | bash -s -- -d /usr/local/bin

EOF
}

detect_platform() {
    local os arch
    
    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        MINGW*|CYGWIN*|MSYS*) 
            error "Native Windows is not supported. Please use WSL (Windows Subsystem for Linux) instead." ;;
        *)          error "Unsupported operating system: $(uname -s). Supported: Linux, macOS, WSL" ;;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        arm64|aarch64)  arch="arm64" ;;
        armv7l)         arch="arm" ;;
        i386|i686)      arch="386" ;;
        *)              error "Unsupported architecture: $(uname -m)" ;;
    esac
    
    echo "${os}/${arch}"
}

get_latest_version() {
    log "Fetching latest release information..."
    
    local latest_url="https://api.github.com/repos/${REPO}/releases/latest"
    local version
    
    # Try different methods to get version
    if command -v curl >/dev/null 2>&1; then
        version=$(curl -s "$latest_url" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        version=$(wget -qO- "$latest_url" | grep '"tag_name":' | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
    
    if [[ -z "$version" ]]; then
        error "Could not fetch latest version. Please check your internet connection."
    fi
    
    echo "$version"
}

check_existing_installation() {
    local install_path="$1"
    
    if [[ -f "$install_path" ]]; then
        if [[ "$FORCE_INSTALL" != "true" ]]; then
            local existing_version
            existing_version=$("$install_path" version 2>/dev/null | head -n1 || echo "unknown")
            
            echo
            warn "phpier is already installed at: $install_path"
            echo "  Current version: $existing_version"
            echo "  Target version: $VERSION"
            echo
            echo "Options:"
            echo "  1. Continue and overwrite (y)"
            echo "  2. Cancel installation (n)"
            echo "  3. Use --force flag to skip this prompt"
            echo
            read -p "Continue with installation? (y/n): " -n 1 -r
            echo
            
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                log "Installation cancelled."
                exit 0
            fi
        fi
    fi
}

download_binary() {
    local version="$1"
    local platform="$2"
    local temp_dir="$3"
    
    local os arch binary_name download_url
    IFS='/' read -r os arch <<< "$platform"
    
    binary_name="${BINARY_NAME}-${os}-${arch}"
    # Note: Windows .exe extension not needed since we only support WSL
    
    download_url="https://github.com/${REPO}/releases/download/${version}/${binary_name}"
    local binary_path="${temp_dir}/${binary_name}"
    
    log "Downloading ${binary_name}..."
    log "URL: $download_url"
    
    # Download binary
    if command -v curl >/dev/null 2>&1; then
        if ! curl -L -o "$binary_path" "$download_url"; then
            error "Failed to download binary from $download_url"
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget -O "$binary_path" "$download_url"; then
            error "Failed to download binary from $download_url"
        fi
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
    
    if [[ ! -f "$binary_path" ]]; then
        error "Downloaded binary not found at $binary_path"
    fi
    
    # Verify download
    local size
    size=$(ls -lh "$binary_path" | awk '{print $5}' 2>/dev/null)
    log "Downloaded $binary_name ($size) ✓"
    
    echo "$binary_path"
}

verify_checksum() {
    local version="$1"
    local binary_path="$2"
    local temp_dir="$3"
    
    if [[ "$VERIFY_CHECKSUM" == "false" ]]; then
        warn "Skipping checksum verification (--no-verify)"
        return 0
    fi
    
    log "Verifying checksum..."
    
    local checksums_url="https://github.com/${REPO}/releases/download/${version}/checksums.txt"
    local checksums_path="${temp_dir}/checksums.txt"
    
    # Download checksums
    if command -v curl >/dev/null 2>&1; then
        curl -sL -o "$checksums_path" "$checksums_url" || {
            warn "Could not download checksums.txt, skipping verification"
            return 0
        }
    elif command -v wget >/dev/null 2>&1; then
        wget -qO "$checksums_path" "$checksums_url" || {
            warn "Could not download checksums.txt, skipping verification"
            return 0
        }
    fi
    
    # Verify checksum
    local binary_name
    binary_name=$(basename "$binary_path")
    local expected_checksum actual_checksum
    
    expected_checksum=$(grep "$binary_name" "$checksums_path" | awk '{print $1}' 2>/dev/null)
    
    if [[ -z "$expected_checksum" ]]; then
        warn "Checksum not found for $binary_name, skipping verification"
        return 0
    fi
    
    if command -v sha256sum >/dev/null 2>&1; then
        actual_checksum=$(sha256sum "$binary_path" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        actual_checksum=$(shasum -a 256 "$binary_path" | awk '{print $1}')
    else
        warn "No checksum utility found, skipping verification"
        return 0
    fi
    
    if [[ "$expected_checksum" == "$actual_checksum" ]]; then
        log "Checksum verification passed ✓"
    else
        error "Checksum verification failed!"
    fi
}

install_binary() {
    local binary_path="$1"
    local install_dir="$2"
    
    log "Installing phpier to $install_dir..."
    
    # Create install directory
    mkdir -p "$install_dir"
    
    # Install binary
    local install_path="${install_dir}/${BINARY_NAME}"
    cp "$binary_path" "$install_path"
    chmod +x "$install_path"
    
    if [[ ! -f "$install_path" ]]; then
        error "Installation failed: binary not found at $install_path"
    fi
    
    # Test installation
    if ! "$install_path" version >/dev/null 2>&1; then
        error "Installation failed: binary is not executable or corrupted"
    fi
    
    echo "$install_path"
}

setup_path() {
    local install_dir="$1"
    
    # Check if install directory is in PATH
    case ":$PATH:" in
        *":$install_dir:"*) 
            log "Installation directory is already in PATH ✓"
            return 0
            ;;
    esac
    
    warn "Installation directory is not in PATH: $install_dir"
    echo
    echo "To use phpier from anywhere, add this to your shell profile:"
    echo
    echo -e "${BLUE}export PATH=\"$install_dir:\$PATH\"${NC}"
    echo
    
    # Detect shell and suggest profile file
    local shell_profile=""
    case "$SHELL" in
        */bash) shell_profile="$HOME/.bashrc or $HOME/.bash_profile" ;;
        */zsh)  shell_profile="$HOME/.zshrc" ;;
        */fish) shell_profile="$HOME/.config/fish/config.fish" ;;
        *)      shell_profile="your shell's configuration file" ;;
    esac
    
    echo "Add the above line to $shell_profile, then restart your terminal or run:"
    echo -e "${BLUE}source $shell_profile${NC}"
    echo
}

cleanup() {
    if [[ -d "$TEMP_DIR" ]]; then
        rm -rf "$TEMP_DIR"
    fi
}

show_completion() {
    local version="$1"
    local install_path="$2"
    
    echo
    success "phpier $version installed successfully! 🎉"
    echo
    echo -e "${BLUE}Installation details:${NC}"
    echo "  Binary: $install_path"
    echo "  Version: $version"
    local size
    size=$(ls -lh "$install_path" | awk '{print $5}' 2>/dev/null)
    echo "  Size: $size"
    echo
    echo -e "${BLUE}Quick start:${NC}"
    echo "  phpier --help              # Show help"
    echo "  phpier version             # Show version"
    echo "  phpier init 8.3            # Initialize PHP 8.3 project"
    echo
    echo -e "${BLUE}Documentation:${NC}"
    echo "  https://github.com/$REPO"
    echo
}

# Parse command line arguments
VERSION=""
INSTALL_DIR="$HOME/.local/bin"
VERIFY_CHECKSUM="true"
FORCE_INSTALL="false"

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -d|--dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        --no-verify)
            VERIFY_CHECKSUM="false"
            shift
            ;;
        --force)
            FORCE_INSTALL="true"
            shift
            ;;
        -*)
            error "Unknown option: $1"
            ;;
        *)
            error "Unexpected argument: $1"
            ;;
    esac
done

# Main installation process
main() {
    # Set trap for cleanup
    trap cleanup EXIT
    
    banner
    
    log "Starting phpier installation..."
    
    # Detect platform
    local platform
    platform=$(detect_platform)
    log "Detected platform: $platform"
    
    # Get version
    if [[ -z "$VERSION" ]]; then
        VERSION=$(get_latest_version)
    fi
    log "Target version: $VERSION"
    
    # Expand install directory
    INSTALL_DIR=$(eval echo "$INSTALL_DIR")
    log "Install directory: $INSTALL_DIR"
    
    # Check existing installation
    check_existing_installation "${INSTALL_DIR}/${BINARY_NAME}"
    
    # Create temporary directory
    rm -rf "$TEMP_DIR"
    mkdir -p "$TEMP_DIR"
    
    # Download binary
    local binary_path
    binary_path=$(download_binary "$VERSION" "$platform" "$TEMP_DIR")
    
    # Verify checksum
    verify_checksum "$VERSION" "$binary_path" "$TEMP_DIR"
    
    # Install binary
    local install_path
    install_path=$(install_binary "$binary_path" "$INSTALL_DIR")
    
    # Setup PATH
    setup_path "$INSTALL_DIR"
    
    # Show completion message
    show_completion "$VERSION" "$install_path"
}

# Run main function
main "$@"