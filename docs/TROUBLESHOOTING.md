# Troubleshooting Guide

dotfiles セットアップ時のよくある問題と解決方法

## 目次

- [シンボリックリンク関連](#シンボリックリンク関連)
- [Homebrew関連](#homebrew関連)
- [シェル設定関連](#シェル設定関連)
- [Neovim関連](#neovim関連)
- [SSH関連](#ssh関連)
- [権限関連](#権限関連)

---

## シンボリックリンク関連

### 既存ファイルがあってリンクが作成できない

**症状**: `ln: /path/to/file: File exists` エラー

**解決策**:
```bash
# バックアップを取ってから削除
mv ~/.zshrc ~/.zshrc.backup
./setup.sh --links-only
```

または、setup.sh が自動的にバックアップを作成するはずです。

### リンク先が間違っている

**確認方法**:
```bash
ls -la ~/.zshrc
# 期待: lrwxr-xr-x  ... .zshrc -> /path/to/dotfiles/.zshrc
```

**解決策**:
```bash
rm ~/.zshrc
./setup.sh --links-only
```

---

## Homebrew関連

### brew コマンドが見つからない

**症状**: `command not found: brew`

**解決策**:
```bash
# Apple Silicon Mac
eval "$(/opt/homebrew/bin/brew shellenv)"

# Intel Mac
eval "$(/usr/local/bin/brew shellenv)"
```

~/.zshrc に上記が含まれていることを確認してください。

### Brewfile からのインストールが失敗する

**症状**: `Error: No such file or directory @ rb_sysopen - Brewfile`

**解決策**:
```bash
# Brewfile の場所を確認
ls -la packages/Brewfile

# 正しいディレクトリから実行
cd ~/dotfiles
brew bundle --file=packages/Brewfile
```

### 特定のパッケージがインストールできない

**症状**: cask や formula のインストールエラー

**解決策**:
```bash
# キャッシュをクリア
brew cleanup
brew update

# 個別にインストールを試す
brew install <package-name>
```

---

## シェル設定関連

### Starship プロンプトが表示されない

**症状**: プロンプトが通常の `%` や `$` のまま

**解決策**:
1. Starship がインストールされているか確認:
   ```bash
   which starship
   ```

2. ~/.zshrc に初期化コードがあるか確認:
   ```bash
   grep starship ~/.zshrc
   # 期待: eval "$(starship init zsh)"
   ```

3. Nerd Font がインストールされているか確認:
   - ターミナルのフォント設定で Nerd Font を選択

### エイリアスが機能しない

**症状**: `vim` と打っても Neovim が起動しない等

**解決策**:
```bash
# シェルをリロード
source ~/.zshrc

# エイリアスを確認
alias | grep vim
```

### パス設定が反映されない

**症状**: インストールしたツールが見つからない

**解決策**:
```bash
# 現在のPATHを確認
echo $PATH | tr ':' '\n'

# シェルを完全に再起動
exec zsh
```

---

## Neovim関連

### プラグインがインストールされない

**症状**: 起動時にエラー、プラグインが動作しない

**解決策**:
```bash
# Lazy.nvim の同期
nvim --headless "+Lazy! sync" +qa

# または Neovim を起動して
:Lazy sync
```

### LSP が動作しない

**症状**: 補完やエラー表示が機能しない

**解決策**:
```bash
# mason.nvim 経由でLSPをインストール
nvim
:Mason
# 必要なLSPをインストール（lua_ls, pyright 等）
```

### render-markdown エラー

**症状**: `module 'render-markdown' not found`

**解決策**:
```bash
# プラグインの再インストール
nvim
:Lazy update render-markdown.nvim
```

---

## SSH関連

### 鍵が生成されない

**症状**: `ssh-keygen` がエラーを返す

**確認事項**:
```bash
# ~/.ssh ディレクトリの権限
ls -la ~ | grep .ssh
# 期待: drwx------

# 修正が必要な場合
chmod 700 ~/.ssh
```

### GitHub に接続できない

**症状**: `Permission denied (publickey)`

**解決策**:
1. SSH エージェントに鍵を追加:
   ```bash
   eval "$(ssh-agent -s)"
   ssh-add ~/.ssh/<your-key>
   ```

2. GitHub に公開鍵を登録:
   ```bash
   cat ~/.ssh/<your-key>.pub
   # 出力を GitHub Settings > SSH Keys に追加
   ```

3. 接続テスト:
   ```bash
   ssh -T git@github.com
   # 期待: Hi username! You've successfully authenticated...
   ```

### ~/.ssh/config の設定が間違っている

**確認方法**:
```bash
cat ~/.ssh/config
```

**正しい形式**:
```
Host github.com
  User git
  IdentityFile ~/.ssh/<your-key>
```

---

## 権限関連

### スクリプトが実行できない

**症状**: `Permission denied`

**解決策**:
```bash
chmod +x setup.sh
chmod +x scripts/*.sh
chmod +x scripts/setup/*.sh
```

### sudo が必要なコマンドで失敗する

**解決策**:
- Linux でパッケージインストールする場合は sudo 権限が必要です
- macOS では Homebrew が sudo なしで動作するように設計されています

---

## セットアップ後の確認項目

以下のコマンドで正常にセットアップされたか確認できます:

```bash
# 1. シンボリックリンクの確認
ls -la ~/.zshrc ~/.config/nvim ~/.config/starship.toml

# 2. エイリアスの動作確認
vim --version | head -1  # Neovim が起動するか

# 3. 環境変数の確認
echo $PNPM_HOME
echo $BUN_INSTALL

# 4. メトリクス収集の確認（有効な場合）
cat ~/.command_metrics.log | tail -5
```

---

## サポート

問題が解決しない場合:
1. GitHub Issues で報告: https://github.com/a2chub/profile/issues
2. スクリプトを `-x` オプション付きで実行してデバッグ:
   ```bash
   bash -x setup.sh
   ```
