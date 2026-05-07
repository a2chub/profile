# Dotfiles リポジトリ

## プロジェクト構成
- エントリーポイント: `setup.sh`（すべてのセットアップはこのスクリプトのフラグで制御）
- `scripts/install-packages.sh` - Homebrew/apt/dnf パッケージインストール
- `scripts/install-tools.sh` - 追加ツール（Starship, lua magick バインディング等）
- `scripts/link-dotfiles.sh` - ホームディレクトリへのシンボリックリンク作成
- `packages/Brewfile` - Homebrew パッケージ定義（leaves-only 方式）
- `config/` - アプリ設定ファイル（nvim, starship, aerospace, borders, sketchybar）
- ルート直下の dotfiles: `.zshrc`, `.tmux.conf`, `.vimrc`

## 規約
- コミットメッセージ: conventional commits 形式（`feat:`, `fix:`, `refactor:` 等）
- ドキュメント・コメントは日本語で記述
- 対応 OS: macOS (Homebrew), Ubuntu/Debian (apt), Fedora/RHEL (dnf)

## 注意点
- `brew bundle` に `--no-lock` オプションは存在しない
- Brewfile の配置場所は `packages/Brewfile`（プロジェクトルートではない）
- Neovim プラグインは lazy.nvim で管理（`config/nvim/init.lua`）
- `.zshrc` の `python3` alias は Apple Silicon (`/opt/homebrew/bin/python3`) ハードコード — Intel Mac では要書き換え
- `GITHUB_PERSONAL_ACCESS_TOKEN` は起動時オートロードしない。`github_pat()` 関数経由で遅延取得

## 編集ガイドライン

### シンボリックリンクを追加するとき
`scripts/link-dotfiles.sh` の `LINKS=()` 配列に `"<dotfiles配下の相対パス>:<リンク先絶対パス>"` 形式で1行追加するだけ。順序は機能的に無関係だが、関連リンク同士はまとめて配置する。

### パッケージを追加するとき
- macOS: `packages/Brewfile`（必須）または `packages/Brewfile.optional`（`--full` 時のみ）に追加
- Debian/RHEL: `scripts/install-packages.sh` の `DEBIAN_PACKAGES` / `REDHAT_PACKAGES` 配列に追加。実体は共通の `_install_packages <check_cmd> <install_cmd> <pkg...>` ヘルパーが処理する

### シェルスクリプトの規約
- 全シェルスクリプトで `set -euo pipefail` を使用
- 共通ユーティリティは `scripts/lib/colors.sh` に集約（`print_success`/`print_warning`/`print_error`/`print_info`/`print_header`）

## Skill routing

When the user's request matches an available skill, ALWAYS invoke it using the Skill
tool as your FIRST action. Do NOT answer directly, do NOT use other tools first.
The skill has specialized workflows that produce better results than ad-hoc answers.

Key routing rules:
- Product ideas, "is this worth building", brainstorming → invoke office-hours
- Bugs, errors, "why is this broken", 500 errors → invoke investigate
- Ship, deploy, push, create PR → invoke ship
- QA, test the site, find bugs → invoke qa
- Code review, check my diff → invoke review
- Update docs after shipping → invoke document-release
- Weekly retro → invoke retro
- Design system, brand → invoke design-consultation
- Visual audit, design polish → invoke design-review
- Architecture review → invoke plan-eng-review
- Save progress, checkpoint, resume → invoke checkpoint
- Code quality, health check → invoke health
