# PATH設定
typeset -U path PATH
path=(
  /opt/homebrew/bin(N-/)
  /opt/homebrew/sbin(N-/)
  /usr/bin
  /usr/sbin
  /bin
  /sbin
  /usr/local/bin(N-/)
  /usr/local/sbin(N-/)
  /Library/Apple/usr/bin
)

# 補完待機中に赤い点を表示
COMPLETION_WAITING_DOTS="true"

# エイリアス
alias vim=nvim
alias python3=/opt/homebrew/bin/python3.14
alias python=python3

# Rye
source "$HOME/.rye/env"

# Google Cloud SDK
if [ -f '/Users/atusi/google-cloud-sdk/path.zsh.inc' ]; then . '/Users/atusi/google-cloud-sdk/path.zsh.inc'; fi
if [ -f '/Users/atusi/google-cloud-sdk/completion.zsh.inc' ]; then . '/Users/atusi/google-cloud-sdk/completion.zsh.inc'; fi

. "$HOME/.local/bin/env"

# pnpm
export PNPM_HOME="/Users/atusi/Library/pnpm"
case ":$PATH:" in
  *":$PNPM_HOME:"*) ;;
  *) export PATH="$PNPM_HOME:$PATH" ;;
esac

export STM32_PRG_PATH=/Applications/STMicroelectronics/STM32Cube/STM32CubeProgrammer/STM32CubeProgrammer.app/Contents/MacOs/bin

# LM Studio CLI
export PATH="$PATH:/Users/atusi/.lmstudio/bin"

export PATH="/opt/homebrew/opt/curl/bin:$PATH"
export LDFLAGS="-L/opt/homebrew/opt/curl/lib"
export CPPFLAGS="-I/opt/homebrew/opt/curl/include"

# Docker CLI completions
fpath=(/Users/atusi/.docker/completions $fpath)
autoload -Uz compinit
compinit

# Claude Code
alias cc='claude --dangerously-skip-permissions'
alias claude="/Users/atusi/.claude/local/claude"

# Kiro shell integration
[[ "$TERM_PROGRAM" == "kiro" ]] && . "$(kiro --locate-shell-integration-path zsh)"

# Windsurf
export PATH="/Users/atusi/.codeium/windsurf/bin:$PATH"
export PATH="$HOME/bin:$PATH"
export PATH="/opt/X11/bin:$PATH"

# Homebrew設定
export HOMEBREW_NO_INSTALL_CLEANUP=1

# Starship prompt
eval "$(starship init zsh)"

# Antigravity
export PATH="/Users/atusi/.antigravity/antigravity/bin:$PATH"

# The next line updates PATH for the Google Cloud SDK.
if [ -f '/Users/atusi/google-cloud-sdk/path.zsh.inc' ]; then . '/Users/atusi/google-cloud-sdk/path.zsh.inc'; fi

# The next line enables shell command completion for gcloud.
if [ -f '/Users/atusi/google-cloud-sdk/completion.zsh.inc' ]; then . '/Users/atusi/google-cloud-sdk/completion.zsh.inc'; fi

# bun completions
[ -s "/Users/atusi/.bun/_bun" ] && source "/Users/atusi/.bun/_bun"

# bun
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"

# opencode
export PATH=/Users/atusi/.opencode/bin:$PATH
