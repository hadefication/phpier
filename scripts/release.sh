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
)

# Auto-detect repository from git remote
if command -v git >/dev/null 2>&1 && [[ -d .git ]]; then
    REPO_URL=$(git remote get-url origin 2>/dev/null || echo "")
    if [[ "$REPO_URL" =~ github\.com[/:]([^/]+)/([^/.]+) ]]; then
        REPO="${BASH_REMATCH[1]}/${BASH_REMATCH[2]}"
    else
        REPO="hadefication/phpier"  # Fallback
    fi
else
    REPO="hadefication/phpier"  # Fallback
fi

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

Build and release phpier for multiple platforms with full automation.

Arguments:
    version     Release version (e.g., v1.0.0, v1.2.3-beta.1)

Options:
    -h, --help          Show this help message
    -d, --dry-run       Build binaries but don't create release or git operations
    -c, --checksums     Generate checksums for binaries (default: true)
    -z, --zip           Create zip archives for binaries (default: true)
    --no-clean          Don't clean build directory before building
    --github            Create GitHub release (requires gh CLI) (default: true)
    --local-only        Build for current platform only
    --skip-tests        Skip running tests before release
    --skip-git          Skip git operations (tagging, pushing)
    --auto-commit       Automatically commit version changes
    --force             Force release even with uncommitted changes

Examples:
    $0 v1.0.0                           # Full automated release
    $0 v1.0.0 --dry-run                 # Build and test without releasing
    $0 v1.0.0 --skip-git --local-only   # Build locally without git operations
    $0 v1.0.0 --force --auto-commit     # Force release with auto-commit

Automated Release Process:
    1. Validate version and check git status
    2. Run tests and linting
    3. Update version in source files
    4. Build cross-platform binaries
    5. Generate checksums and archives
    6. Commit version changes and create git tag
    7. Push to GitHub and create release

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
    
    # Note: No Windows .exe needed since we only support WSL
    
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

check_git_status() {
    log "Checking git repository status..."
    
    # Check if we're in a git repository
    if ! git rev-parse --git-dir >/dev/null 2>&1; then
        error "Not in a git repository. Please run this script from the phpier root directory."
    fi
    
    # Check for uncommitted changes
    if [[ -n $(git status --porcelain) ]]; then
        if [[ "$FORCE_RELEASE" != "true" ]]; then
            warn "Working directory has uncommitted changes:"
            git status --short
            echo
            if [[ "$AUTO_COMMIT" == "true" ]]; then
                log "Auto-committing changes before release..."
                git add .
                git commit -m "Prepare release $VERSION"
            else
                echo "Options:"
                echo "  1. Commit changes manually and retry"
                echo "  2. Use --auto-commit to commit automatically"  
                echo "  3. Use --force to proceed anyway"
                echo "  4. Use --dry-run to test without releasing"
                exit 1
            fi
        else
            warn "Proceeding with uncommitted changes (--force)"
        fi
    fi
    
    # Check if tag already exists
    if git rev-parse "$VERSION" >/dev/null 2>&1; then
        if [[ "$FORCE_RELEASE" != "true" ]]; then
            error "Git tag $VERSION already exists. Use --force to overwrite."
        else
            warn "Git tag $VERSION exists, will be overwritten (--force)"
            git tag -d "$VERSION" || true
        fi
    fi
    
    log "Git repository status check passed ✓"
}

run_tests() {
    if [[ "$SKIP_TESTS" == "true" ]]; then
        warn "Skipping tests (--skip-tests)"
        return 0
    fi
    
    log "Running tests and linting..."
    
    # Run Go tests
    if ! go test -v ./...; then
        error "Tests failed. Fix tests before releasing."
    fi
    
    # Run Go vet
    if ! go vet ./...; then
        error "go vet failed. Fix issues before releasing."
    fi
    
    # Run Go fmt check
    if [[ -n $(gofmt -l .) ]]; then
        error "Code is not properly formatted. Run 'go fmt ./...' before releasing."
    fi
    
    log "Tests and linting passed ✓"
}

update_version_files() {
    local version=$1
    
    if [[ "$DRY_RUN" == "true" ]]; then
        log "Would update version to $version in source files (dry run)"
        return 0
    fi
    
    log "Updating version in source files..."
    
    # Update version in main.go (if it exists)
    if [[ -f "main.go" ]]; then
        # Look for version variable and update it
        if grep -q 'var version' main.go; then
            sed -i.bak "s/var version = \".*\"/var version = \"$version\"/" main.go
            rm -f main.go.bak
            log "Updated version in main.go ✓"
        fi
    fi
    
    # Update version in cmd/version.go (if it exists)
    if [[ -f "cmd/version.go" ]]; then
        if grep -q 'Version.*=' cmd/version.go; then
            sed -i.bak "s/Version.*=.*/Version = \"$version\"/" cmd/version.go
            rm -f cmd/version.go.bak
            log "Updated version in cmd/version.go ✓"
        fi
    fi
    
    log "Version files updated ✓"
}

