# Git履歴表示機能の追加

## 概要

設定ファイルの過去の修正履歴（Gitコミット履歴とDiff）を表示する機能を追加する。

## 機能要件

1. 「History」ボタンをUIに追加
2. クリックで過去5件のコミット履歴を表示
3. 各コミットをクリックするとDiff内容を表示

## UI設計

```
+------------------------------------------------------------------+
| File: .zshrc                                                      |
| 📁 dotfiles/.zshrc                                                |
| [Zsh] [shell] | [Docs] [GitHub] [History]  ← 新規ボタン           |
+------------------------------------------------------------------+
| [View]  [Edit]                                                    |
+------------------------------------------------------------------+
|                                    | History Panel (右側スライド) |
|  コードビューワー                  | +-------------------------+ |
|                                    | | Commit History          | |
|                                    | | > abc123 - Fix alias    | |
|                                    | | > def456 - Add PATH     | |
|                                    | | > ghi789 - Initial      | |
|                                    | +-------------------------+ |
|                                    | | Diff View               | |
|                                    | | - old line              | |
|                                    | | + new line              | |
|                                    | +-------------------------+ |
+------------------------------------------------------------------+
```

## 実装計画

### 1. サーバー側 (server.py)

**新規APIエンドポイント:**
- `GET /api/history/<config_id>` - ファイルの履歴一覧
- `GET /api/history/<config_id>/<commit>` - 特定コミットのDiff

**新規関数:**
```python
import subprocess

def get_file_history(config_id: str, limit: int = 5) -> list:
    # git log --format="%H|%ai|%an|%s" -n limit -- <file>
    # 戻り値: [{"hash", "date", "author", "message"}, ...]

def get_commit_diff(config_id: str, commit: str) -> str:
    # git diff <commit>~1 <commit> -- <file>
    # 戻り値: diff文字列
```

### 2. フロントエンド (index.html)

**追加要素:**
- 「History」ボタン（linksセクションに追加）
- 履歴パネル（サイドパネル方式）

### 3. スタイル (style.css)

**追加スタイル:**
- `.history-panel` - 右からスライドインするパネル
- `.history-list` - コミット一覧
- `.history-item` - 各コミット行
- `.diff-view` - Diff表示領域
- `.diff-add` / `.diff-del` - 追加/削除行の色分け

### 4. JavaScript (app.js)

**追加機能:**
- `state.historyOpen` - パネル開閉状態
- `state.historyData` - 履歴データキャッシュ
- `loadHistory(configId)` - 履歴API呼び出し
- `loadDiff(configId, commit)` - Diff API呼び出し
- `renderHistoryPanel()` - パネル描画
- `toggleHistoryPanel()` - パネル開閉

---

## 修正ファイル

| ファイル | 変更内容 |
|---------|---------|
| `viewer/server.py` | subprocess import追加、履歴API追加（約50行） |
| `viewer/static/index.html` | Historyボタン、履歴パネル追加（約20行） |
| `viewer/static/style.css` | パネル・Diff表示スタイル追加（約80行） |
| `viewer/static/app.js` | 履歴機能のJS追加（約100行） |

---

## セキュリティ考慮

- `validate_path()` でファイルパス検証（既存関数）
- コミットハッシュは英数字のみ許可（正規表現チェック）
- subprocess実行時はリスト形式でコマンド渡し（インジェクション防止）

---

## 検証方法

```bash
./start-viewer.sh
# ブラウザで設定ファイルを選択
# 「History」ボタンをクリック
# コミット履歴が表示されることを確認
# 各コミットをクリックしてDiffが表示されることを確認
```
