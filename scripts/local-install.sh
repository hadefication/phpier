#!/bin/bash

# Local installation script for phpier CLI
# This script builds the binary and installs it to /usr/local/bin

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="phpier"
INSTALL_PATH="/usr/local/bin"
BUILD_PATH="./phpier"

echo -e "${BLUE}🔨 Building phpier CLI...${NC}"

# Get version information from git if available
VERSION="dev"
COMMIT="unknown"
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

if command -v git &> /dev/null && git rev-parse --is-inside-work-tree &> /dev/null; then
    # Get version from git tags or use dev
    if git describe --tags --exact-match 2>/dev/null; then
        VERSION=$(git describe --tags --exact-match 2>/dev/null)
    elif git describe --tags 2>/dev/null; then
        VERSION=$(git describe --tags 2>/dev/null)
    else
        VERSION="dev-$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
    fi
    
    COMMIT=$(git rev-parse HEAD 2>/dev/null || echo 'unknown')
fi

echo -e "${YELLOW}📦 Version: ${VERSION}${NC}"
echo -e "${YELLOW}📦 Commit: ${COMMIT:0:8}${NC}"
echo -e "${YELLOW}📦 Date: ${DATE}${NC}"

# Check if we need sudo for installation
if [ ! -w "${INSTALL_PATH}" ]; then
    echo -e "${YELLOW}🔒 Installation requires sudo privileges${NC}"
    SUDO_CMD="sudo"
else
    SUDO_CMD=""
fi

# Step 1: Check if phpier already exists and remove it first
if [ -f "${INSTALL_PATH}/${BINARY_NAME}" ]; then
    echo -e "${YELLOW}⚠️  Existing phpier installation found${NC}"
    
    # Show current version if possible
    if [ -x "${INSTALL_PATH}/${BINARY_NAME}" ]; then
        CURRENT_VERSION=$("${INSTALL_PATH}/${BINARY_NAME}" version 2>/dev/null || "${INSTALL_PATH}/${BINARY_NAME}" --version 2>/dev/null || echo "Unknown version")
        echo -e "${YELLOW}   Current version: ${CURRENT_VERSION}${NC}"
    fi
    
    echo -e "${BLUE}🗑️  Removing existing installation...${NC}"
    ${SUDO_CMD} rm "${INSTALL_PATH}/${BINARY_NAME}"
    
    if [ -f "${INSTALL_PATH}/${BINARY_NAME}" ]; then
        echo -e "${RED}❌ Failed to remove existing installation${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ Existing installation removed${NC}"
fi

# Step 2: Build with version information
echo -e "${BLUE}⚙️  Compiling binary...${NC}"
go build -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" -o "${BUILD_PATH}"

if [ ! -f "${BUILD_PATH}" ]; then
    echo -e "${RED}❌ Build failed - binary not found${NC}"
    exit 1
fi

# Step 3: Add executable permissions
echo -e "${BLUE}🔐 Setting executable permissions...${NC}"
chmod +x "${BUILD_PATH}"

# Step 4: Move binary to /usr/local/bin
echo -e "${BLUE}📦 Installing to ${INSTALL_PATH}...${NC}"
${SUDO_CMD} mv "${BUILD_PATH}" "${INSTALL_PATH}/${BINARY_NAME}"

# Verify installation
if [ -f "${INSTALL_PATH}/${BINARY_NAME}" ]; then
    echo -e "${GREEN}✅ Successfully installed phpier to ${INSTALL_PATH}/${BINARY_NAME}${NC}"
    
    # Check if install path is in PATH
    if echo "$PATH" | grep -q "${INSTALL_PATH}"; then
        echo -e "${GREEN}✅ ${INSTALL_PATH} is in your PATH${NC}"
        
        # Test the installation
        echo -e "${BLUE}🧪 Testing installation...${NC}"
        if command -v phpier &> /dev/null; then
            VERSION_OUTPUT=$(phpier version 2>/dev/null || phpier --version 2>/dev/null || echo "Version command not available")
            echo -e "${GREEN}✅ Installation test successful${NC}"
            echo -e "${BLUE}📋 Installed version: ${VERSION_OUTPUT}${NC}"
        else
            echo -e "${YELLOW}⚠️  Command 'phpier' not found in PATH, try reloading your shell${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  ${INSTALL_PATH} is not in your PATH${NC}"
        echo -e "${YELLOW}   Add it to your shell profile:${NC}"
        echo -e "${YELLOW}   export PATH=\"${INSTALL_PATH}:\$PATH\"${NC}"
    fi
    
    echo ""
    echo -e "${GREEN}🎉 Installation complete!${NC}"
    echo -e "${BLUE}📚 Usage examples:${NC}"
    echo -e "   phpier init 8.3"
    echo -e "   phpier up"
    echo -e "   phpier php -v"
    echo -e "   phpier --help"
    
else
    echo -e "${RED}❌ Installation failed - binary not found in ${INSTALL_PATH}${NC}"
    exit 1
fi