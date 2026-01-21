#!/usr/bin/env python3
"""
Dotfiles Dashboard Server
Python標準ライブラリのみを使用した軽量Webサーバー
"""

import json
import os
import re
import shutil
import subprocess
from datetime import datetime
from http.server import HTTPServer, SimpleHTTPRequestHandler
from pathlib import Path
from urllib.parse import urlparse, parse_qs

# サーバー設定
PORT = int(os.environ.get("DOTFILES_VIEWER_PORT", 8765))
HOST = "127.0.0.1"

# パス設定
SCRIPT_DIR = Path(__file__).parent.resolve()
DOTFILES_DIR = SCRIPT_DIR.parent.resolve()
STATIC_DIR = SCRIPT_DIR / "static"
BACKUPS_DIR = DOTFILES_DIR / ".backups"

# 管理対象の設定ファイル
CONFIG_FILES = {
    "zshrc": {
        "id": "zshrc",
        "displayName": ".zshrc",
        "software": "Zsh",
        "category": "Shell",
        "format": "shell",
        "sourcePath": ".zshrc",
        "docs": "https://zsh.sourceforge.io/Doc/",
        "repo": "https://github.com/zsh-users/zsh"
    },
    "tmux": {
        "id": "tmux",
        "displayName": ".tmux.conf",
        "software": "tmux",
        "category": "Terminal",
        "format": "conf",
        "sourcePath": ".tmux.conf",
        "docs": "https://github.com/tmux/tmux/wiki",
        "repo": "https://github.com/tmux/tmux"
    },
    "vimrc": {
        "id": "vimrc",
        "displayName": ".vimrc",
        "software": "Vim",
        "category": "Editor",
        "format": "vim",
        "sourcePath": ".vimrc",
        "docs": "https://vimdoc.sourceforge.net/",
        "repo": "https://github.com/vim/vim"
    },
    "nvim": {
        "id": "nvim",
        "displayName": "nvim/init.lua",
        "software": "Neovim",
        "category": "Editor",
        "format": "lua",
        "sourcePath": "config/nvim/init.lua",
        "docs": "https://neovim.io/doc/user/",
        "repo": "https://github.com/neovim/neovim"
    },
    "starship": {
        "id": "starship",
        "displayName": "starship.toml",
        "software": "Starship",
        "category": "Shell",
        "format": "toml",
        "sourcePath": "config/starship.toml",
        "docs": "https://starship.rs/config/",
        "repo": "https://github.com/starship/starship"
    },
    "aerospace": {
        "id": "aerospace",
        "displayName": "aerospace.toml",
        "software": "AeroSpace",
        "category": "WM",
        "format": "toml",
        "sourcePath": "config/aerospace/aerospace.toml",
        "docs": "https://nikitabobko.github.io/AeroSpace/guide",
        "repo": "https://github.com/nikitabobko/AeroSpace"
    },
    "borders": {
        "id": "borders",
        "displayName": "bordersrc",
        "software": "JankyBorders",
        "category": "WM",
        "format": "bash",
        "sourcePath": "config/borders/bordersrc",
        "docs": "https://github.com/FelixKratz/JankyBorders",
        "repo": "https://github.com/FelixKratz/JankyBorders"
    },
    "brewfile": {
        "id": "brewfile",
        "displayName": "Brewfile",
        "software": "Homebrew",
        "category": "Package",
        "format": "ruby",
        "sourcePath": "packages/Brewfile",
        "docs": "https://docs.brew.sh/",
        "repo": "https://github.com/Homebrew/brew"
    }
}


def get_config_path(config_id: str) -> Path:
    """設定ファイルの絶対パスを取得"""
    if config_id not in CONFIG_FILES:
        return None
    return DOTFILES_DIR / CONFIG_FILES[config_id]["sourcePath"]


def validate_path(path: Path) -> bool:
    """パスがdotfilesディレクトリ内にあることを検証"""
    try:
        path.resolve().relative_to(DOTFILES_DIR)
        return True
    except ValueError:
        return False


def create_backup(config_id: str) -> str:
    """設定ファイルのバックアップを作成"""
    config_path = get_config_path(config_id)
    if not config_path or not config_path.exists():
        return None

    BACKUPS_DIR.mkdir(exist_ok=True)
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    backup_name = f"{config_id}_{timestamp}"
    backup_path = BACKUPS_DIR / backup_name

    shutil.copy2(config_path, backup_path)
    return str(backup_path)


