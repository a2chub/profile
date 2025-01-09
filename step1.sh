#!/bin/bash

# OSの判別
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="Linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="MacOS"
else
    echo "Unsupported OS"
    exit 1
fi

# MacOSの場合、brewの存在を確認し、なければインストール
check_and_install_brew() {
    if ! command -v brew &> /dev/null; then
        echo "Homebrewが見つかりません。インストールを行います。"
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
    else
        echo "Homebrewは既にインストールされています。"
    fi
}

# 必要なアプリのインストール
install_apps() {
    if [ "$OS" == "Linux" ]; then
        sudo apt update
        sudo apt install -y vim tmux
        # miseのインストールまたはアップデート
        if ! command -v mise &> /dev/null; then
            echo "miseが見つかりません。インストールを行います。"
            curl https://mise.run | sh
        else
            echo "miseは既にインストールされています。アップデートを行います。"
            mise self-update
        fi
    elif [ "$OS" == "MacOS" ]; then
        check_and_install_brew
        brew install vim tmux
        # miseのインストールまたはアップデート
        if ! command -v mise &> /dev/null; then
            echo "miseが見つかりません。インストールを行います。"
            curl https://mise.run | sh
        else
            echo "miseは既にインストールされています。アップデートを行います。"
            mise self-update
        fi
    fi
}

# vim-jetpackのインストール
install_vim_jetpack() {
    if [ ! -d "~/.vim/pack/jetpack/opt/vim-jetpack" ]; then
        echo "vim-jetpackをインストールします。"
        git clone https://github.com/tani/vim-jetpack.git ~/.vim/pack/jetpack/opt/vim-jetpack
    else
        echo "vim-jetpackは既にインストールされています。"
    fi
}

# gitのインストール
install_git() {
    if ! command -v git &> /dev/null; then
        echo "gitが見つかりません。インストールを行います。"
        if [ "$OS" == "Linux" ]; then
            sudo apt install -y git
        elif [ "$OS" == "MacOS" ]; then
            brew install git
        fi
    else
        echo "gitは既にインストールされています。"
    fi
}

# 設定ファイルのクローン
setup_configs() {
    install_git
    install_vim_jetpack
    git clone https://github.com/yourusername/dotfiles.git ~/dotfiles
    ln -s ~/dotfiles/.vimrc ~/.vimrc
    ln -s ~/dotfiles/.tmux.conf ~/.tmux.conf
    # 他の設定ファイルも同様にリンクを作成
}

# メイン処理
install_apps
setup_configs

# SSH設定はオプション
read -p "SSH設定を行いますか？ (y/n): " setup_ssh_choice
if [ "$setup_ssh_choice" == "y" ]; then
    ./step2.sh
fi

echo "セットアップが完了しました。"