# dotsync — Engineering Review & Implementation Plan

## Context

/office-hours で設計した dotsync（dotfiles reverse-sync TUI ツール）のエンジニアリングレビュー。
Go + Bubble Tea v2 で、ローカルマシンの未管理設定ファイルを検知し、dotfiles リポに取り込む TUI ツール。
デザインドキュメント: `~/.gstack/projects/a2chub-profile/atusi-master-design-20260402-231130.md`

## Data Flow

```
                     ┌──────────────┐
                     │   main.go    │
                     │  parse flags │
                     └──────┬───────┘
                            │
                     ┌──────▼───────┐
                     │  tui/app.go  │
                     │  MVU model   │
                     └──────┬───────┘
                            │
         ┌──────────────────┼──────────────────┐
         │                  │                  │
   ┌─────▼──────┐   ┌──────▼──────┐   ┌──────▼──────┐
   │ config.go  │   │  brew.go    │   │ secrets.go  │
   │ parse      │   │ brew leaves │   │ filename +  │
   │ link-script│   │ vs Brewfile │   │ entropy     │
   │ walk .cfg  │   │             │   │             │
   └─────┬──────┘   └──────┬──────┘   └──────┬──────┘
         │                  │                  │
         └────────┬─────────┘                  │
                  │                            │
         ┌────────▼─────────┐                  │
         │ manifest.go      │◄─────────────────┘
         │ []ScanItem       │
         └────────┬─────────┘
                  │
         ┌────────▼─────────┐
         │ importer/        │
         │  copy files      │
         │  update script   │
         │  update Brewfile │
         │  gen commit msg  │
         └──────────────────┘
```

## Critical Discovery: Three Symlink Patterns

link-dotfiles.sh は3種類のシンボリックリンクを作成:
1. **ディレクトリ単位**: `~/.config/borders` → dotfiles (ディレクトリ全体)
2. **ファイル単位**: `~/.config/aerospace/aerospace.toml` → dotfiles (ファイルのみ)
3. **ルート dotfile**: `~/.zshrc` → dotfiles (ホーム直下)

さらにドリフト検知が必要: zellij, jj は link-dotfiles.sh に記載があるのに実体ディレクトリのまま。

## Dependencies

```
require (
    charm.land/bubbletea/v2    // Bubble Tea v2 (2026年3月)
    github.com/charmbracelet/lipgloss/v2
    github.com/charmbracelet/bubbles/v2
)
```

外部依存はこの3つのみ。Shannon entropy は自前実装 (log2, byte-level)。コマンドライン引数は stdlib `flag`。

## Build Order

```
Phase 1: scanner/secrets.go     ← 依存ゼロ、純粋ロジック
Phase 2: scanner/config.go      ← os, regexp のみ
Phase 3: scanner/brew.go        ← os/exec
Phase 4: manifest/manifest.go   ← scanner パッケージに依存
Phase 5: importer/*             ← manifest に依存
Phase 6: tui/styles.go          ← lipgloss のみ
Phase 7: tui/app.go + 画面      ← manifest, importer, bubbletea
Phase 8: cmd/dotsync/main.go    ← tui
```

## Review Decisions

- **#1 リポ配置**: dotfiles リポ内 `tools/dotsync/` に配置（Phase 3 で別リポに切り出し可能）
- **#2 Bubble Tea バージョン**: v2 を採用（学習目的、最新機能活用）
- **#3 manifest/brewfile.go**: 削除。Brewfile 関連は scanner + importer の2箇所に統合（DRY）
- **#4 テスト**: ハッピーパスのみ。エッジケースは後から追加
- **#5 冪等性**: importer に追記前の既存エントリチェックを追加（Codex 指摘）

## File Structure

```
tools/dotsync/
├── go.mod
├── go.sum
├── cmd/dotsync/main.go
├── internal/
│   ├── scanner/
│   │   ├── config.go          # link-dotfiles.sh パース + ~/.config 走査
│   │   ├── config_test.go
│   │   ├── brew.go            # brew leaves vs Brewfile 差分
│   │   ├── brew_test.go
│   │   ├── secrets.go         # 秘密情報検知
│   │   └── secrets_test.go
│   ├── manifest/
│   │   ├── manifest.go        # 統合データ構造
│   │   └── manifest_test.go
│   ├── importer/
│   │   ├── importer.go        # インポート実行
│   │   ├── importer_test.go
│   │   ├── linker.go          # link-dotfiles.sh 追記
│   │   ├── linker_test.go
│   │   ├── brewfile.go        # Brewfile 追記
│   │   ├── brewfile_test.go
│   │   ├── commit.go          # コミットメッセージ生成
│   │   └── commit_test.go
│   └── tui/
│       ├── app.go             # Bubble Tea ルートモデル
│       ├── app_test.go
│       ├── scan.go            # スキャン画面
│       ├── detail.go          # 詳細画面
│       ├── import_view.go     # インポート確認画面
│       └── styles.go          # Lip Gloss スタイル
└── testdata/
    ├── link-dotfiles.sh       # テスト用フィクスチャ
    ├── Brewfile               # テスト用フィクスチャ
    └── sample-config/         # モック ~/.config ツリー
```

## Key Implementation Details

### scanner/config.go — ParseLinkScript
- 正規表現: `create_link\s+"([^"]+)"\s+"([^"]+)"`
- `$DOTFILES_DIR` → 実パス置換、`$HOME` → 実パス置換
- `$1` も `$DOTFILES_DIR` と同様に扱う（スクリプト冒頭で `DOTFILES_DIR="$1"` のため）
- `#` で始まる行はスキップ
- パースエラー時は追記専用モードにフォールバック

