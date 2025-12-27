" =============================================================================
" 文字コード自動認識（EUC-JP, Shift_JIS, ISO-2022-JP 対応）
" =============================================================================

if has('iconv')
  let s:enc_euc = 'euc-jp'
  let s:enc_jis = 'iso-2022-jp'
  if iconv("\x87\x64\x87\x6a", 'cp932', 'eucjp-ms') ==# "\xad\xc5\xad\xcb"
    let s:enc_euc = 'eucjp-ms'
    let s:enc_jis = 'iso-2022-jp-3'
  elseif iconv("\x87\x64\x87\x6a", 'cp932', 'euc-jisx0213') ==# "\xad\xc5\xad\xcb"
    let s:enc_euc = 'euc-jisx0213'
    let s:enc_jis = 'iso-2022-jp-3'
  endif
  let &fileencodings = s:enc_jis .','. s:enc_euc .',cp932,utf-8'
  unlet s:enc_euc
  unlet s:enc_jis
endif

" 日本語を含まない場合は fileencoding に encoding を使うようにする
augroup encoding_check
  autocmd!
  autocmd BufReadPost * if &fileencoding =~# 'iso-2022-jp' && search("[^\x01-\x7e]", 'n') == 0 | let &fileencoding=&encoding | endif
augroup END

" □とか◯の文字が有ってもカーソルの位置がズレないようにする
if exists('&ambiwidth')
  set ambiwidth=double
endif

" =============================================================================
" Vim専用設定（Neovimでは init.lua を使用）
" =============================================================================
if !has('nvim')
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
  Jetpack 'motemen/hatena-vim'
  call jetpack#end()

  " 基本設定
  set number
  set scrolloff=5
  set noswapfile
  set hidden
  set backspace=indent,eol,start
  set vb t_vb=
  set showcmd
  set laststatus=2
  set showmatch
  set wildmode=list:longest
  set list
  set listchars=tab:>-
  set cursorline
  set autoindent
  set shiftwidth=2
  set tabstop=2
  set expandtab
  set smartindent

  filetype plugin indent on
  syntax on

  " ;と:を入れ替え
  noremap ; :
  noremap : ;

  " NERDTree
  nnoremap ;t :NERDTreeToggle<CR>
  let g:NERDTreeMapOpenSplit = "-"

  " QuickBuf
  let g:qb_hotkey = ":;"
endif

" =============================================================================
" hatena.vim 設定
" =============================================================================
let g:hatena_user='a2c'
