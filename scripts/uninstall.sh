#!/bin/bash

# phpier Uninstallation Script
# This script removes phpier from your system

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
COMMON_INSTALL_PATHS=(
    "$HOME/.local/bin"
    "/usr/local/bin"
    "/opt/homebrew/bin"
    "/usr/bin"
    "$HOME/bin"
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

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

banner() {
    echo -e "${PURPLE}"
    cat << 'EOF'
 _   _       _           _        _ _ 
| | | |_ __ (_)_ __  ___| |_ __ _| | |
| | | | '_ \| | '_ \/ __| __/ _` | | |
| |_| | | | | | | | \__ \ || (_| | | |
 \___/|_| |_|_|_| |_|___/\__\__,_|_|_|
                                     
EOF
    echo -e "${NC}"
    echo -e "${BLUE}PHPier Uninstaller${NC}"
    echo
}

usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Uninstall phpier from your system.

Options:
    -h, --help              Show this help message
    -f, --force             Force removal without confirmation
    --dry-run               Show what would be removed without actually removing
    --path PATH             Only check specific path for phpier binary
    --all                   Remove all found phpier installations

Examples:
    # Interactive uninstall
    $0

    # Force uninstall without confirmation
    $0 --force

    # Show what would be removed
    $0 --dry-run

    # Remove from specific path
    $0 --path /usr/local/bin

EOF
}

find_phpier_installations() {
    local found_paths=()
    
    # Check common installation paths
    for path in "${COMMON_INSTALL_PATHS[@]}"; do
        local full_path="$path/$BINARY_NAME"
        if [[ -f "$full_path" ]] && [[ -x "$full_path" ]]; then
            # Verify it's actually phpier by checking version output
            if "$full_path" version &>/dev/null || "$full_path" --help | grep -q "phpier\|PHP.*development"; then
                found_paths+=("$full_path")
            fi
        fi
    done
    
    # Check if phpier is in PATH but not in common paths
    if command -v phpier &>/dev/null; then
        local which_path
        which_path=$(command -v phpier)
        
        # Check if this path is already in our found_paths
        local already_found=false
        for found_path in "${found_paths[@]}"; do
            if [[ "$found_path" == "$which_path" ]]; then
                already_found=true
                break
            fi
        done
        
        if [[ "$already_found" == "false" ]]; then
            found_paths+=("$which_path")
        fi
    fi
    
    # Return found paths as newline-separated string
    printf '%s\n' "${found_paths[@]}"
}

show_installation_info() {
    local install_path="$1"
    
    echo -e "${BLUE}Installation Details:${NC}"
    echo "  Path: $install_path"
    
    # Try to get version
    if version_output=$("$install_path" version 2>/dev/null); then
        echo "  Version: $version_output"
    else
        echo "  Version: Unable to determine"
    fi
    
    # Get file size
    if [[ -f "$install_path" ]]; then
        local size
        size=$(ls -lh "$install_path" 2>/dev/null | awk '{print $5}')
        echo "  Size: $size"
    fi
    
    # Check if it's in PATH
    if command -v phpier &>/dev/null && [[ "$(command -v phpier)" == "$install_path" ]]; then
        echo "  Status: In PATH (active)"
    else
        echo "  Status: Not in current PATH"
    fi
    
    echo
}

remove_phpier() {
    local install_path="$1"
    local dry_run="$2"
    
    if [[ "$dry_run" == "true" ]]; then
        log "Would remove: $install_path"
        return 0
    fi
    
    log "Removing phpier from: $install_path"
    
    # Check if we need sudo
    if [[ ! -w "$(dirname "$install_path")" ]]; then
        warn "Removing $install_path requires sudo privileges"
        if ! sudo rm -f "$install_path"; then
            error "Failed to remove $install_path (permission denied)"
        fi
    else
        if ! rm -f "$install_path"; then
            error "Failed to remove $install_path"
        fi
    fi
    
    # Verify removal
    if [[ -f "$install_path" ]]; then
        error "Failed to remove $install_path (file still exists)"
    fi
    
    success "Removed phpier from: $install_path"
}

cleanup_phpier_data() {
    local dry_run="$1"
    
    log "Checking for phpier data directories..."
    
    local data_paths=(
        "$HOME/.phpier"
        "$HOME/.config/phpier"
        "/tmp/phpier-*"
    )
    
    local found_data=false
    
    for pattern in "${data_paths[@]}"; do
        # Use find for glob patterns, direct check for specific paths
        if [[ "$pattern" =~ \* ]]; then
            # Handle glob patterns
            while IFS= read -r -d '' path; do
                if [[ "$dry_run" == "true" ]]; then
                    log "Would remove data: $path"
                else
                    log "Removing data: $path"
                    rm -rf "$path"
                fi
                found_data=true
            done < <(find /tmp -maxdepth 1 -name "phpier-*" -print0 2>/dev/null || true)
        else
            # Handle specific paths
            if [[ -e "$pattern" ]]; then
                if [[ "$dry_run" == "true" ]]; then
                    log "Would remove data: $pattern"
                else
                    log "Removing data: $pattern"
                    rm -rf "$pattern"
                fi
                found_data=true
            fi
        fi
    done
    
    if [[ "$found_data" == "false" ]]; then
        log "No phpier data directories found"
    fi
}

cleanup_phpier_docker() {
    local dry_run="$1"
    
    log "Checking for phpier Docker containers..."
    
    # Check if Docker is running
    if ! command -v docker &> /dev/null; then
        log "Docker not found, skipping Docker cleanup"
        return 0
    fi
    
    if ! docker info &> /dev/null; then
        warn "Docker daemon not running, skipping Docker cleanup"
        return 0
    fi
    
    # Stop and remove phpier containers
    local containers
    containers=$(docker ps -aq --filter "name=phpier" 2>/dev/null || true)
    
    if [[ -n "$containers" ]]; then
        if [[ "$dry_run" == "true" ]]; then
            log "Would stop and remove phpier containers:"
            docker ps -a --filter "name=phpier" --format "  - {{.Names}} ({{.Status}})" 2>/dev/null || true
        else
            log "Stopping and removing phpier containers..."
            echo "$containers" | xargs docker stop &>/dev/null || true
            echo "$containers" | xargs docker rm &>/dev/null || true
            success "Removed phpier containers"
        fi
    else
        log "No phpier containers found"
    fi
    
    # Remove phpier networks
    local networks
    networks=$(docker network ls --filter "name=phpier" -q 2>/dev/null || true)
    
    if [[ -n "$networks" ]]; then
        if [[ "$dry_run" == "true" ]]; then
            log "Would remove phpier networks:"
            docker network ls --filter "name=phpier" --format "  - {{.Name}}" 2>/dev/null || true
        else
            log "Removing phpier networks..."
            echo "$networks" | xargs docker network rm &>/dev/null || true
            success "Removed phpier networks"
        fi
    else
        log "No phpier networks found"
    fi
    
    # Remove phpier images (built project images)
    local images
    images=$(docker images --filter "reference=phpier-*" -q 2>/dev/null || true)
    
    if [[ -n "$images" ]]; then
        if [[ "$dry_run" == "true" ]]; then
            log "Would remove phpier project images:"
            docker images --filter "reference=phpier-*" --format "  - {{.Repository}}:{{.Tag}} ({{.Size}})" 2>/dev/null || true
        else
            log "Removing phpier project images..."
            echo "$images" | xargs docker rmi -f &>/dev/null || true
            success "Removed phpier project images"
        fi
    else
        log "No phpier project images found"
    fi
    
    # Note about preserved volumes
    local volumes
    volumes=$(docker volume ls --filter "name=phpier" -q 2>/dev/null || true)
    
    if [[ -n "$volumes" ]]; then
        log "Preserving phpier Docker volumes for fresh installation testing:"
        docker volume ls --filter "name=phpier" --format "  - {{.Name}}" 2>/dev/null || true
        echo
        warn "Docker volumes are preserved. To remove them manually if needed:"
        echo "  docker volume rm \$(docker volume ls --filter \"name=phpier\" -q)"
    else
        log "No phpier Docker volumes found"
    fi
}

show_summary() {
    local removed_count="$1"
    local dry_run="$2"
    
    echo
    if [[ "$dry_run" == "true" ]]; then
        log "Dry run completed - no files were actually removed"
        echo
        echo "To perform the actual uninstall, run:"
        echo "  $0 --force"
    elif [[ "$removed_count" -gt 0 ]]; then
        success "Uninstallation completed! 🎉"
        echo
        echo "Removed $removed_count phpier installation(s)"
        echo
        echo "To verify uninstallation:"
        echo "  command -v phpier  # Should return nothing"
        echo "  phpier --help      # Should show 'command not found'"
        echo
        echo "To reinstall phpier in the future:"
        echo "  curl -sSL https://raw.githubusercontent.com/hadefication/phpier/main/scripts/install.sh | bash"
    else
        warn "No phpier installations were removed"
    fi
}

# Parse command line arguments
FORCE="false"
DRY_RUN="false"
SPECIFIC_PATH=""
REMOVE_ALL="false"

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            usage
            exit 0
            ;;
        -f|--force)
            FORCE="true"
            shift
            ;;
        --dry-run)
            DRY_RUN="true"
            shift
            ;;
        --path)
            SPECIFIC_PATH="$2"
            shift 2
            ;;
        --all)
            REMOVE_ALL="true"
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

# Main execution
main() {
    banner
    
    log "Starting phpier uninstallation..."
    
    # Find installations
    local found_installations=()
    if [[ -n "$SPECIFIC_PATH" ]]; then
        # Check specific path
        local full_path="$SPECIFIC_PATH/$BINARY_NAME"
        if [[ -f "$full_path" ]] && [[ -x "$full_path" ]]; then
            found_installations=("$full_path")
        else
            error "phpier not found at: $full_path"
        fi
    else
        # Find all installations
        log "Searching for phpier installations..."
        local installations_output
        installations_output=$(find_phpier_installations)
        if [[ -n "$installations_output" ]]; then
            while IFS= read -r line; do
                if [[ -n "$line" ]]; then
                    found_installations+=("$line")
                    log "Found phpier at: $line"
                fi
            done <<< "$installations_output"
        fi
    fi
    
    # Check if any installations were found
    if [[ ${#found_installations[@]} -eq 0 ]]; then
        warn "No phpier installations found on this system"
        echo
        echo "phpier may have been:"
        echo "  - Already uninstalled"
        echo "  - Installed in a non-standard location"
        echo "  - Installed with a different name"
        echo
        echo "To find phpier manually:"
        echo "  find / -name 'phpier*' -type f 2>/dev/null"
        exit 0
    fi
    
    # Show found installations
    log "Found ${#found_installations[@]} phpier installation(s):"
    echo
    
    for installation in "${found_installations[@]}"; do
        show_installation_info "$installation"
    done
    
    # Confirm removal unless forced or dry run
    if [[ "$FORCE" == "false" ]] && [[ "$DRY_RUN" == "false" ]]; then
        echo "This will remove phpier from your system."
        echo
        if [[ ${#found_installations[@]} -eq 1 ]]; then
            echo "Remove phpier from: ${found_installations[0]}"
        else
            echo "Remove all ${#found_installations[@]} phpier installations:"
            for installation in "${found_installations[@]}"; do
                echo "  - $installation"
            done
        fi
        echo
        read -p "Continue with uninstallation? (y/N): " -n 1 -r
        echo
        
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log "Uninstallation cancelled by user"
            exit 0
        fi
    fi
    
    # Remove installations
    local removed_count=0
    for installation in "${found_installations[@]}"; do
        if [[ "$REMOVE_ALL" == "true" ]] || [[ ${#found_installations[@]} -eq 1 ]] || [[ "$FORCE" == "true" ]]; then
            remove_phpier "$installation" "$DRY_RUN"
            ((removed_count++))
        else
            echo
            read -p "Remove phpier from $installation? (y/N): " -n 1 -r
            echo
            if [[ $REPLY =~ ^[Yy]$ ]]; then
                remove_phpier "$installation" "$DRY_RUN"
                ((removed_count++))
            fi
        fi
    done
    
    # Clean up data directories
    cleanup_phpier_data "$DRY_RUN"
    
    # Clean up Docker containers and networks (but preserve volumes)
    cleanup_phpier_docker "$DRY_RUN"
    
    # Show summary
    show_summary "$removed_count" "$DRY_RUN"
}

# Run main function
main "$@"