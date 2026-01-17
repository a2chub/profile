# dotfiles リポジトリ リファクタリング計画

> **ステータス**: Phase 1 + Phase 2 実装完了（2026-01-17）

## 現状分析サマリー

### 強み
- XDG準拠のconfig構成
- OS判別機能を備えたスクリプト（setup.sh）
- シンボリックリンク自動管理とバックアップ機能
- コマンドメトリクス収集基盤（preexec/precmd）
- 詳細な戦略ドキュメント（STRATEGY.md, ARCHITECTURE.md）

### 課題
- スクリプト・設定ファイルの配置が一貫していない
- エラーハンドリング・セキュリティが不十分
- ドキュメントと実装の乖離
- 長期運用に必要な機能（ロールバック、ログ、テスト）が不足

---

## 優先度別 改善リスト

### P0: Critical（セキュリティ・安定性に直結）

| # | 問題 | 改善内容 | 影響ファイル |
|---|------|---------|-------------|
| 1 | **エラーハンドリング不足** | `step1_app-install.sh`に`set -e`とエラートラップ追加 | `step1_app-install.sh` |
| 2 | **SSH設定の冪等性違反** | 空行追加ロジックを修正、ファイル存在確認を追加 | `step2_make-ssh-key.sh` |
| 3 | **外部スクリプト実行リスク** | Homebrew/mise等のインストールにチェックサム検証追加 | `step1_app-install.sh` |
| 4 | **.gitignore不完全** | `*.log`, `firebase-debug.log`, `.DS_Store`を追加 | `.gitignore` |

### P1: High（構造的問題・長期メンテナンス性）

| # | 問題 | 改善内容 | 影響 |
|---|------|---------|------|
| 5 | **スクリプト配置の不統一** | トップレベルの`step*.sh`を`scripts/setup/`に移動 | ディレクトリ構造 |
| 6 | **セットアップフローの曖昧さ** | `setup.sh`を唯一のエントリーポイントに統一 | `setup.sh`, `README.md` |
| 7 | **ドキュメントと実装の乖離** | `agent/docs/ROADMAP.md`を現状に合わせて更新 | `ROADMAP.md` |
| 8 | **トラブルシューティング不足** | `docs/TROUBLESHOOTING.md`を新規作成 | 新規ファイル |
| 9 | **設定ファイルの散在** | `.zshrc`等を`config/shell/`に移動し、リンク構成を更新 | 複数ファイル |

### P2: Medium（運用効率・可視性向上）

| # | 問題 | 改善内容 | 影響 |
|---|------|---------|------|
| 10 | **ロールバック機能なし** | `--dry-run`オプションと復元スクリプトを追加 | セットアップスクリプト |
| 11 | **ログ出力なし** | 統合ログファイル（`~/.dotfiles/logs/`）を導入 | セットアップスクリプト |
| 12 | **セットアップ検証なし** | `scripts/verify-setup.sh`を新規作成 | 新規ファイル |
| 13 | **環境要件が不明確** | `docs/REQUIREMENTS.md`を新規作成 | 新規ファイル |
| 14 | **OS別設定の分離不足** | `config/`をcommon/macos/linux構造に再編 | ディレクトリ構造 |

### P3: Low（最適化・将来対応）

| # | 問題 | 改善内容 | 影響 |
|---|------|---------|------|
| 15 | **環境変数の一元管理** | `.zshrc`から環境変数を`config/shell/env.sh`に分離 | `.zshrc` |
| 16 | **agent/の空実装** | Phase 1スクリプト（history-stats等）を実装または計画修正 | `agent/scripts/` |
| 17 | **メトリクスログ管理** | ログローテーション機構を追加 | `.zshrc` |
| 18 | **マルチマシン同期** | `docs/MULTI_MACHINE_SETUP.md`を新規作成 | 新規ファイル |
| 19 | **テスト自動化** | Docker/CIでのセットアップテストを整備 | `test/`, CI設定 |

---

## 推奨ディレクトリ構造（リファクタリング後）