def validate_toml(content: str) -> tuple[bool, str]:
    """簡易的なTOML構文検証"""
    lines = content.split('\n')
    bracket_stack = []

    for i, line in enumerate(lines, 1):
        stripped = line.strip()
        if not stripped or stripped.startswith('#'):
            continue

        # セクションヘッダーのチェック
        if stripped.startswith('['):
            if not stripped.endswith(']'):
                return False, f"Line {i}: Unclosed bracket in section header"

        # 基本的な key = value 形式のチェック
        if '=' in stripped and not stripped.startswith('['):
            parts = stripped.split('=', 1)
            key = parts[0].strip()
            if not key or key.startswith('#'):
                continue
            if not re.match(r'^[a-zA-Z_][a-zA-Z0-9_\-\.]*$', key.split('.')[0]):
                return False, f"Line {i}: Invalid key format"

    return True, "OK"


def parse_brewfile(content: str) -> dict:
    """Brewfileをパースしてカテゴリ別に分類"""
    result = {
        "taps": [],
        "formulae": [],
        "casks": []
    }

    for line in content.split('\n'):
        line = line.strip()
        if not line or line.startswith('#'):
            continue

        if line.startswith('tap '):
            match = re.match(r'tap\s+["\']([^"\']+)["\']', line)
            if match:
                result["taps"].append({"name": match.group(1)})

        elif line.startswith('brew '):
            match = re.match(r'brew\s+["\']([^"\']+)["\']', line)
            if match:
                result["formulae"].append({"name": match.group(1)})

        elif line.startswith('cask '):
            match = re.match(r'cask\s+["\']([^"\']+)["\']', line)
            if match:
                result["casks"].append({"name": match.group(1)})

    return result


def validate_commit_hash(commit: str) -> bool:
    """コミットハッシュが有効な形式かチェック"""
    return bool(re.match(r'^[a-fA-F0-9]{7,40}$', commit))


def get_file_history(config_id: str, limit: int = 5) -> list:
    """ファイルのGit履歴を取得"""
    config_path = get_config_path(config_id)
    if not config_path:
        return []

    relative_path = CONFIG_FILES[config_id]["sourcePath"]

    try:
        result = subprocess.run(
            ["git", "log", f"--format=%H|%ai|%an|%s", f"-n{limit}", "--", relative_path],
            cwd=str(DOTFILES_DIR),
            capture_output=True,
            text=True,
            timeout=10
        )

        if result.returncode != 0:
            return []

        history = []
        for line in result.stdout.strip().split('\n'):
            if not line:
                continue
            parts = line.split('|', 3)
            if len(parts) == 4:
                history.append({
                    "hash": parts[0][:7],
                    "fullHash": parts[0],
                    "date": parts[1],
                    "author": parts[2],
                    "message": parts[3]
                })
        return history
    except Exception:
        return []


def get_commit_diff(config_id: str, commit: str) -> str:
    """指定コミットでのファイル変更差分を取得"""
    if not validate_commit_hash(commit):
        return ""

    relative_path = CONFIG_FILES[config_id]["sourcePath"]

    try:
        result = subprocess.run(
            ["git", "diff", f"{commit}~1", commit, "--", relative_path],
            cwd=str(DOTFILES_DIR),
            capture_output=True,
            text=True,
            timeout=10
        )

        if result.returncode != 0:
            # 最初のコミットの場合はshowを使用
            result = subprocess.run(
                ["git", "show", f"{commit}:{relative_path}"],
                cwd=str(DOTFILES_DIR),
                capture_output=True,
                text=True,
                timeout=10
            )
            if result.returncode == 0:
                return f"(Initial commit)\n\n{result.stdout}"
            return ""

        return result.stdout
    except Exception:
        return ""


