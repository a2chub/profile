# Agent Automation Project

ヒストリーデータを活用した作業自動化エージェントシステム

> **現在のステータス**: Phase 0 完了、Phase 1 メトリクス収集実装済み。
> 分析スクリプト（history-stats等）は延期中。
> 詳細は [ROADMAP.md](docs/ROADMAP.md) を参照。

## 概要

コマンドヒストリーの分析により繰り返しパターンを検出し、
ルーティーン作業を自動化するエージェントを構築・運用するプロジェクト。

## 目的

1. **作業効率化** - 反復的なコマンド操作の自動化
2. **コスト最適化** - 適切なツール選択（シェル vs AI）
3. **品質向上** - 標準化されたワークフローの実行
4. **知識蓄積** - 作業パターンの可視化と改善

## ディレクトリ構造

```
agent/
├── README.md           # このファイル
├── docs/               # ドキュメント
│   ├── STRATEGY.md     # 戦略・方針
│   ├── ROADMAP.md      # 実装ロードマップ
│   └── ARCHITECTURE.md # アーキテクチャ設計
├── scripts/            # 分析・ユーティリティスクリプト（将来）
├── agents/             # エージェント定義（将来）
└── skills/             # Claude Code スキル（将来）
```

## クイックスタート

```bash
# ヒストリー統計を確認（将来実装）
./scripts/history-stats

# パターン分析レポート（将来実装）
./scripts/history-patterns
```

## ドキュメント

- [戦略・方針](docs/STRATEGY.md) - プロジェクトの戦略と原則
- [実装ロードマップ](docs/ROADMAP.md) - フェーズ別の実装計画
- [アーキテクチャ](docs/ARCHITECTURE.md) - 技術設計と構成

## 関連設定

- ヒストリー設定: `~/.zshrc`（EXTENDED_HISTORY有効）
- アーカイブ先: `~/.zsh_history_archive/`

## ライセンス

Private - Personal Use
