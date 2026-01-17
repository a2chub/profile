# Dotfiles

macOS / Linux 用の設定ファイル管理リポジトリ

## クイックスタート

```bash
git clone https://github.com/a2chub/profile.git ~/dotfiles
cd ~/dotfiles
./setup.sh
```

## セットアップオプション

| オプション | 説明 |
|-----------|------|
| `--skip-packages` | パッケージインストールをスキップ |
| `--skip-links` | シンボリックリンク作成をスキップ |
| `--skip-tools` | ツールインストールをスキップ |
| `--links-only` | シンボリックリンクのみ作成 |
| `--setup-ssh` | GitHub用SSH鍵を設定 |
| `--setup-docker` | Dockerをインストール |
| `--install-apps` | レガシーアプリインストールスクリプトを実行 |

```bash
# 例: シンボリックリンクのみ作成
./setup.sh --links-only

# 例: SSH設定も実行
./setup.sh --setup-ssh
```

## 対応OS

- macOS (Homebrew)
- Ubuntu / Debian (apt)
- Fedora / RHEL (dnf)

## 含まれる設定ファイル

| ファイル | 説明 |
|----------|------|
| `.zshrc` | Zsh設定 (Starship, エイリアス, PATH) |
| `.tmux.conf` | tmux設定 (vi風キーバインド) |
| `.vimrc` | Vim設定 (日本語エンコーディング対応) |
| `config/nvim/` | Neovim設定 (Lua) |
| `config/starship.toml` | Starshipプロンプト設定 (Gruvbox Dark) |

## インストールされるパッケージ

### 基本ツール
- git, curl, wget
- tmux, neovim
- ripgrep, fd, jq, tree, htop

### 追加ツール
- [Starship](https://starship.rs/) - クロスシェルプロンプト
- [vim-jetpack](https://github.com/tani/vim-jetpack) - Vimプラグインマネージャー
- Nerd Fonts (Hack) - アイコン表示用フォント

## ディレクトリ構成

```
dotfiles/
├── setup.sh                     # メインセットアップスクリプト（唯一のエントリーポイント）
├── scripts/
│   ├── install-packages.sh      # パッケージインストール
│   ├── link-dotfiles.sh         # シンボリックリンク作成
│   ├── install-tools.sh         # 追加ツールインストール
│   └── setup/                   # オプションセットアップスクリプト
│       ├── install-apps.sh      # アプリインストール（旧step1）
│       ├── setup-ssh.sh         # SSH鍵設定（旧step2）
│       └── install-docker.sh    # Dockerインストール
├── packages/
│   └── Brewfile                 # Homebrew パッケージ定義
├── config/
│   ├── nvim/                    # Neovim設定
│   ├── starship.toml            # Starship設定
│   ├── aerospace/               # AeroSpace (macOS)
│   ├── borders/                 # Borders (macOS)
│   └── sketchybar/              # Sketchybar (macOS)
├── docs/                        # ドキュメント
│   ├── TROUBLESHOOTING.md       # トラブルシューティング
│   └── plans/                   # 計画ドキュメント
├── agent/                       # エージェント自動化プロジェクト
├── .zshrc                       # Zsh設定
├── .tmux.conf                   # tmux設定
└── .vimrc                       # Vim設定
```

## 手動セットアップ

個別にリンクを作成する場合:

```bash
ln -sf ~/dotfiles/.vimrc ~/.vimrc
ln -sf ~/dotfiles/.zshrc ~/.zshrc
ln -sf ~/dotfiles/.tmux.conf ~/.tmux.conf
ln -sf ~/dotfiles/config/nvim ~/.config/nvim
ln -sf ~/dotfiles/config/starship.toml ~/.config/starship.toml
```

## キーバインド

### Neovim / Vim
| キー | 動作 |
|------|------|
| `;` / `:` | 入れ替え |
| `Ctrl+h/j/k/l` | ペイン移動 |
| `Ctrl+n` | Neo-tree トグル |
| `<Leader>t` | ターミナルを下部に開く |
| `<Leader>ff` | Telescope ファイル検索 |

### tmux
| キー | 動作 |
|------|------|
| `Ctrl+h/j/k/l` | ペイン移動 |
| `Ctrl+b "` | 横分割 |
| `Ctrl+b %` | 縦分割 |
| `v` (コピーモード) | 選択開始 |
| `y` (コピーモード) | コピー |