create_git_tag() {
    local version=$1
    
    if [[ "$DRY_RUN" == "true" ]] || [[ "$SKIP_GIT" == "true" ]]; then
        log "Would create git tag $version (dry run or skip-git)"
        return 0
    fi
    
    log "Creating git tag $version..."
    
    # Commit version changes if any
    if [[ -n $(git status --porcelain) ]]; then
        log "Committing version update..."
        git add .
        git commit -m "Release $version"
    fi
    
    # Create annotated tag
    git tag -a "$version" -m "Release $version"
    
    log "Git tag $version created ✓"
}

push_to_github() {
    local version=$1
    
    if [[ "$DRY_RUN" == "true" ]] || [[ "$SKIP_GIT" == "true" ]]; then
        log "Would push to GitHub (dry run or skip-git)"
        return 0
    fi
    
    log "Pushing to GitHub..."
    
    # Push commits
    if ! git push; then
        error "Failed to push commits to GitHub"
    fi
    
    # Push tag
    if ! git push origin "$version"; then
        error "Failed to push tag $version to GitHub"
    fi
    
    log "Pushed to GitHub ✓"
}

show_summary() {
    local version=$1
    
    echo
    log "Release process completed! 🎉"
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
        echo
        echo "To create the actual release, run:"
        echo "  $0 $version --checksums --zip --github"
    else
        echo -e "${GREEN}Release completed successfully!${NC}"
        echo
        echo "What happened:"
        if [[ "$SKIP_TESTS" != "true" ]]; then
            echo "  ✅ Tests and linting passed"
        fi
        echo "  ✅ Built binaries for all platforms"
        echo "  ✅ Generated checksums and archives"
        if [[ "$SKIP_GIT" != "true" ]]; then
            echo "  ✅ Created git tag and pushed to GitHub"
        fi
        if [[ "$GITHUB_RELEASE" == "true" ]]; then
            echo "  ✅ Created GitHub release with assets"
        fi
        echo
        echo "Next steps:"
        echo "  1. Verify the GitHub release page"
        echo "  2. Test installation with: curl -sSL https://raw.githubusercontent.com/$REPO/main/scripts/install.sh | bash"
        echo "  3. Update documentation if needed"
        echo "  4. Announce the release"
    fi
}

# Parse command line arguments
VERSION=""
DRY_RUN="false"
GENERATE_CHECKSUMS="true"   # Default to true for full releases
CREATE_ARCHIVES="true"      # Default to true for full releases
NO_CLEAN="false"
GITHUB_RELEASE="true"       # Default to true for full releases
LOCAL_ONLY="false"
SKIP_TESTS="false"
SKIP_GIT="false"
AUTO_COMMIT="false"
FORCE_RELEASE="false"

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
        --skip-tests)
            SKIP_TESTS="true"
            shift
            ;;
        --skip-git)
            SKIP_GIT="true"
            shift
            ;;
        --auto-commit)
            AUTO_COMMIT="true"
            shift
            ;;
        --force)
            FORCE_RELEASE="true"
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
    log "Starting automated phpier release process..."
    
    # Step 1: Validate requirements and version
    check_requirements
    validate_version "$VERSION"
    
    # Step 2: Check git status and handle uncommitted changes
    check_git_status
    
    # Step 3: Run tests and linting
    run_tests
    
    # Step 4: Update version in source files
    update_version_files "$VERSION"
    
    # Step 5: Build binaries
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
    
    # Step 6: Generate checksums and archives
    if [[ "$GENERATE_CHECKSUMS" == "true" ]]; then
        create_checksums
    fi
    
    if [[ "$CREATE_ARCHIVES" == "true" ]]; then
        create_archives
    fi
    
    # Step 7: Create git tag and push to GitHub
    create_git_tag "$VERSION"
    push_to_github "$VERSION"
    
    # Step 8: Create GitHub release
    if [[ "$GITHUB_RELEASE" == "true" ]] && [[ "$DRY_RUN" == "false" ]]; then
        # Wait a moment for the tag to be available on GitHub
        log "Waiting for tag to be available on GitHub..."
        sleep 3
        create_github_release "$VERSION"
    fi
    
    show_summary "$VERSION"
}

# Run main function
main