```
dotfiles/
├── README.md                      # クイックスタート
├── setup.sh                       # 唯一のエントリーポイント
├── .gitignore                     # 更新版
│
├── config/                        # 設定ファイル（XDG準拠）
│   ├── shell/                     # シェル設定
│   │   ├── zshrc                  # .zshrc → ここに移動
│   │   ├── env.sh                 # 環境変数（分離）
│   │   └── aliases.sh             # エイリアス（分離）
│   ├── editor/                    # エディタ設定
│   │   ├── nvim/
│   │   └── vim/
│   ├── terminal/                  # ターミナル関連
│   │   ├── starship.toml
│   │   └── tmux.conf
│   └── macos/                     # macOS専用
│       ├── aerospace/
│       ├── borders/
│       └── sketchybar/
│
├── scripts/                       # セットアップ・ユーティリティ
│   ├── setup/                     # セットアップ関連
│   │   ├── install-packages.sh
│   │   ├── install-tools.sh
│   │   ├── link-dotfiles.sh
│   │   ├── setup-ssh.sh           # step2 → リネーム
│   │   └── install-docker.sh      # stepx → リネーム
│   ├── verify-setup.sh            # 新規：検証スクリプト
│   └── utils/                     # ユーティリティ
│
├── packages/                      # パッケージ定義
│   ├── Brewfile                   # mac/ → ここに移動
│   ├── apt-packages.txt           # 新規：Debian/Ubuntu用
│   └── dnf-packages.txt           # 新規：Fedora用
│
├── docs/                          # ドキュメント
│   ├── TROUBLESHOOTING.md         # 新規
│   ├── REQUIREMENTS.md            # 新規
│   ├── MULTI_MACHINE_SETUP.md     # 新規
│   └── plans/
│
├── agent/                         # エージェント自動化（既存）
│   ├── README.md
│   ├── docs/
│   └── scripts/
│
└── test/                          # テスト
    ├── Dockerfile
    └── test_setup.sh
```

---

## 実装フェーズ

### Phase 1: 安定化（P0対応）
1. `set -e`とエラートラップの追加
2. SSH設定スクリプトの冪等性修正
3. `.gitignore`の更新
4. 不要ファイル（firebase-debug.log等）の削除

### Phase 2: 構造整理（P1対応）
1. スクリプトを`scripts/setup/`に移動
2. `setup.sh`の統合・更新
3. ドキュメント更新（ROADMAP.md, TROUBLESHOOTING.md）
4. 設定ファイルの再配置

### Phase 3: 運用機能追加（P2対応）
1. ログ機能の追加
2. `--dry-run`オプションの実装
3. 検証スクリプトの作成
4. 環境要件ドキュメントの作成

### Phase 4: 最適化（P3対応）
1. 環境変数の分離
2. agent/スクリプトの実装
3. CI/CD統合

---

## 検証方法

### セットアップ後の確認項目
```bash
# 1. シンボリックリンクの確認
ls -la ~/.zshrc ~/.config/nvim ~/.config/starship.toml

# 2. エイリアスの動作確認
vim --version | head -1  # Neovimが起動するか
python --version         # Python 3.14か

# 3. 環境変数の確認
echo $PNPM_HOME
echo $BUN_INSTALL

# 4. メトリクス収集の確認
cat ~/.command_metrics.log | tail -5
```

### Dockerでのテスト
```bash
cd test/
docker build -t dotfiles-test .
docker run -it dotfiles-test ./setup.sh
```

---

## 決定事項

- **リファクタリング範囲**: P0 + P1（安定化 + 構造整理まで）
- **step*.shスクリプト**: `scripts/setup/`に移動して保持
- **agent/プロジェクト**: ドキュメントを現状に合わせる（未実装部分を延期として記載）

---

## 実装タスク一覧

### Phase 1: 安定化（P0） ✅ 完了

- [x] `step1_app-install.sh`: `set -e`とエラートラップ追加
- [x] `step2_make-ssh-key.sh`: 冪等性修正（ファイル存在確認、空行問題）
- [x] `.gitignore`: `*.log`, `firebase-debug.log`追加
- [x] 不要ファイル削除: `firebase-debug.log`, `.zshrc.org`

### Phase 2: 構造整理（P1） ✅ 完了

- [x] ディレクトリ作成: `scripts/setup/`, `packages/`
- [x] スクリプト移動・リネーム:
  - `step1_app-install.sh` → `scripts/setup/install-apps.sh`
  - `step2_make-ssh-key.sh` → `scripts/setup/setup-ssh.sh`
  - `stepx_install_docker.sh` → `scripts/setup/install-docker.sh`
  - `mac/Brewfile` → `packages/Brewfile`
- [x] `setup.sh`更新: 新しいパス構造に対応
- [x] `scripts/link-dotfiles.sh`更新: 新しい配置に対応
- [x] `README.md`更新: 新しい構造を反映
- [x] `docs/TROUBLESHOOTING.md`新規作成
- [x] `agent/docs/ROADMAP.md`更新: 未実装部分を延期として記載

### 修正対象ファイル

| ファイル | 変更内容 |
|---------|---------|
| `step1_app-install.sh` | エラーハンドリング追加 → 移動 |
| `step2_make-ssh-key.sh` | 冪等性修正 → 移動 |
| `stepx_install_docker.sh` | 移動のみ |
| `.gitignore` | パターン追加 |
| `setup.sh` | パス更新 |
| `scripts/link-dotfiles.sh` | パス更新 |
| `README.md` | 構造説明更新 |
| `agent/docs/ROADMAP.md` | 進捗更新 |

### 新規作成ファイル

| ファイル | 内容 |
|---------|------|
| `docs/TROUBLESHOOTING.md` | よくある問題と対処法 |
| `scripts/setup/` | ディレクトリ |
| `packages/` | ディレクトリ |
