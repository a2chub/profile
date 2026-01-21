# ==============================================================================
#                                    .zshrc
# ==============================================================================
# macOS用のZsh設定ファイル
# 最終更新: 2024年12月
# ==============================================================================


# ==============================================================================
# 1. 基本設定
# ==============================================================================

# PATH設定（重複を自動排除）
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


# ==============================================================================
# 2. ヒストリー設定
# ==============================================================================

# 基本設定
HISTFILE=~/.zsh_history
HISTSIZE=50000              # メモリ上に保持する行数
SAVEHIST=1000000            # ファイルに保存する行数（実質無制限）

# ヒストリーオプション
setopt HIST_IGNORE_DUPS      # 連続する重複コマンドを無視
setopt HIST_IGNORE_ALL_DUPS  # 古い重複を削除
setopt SHARE_HISTORY         # 複数ターミナル間で共有
setopt EXTENDED_HISTORY      # タイムスタンプを記録
setopt HIST_REDUCE_BLANKS    # 余分な空白を削除
setopt INC_APPEND_HISTORY    # 即時にファイルに追記

# ------------------------------------------------------------------------------
# ヒストリー月次アーカイブ
# ------------------------------------------------------------------------------
# 動作: 3ヶ月以上前のエントリを月別ファイルに自動移動
# 保存先: ~/.zsh_history_archive/zsh_history_YYYY-MM.txt
# 実行タイミング: シェル起動時（月に1回のみ）
# ------------------------------------------------------------------------------

HIST_ARCHIVE_DIR=~/.zsh_history_archive
[[ -d "$HIST_ARCHIVE_DIR" ]] || mkdir -p "$HIST_ARCHIVE_DIR"

_current_month=$(date +%Y-%m)
_last_archive_file="$HIST_ARCHIVE_DIR/.last_archive"

if [[ -f "$_last_archive_file" ]]; then
    _last_archived=$(cat "$_last_archive_file")
else
    _last_archived=""
fi

if [[ "$_last_archived" != "$_current_month" && -f "$HISTFILE" && -s "$HISTFILE" ]]; then
    _cutoff_timestamp=$(date -v-2m -v1d -v0H -v0M -v0S +%s)
    _temp_recent=$(mktemp)

    # EXTENDED_HISTORY形式: ": timestamp:0;command"
    while IFS= read -r line || [[ -n "$line" ]]; do
        if [[ "$line" == ": "* ]]; then
            _ts=${line#: }
            _ts=${_ts%%:*}
            if [[ "$_ts" == <-> ]] && (( _ts >= _cutoff_timestamp )); then
                echo "$line" >> "$_temp_recent"
            elif [[ "$_ts" == <-> ]]; then
                _entry_month=$(date -r $_ts +%Y-%m)
                echo "$line" >> "$HIST_ARCHIVE_DIR/zsh_history_${_entry_month}.txt"
            else
                echo "$line" >> "$_temp_recent"
            fi
        else
            echo "$line" >> "$_temp_recent"
        fi
    done < "$HISTFILE"

    mv "$_temp_recent" "$HISTFILE"
    echo "$_current_month" > "$_last_archive_file"
fi

unset _current_month _last_archive_file _last_archived _cutoff_timestamp _temp_recent _ts _entry_month


# ==============================================================================
# 3. エイリアス
# ==============================================================================

# エディタ
alias vim=nvim

# Python
alias python3=/opt/homebrew/bin/python3.14
alias python=python3

# AI ツール
alias cc='claude --dangerously-skip-permissions'

# Git ワークフロー
alias gsync='git fetch --all && git pull && git status'
alias gac='git add . && git commit'

# Docker ワークフロー
alias drestart='docker-compose down && docker-compose up -d'
alias dclean='docker system prune -f'
alias dlogs='docker-compose logs -f'


# ==============================================================================
# 4. 補完設定
# ==============================================================================

# Docker CLI補完
fpath=(/Users/atusi/.docker/completions $fpath)
autoload -Uz compinit
compinit

# Bun補完
[ -s "/Users/atusi/.bun/_bun" ] && source "/Users/atusi/.bun/_bun"


# ==============================================================================
# 5. 開発ツール - 言語・ランタイム
# ==============================================================================

# Python (Rye)
[[ -f "$HOME/.rye/env" ]] && source "$HOME/.rye/env"

# Node.js (pnpm)
export PNPM_HOME="/Users/atusi/Library/pnpm"
case ":$PATH:" in
  *":$PNPM_HOME:"*) ;;
  *) export PATH="$PNPM_HOME:$PATH" ;;
