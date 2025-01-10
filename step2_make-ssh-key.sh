#!/bin/bash

# SSH設定の自動化
setup_ssh() {
    read -p "メールアドレスを入力してください: " email
    username=$(whoami)
    hostname=$(hostname)
    ssh-keygen -t rsa -b 4096 -C "$email" -f ~/.ssh/${username}_${hostname} -N ""

    ls -la ~/.ssh

    #ssh-copy-id -i "~/.ssh/${username}_${hostname}.pub" user@your-ssh-server
    echo "Host your-ssh-server" >> ~/.ssh/config
    echo "  User yourusername" >> ~/.ssh/config
    echo "  IdentityFile ~/.ssh/${username}_${hostname}" >> ~/.ssh/config
}

setup_ssh 