### scanner/config.go — ScanConfigDir
分類ロジック:
1. `os.Lstat()` でシンボリックリンクチェック
2. シンボリックリンク → ターゲットが dotfiles 配下? → MANAGED
3. 実ディレクトリ → managed map に存在? → DRIFTED
4. 実ディレクトリ → 中のファイルにシンボリックリンクあり? → PARTIALLY_MANAGED
5. 上記いずれでもない → UNMANAGED

### scanner/secrets.go — ShannonEntropy
```
1. 256バケットでバイト頻度をカウント
2. total = len(data)
3. count > 0 の各バケット: p = count/total, entropy -= p * log2(p)
4. 戻り値: 0.0-8.0 (bits)
```
閾値 4.5。バイナリファイル、空ファイル、16バイト未満、1MB超はスキップ。

### importer/linker.go — AppendLinkEntry
- `print_success` 行の前に挿入
- 追加前に既存エントリチェック（冪等性）
- アトミック書き込み（temp → rename）

### tui/app.go — State Machine
```
ScreenScan  --[完了]--> ScreenList
ScreenList  --[Enter]--> ScreenDetail
ScreenList  --[i]------> ScreenImport
ScreenDetail --[Esc]----> ScreenList
ScreenImport --[確認]---> 実行 → ScreenDone
ScreenDone   --[q]------> 終了、commit msg を stdout に出力
```

## Test Strategy

### 高優先テスト
- `ShannonEntropy`: 既知値テスト (all-zero=0.0, uniform=8.0, english~4.0, base64~5.5)
- `ParseLinkScript`: 実際の link-dotfiles.sh をフィクスチャとして使用
- `ScanConfigDir`: t.TempDir() でモック ~/.config を作成、3パターンのシンボリックリンクをテスト
- `ParseBrewfile`: 実際の Brewfile をフィクスチャとして使用
- `AppendLinkEntry`: 冪等性テスト

### テスタビリティ設計
- `RunBrewLeaves` は os/exec を呼ぶため、パース部分を `ParseLeaves(output string) []string` に分離してテスト可能にする
- TUI テストは最小限: `Update` 関数のステート遷移のみ

## Verification

```bash
# Phase 1-3: 各パッケージのテスト
cd dotsync && go test ./internal/scanner/...
cd dotsync && go test ./internal/manifest/...
cd dotsync && go test ./internal/importer/...

# Phase 7-8: 統合テスト
cd dotsync && go build ./cmd/dotsync
./dotsync -dotfiles ~/dotfiles -dry-run  # dry-run で動作確認

# 手動 E2E: 実際にインポートしてみる
./dotsync -dotfiles ~/dotfiles
# TUI で ghostty を選択してインポート
# 確認: config/ghostty/ が存在、link-dotfiles.sh にエントリ追加、git diff で確認
```

## NOT in scope (Phase 1)

- `.dotsync.toml` 設定ファイル対応（Phase 2）
- dry-run モード（Phase 2）
- コミットメッセージ生成の stdout 出力（Phase 2）
- 壊れたシンボリックリンクの修復（Phase 2）
- goreleaser / GitHub Actions（Phase 3）
- 汎用化（他のリポ構造対応）（Phase 3）
- AI による取込み推奨（Phase 4+）

## What already exists

- `scripts/link-dotfiles.sh` — dotsync はこれをパースして管理対象を把握し、インポート時に追記する
- `packages/Brewfile` — dotsync はこれをパースして差分を検知し、インポート時に追記する
- `scripts/lib/colors.sh` — dotsync は Go なので直接使わないが、コメント規約（日本語）は踏襲する
- `.gitignore` — 秘密情報パターンが定義済み。dotsync の DangerousPatterns と同期すべき

## Parallelization Strategy

| Step | Modules | Depends on |
|------|---------|------------|
| scanner/secrets | scanner/ | — |
| scanner/config | scanner/ | — |
| scanner/brew | scanner/ | — |
| manifest | manifest/ | scanner |
| importer | importer/ | manifest |
| tui | tui/ | manifest, importer |
| main | cmd/ | tui |

Lane A: scanner/secrets + scanner/config (parallel, same package but independent files)
Lane B: scanner/brew (parallel with Lane A, independent)
→ merge → manifest → importer → tui → main (sequential)

3 scanner ファイルは並列実装可能。それ以降は順次。

## Completion Summary

- Step 0: Scope Challenge — scope accepted as-is (新規プロジェクト、12ファイルは Go の慣習として妥当)
- Architecture Review: 2 issues found (リポ配置, BT バージョン) → both resolved
- Code Quality Review: 1 issue found (DRY 違反) → resolved
- Test Review: diagram produced, happy path のみ (ユーザー選択)
- Performance Review: 0 issues found
- NOT in scope: written
- What already exists: written
- TODOS.md updates: 2 items proposed (watch_dirs 拡張, ドリフト修復)
- Failure modes: 0 critical gaps (CopyDir の既存ファイル衝突は error で処理)
- Outside voice: ran (codex) → 12 findings, 2 accepted (冪等性, パスエスケープ)
- Parallelization: 2 lanes parallel (scanner), 5 sequential (manifest→tui)
- Lake Score: 4/5 recommendations chose complete option

## GSTACK REVIEW REPORT

| Review | Trigger | Why | Runs | Status | Findings |
|--------|---------|-----|------|--------|----------|
| CEO Review | `/plan-ceo-review` | Scope & strategy | 0 | — | — |
| Codex Review | `/codex review` | Independent 2nd opinion | 1 | ISSUES_FOUND | 12 findings, 2 accepted |
| Eng Review | `/plan-eng-review` | Architecture & tests (required) | 1 | CLEAR (PLAN) | 3 issues, 0 critical gaps |
| Design Review | `/plan-design-review` | UI/UX gaps | 0 | — | — |

**VERDICT:** ENG REVIEW CLEARED — ready to implement
