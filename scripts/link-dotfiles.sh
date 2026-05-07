#!/bin/bash
#
# Symlink Creation Script
#
set -euo pipefail

DOTFILES_DIR="$1"

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/lib/colors.sh"

if [[ -z "$DOTFILES_DIR" ]]; then
    print_error "Usage: link-dotfiles.sh <dotfiles-dir>"
    exit 1
fi

# シンボリックリンク作成（既存ファイルはバックアップ、既存シンボリックリンクは置換）
create_link() {
    local src="$1"
    local dest="$2"

    if [[ ! -e "$src" ]]; then
        print_warning "Source not found: $src"
        return 0
    fi

    local parent_dir
    parent_dir="$(dirname "$dest")"
    if [[ ! -d "$parent_dir" ]]; then
        mkdir -p "$parent_dir"
        print_info "Created directory: $parent_dir"
    fi

    if [[ -e "$dest" && ! -L "$dest" ]]; then
        local backup="${dest}.backup.$(date +%Y%m%d%H%M%S)"
        mv "$dest" "$backup"
        print_warning "Backed up existing file to: $backup"
    fi

    if [[ -L "$dest" ]]; then
        rm "$dest"
    fi

    ln -s "$src" "$dest"
    print_success "Linked: $dest -> $src"
}

# シンボリックリンク定義: "<dotfiles配下の相対パス>:<リンク先絶対パス>"
LINKS=(
    ".zshrc:$HOME/.zshrc"
    ".tmux.conf:$HOME/.tmux.conf"
    ".vimrc:$HOME/.vimrc"
    "config/nvim:$HOME/.config/nvim"
    "config/starship.toml:$HOME/.config/starship.toml"
    "config/aerospace:$HOME/.config/aerospace"
    "config/borders:$HOME/.config/borders"
    "config/zellij:$HOME/.config/zellij"
    "config/jj:$HOME/.config/jj"
    "config/gh:$HOME/.config/gh"
    "config/ghostty:$HOME/.config/ghostty"
    "config/openspec:$HOME/.config/openspec"
    "config/uv:$HOME/.config/uv"
)

echo "Creating symlinks..."
for entry in "${LINKS[@]}"; do
    src="${entry%%:*}"
    dest="${entry#*:}"
    create_link "$DOTFILES_DIR/$src" "$dest"
done

print_success "All symlinks created"
