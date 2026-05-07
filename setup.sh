#!/bin/bash
#
# Dotfiles Setup Script
# Supports: macOS, Ubuntu/Debian, Fedora/RHEL
#
set -euo pipefail

DOTFILES_DIR="$(cd "$(dirname "$0")" && pwd)"
SCRIPTS_DIR="$DOTFILES_DIR/scripts"
SETUP_DIR="$SCRIPTS_DIR/setup"

source "$SCRIPTS_DIR/lib/colors.sh"

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
    FULL_INSTALL=false
    SETUP_SSH=false
    SETUP_DOCKER=false
    INSTALL_APPS=false

    while [[ $# -gt 0 ]]; do
        case $1 in
            --skip-packages) SKIP_PACKAGES=true; shift ;;
            --skip-links) SKIP_LINKS=true; shift ;;
            --skip-tools) SKIP_TOOLS=true; shift ;;
            --full) FULL_INSTALL=true; shift ;;
            --links-only) SKIP_PACKAGES=true; SKIP_TOOLS=true; shift ;;
            --setup-ssh) SETUP_SSH=true; shift ;;
            --setup-docker) SETUP_DOCKER=true; shift ;;
            --install-apps) INSTALL_APPS=true; shift ;;
            -h|--help)
                echo "Usage: ./setup.sh [options]"
                echo ""
                echo "Options:"
                echo "  --skip-packages   Skip package installation"
                echo "  --skip-links      Skip symlink creation"
                echo "  --skip-tools      Skip tool installation (Starship, etc.)"
                echo "  --full            Install optional heavy packages (octave, whisper, etc.)"
                echo "  --links-only      Only create symlinks"
                echo "  --setup-ssh       Setup SSH key for GitHub"
                echo "  --setup-docker    Install Docker"
                echo "  --install-apps    Run legacy app installation script"
                echo "  -h, --help        Show this help"
                exit 0
                ;;
            *) print_error "Unknown option: $1"; exit 1 ;;
        esac
    done

    # Step 1: Install packages
    if [[ "$SKIP_PACKAGES" == false ]]; then
        print_header "Step 1: Installing Packages"
        bash "$SCRIPTS_DIR/install-packages.sh" "$OS" "$FULL_INSTALL"
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

    # Optional: Legacy app installation
    if [[ "$INSTALL_APPS" == true ]]; then
        print_header "Installing Apps (Legacy Script)"
        if [[ -f "$SETUP_DIR/install-apps.sh" ]]; then
            bash "$SETUP_DIR/install-apps.sh"
        else
            print_error "install-apps.sh not found in $SETUP_DIR"
        fi
    fi

    # Optional: SSH setup
    if [[ "$SETUP_SSH" == true ]]; then
        print_header "Setting up SSH"
        if [[ -f "$SETUP_DIR/setup-ssh.sh" ]]; then
            bash "$SETUP_DIR/setup-ssh.sh"
        else
            print_error "setup-ssh.sh not found in $SETUP_DIR"
        fi
    fi

    # Optional: Docker installation
    if [[ "$SETUP_DOCKER" == true ]]; then
        print_header "Installing Docker"
        if [[ -f "$SETUP_DIR/install-docker.sh" ]]; then
            bash "$SETUP_DIR/install-docker.sh"
        else
            print_error "install-docker.sh not found in $SETUP_DIR"
        fi
    fi

    print_header "Setup Complete!"
    echo "Please restart your terminal or run: source ~/.zshrc"
}

main "$@"
