#!/bin/bash

# Quick local build script for development
# This script builds phpier for the current platform only

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

BINARY_NAME="phpier"
VERSION=${1:-"dev-$(date +%Y%m%d-%H%M%S)"}

log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

# Check if we're in the right directory
if [[ ! -f "main.go" ]]; then
    echo "❌ main.go not found. Please run this script from the phpier root directory."
    exit 1
fi

log "Building phpier locally..."

# Get build info
COMMIT="unknown"
DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

if command -v git &> /dev/null && [[ -d .git ]]; then
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
fi

# Build flags
LDFLAGS="-X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE"

# Build binary
go build -ldflags="$LDFLAGS" -o "$BINARY_NAME" main.go

if [[ -f "$BINARY_NAME" ]]; then
    SIZE=$(ls -lh "$BINARY_NAME" | awk '{print $5}')
    log "Built $BINARY_NAME ($SIZE) ✓"
    
    echo
    echo -e "${BLUE}Binary built successfully!${NC}"
    echo "  File: ./$BINARY_NAME"
    echo "  Version: $VERSION"
    echo "  Commit: $COMMIT"
    echo "  Size: $SIZE"
    echo
    echo "Test it with: ./$BINARY_NAME --help"
else
    echo "❌ Failed to build $BINARY_NAME"
    exit 1
fi