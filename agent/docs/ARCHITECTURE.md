# アーキテクチャ設計

## システム概要

```
┌─────────────────────────────────────────────────────────────────┐
│                        Agent System                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │   データ層   │───▶│   分析層    │───▶│  実行層      │      │
│  │              │    │              │    │              │      │
│  │ - History    │    │ - Scripts   │    │ - Agents    │      │
│  │ - Metrics    │    │ - Reports   │    │ - Skills    │      │
│  │ - Archives   │    │ - Patterns  │    │ - Functions │      │
│  └──────────────┘    └──────────────┘    └──────────────┘      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## ディレクトリ構造（将来像）

```
agent/
├── README.md
├── docs/
│   ├── STRATEGY.md
│   ├── ROADMAP.md
│   └── ARCHITECTURE.md
│
├── scripts/                    # 分析スクリプト
│   ├── history-stats          # 基本統計
│   ├── history-patterns       # パターン検出
│   └── weekly-report          # 週次レポート
│
├── agents/                     # エージェント定義
│   ├── general/               # 汎用（シェルスクリプト）
│   │   ├── git-sync.sh
│   │   ├── docker-refresh.sh
│   │   └── dev-env-check.sh
│   │
│   └── specialized/           # 専門（AI駆動）
│       ├── project-setup/
│       └── dependency-update/
│
├── skills/                     # Claude Code スキル
│   ├── analyze-history/
│   └── suggest-automation/
│
└── config/                     # 設定ファイル
    ├── patterns.yaml          # 検出パターン定義
    └── agents.yaml            # エージェント設定
```

---

## データフロー

### 1. データ収集

```
┌─────────────┐
│ Terminal    │
│ (zsh)       │
└──────┬──────┘
       │ コマンド実行
       ▼
┌─────────────┐     ┌─────────────────────┐
│ preexec()   │────▶│ ~/.command_metrics  │
│ precmd()    │     │ (実行メトリクス)     │
└──────┬──────┘     └─────────────────────┘
       │
       ▼
┌─────────────┐     ┌─────────────────────┐
│ HISTFILE    │────▶│ ~/.zsh_history      │
│             │     │ (EXTENDED形式)       │
└─────────────┘     └──────────┬──────────┘
                               │ 月次
                               ▼
                    ┌─────────────────────┐
                    │ ~/.zsh_history_     │
                    │ archive/            │
                    └─────────────────────┘
```

### 2. 分析処理

```
┌─────────────────────┐
│ ~/.zsh_history      │
│ ~/.command_metrics  │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ history-stats       │───▶ 頻出コマンド
│ history-patterns    │───▶ 連続パターン
│ weekly-report       │───▶ 週次サマリー
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│ 自動化候補リスト    │
└─────────────────────┘
```

### 3. エージェント実行

```
┌─────────────┐
│ ユーザー    │
│ トリガー    │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────────┐
│           エージェント選択               │
├─────────────────────────────────────────┤
│                                          │
│  パターンマッチ?                         │
│       │                                  │
│  ┌────┴────┐                            │
│  ▼         ▼                            │
│ Yes       No                            │
│  │         │                            │
│  ▼         ▼                            │
│ Shell    AI Agent                       │
│ Script   (Haiku/Sonnet/Opus)           │
│                                          │
└─────────────────────────────────────────┘
```

---

## コンポーネント詳細

### データ層

#### ヒストリーファイル (`~/.zsh_history`)

**形式:** EXTENDED_HISTORY
```
: 1703846400:0;git status
: 1703846410:0;git add .
: 1703846420:0;git commit -m "fix bug"
```

**フィールド:**
- タイムスタンプ（Unix時間）
- 実行時間（秒）
- コマンド

#### メトリクスログ (`~/.command_metrics.log`)（将来）

**形式:**
```
timestamp:duration:pwd:command
1703846400:0:/Users/atusi/project:git status
1703846410:2:/Users/atusi/project:npm test
```

**フィールド:**
- 実行時刻
- 実行時間
- 作業ディレクトリ
- コマンド

### 分析層

#### history-stats

**入力:** ヒストリーファイル
**出力:**
```
=== 頻出コマンド TOP 20 ===
  145 git
   98 cd
   87 ls
   ...

=== 日別推移 ===
2024-12-28: 89 commands
2024-12-29: 124 commands
```

#### history-patterns

**入力:** ヒストリーファイル
**出力:**
```
=== 連続パターン ===
git add → git commit: 45回
cd → ls: 38回
npm install → npm run: 22回

=== 自動化候補 ===
1. git add + commit パターン → alias推奨
2. docker down + up パターン → スクリプト推奨
```

### 実行層

#### 汎用エージェント（シェルスクリプト）

```bash
# git-sync.sh
#!/bin/bash
set -e

echo "=== Git Sync ==="
git fetch --all
git pull
git status

echo "=== Complete ==="
```

#### 専門エージェント（AI駆動）

```yaml
# agents/specialized/project-setup/config.yaml
name: project-setup
description: 新規プロジェクトの初期化
model: haiku  # コスト優先
triggers:
  - "新規プロジェクト"
  - "プロジェクト作成"
steps:
  - analyze_requirements
  - create_structure
  - init_git
  - setup_dependencies
```

---

## インターフェース

### CLI

```bash
# 分析スクリプト
./agent/scripts/history-stats
./agent/scripts/history-patterns
./agent/scripts/weekly-report

# エージェント実行
./agent/agents/general/git-sync.sh
```

### Claude Code スキル

```
# ヒストリー分析
/analyze-history

# 自動化提案
/suggest-automation
```

### シェル関数/エイリアス

```bash
# .zshrc に追加
alias gsync='git fetch --all && git pull && git status'
alias drestart='docker-compose down && docker-compose up -d'

# 関数
agent-run() {
    local agent=$1
    shift
    ~/dotfiles/agent/agents/general/${agent}.sh "$@"
}
```

---

## 設定ファイル

### patterns.yaml（将来）

```yaml
# 検出するパターン定義
patterns:
  git-workflow:
    sequence: ["git add", "git commit"]
    threshold: 10  # 10回以上で検出
    suggest: "alias gac='git add . && git commit'"

  docker-restart:
    sequence: ["docker-compose down", "docker-compose up"]
    threshold: 5
    suggest: "alias drestart='docker-compose down && docker-compose up -d'"
```

### agents.yaml（将来）

```yaml
# エージェント設定
agents:
  git-sync:
    type: shell
    path: agents/general/git-sync.sh
    description: Git同期

  project-setup:
    type: ai
    model: haiku
    path: agents/specialized/project-setup/
    description: 新規プロジェクト初期化
```

---

## セキュリティ考慮事項

| 項目 | リスク | 対策 |
|------|--------|------|
| ヒストリー | 機密情報の記録 | 環境変数経由でシークレット管理 |
| メトリクス | 行動追跡 | ローカル保存のみ |
| エージェント | 意図しない実行 | 確認プロンプト、dry-run |

---

## 拡張ポイント

### MCP サーバー化（将来）

```
┌─────────────┐     ┌─────────────┐
│ Claude Code │────▶│ MCP Server  │
│             │◀────│ (Agent)     │
└─────────────┘     └──────┬──────┘
                           │
                    ┌──────┴──────┐
                    ▼             ▼
              ┌─────────┐  ┌─────────┐
              │ Scripts │  │ Agents  │
              └─────────┘  └─────────┘
```

### IDE 統合（将来）

- VS Code 拡張
- Cursor 統合
- ステータスバー表示

---

## 更新履歴

| 日付 | 変更内容 |
|------|---------|
| 2024-12-29 | 初版作成 |
