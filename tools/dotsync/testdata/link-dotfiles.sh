#!/bin/bash
#
# Symlink Creation Script (テスト用フィクスチャ)
#
DOTFILES_DIR="$1"

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
