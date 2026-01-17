#!/bin/bash
#
# Symlink Creation Script
#
set -e

DOTFILES_DIR="$1"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_success() { echo -e "${GREEN}[OK]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARN]${NC} $1"; }
print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }

# Create symlink with backup
create_link() {
    local src="$1"
    local dest="$2"

    # Check if source exists
    if [[ ! -e "$src" ]]; then
        print_warning "Source not found: $src"
        return 1
    fi

    # Create parent directory if needed
    local parent_dir
    parent_dir="$(dirname "$dest")"
    if [[ ! -d "$parent_dir" ]]; then
        mkdir -p "$parent_dir"
        print_info "Created directory: $parent_dir"
    fi

    # If destination exists and is not a symlink, back it up
    if [[ -e "$dest" && ! -L "$dest" ]]; then
        local backup="${dest}.backup.$(date +%Y%m%d%H%M%S)"
        mv "$dest" "$backup"
        print_warning "Backed up existing file to: $backup"
    fi

    # Remove existing symlink
    if [[ -L "$dest" ]]; then
        rm "$dest"
    fi

    # Create symlink
    ln -s "$src" "$dest"
    print_success "Linked: $dest -> $src"
}

# Main
echo "Creating symlinks..."

# Shell configuration
create_link "$DOTFILES_DIR/.zshrc" "$HOME/.zshrc"
create_link "$DOTFILES_DIR/.tmux.conf" "$HOME/.tmux.conf"

# Vim configuration
create_link "$DOTFILES_DIR/.vimrc" "$HOME/.vimrc"

# Neovim configuration
create_link "$DOTFILES_DIR/config/nvim" "$HOME/.config/nvim"

# Starship prompt
create_link "$DOTFILES_DIR/config/starship.toml" "$HOME/.config/starship.toml"

# AeroSpace window manager
create_link "$DOTFILES_DIR/config/aerospace/aerospace.toml" "$HOME/.config/aerospace/aerospace.toml"

# Borders (window border styling)
create_link "$DOTFILES_DIR/config/borders" "$HOME/.config/borders"

# Sketchybar
create_link "$DOTFILES_DIR/config/sketchybar" "$HOME/.config/sketchybar"

print_success "All symlinks created"
