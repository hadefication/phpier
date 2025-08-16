#!/bin/bash

# phpier Release Script
# This script builds and releases phpier for multiple platforms

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="phpier"
BUILD_DIR="dist"
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

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

usage() {
    cat << EOF
Usage: $0 [OPTIONS] <version>

Build and release phpier for multiple platforms.

Arguments:
    version     Release version (e.g., v1.0.0, v1.2.3-beta.1)

Options:
    -h, --help          Show this help message
    -d, --dry-run       Build binaries but don't create release
    -c, --checksums     Generate checksums for binaries
    -z, --zip           Create zip archives for binaries
    --no-clean          Don't clean build directory before building
    --github            Create GitHub release (requires gh CLI)
    --local-only        Build for current platform only

Examples:
    $0 v1.0.0                    # Build and release v1.0.0
    $0 v1.0.0 --dry-run          # Build v1.0.0 without releasing
    $0 v1.0.0 --checksums --zip  # Build with checksums and zip files
    $0 v1.0.0 --local-only       # Build for current platform only

EOF
}

check_requirements() {
    log "Checking requirements..."
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        error "Go is not installed. Please install Go 1.20 or later."
    fi
    
    # Check Go version
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [[ "$GO_VERSION" < "1.20" ]]; then
        error "Go version $GO_VERSION is too old. Please install Go 1.20 or later."
    fi
    
    # Check if we're in the right directory
    if [[ ! -f "main.go" ]]; then
        error "main.go not found. Please run this script from the phpier root directory."
    fi
    
    if [[ ! -f "go.mod" ]]; then
        error "go.mod not found. Please run this script from the phpier root directory."
    fi
    
    # Check if git is available and repository is clean
    if command -v git &> /dev/null; then
        if [[ -n $(git status --porcelain) ]]; then
            warn "Working directory is not clean. Consider committing changes before release."
        fi
    fi
    
    log "Requirements check passed ✓"
}

validate_version() {
    local version=$1
    
    # Check if version follows semantic versioning
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        error "Invalid version format: $version. Expected format: v1.2.3 or v1.2.3-beta.1"
    fi
    
    log "Version $version is valid ✓"
}

clean_build_dir() {
    if [[ "$NO_CLEAN" != "true" ]]; then
        log "Cleaning build directory..."
        rm -rf "$BUILD_DIR"
    fi
    mkdir -p "$BUILD_DIR"
}

get_ldflags() {
    local version=$1
    local commit=""
    local date=""
    
    if command -v git &> /dev/null && [[ -d .git ]]; then
        commit=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
        date=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    else
        commit="unknown"
        date=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    fi
    
    echo "-X main.version=$version -X main.commit=$commit -X main.date=$date"
}

build_binary() {
    local platform=$1
    local version=$2
    local ldflags=$3
    
    local os=$(echo $platform | cut -d'/' -f1)
    local arch=$(echo $platform | cut -d'/' -f2)
    local output_name="${BINARY_NAME}-${os}-${arch}"
    
    if [[ "$os" == "windows" ]]; then
        output_name="${output_name}.exe"
    fi
    
    local output_path="${BUILD_DIR}/${output_name}"
    
    log "Building $output_name..."
    
    GOOS=$os GOARCH=$arch CGO_ENABLED=0 go build \
        -ldflags="$ldflags" \
        -o "$output_path" \
        main.go
    
    if [[ ! -f "$output_path" ]]; then
        error "Failed to build $output_name"
    fi
    
    # Get file size
    local size=$(ls -lh "$output_path" | awk '{print $5}')
    log "Built $output_name ($size) ✓"
    
    echo "$output_path"
}

create_checksums() {
    log "Generating checksums..."
    
    local checksum_file="${BUILD_DIR}/checksums.txt"
    
    cd "$BUILD_DIR"
    
    # Generate SHA256 checksums
    if command -v sha256sum &> /dev/null; then
        sha256sum ${BINARY_NAME}-* > checksums.txt
    elif command -v shasum &> /dev/null; then
        shasum -a 256 ${BINARY_NAME}-* > checksums.txt
    else
        warn "No checksum utility found. Skipping checksums."
        cd ..
        return
    fi
    
    cd ..
    
    log "Checksums generated in $checksum_file ✓"
}

create_archives() {
    log "Creating zip archives..."
    
    cd "$BUILD_DIR"
    
    for binary in ${BINARY_NAME}-*; do
        # Skip already compressed files
        if [[ "$binary" == *.zip ]] || [[ "$binary" == *.txt ]]; then
            continue
        fi
        
        local archive_name="${binary}.zip"
        zip -q "$archive_name" "$binary"
        
        if [[ -f "$archive_name" ]]; then
            log "Created $archive_name ✓"
        else
            warn "Failed to create $archive_name"
        fi
    done
    
    cd ..
}

