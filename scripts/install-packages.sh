#!/bin/bash
#
# Package Installation Script
#
set -e

OS="$1"

# Script directory and paths
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PACKAGES_DIR="$(dirname "$SCRIPT_DIR")/packages"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() { echo -e "${GREEN}[OK]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARN]${NC} $1"; }

# Common packages to install
COMMON_PACKAGES=(
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

# macOS specific packages
MACOS_PACKAGES=(
    git
    curl
    wget
    tmux
    neovim
    ripgrep
    fd
    jq
    tree
    htop
    coreutils
    gnu-sed
)

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
    if [[ -f "$BREWFILE" ]]; then
        echo "Installing packages via Brewfile..."
        brew bundle --file="$BREWFILE"
        print_success "Brewfile packages installed"
    else
        # Fallback: install basic packages individually
        echo "Brewfile not found, installing basic packages via Homebrew..."
        for pkg in "${MACOS_PACKAGES[@]}"; do
            if brew list "$pkg" &>/dev/null; then
                print_success "$pkg already installed"
            else
                echo "Installing $pkg..."
                brew install "$pkg" || print_warning "Failed to install $pkg"
            fi
        done
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

    for pkg in "${DEBIAN_PACKAGES[@]}"; do
        if dpkg -l "$pkg" &>/dev/null; then
            print_success "$pkg already installed"
        else
            echo "Installing $pkg..."
            sudo apt install -y "$pkg" || print_warning "Failed to install $pkg"
        fi
    done

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

    for pkg in "${REDHAT_PACKAGES[@]}"; do
        if rpm -q "$pkg" &>/dev/null; then
            print_success "$pkg already installed"
        else
            echo "Installing $pkg..."
            sudo dnf install -y "$pkg" || print_warning "Failed to install $pkg"
        fi
    done
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