esac

# Bun
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"

# curl (Homebrew版)
export PATH="/opt/homebrew/opt/curl/bin:$PATH"
export LDFLAGS="-L/opt/homebrew/opt/curl/lib"
export CPPFLAGS="-I/opt/homebrew/opt/curl/include"


# ==============================================================================
# 6. 開発ツール - クラウド・組み込み
# ==============================================================================

# Google Cloud SDK
if [ -f '/Users/atusi/google-cloud-sdk/path.zsh.inc' ]; then
    . '/Users/atusi/google-cloud-sdk/path.zsh.inc'
fi
if [ -f '/Users/atusi/google-cloud-sdk/completion.zsh.inc' ]; then
    . '/Users/atusi/google-cloud-sdk/completion.zsh.inc'
fi

# STM32 Programmer
export STM32_PRG_PATH=/Applications/STMicroelectronics/STM32Cube/STM32CubeProgrammer/STM32CubeProgrammer.app/Contents/MacOs/bin


# ==============================================================================
# 7. AI・エディタツール
# ==============================================================================

# GitHub Personal Access Token (Keychainから取得)
export GITHUB_PERSONAL_ACCESS_TOKEN=$(security find-generic-password -a "$USER" -s "github-pat" -w 2>/dev/null)

# LM Studio
export PATH="$PATH:/Users/atusi/.lmstudio/bin"

# Windsurf (Codeium)
export PATH="/Users/atusi/.codeium/windsurf/bin:$PATH"

# Kiro shell integration
[[ "$TERM_PROGRAM" == "kiro" ]] && . "$(kiro --locate-shell-integration-path zsh)"

# OpenCode
export PATH=/Users/atusi/.opencode/bin:$PATH

# Antigravity
export PATH="/Users/atusi/.antigravity/antigravity/bin:$PATH"


# ==============================================================================
# 8. その他
# ==============================================================================

# ローカルバイナリ
[[ -f "$HOME/.local/bin/env" ]] && . "$HOME/.local/bin/env"
export PATH="$HOME/bin:$PATH"

# X11
export PATH="/opt/X11/bin:$PATH"

# Homebrew設定
export HOMEBREW_NO_INSTALL_CLEANUP=1


# ==============================================================================
# 9. エージェント・自動化
# ==============================================================================

# エージェントスクリプト実行ヘルパー
agent-run() {
    local agent=$1
    shift
    local script="$HOME/dotfiles/agent/agents/general/${agent}.sh"
    if [[ -x "$script" ]]; then
        "$script" "$@"
    else
        echo "Agent not found: $agent" >&2
        return 1
    fi
}

# ------------------------------------------------------------------------------
# コマンドメトリクス収集
# ------------------------------------------------------------------------------
# 保存先: ~/.command_metrics.log
# 形式: timestamp:duration:pwd:command
# 用途: 作業パターン分析、自動化候補の検出
# ------------------------------------------------------------------------------

METRICS_LOG="$HOME/.command_metrics.log"

preexec() {
    _cmd_start_time=$EPOCHSECONDS
    _last_cmd=$1
}

precmd() {
    local exit_code=$?
    if [[ -n "$_cmd_start_time" && -n "$_last_cmd" ]]; then
        local duration=$(( EPOCHSECONDS - _cmd_start_time ))
        # 形式: timestamp:duration:exit_code:pwd:command
        echo "${_cmd_start_time}:${duration}:${exit_code}:${PWD}:${_last_cmd}" >> "$METRICS_LOG"
    fi
    unset _cmd_start_time _last_cmd
}


# ==============================================================================
# 10. プロンプト
# ==============================================================================

# Starship prompt（最後に読み込む）
eval "$(starship init zsh)"

alias claude-mem='bun "/Users/atusi/.claude/plugins/marketplaces/thedotmack/plugin/scripts/worker-service.cjs"'
