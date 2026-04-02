#!/bin/bash
#
# Symlink Creation Script
#
set -e

DOTFILES_DIR="$1"

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/colors.sh"

if [[ -z "$DOTFILES_DIR" ]]; then
    print_error "Usage: link-dotfiles.sh <dotfiles-dir>"
    exit 1
fi

# Create symlink with backup
create_link() {
    local src="$1"
    local dest="$2"

    # Check if source exists
    if [[ ! -e "$src" ]]; then
        print_warning "Source not found: $src"
        return 0
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

# Zellij terminal multiplexer
create_link "$DOTFILES_DIR/config/zellij" "$HOME/.config/zellij"

# jj (Jujutsu) version control
create_link "$DOTFILES_DIR/config/jj" "$HOME/.config/jj"

create_link "$DOTFILES_DIR/config/.starship.toml" "$HOME/.config/.starship.toml"

create_link "$DOTFILES_DIR/config/aerospace" "$HOME/.config/aerospace"

create_link "$DOTFILES_DIR/config/gh" "$HOME/.config/gh"

create_link "$DOTFILES_DIR/config/ghostty" "$HOME/.config/ghostty"

create_link "$DOTFILES_DIR/config/openspec" "$HOME/.config/openspec"

create_link "$DOTFILES_DIR/config/uv" "$HOME/.config/uv"

print_success "All symlinks created"
