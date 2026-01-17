#!/bin/bash

set -euo pipefail

# エラートラップ
error_handler() {
    local line_no=$1
    local error_code=$2
    echo "エラー: 行 $line_no でコマンドが失敗しました (終了コード: $error_code)" >&2
    exit "$error_code"
}
trap 'error_handler ${LINENO} $?' ERR

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
    if [ ! -d "$HOME/.vim/pack/jetpack/opt/vim-jetpack" ]; then
        echo "vim-jetpackをインストールします。"
        curl -fLo ~/.vim/pack/jetpack/opt/vim-jetpack/plugin/jetpack.vim --create-dirs https://raw.githubusercontent.com/tani/vim-jetpack/master/plugin/jetpack.vim
    else
        echo "vim-jetpackは既にインストールされています。"
    fi
    # Powerline fontsのインストール
    if [ ! -d "$HOME/.local/share/fonts" ] || [ -z "$(ls -A $HOME/.local/share/fonts 2>/dev/null | grep -i powerline)" ]; then
        echo "Powerline fontsをインストールします。"
        TEMP_FONTS=$(mktemp -d)
        git clone https://github.com/powerline/fonts.git --depth=1 "$TEMP_FONTS"
        "$TEMP_FONTS/install.sh"
        rm -rf "$TEMP_FONTS"
    else
        echo "Powerline fontsは既にインストールされています。"
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

    # dotfilesが既に存在する場合はスキップ
    if [ -d ~/dotfiles ]; then
        echo "dotfilesは既に存在します。シンボリックリンクの作成には ./setup.sh --links-only を使用してください。"
    else
        echo "git clone https://github.com/a2chub/profile.git ~/dotfiles"
        git clone https://github.com/a2chub/profile.git ~/dotfiles
    fi

    # シンボリックリンクの作成（setup.shを使用することを推奨）
    if [ -f ~/dotfiles/setup.sh ]; then
        echo "シンボリックリンクの作成には ~/dotfiles/setup.sh --links-only を実行してください。"
    fi
}

# メイン処理
install_apps
setup_configs

# SSH設定はオプション
#read -p "SSH設定を行いますか？ (y/n): " setup_ssh_choice
#if [ "$setup_ssh_choice" == "y" ]; then
#    ./step2.sh
#fi

echo "セットアップが完了しました。"



