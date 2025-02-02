" 文字コードの自動認識
if &encoding !=# 'utf-8'
  set encoding=japan
  set fileencoding=japan
endif

if has('iconv')
  let s:enc_euc = 'euc-jp'
  let s:enc_jis = 'iso-2022-jp'
  " iconvがeucJP-msに対応しているかをチェック
  if iconv("\x87\x64\x87\x6a", 'cp932', 'eucjp-ms') ==# "\xad\xc5\xad\xcb"
    let s:enc_euc = 'eucjp-ms'
    let s:enc_jis = 'iso-2022-jp-3'
  " iconvがJISX0213に対応しているかをチェック
  elseif iconv("\x87\x64\x87\x6a", 'cp932', 'euc-jisx0213') ==# "\xad\xc5\xad\xcb"
    let s:enc_euc = 'euc-jisx0213'
    let s:enc_jis = 'iso-2022-jp-3'
  endif
  " fileencodingsを構築
  if &encoding ==# 'utf-8'
    let s:fileencodings_default = &fileencodings
    let &fileencodings = s:enc_jis .','. s:enc_euc .',cp932'
    let &fileencodings = &fileencodings .','. s:fileencodings_default
    unlet s:fileencodings_default
  else
    let &fileencodings = &fileencodings .','. s:enc_jis
    set fileencodings+=utf-8,ucs-2le,ucs-2
    if &encoding =~# '^\(euc-jp\|euc-jisx0213\|eucjp-ms\)$'
      set fileencodings+=cp932
      set fileencodings-=euc-jp
      set fileencodings-=euc-jisx0213
      set fileencodings-=eucjp-ms
      let &encoding = s:enc_euc
      let &fileencoding = s:enc_euc
    else
      let &fileencodings = &fileencodings .','. s:enc_euc
    endif
  endif
  " 定数を処分
  unlet s:enc_euc
  unlet s:enc_jis
endif

" 日本語を含まない場合は fileencoding に encoding を使うようにする
if has('autocmd')
  function! AU_ReCheck_FENC()
    if &fileencoding =~# 'iso-2022-jp' && search("[^\x01-\x7e]", 'n') == 0
      let &fileencoding=&encoding
    endif
  endfunction
  autocmd BufReadPost * call AU_ReCheck_FENC()
endif

" 改行コードの自動認識
set fileformats=unix,dos,mac
" □とか◯の文字が有ってもカーソルの位置がズレないようにする
if exists('&ambiwidth')
  set ambiwidth=double
endif

" Jetpack Plugin インストール
packadd vim-jetpack
call jetpack#begin()
Jetpack 'tani/vim-jetpack'
Jetpack 'preservim/nerdtree'
Jetpack 'preservim/nerdcommenter'
Jetpack 'vim-airline/vim-airline'
Jetpack 'vim-airline/vim-airline-themes'
Jetpack 'vim-jp/vimdoc-ja'
Jetpack 'mattn/gist-vim'
Jetpack 'Shougo/qbuf.vim'
call jetpack#end()

" ショートカット
" Ev/Rvでvimrcの編集と反映
command! Ev edit $MYVIMRC
command! Rv source $MYVIMRC

" ファイルタイプの自動認識
filetype plugin indent on


"-------------------------------------------------------------------------------
" 基本設定 Basics
"-------------------------------------------------------------------------------
"let mapleader = "\"             " キーマップリーダー
set scrolloff=5                  " スクロール時の余白確保
""set autoread                     " 他で書き換えられたら自動で読み直す
set noswapfile                   " スワップファイル作らない
set hidden                       " 編集中でも他のファイルを開けるようにする
set backspace=indent,eol,start   " バックスペースでなんでも消せるように
set formatoptions=lmoq           " テキスト整形オプション，マルチバイト系を追加
set vb t_vb=                     " ビープをならさない
set browsedir=buffer             " Exploreの初期ディレクトリ
"set whichwrap=b,s,h,l,<,>,[,]    " カーソルを行頭、行末で止まらないようにする
set showcmd                      " コマンドをステータス行に表示
"set showmode                     " 現在のモードを表示
set viminfo='50,<1000,s100,\"50  " viminfoファイルの設定
set modelines=0                  " モードラインは無効
set number                       " 行番号を表示する