create_github_release() {
    local version=$1
    
    if ! command -v gh &> /dev/null; then
        error "GitHub CLI (gh) is not installed. Please install it to create GitHub releases."
    fi
    
    # Check if logged in to GitHub
    if ! gh auth status &> /dev/null; then
        error "Not authenticated with GitHub. Please run 'gh auth login' first."
    fi
    
    log "Creating GitHub release $version..."
    
    # Create release notes
    local release_notes_file=$(mktemp)
    cat > "$release_notes_file" << EOF
# phpier $version

## What's Changed

- Bug fixes and improvements
- Updated dependencies
- Performance enhancements

## Installation

### Download Binary
Download the appropriate binary for your platform from the assets below.

### Build from Source
\`\`\`bash
git clone <repository-url>
cd phpier
go build -o phpier main.go
\`\`\`

### Verify Checksums
\`\`\`bash
# Download checksums.txt and verify
sha256sum -c checksums.txt
\`\`\`

## Full Changelog
**Full Changelog**: https://github.com/your-org/phpier/compare/v1.0.0...$version
EOF
    
    # Create the release
    gh release create "$version" \
        --title "Release $version" \
        --notes-file "$release_notes_file" \
        "$BUILD_DIR"/*
    
    rm "$release_notes_file"
    
    log "GitHub release $version created ✓"
}

show_summary() {
    local version=$1
    
    echo
    log "Release build completed! 🎉"
    echo
    echo -e "${BLUE}Version:${NC} $version"
    echo -e "${BLUE}Build directory:${NC} $BUILD_DIR"
    echo
    echo -e "${BLUE}Built binaries:${NC}"
    
    cd "$BUILD_DIR"
    for file in ${BINARY_NAME}-*; do
        if [[ -f "$file" ]]; then
            local size=$(ls -lh "$file" | awk '{print $5}')
            echo "  - $file ($size)"
        fi
    done
    cd ..
    
    echo
    if [[ "$DRY_RUN" == "true" ]]; then
        echo -e "${YELLOW}This was a dry run. No release was created.${NC}"
    else
        echo -e "${GREEN}Release assets are ready in $BUILD_DIR/${NC}"
    fi
    
    echo
    echo "Next steps:"
    echo "  1. Test the binaries on different platforms"
    echo "  2. Update documentation if needed"
    echo "  3. Announce the release"
}

# Parse command line arguments
VERSION=""
DRY_RUN="false"
GENERATE_CHECKSUMS="false"
CREATE_ARCHIVES="false"
NO_CLEAN="false"
GITHUB_RELEASE="false"
LOCAL_ONLY="false"

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -d|--dry-run)
            DRY_RUN="true"
            shift
            ;;
        -c|--checksums)
            GENERATE_CHECKSUMS="true"
            shift
            ;;
        -z|--zip)
            CREATE_ARCHIVES="true"
            shift
            ;;
        --no-clean)
            NO_CLEAN="true"
            shift
            ;;
        --github)
            GITHUB_RELEASE="true"
            shift
            ;;
        --local-only)
            LOCAL_ONLY="true"
            shift
            ;;
        -*)
            error "Unknown option: $1"
            ;;
        *)
            if [[ -z "$VERSION" ]]; then
                VERSION="$1"
            else
                error "Too many arguments. Version already set to $VERSION"
            fi
            shift
            ;;
    esac
done

# Validate arguments
if [[ -z "$VERSION" ]]; then
    error "Version is required. Use $0 --help for usage information."
fi

# Main execution
main() {
    log "Starting phpier release process..."
    
    check_requirements
    validate_version "$VERSION"
    clean_build_dir
    
    local ldflags=$(get_ldflags "$VERSION")
    local built_binaries=()
    
    if [[ "$LOCAL_ONLY" == "true" ]]; then
        # Build for current platform only
        local current_os=$(go env GOOS)
        local current_arch=$(go env GOARCH)
        local current_platform="${current_os}/${current_arch}"
        
        log "Building for current platform only: $current_platform"
        local binary_path=$(build_binary "$current_platform" "$VERSION" "$ldflags")
        built_binaries+=("$binary_path")
    else
        # Build for all platforms
        log "Building for ${#PLATFORMS[@]} platforms..."
        
        for platform in "${PLATFORMS[@]}"; do
            local binary_path=$(build_binary "$platform" "$VERSION" "$ldflags")
            built_binaries+=("$binary_path")
        done
    fi
    
    # Generate checksums if requested
    if [[ "$GENERATE_CHECKSUMS" == "true" ]]; then
        create_checksums
    fi
    
    # Create archives if requested
    if [[ "$CREATE_ARCHIVES" == "true" ]]; then
        create_archives
    fi
    
    # Create GitHub release if requested and not dry run
    if [[ "$GITHUB_RELEASE" == "true" ]] && [[ "$DRY_RUN" == "false" ]]; then
        create_github_release "$VERSION"
    fi
    
    show_summary "$VERSION"
}

# Run main function
main