#!/bin/bash
#
# Additional Tools Installation Script
#
set -e

OS="$1"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_success() { echo -e "${GREEN}[OK]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARN]${NC} $1"; }

# Install Starship prompt
install_starship() {
    if command -v starship &>/dev/null; then
        print_success "Starship already installed"
    else
        echo "Installing Starship..."
        curl -sS https://starship.rs/install.sh | sh -s -- -y
        print_success "Starship installed"
    fi
}

# Install vim-jetpack for Neovim
install_vim_jetpack() {
    local jetpack_path="$HOME/.local/share/nvim/site/pack/jetpack/opt/vim-jetpack"

    if [[ -d "$jetpack_path" ]]; then
        print_success "vim-jetpack already installed"
    else
        echo "Installing vim-jetpack..."
        mkdir -p "$(dirname "$jetpack_path")"
        git clone --depth 1 https://github.com/tani/vim-jetpack.git "$jetpack_path"
        print_success "vim-jetpack installed"
    fi
}

# Install vim-jetpack for Vim
install_vim_jetpack_vim() {
    local jetpack_path="$HOME/.vim/pack/jetpack/opt/vim-jetpack"

    if [[ -d "$jetpack_path" ]]; then
        print_success "vim-jetpack (vim) already installed"
    else
        echo "Installing vim-jetpack for Vim..."
        mkdir -p "$(dirname "$jetpack_path")"
        git clone --depth 1 https://github.com/tani/vim-jetpack.git "$jetpack_path"
        print_success "vim-jetpack (vim) installed"
    fi
}

# Install Lua magick binding for image.nvim
install_lua_magick() {
    if command -v luarocks &>/dev/null; then
        if luarocks --local list 2>/dev/null | grep -q "magick"; then
            print_success "Lua magick already installed"
        else
            echo "Installing Lua magick binding..."
            luarocks --local --lua-version=5.1 install magick
            print_success "Lua magick installed"
        fi
    else
        print_warning "luarocks not found, skipping magick installation"
    fi
}

# Install Nerd Fonts (optional, for icons)
install_nerd_fonts() {
    if [[ "$OS" == "macos" ]]; then
        if brew list --cask font-hack-nerd-font &>/dev/null; then
            print_success "Nerd Fonts already installed"
        else
            echo "Installing Nerd Fonts..."
            brew tap homebrew/cask-fonts 2>/dev/null || true
            brew install --cask font-hack-nerd-font || print_warning "Failed to install Nerd Fonts"
        fi
    else
        # Linux: Download manually
        local fonts_dir="$HOME/.local/share/fonts"
        if [[ -f "$fonts_dir/HackNerdFont-Regular.ttf" ]]; then
            print_success "Nerd Fonts already installed"
        else
            echo "Installing Nerd Fonts..."
            mkdir -p "$fonts_dir"
            cd /tmp
            curl -fLo "Hack.zip" https://github.com/ryanoasis/nerd-fonts/releases/latest/download/Hack.zip
            unzip -o Hack.zip -d "$fonts_dir"
            rm Hack.zip
            fc-cache -fv "$fonts_dir" 2>/dev/null || true
            print_success "Nerd Fonts installed"
        fi
    fi
}

# Sync Neovim plugins
sync_nvim_plugins() {
    if command -v nvim &>/dev/null; then
        echo "Syncing Neovim plugins..."
        nvim --headless "+Jetpack sync" +qa 2>/dev/null || print_warning "Jetpack sync may need manual run"
        print_success "Neovim plugins synced (or run :Jetpack sync manually)"
    fi
}

# Main
echo "Installing additional tools..."

install_starship
install_vim_jetpack
install_vim_jetpack_vim
install_lua_magick
install_nerd_fonts
sync_nvim_plugins

print_success "Tool installation complete"
