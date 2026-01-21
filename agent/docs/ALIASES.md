# エイリアス・関数リファレンス

## 概要

`.zshrc` に定義されたエイリアスと関数の一覧。
すべて手動実行（ユーザーがコマンドを入力した時のみ動作）。

---

## 基本ツール

### vim

**Neovimをvimとして起動**

```bash
vim
```

**展開後:**
```bash
nvim
```

**用途:**
- Vimコマンドを入力するとNeovimが起動
- 既存のワークフローを変更せずにNeovimを使用

---

### python / python3

**Homebrew版Python 3.14を使用**

```bash
python script.py
python3 script.py
```

**展開後:**
```bash
/opt/homebrew/bin/python3.14 script.py
```

**用途:**
- システムPythonではなくHomebrew版を使用
- 最新のPython機能を利用

---

### cc

**Claude CLIを権限スキップモードで起動**

```bash
cc
```

**展開後:**
```bash
claude --dangerously-skip-permissions
```

**用途:**
- Claude Codeの対話モードを素早く起動
- 開発作業中のAI支援

**注意:**
- `--dangerously-skip-permissions`は開発環境でのみ使用推奨
- 本番環境では通常の`claude`コマンドを使用

---

## Git ワークフロー

### gsync

**リモートと同期して状態を確認する**

```bash
gsync
```

**展開後:**
```bash
git fetch --all && git pull && git status
```

**用途:**
- 作業開始時にリモートの最新状態を取得
- 他のメンバーの変更を取り込む

**運用例:**
```bash
# 朝一番、作業開始時
cd ~/projects/my-app
gsync

# 出力例:
# Fetching origin
# Already up to date.
# On branch main
# nothing to commit, working tree clean
```

---

### gac

**変更を全てステージングしてコミット**

```bash
gac
```

**展開後:**
```bash
git add . && git commit
```

**用途:**
- 全ての変更をまとめてコミット
- コミットメッセージはエディタで入力

**運用例:**
```bash
# 機能実装後
gac
# → エディタが開く → メッセージ入力 → 保存して閉じる

# メッセージを直接指定したい場合は通常のgitコマンドを使用
git add . && git commit -m "feat: add login feature"
```

**注意:**
- `.gitignore` に含まれないファイルは全てステージングされる
- 意図しないファイルが含まれていないか `git status` で事前確認推奨

---

## Docker ワークフロー

### drestart

**Docker Compose を再起動**

```bash
drestart
```

**展開後:**
```bash
docker-compose down && docker-compose up -d
```

**用途:**
- コンテナの状態をリセット
- 設定変更後の再起動
- 問題発生時のリフレッシュ

**運用例:**
```bash
# docker-compose.yml を編集後
cd ~/projects/my-app
vim docker-compose.yml
drestart

# 出力例:
# Stopping my-app_web_1 ... done
# Stopping my-app_db_1  ... done
# Creating my-app_db_1  ... done
# Creating my-app_web_1 ... done
```

---

### dclean

**未使用の Docker リソースを削除**

```bash
dclean
```

**展開後:**
```bash
docker system prune -f
```

**用途:**
- 停止中のコンテナを削除
- 未使用のネットワークを削除
- ダングリングイメージを削除
- ディスク容量の解放

**運用例:**
```bash
# ディスク容量が逼迫している時
dclean

# 出力例:
# Deleted Containers:
# 4a7f...
# Deleted Networks:
# my-app_default
# Total reclaimed space: 1.2GB

# より強力なクリーンアップが必要な場合
docker system prune -a -f  # 全ての未使用イメージも削除
```

**注意:**
- 停止中のコンテナは削除される
- 必要なコンテナは事前に起動しておくこと

---

### dlogs

**Docker Compose のログをリアルタイム表示**

```bash
dlogs
```

**展開後:**
```bash
docker-compose logs -f
```

**用途:**
- アプリケーションログの確認
- エラー調査
- デバッグ

**運用例:**
```bash
# サービス起動後にログを監視
drestart
dlogs

# 特定のサービスのみ見たい場合は通常のコマンドを使用
docker-compose logs -f web
docker-compose logs -f --tail=100 db
```

**終了:** `Ctrl + C`

---

## エージェント関連

### agent-run

**エージェントスクリプトを実行**

```bash
agent-run <agent-name> [args...]
```

**用途:**
- `~/dotfiles/agent/agents/general/` 内のスクリプトを実行
- 将来的に汎用エージェントを追加した際に使用

**運用例:**
```bash
# git-sync エージェントを実行（将来実装予定）
agent-run git-sync

# 引数付きで実行
agent-run deploy production
```

**現在の状態:**
スクリプトはまだ作成されていないため、使用するとエラーになる。

```bash
agent-run test
# Agent not found: test
```

---

## 運用シナリオ

### シナリオ1: 朝の作業開始

```bash
# 1. プロジェクトディレクトリに移動
cd ~/projects/my-app

# 2. リモートと同期
gsync

# 3. Docker環境を起動
drestart

# 4. ログを確認（別ターミナルで）
dlogs
```

### シナリオ2: 機能開発〜コミット

```bash
# 1. 機能を実装...
vim src/feature.ts

# 2. 変更を確認
git status
git diff

# 3. コミット
gac
# → エディタでメッセージ入力

# 4. プッシュ
git push
```

### シナリオ3: Docker 環境のトラブルシューティング

```bash
# 1. 現在の状態を確認
docker ps -a

# 2. ログを確認
dlogs

# 3. 再起動してみる
drestart

# 4. それでもダメなら完全リセット
docker-compose down -v  # ボリュームも削除
dclean
docker-compose up -d
```

### シナリオ4: 週末のクリーンアップ

```bash
# 1. 不要なリソースを削除
dclean

# 2. ディスク使用量を確認
docker system df

# 3. さらに容量を確保したい場合
docker system prune -a -f
docker volume prune -f
```

---

## カスタマイズ

### エイリアスの追加・変更

`.zshrc` の「3. エイリアス」セクションを編集:

```bash
vim ~/dotfiles/.zshrc
```

変更後は再読み込み:

```bash
source ~/.zshrc
```

### よく使うカスタムエイリアスの例

```bash
# ブランチ作成 + チェックアウト
alias gcb='git checkout -b'

# 直前のコミットを修正
alias gca='git commit --amend'

# Docker Compose のビルド + 起動
alias dbuild='docker-compose build && docker-compose up -d'

# 特定プロジェクトへの移動
alias proj='cd ~/projects/my-main-project'
```

---

## 更新履歴

| 日付 | 変更内容 |
|------|---------|
| 2024-12-29 | 初版作成 |