highlight WhitespaceEOL ctermbg=gray guibg=red
match WhitespaceEOL /¥s¥+$/
match WhitespaceEOL /\s\+$/
autocmd WinEnter * match WhitespaceEOL /¥s¥+$/

filetype on
filetype plugin on
syntax on

set enc=utf-8
set fenc=utf-8
set fencs=iso-2022-jp,euc-jp,cp932
"set statusline=%<[%n]%m%r%h%w%{'['.(&fenc!=''?&fenc:&enc).':'.&ff.']['.&ft.']'}\ %F%=%l,%c%V%8P
"set statusline=%<%m%w%%c%V%8P
set laststatus=2
set fileformats=unix,dos,mac
set showmatch
set wildmode=list:longest
set list
set listchars:tab:>-
set cursorline
set autoindent
set shiftwidth=2
set tabstop=2
set smartindent
set history=200     " keep 50 lines of command line history
set ruler           " show the cursor position all the time
set backupdir=$HOME/.vim_backup
let &directory = &backupdir


"autocmd FileType python set tabstop=8 shiftwidth=2 expandtab fenc=utf-8
autocmd FileType perl set isfname-=-
autocmd FileType yaml set expandtab ts=2 sw=2 enc=utf-8 fenc=utf-8

au  BufEnter *.py setlocal indentkeys+=0#

" CTRL+p でファイルを実行
nnoremap <silent> <C-p> :<C-u>execute '!' &l:filetype '%'<Return>

" カレントウィンドウにのみ罫線を引く
augroup cch
  autocmd! cch
  autocmd WinLeave * set nocursorline
  autocmd WinEnter,BufRead * set cursorline
augroup END


" For Python
if has('autocmd')
  autocmd FileType python set tabstop=8 shiftwidth=2 expandtab fenc=utf-8 smarttab
  "autocmd FileType python set omnifunc=pythoncomplete#Complete
  "" http://vim.sourceforge.net/scripts/script.php?script_id=30
  autocmd FileType python source $HOME/.vim/plugin/python.vim

  autocmd BufNewFile  *.py 0r ~/.vim/template/python.txt
endif

" ;でコマンド入力( ;と:を入れ替)
noremap ; :
noremap : ;


"===================================================================
" Plugin 関係
"===================================================================

" NERD_tree.vim
nnoremap ;t :NERDTreeToggle<CR>
let g:NERDTreeMapOpenSplit = "-"

" NERD commenter
map <Leader>x ,c<space>

" zen coding
let g:user_zen_expandabbr_key = '<c-e>'
let g:user_zen_settings = { 'indentation': ' ' }

"minibufexpl.vim
"set minibfexp
let g:miniBufExplMapWindowNavVim=1
let g:miniBufExplSplitBelow=0
let g:miniBufExplMapWindowNavArrows=1
let g:miniBufExplMapCTabSwitchBufs=1
let g:miniBufExplModSelTarget=1
let g:miniBufExplSplitToEdge=1

" QuickBuf: qbuf.vim
let g:qb_hotkey = ":;"

" changelogの記入設定
" 新規エントリーは <Leader> + o
let g:changelog_username = "a2c <atusi@a2c.biz>"
let g:changelog_timeformat = "# %Y-%m-%d (%a)"
nnoremap <Leader><Leader><Leader> :new ~atusi/Dropbox/changelog<cr>

"set runtimepath+=$HOME/.vim/plugin/hatena
let g:hatena_user='a2c'

"===================================================================
" 編集 関係
"===================================================================
" Tabキーを空白に変換
set expandtab

" コンマの後に自動的にスペースを挿入
inoremap , ,<Space>

" Tohtml setting
let html_number_lines = 0
let html_use_css = 1
let use_xhtml = 1

" 外部で変更されたら再度読み込む（Window切替時)
"augroup vimrc-checktime
  "autocmd!
  "autocmd WinEnter * checktime
"augroup END

" 保存時に行末の空白を除去する
autocmd BufWritePre * :%s/\s\+$//ge
" " 保存時にtabをスペースに変換する
autocmd BufWritePre * :%s/\t/  /ge

" ヤンク時にクリップボードにもコピーする
set clipboard=unnamed
