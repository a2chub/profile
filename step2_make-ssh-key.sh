#!/bin/bash

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
    else
        echo "SSH鍵は既に存在します: ~/.ssh/${username}_${hostname}"
    fi

    ls -la ~/.ssh

    # ~/.ssh/config にgithub.com設定が存在しない場合のみ追加
    if ! grep -q "Host github.com" ~/.ssh/config 2>/dev/null; then
        echo "" >> ~/.ssh/config
        echo "Host github.com" >> ~/.ssh/config
        echo "  User git" >> ~/.ssh/config
        echo "  IdentityFile ~/.ssh/${username}_${hostname}" >> ~/.ssh/config
        chmod 600 ~/.ssh/config
    else
        echo "GitHub用のSSH設定は既に存在します"
    fi
}

setup_ssh