class DotfilesHandler(SimpleHTTPRequestHandler):
    """dotfilesダッシュボード用HTTPリクエストハンドラー"""

    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=str(STATIC_DIR), **kwargs)

    def do_GET(self):
        parsed = urlparse(self.path)
        path = parsed.path

        # ルートパス -> index.html
        if path == "/" or path == "":
            self.path = "/index.html"
            return super().do_GET()

        # 静的ファイル
        if path.startswith("/static/"):
            self.path = path[7:]  # /static/ を除去
            return super().do_GET()

        # API: 設定ファイル一覧
        if path == "/api/configs":
            self.send_json_response(list(CONFIG_FILES.values()))
            return

        # API: 設定ファイル内容取得
        if path.startswith("/api/configs/"):
            config_id = path.split("/")[-1]
            self.handle_get_config(config_id)
            return

        # API: Brewパッケージ
        if path == "/api/brew/formulae":
            self.handle_brew_category("formulae")
            return

        if path == "/api/brew/casks":
            self.handle_brew_category("casks")
            return

        if path == "/api/brew/taps":
            self.handle_brew_category("taps")
            return

        # API: Git履歴
        if path.startswith("/api/history/"):
            parts = path.split("/")
            if len(parts) == 4:
                # /api/history/<config_id>
                config_id = parts[3]
                self.handle_get_history(config_id)
                return
            elif len(parts) == 5:
                # /api/history/<config_id>/<commit>
                config_id = parts[3]
                commit = parts[4]
                self.handle_get_diff(config_id, commit)
                return

        # その他は静的ファイルとして処理
        return super().do_GET()

    def do_PUT(self):
        parsed = urlparse(self.path)
        path = parsed.path

        # API: 設定ファイル保存
        if path.startswith("/api/configs/"):
            config_id = path.split("/")[-1]
            self.handle_put_config(config_id)
            return

        self.send_error(404, "Not Found")

    def handle_get_config(self, config_id: str):
        """設定ファイル内容を取得"""
        if config_id not in CONFIG_FILES:
            self.send_error(404, f"Config '{config_id}' not found")
            return

        config_path = get_config_path(config_id)
        if not config_path.exists():
            self.send_error(404, f"File not found: {config_path}")
            return

        try:
            content = config_path.read_text(encoding="utf-8")
            response = {
                **CONFIG_FILES[config_id],
                "content": content,
                "size": config_path.stat().st_size,
                "modified": datetime.fromtimestamp(
                    config_path.stat().st_mtime
                ).isoformat()
            }
            self.send_json_response(response)
        except Exception as e:
            self.send_error(500, str(e))

    def handle_put_config(self, config_id: str):
        """設定ファイルを保存"""
        if config_id not in CONFIG_FILES:
            self.send_error(404, f"Config '{config_id}' not found")
            return

        config_path = get_config_path(config_id)
        if not validate_path(config_path):
            self.send_error(403, "Invalid path")
            return

        try:
            content_length = int(self.headers.get("Content-Length", 0))
            if content_length > 1024 * 1024:  # 1MB制限
                self.send_error(413, "File too large")
                return

            body = self.rfile.read(content_length)
            data = json.loads(body.decode("utf-8"))
            content = data.get("content", "")

            # TOML形式の場合は構文検証
            if CONFIG_FILES[config_id]["format"] == "toml":
                valid, msg = validate_toml(content)
                if not valid:
                    self.send_json_response(
                        {"success": False, "error": f"TOML validation failed: {msg}"},
                        status=400
                    )
                    return

            # バックアップ作成
            backup_path = create_backup(config_id)

            # ファイル保存
            config_path.write_text(content, encoding="utf-8")

            self.send_json_response({
                "success": True,
                "backup": backup_path,
                "message": f"Saved {config_id}"
            })
        except json.JSONDecodeError:
            self.send_error(400, "Invalid JSON")
        except Exception as e:
            self.send_error(500, str(e))

    def handle_brew_category(self, category: str):
        """Brewパッケージカテゴリを取得"""
        brewfile_path = get_config_path("brewfile")
        if not brewfile_path.exists():
            self.send_json_response([])
            return

        try:
            content = brewfile_path.read_text(encoding="utf-8")
            packages = parse_brewfile(content)
            self.send_json_response(packages.get(category, []))
        except Exception as e:
            self.send_error(500, str(e))

    def handle_get_history(self, config_id: str):
        """ファイルのGit履歴を取得"""
        if config_id not in CONFIG_FILES:
            self.send_error(404, f"Config '{config_id}' not found")
            return

        history = get_file_history(config_id)
        self.send_json_response(history)

    def handle_get_diff(self, config_id: str, commit: str):
        """特定コミットのDiffを取得"""
        if config_id not in CONFIG_FILES:
            self.send_error(404, f"Config '{config_id}' not found")
            return

        if not validate_commit_hash(commit):
            self.send_error(400, "Invalid commit hash")
            return

        diff = get_commit_diff(config_id, commit)
        self.send_json_response({"diff": diff})

    def send_json_response(self, data, status=200):
        """JSONレスポンスを送信"""
        response = json.dumps(data, ensure_ascii=False, indent=2)
        self.send_response(status)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", len(response.encode("utf-8")))
        self.send_header("Access-Control-Allow-Origin", "*")
        self.end_headers()
        self.wfile.write(response.encode("utf-8"))

    def log_message(self, format, *args):
        """ログフォーマットをカスタマイズ"""
        print(f"[{datetime.now().strftime('%H:%M:%S')}] {args[0]}")


def main():
    """サーバーを起動"""
    print(f"Dotfiles Dashboard Server")
    print(f"=" * 40)
    print(f"Dotfiles: {DOTFILES_DIR}")
    print(f"Static:   {STATIC_DIR}")
    print(f"Backups:  {BACKUPS_DIR}")
    print(f"=" * 40)
    print(f"Starting server at http://{HOST}:{PORT}")
    print(f"Press Ctrl+C to stop")
    print()

    server = HTTPServer((HOST, PORT), DotfilesHandler)
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\nShutting down...")
        server.shutdown()


if __name__ == "__main__":
    main()
