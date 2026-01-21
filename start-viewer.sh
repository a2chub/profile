#!/bin/bash
#
# Dotfiles Dashboard起動スクリプト
#

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
VIEWER_DIR="$SCRIPT_DIR/viewer"
PORT="${DOTFILES_VIEWER_PORT:-8765}"

# カラー出力
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Dotfiles Dashboard${NC}"
echo "================================"
echo "Starting server on http://127.0.0.1:$PORT"
echo ""

# ブラウザを開く（macOS）
if command -v open &> /dev/null; then
    (sleep 1 && open "http://127.0.0.1:$PORT") &
fi

# サーバー起動
cd "$VIEWER_DIR"
python3 server.py
