# Dotfiles Dashboard/Catalog Viewer 計画

## 概要

dotfilesリポジトリで管理されている設定ファイルとBrewパッケージを可視化・編集できるWebダッシュボードを作成する。

## 要件

- 設定ファイルの一覧表示（ソフトウェア名、パス、形式）
- Brewパッケージの一覧表示（Formulae/Casks/Taps）
- 各ソフトウェアの本家リポジトリ/ドキュメントへのリンク
- **設定ファイルの編集・保存機能**
- 軽量実装：Python標準ライブラリ + HTML + バニラJS

---

## ディレクトリ構成

```
dotfiles/
├── viewer/                      # ダッシュボードアプリ
│   ├── server.py                # Pythonサーバー（API + 静的ファイル配信）
│   ├── static/
│   │   ├── index.html           # メインHTML
│   │   ├── style.css            # スタイルシート
│   │   └── app.js               # JavaScriptロジック
│   └── README.md                # ビューワー詳細ドキュメント
├── start-viewer.sh              # 起動スクリプト（トップレベル）
└── ...
```

---

## 管理対象の設定ファイル

| ID | ファイル | ソフトウェア | カテゴリ | 形式 | 公式ドキュメント |
|----|---------|-------------|---------|------|-----------------|
| zshrc | `.zshrc` | Zsh | Shell | shell | https://zsh.sourceforge.io/Doc/ |
| tmux | `.tmux.conf` | tmux | Terminal | conf | https://github.com/tmux/tmux/wiki |
| vimrc | `.vimrc` | Vim | Editor | vim | https://vimdoc.sourceforge.net/ |
| nvim | `config/nvim/init.lua` | Neovim | Editor | lua | https://neovim.io/doc/user/ |
| starship | `config/starship.toml` | Starship | Shell | toml | https://starship.rs/config/ |
| aerospace | `config/aerospace/aerospace.toml` | AeroSpace | WM | toml | https://nikitabobko.github.io/AeroSpace/guide |
| borders | `config/borders/bordersrc` | Borders | WM | bash | https://github.com/FelixKratz/JankyBorders |
| brewfile | `packages/Brewfile` | Homebrew | Package | ruby | https://docs.brew.sh/ |

---

## APIエンドポイント

| Method | Path | 説明 |
|--------|------|------|
| GET | `/` | index.html配信 |
| GET | `/static/*` | 静的ファイル配信 |
| GET | `/api/configs` | 設定ファイル一覧 |
| GET | `/api/configs/<id>` | 設定ファイル内容取得 |
| PUT | `/api/configs/<id>` | 設定ファイル保存（バックアップ作成） |
| GET | `/api/brew/formulae` | Formulae一覧 |
| GET | `/api/brew/casks` | Casks一覧 |
| GET | `/api/brew/taps` | Taps一覧 |

---

## 画面構成

```
+------------------------------------------------------------------+
|  Dotfiles Dashboard                                    [Reload]   |
+------------------------------------------------------------------+
| +------------------+  +----------------------------------------+ |
| | Navigation       |  | Content Area                           | |
| |                  |  |                                        | |
| | [Configs]        |  | File: .zshrc                           | |
| |   .zshrc         |  | Software: Zsh  |  Format: Shell        | |
| |   .tmux.conf     |  | [Docs] [GitHub]                        | |
| |   .vimrc         |  +----------------------------------------+ |
| |   nvim/init.lua  |  |                                        | |
| |   starship.toml  |  | +------------------------------------+ | |
| |   aerospace.toml |  | | # Editor                           | | |
| |   bordersrc      |  | |                                    | | |
| |   Brewfile       |  | | export PATH="..."                  | | |
| |                  |  | | alias vim=nvim                     | | |
| | [Packages]       |  | | ...                                | | |
| |   > Formulae     |  | |                                    | | |
| |   > Casks        |  | +------------------------------------+ | |
| |   > Taps         |  | [Save] [Discard] [Download]            | |
| +------------------+  +----------------------------------------+ |
+------------------------------------------------------------------+
```

---

## セキュリティ対策

1. **ホワイトリスト方式**: 編集可能なファイルを事前定義
2. **パス検証**: `..`によるディレクトリトラバーサル防止
3. **localhost限定**: `127.0.0.1`にのみバインド
4. **保存前バックアップ**: `.backups/`に自動バックアップ
5. **TOML構文検証**: starship.toml等は保存前にパース検証

---

## 実装タスク

### Phase 1: 基盤構築
- [ ] `viewer/`ディレクトリ作成
- [ ] `server.py`実装（http.server使用）
- [ ] `static/index.html`作成（基本レイアウト）
- [ ] `static/style.css`作成
- [ ] `static/app.js`作成

### Phase 2: 設定ファイル表示
- [ ] `/api/configs`エンドポイント実装
- [ ] `/api/configs/<id>`エンドポイント実装
- [ ] サイドバーナビゲーション実装
- [ ] 設定ファイル詳細表示（メタデータ + 内容）
- [ ] 外部ドキュメントリンク表示

### Phase 3: Brewパッケージ表示
- [ ] Brewfileパーサー実装
- [ ] `/api/brew/*`エンドポイント実装
- [ ] パッケージ一覧表示（カテゴリ分け）
- [ ] 検索/フィルター機能

### Phase 4: 編集機能
- [ ] PUT `/api/configs/<id>`実装
- [ ] バックアップ機構実装
- [ ] 編集UI実装（textarea）
- [ ] 保存/破棄ボタン実装
- [ ] TOML構文検証

### Phase 5: 仕上げ
- [ ] `start-viewer.sh`作成
- [ ] `viewer/README.md`作成
- [ ] レスポンシブデザイン調整
- [ ] エラーハンドリング強化

---

## 修正・作成ファイル一覧

### 新規作成

| ファイル | 内容 |
|---------|------|
| `viewer/server.py` | Pythonサーバー（約200行） |
| `viewer/static/index.html` | メインHTML（約100行） |
| `viewer/static/style.css` | スタイルシート（約150行） |
| `viewer/static/app.js` | JavaScriptロジック（約300行） |
| `viewer/README.md` | ビューワー詳細ドキュメント |
| `start-viewer.sh` | 起動スクリプト |

### 既存ファイル（参照のみ）

- `.zshrc`, `.tmux.conf`, `.vimrc` - 編集対象
- `config/nvim/init.lua`, `config/starship.toml` - 編集対象
- `config/aerospace/aerospace.toml`, `config/borders/bordersrc` - 編集対象
- `packages/Brewfile` - パース対象・編集対象

---

## 検証方法

```bash
# 1. サーバー起動
./start-viewer.sh

# 2. ブラウザでアクセス
open http://localhost:8765

# 3. 確認項目
# - 設定ファイル一覧が表示される
# - 各ファイルの内容が表示される
# - 外部リンクが正しく開く
# - ファイル編集・保存ができる
# - .backups/にバックアップが作成される
# - Brewパッケージ一覧が表示される
```

---

## 技術仕様

- **Python**: 3.9+（標準ライブラリのみ）
- **ポート**: 8765（環境変数`DOTFILES_VIEWER_PORT`で変更可）
- **対応ブラウザ**: モダンブラウザ（Chrome, Firefox, Safari, Edge）
