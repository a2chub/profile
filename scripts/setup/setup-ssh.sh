#!/bin/bash

set -euo pipefail

# SSH設定の自動化
setup_ssh() {
    read -p "メールアドレスを入力してください: " email
    username=$(whoami)
    hostname=$(hostname)

    # ~/.ssh ディレクトリを作成（パーミッション700）
    mkdir -p ~/.ssh
    chmod 700 ~/.ssh

    # 鍵が存在しない場合のみ生成
    if [ ! -f "$HOME/.ssh/${username}_${hostname}" ]; then
        ssh-keygen -t ed25519 -C "$email" -f ~/.ssh/${username}_${hostname} -N ""
        echo "SSH鍵を生成しました: ~/.ssh/${username}_${hostname}"
    else
        echo "SSH鍵は既に存在します: ~/.ssh/${username}_${hostname}"
    fi

    ls -la ~/.ssh

    # ~/.ssh/config ファイルが存在しない場合は作成
    if [ ! -f ~/.ssh/config ]; then
        touch ~/.ssh/config
        chmod 600 ~/.ssh/config
    fi

    # ~/.ssh/config にgithub.com設定が存在しない場合のみ追加
    if ! grep -q "Host github.com" ~/.ssh/config 2>/dev/null; then
        # ファイルが空でない場合のみ空行を追加
        if [ -s ~/.ssh/config ]; then
            echo "" >> ~/.ssh/config
        fi
        cat >> ~/.ssh/config << EOF
Host github.com
  User git
  IdentityFile ~/.ssh/${username}_${hostname}
EOF
        echo "GitHub用のSSH設定を追加しました"
    else
        echo "GitHub用のSSH設定は既に存在します"
    fi
}

setup_ssh
