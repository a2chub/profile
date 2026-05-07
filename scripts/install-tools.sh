#!/bin/bash
#
# Additional Tools Installation Script
#
set -euo pipefail

OS="$1"

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/colors.sh"

# Install vim-jetpack into the given directory
# Args: <jetpack_path> <label>
install_vim_jetpack() {
    local jetpack_path="$1"
    local label="$2"

    if [[ -d "$jetpack_path" ]]; then
        print_success "vim-jetpack ($label) already installed"
    else
        echo "Installing vim-jetpack ($label)..."
        mkdir -p "$(dirname "$jetpack_path")"
        git clone --depth 1 https://github.com/tani/vim-jetpack.git "$jetpack_path"
        print_success "vim-jetpack ($label) installed"
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

install_vim_jetpack "$HOME/.local/share/nvim/site/pack/jetpack/opt/vim-jetpack" "nvim"
install_vim_jetpack "$HOME/.vim/pack/jetpack/opt/vim-jetpack" "vim"
install_lua_magick
sync_nvim_plugins

print_success "Tool installation complete"
