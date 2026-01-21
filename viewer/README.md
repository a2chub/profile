# Dotfiles Dashboard

dotfilesリポジトリで管理されている設定ファイルとBrewパッケージを可視化・編集できるWebダッシュボード。

## 起動方法

```bash
# リポジトリルートから
./start-viewer.sh

# または直接
python3 viewer/server.py
```

ブラウザで http://localhost:8765 にアクセス。

## 機能

### 設定ファイル管理
- 8つの設定ファイルを一覧表示
- 内容の閲覧・編集
- 保存前に自動バックアップ（`.backups/`に保存）
- 各ソフトウェアの公式ドキュメント/GitHubへのリンク

### Brewパッケージカタログ
- Brewfileをパースして表示
- Formulae/Casks/Tapsのカテゴリ分け
- 検索機能

## 管理対象ファイル

| ファイル | ソフトウェア | 形式 |
|---------|-------------|------|
| `.zshrc` | Zsh | shell |
| `.tmux.conf` | tmux | conf |
| `.vimrc` | Vim | vim |
| `config/nvim/init.lua` | Neovim | lua |
| `config/starship.toml` | Starship | toml |
| `config/aerospace/aerospace.toml` | AeroSpace | toml |
| `config/borders/bordersrc` | JankyBorders | bash |
| `packages/Brewfile` | Homebrew | ruby |

## 設定

環境変数でポートを変更可能：

```bash
DOTFILES_VIEWER_PORT=9000 ./start-viewer.sh
```

## セキュリティ

- `127.0.0.1`のみにバインド（ローカルアクセス限定）
- 編集可能なファイルはホワイトリスト方式
- ディレクトリトラバーサル防止
- TOML形式は保存前に構文検証
- 保存前に自動バックアップ

## 技術スタック

- Python 3.9+（標準ライブラリのみ）
- HTML + CSS + バニラJavaScript
- 外部依存なし
