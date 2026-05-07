#!/bin/bash
#
# Package Installation Script
#
set -euo pipefail

OS="$1"
FULL_INSTALL="${2:-false}"

# Script directory and paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PACKAGES_DIR="$(dirname "$SCRIPT_DIR")/packages"

source "$SCRIPT_DIR/lib/colors.sh"

# Install a list of packages, checking for existence first
# Args: <check_cmd> <install_cmd> <pkg1> <pkg2> ...
_install_packages() {
    local check_cmd="$1"
    local install_cmd="$2"
    shift 2
    for pkg in "$@"; do
        if eval "$check_cmd \"$pkg\"" &>/dev/null; then
            print_success "$pkg already installed"
        else
            echo "Installing $pkg..."
            eval "$install_cmd \"$pkg\"" || print_warning "Failed to install $pkg"
        fi
    done
}

# Install Homebrew (macOS)
install_homebrew() {
    if command -v brew &>/dev/null; then
        print_success "Homebrew already installed"
        brew update
    else
        echo "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

        # Add to PATH for Apple Silicon
        if [[ -f /opt/homebrew/bin/brew ]]; then
            eval "$(/opt/homebrew/bin/brew shellenv)"
        fi
        print_success "Homebrew installed"
    fi
}

# Install packages on macOS
install_macos() {
    install_homebrew

    BREWFILE="$PACKAGES_DIR/Brewfile"
    BREWFILE_OPTIONAL="$PACKAGES_DIR/Brewfile.optional"
    if [[ -f "$BREWFILE" ]]; then
        echo "Installing base packages via Brewfile..."
        if brew bundle --file="$BREWFILE"; then
            print_success "Brewfile packages installed"
        else
            print_warning "Some Brewfile packages failed to install (sudo-requiring casks may need manual installation)"
        fi

        # オプションパッケージ（--full 指定時のみ）
        if [[ "$FULL_INSTALL" == true ]] && [[ -f "$BREWFILE_OPTIONAL" ]]; then
            echo "Installing optional packages via Brewfile.optional..."
            if brew bundle --file="$BREWFILE_OPTIONAL"; then
                print_success "Optional packages installed"
            else
                print_warning "Some optional packages failed to install"
            fi
        fi
    else
        # Fallback: Brewfile が見つからない場合は基本パッケージを個別インストール
        print_warning "Brewfile not found at $BREWFILE"
        local fallback_packages=(git curl wget tmux neovim ripgrep fd jq tree htop)
        echo "Installing basic packages via Homebrew..."
        _install_packages "brew list" "brew install" "${fallback_packages[@]}"
    fi
}

# Install packages on Debian/Ubuntu
install_debian() {
    echo "Updating package lists..."
    sudo apt update

    echo "Installing packages via apt..."
    DEBIAN_PACKAGES=(
        git
        curl
        wget
        tmux
        ripgrep
        fd-find
        jq
        tree
        htop
        build-essential
    )

    _install_packages "dpkg -l" "sudo apt install -y" "${DEBIAN_PACKAGES[@]}"

    # Neovim (get latest from GitHub releases or PPA)
    if ! command -v nvim &>/dev/null; then
        echo "Installing Neovim..."
        sudo apt install -y software-properties-common
        sudo add-apt-repository -y ppa:neovim-ppa/unstable
        sudo apt update
        sudo apt install -y neovim
    else
        print_success "Neovim already installed"
    fi
}

# Install packages on Fedora/RHEL
install_redhat() {
    echo "Installing packages via dnf..."
    REDHAT_PACKAGES=(
        git
        curl
        wget
        tmux
        neovim
        ripgrep
        fd-find
        jq
        tree
        htop
    )

    _install_packages "rpm -q" "sudo dnf install -y" "${REDHAT_PACKAGES[@]}"
}

# Main
case "$OS" in
    macos)
        install_macos
        ;;
    debian)
        install_debian
        ;;
    redhat)
        install_redhat
        ;;
    *)
        echo "Unknown OS: $OS"
        exit 1
        ;;
esac

print_success "Package installation complete"
