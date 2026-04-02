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
