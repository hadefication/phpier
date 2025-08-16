#!/bin/bash

# Local uninstallation script for phpier CLI
# This script removes the phpier binary from /usr/local/bin

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
BINARY_FULL_PATH="${INSTALL_PATH}/${BINARY_NAME}"

echo -e "${BLUE}🗑️  Uninstalling phpier CLI...${NC}"

# Check if binary exists
if [ ! -f "${BINARY_FULL_PATH}" ]; then
    echo -e "${YELLOW}⚠️  phpier not found in ${INSTALL_PATH}${NC}"
    echo -e "${YELLOW}   Binary may not be installed or may be in a different location${NC}"
    
    # Check if phpier is available in PATH but not in expected location
    if command -v phpier &> /dev/null; then
        CURRENT_LOCATION=$(which phpier)
        echo -e "${BLUE}📍 Found phpier at: ${CURRENT_LOCATION}${NC}"
        read -p "Remove phpier from ${CURRENT_LOCATION}? (y/N): " -r
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            if [ ! -w "$(dirname "${CURRENT_LOCATION}")" ]; then
                sudo rm "${CURRENT_LOCATION}"
            else
                rm "${CURRENT_LOCATION}"
            fi
            echo -e "${GREEN}✅ Removed phpier from ${CURRENT_LOCATION}${NC}"
        else
            echo -e "${YELLOW}❌ Uninstallation cancelled${NC}"
        fi
    else
        echo -e "${GREEN}✅ phpier is not installed${NC}"
    fi
    exit 0
fi

# Check if we need sudo for removal
if [ ! -w "${INSTALL_PATH}" ]; then
    echo -e "${YELLOW}🔒 Removal requires sudo privileges${NC}"
    SUDO_CMD="sudo"
else
    SUDO_CMD=""
fi

# Show current version before removal
if [ -x "${BINARY_FULL_PATH}" ]; then
    echo -e "${BLUE}📋 Current version:${NC}"
    CURRENT_VERSION=$("${BINARY_FULL_PATH}" version 2>/dev/null || "${BINARY_FULL_PATH}" --version 2>/dev/null || echo "Unknown version")
    echo -e "${BLUE}   ${CURRENT_VERSION}${NC}"
fi

# Confirm removal
read -p "Are you sure you want to remove phpier? (y/N): " -r
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}❌ Uninstallation cancelled${NC}"
    exit 0
fi

# Remove binary
echo -e "${BLUE}🗑️  Removing binary...${NC}"
${SUDO_CMD} rm "${BINARY_FULL_PATH}"

# Verify removal
if [ ! -f "${BINARY_FULL_PATH}" ]; then
    echo -e "${GREEN}✅ Successfully removed phpier from ${INSTALL_PATH}${NC}"
    
    # Check if phpier is still available (might be in another location)
    if command -v phpier &> /dev/null; then
        REMAINING_LOCATION=$(which phpier)
        echo -e "${YELLOW}⚠️  phpier is still available at: ${REMAINING_LOCATION}${NC}"
        echo -e "${YELLOW}   You may have multiple installations${NC}"
    else
        echo -e "${GREEN}✅ phpier completely removed from system${NC}"
    fi
    
    echo -e "${GREEN}🎉 Uninstallation complete!${NC}"
else
    echo -e "${RED}❌ Uninstallation failed - binary still exists${NC}"
    exit 1
fi