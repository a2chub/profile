#!/bin/bash
#
# Dotfiles Setup Script
# Supports: macOS, Ubuntu/Debian, Fedora/RHEL
#
set -e

DOTFILES_DIR="$(cd "$(dirname "$0")" && pwd)"
SCRIPTS_DIR="$DOTFILES_DIR/scripts"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "\n${BLUE}===================================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}===================================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Detect OS
detect_os() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
        OS="macos"
    elif [[ -f /etc/debian_version ]]; then
        OS="debian"
    elif [[ -f /etc/redhat-release ]]; then
        OS="redhat"
    else
        OS="unknown"
    fi
    echo "$OS"
}

# Main setup
main() {
    print_header "Dotfiles Setup"

    OS=$(detect_os)
    echo "Detected OS: $OS"
    echo "Dotfiles directory: $DOTFILES_DIR"

    if [[ "$OS" == "unknown" ]]; then
        print_error "Unsupported OS. Exiting."
        exit 1
    fi

    # Parse arguments
    SKIP_PACKAGES=false
    SKIP_LINKS=false
    SKIP_TOOLS=false

    while [[ $# -gt 0 ]]; do
        case $1 in
            --skip-packages) SKIP_PACKAGES=true; shift ;;
            --skip-links) SKIP_LINKS=true; shift ;;
            --skip-tools) SKIP_TOOLS=true; shift ;;
            --links-only) SKIP_PACKAGES=true; SKIP_TOOLS=true; shift ;;
            -h|--help)
                echo "Usage: ./setup.sh [options]"
                echo ""
                echo "Options:"
                echo "  --skip-packages   Skip package installation"
                echo "  --skip-links      Skip symlink creation"
                echo "  --skip-tools      Skip tool installation (Starship, etc.)"
                echo "  --links-only      Only create symlinks"
                echo "  -h, --help        Show this help"
                exit 0
                ;;
            *) print_warning "Unknown option: $1"; shift ;;
        esac
    done

    # Step 1: Install packages
    if [[ "$SKIP_PACKAGES" == false ]]; then
        print_header "Step 1: Installing Packages"
        bash "$SCRIPTS_DIR/install-packages.sh" "$OS"
    else
        print_warning "Skipping package installation"
    fi

    # Step 2: Create symlinks
    if [[ "$SKIP_LINKS" == false ]]; then
        print_header "Step 2: Creating Symlinks"
        bash "$SCRIPTS_DIR/link-dotfiles.sh" "$DOTFILES_DIR"
    else
        print_warning "Skipping symlink creation"
    fi

    # Step 3: Install additional tools
    if [[ "$SKIP_TOOLS" == false ]]; then
        print_header "Step 3: Installing Tools"
        bash "$SCRIPTS_DIR/install-tools.sh" "$OS"
    else
        print_warning "Skipping tool installation"
    fi

    print_header "Setup Complete!"
    echo "Please restart your terminal or run: source ~/.zshrc"
}

main "$@"
