-- 行番号表示
vim.opt.number = true
--vim.opt.relativenumber = true

-- -- マウス操作有効
vim.opt.mouse = 'a'

-- -- 検索設定
vim.opt.ignorecase = true
vim.opt.smartcase = true
vim.opt.hlsearch = true

-- -- テキスト折り返し
vim.opt.wrap = true
vim.opt.breakindent = true

-- -- インデント設定
vim.opt.tabstop = 2
vim.opt.shiftwidth = 2
vim.opt.expandtab = true
vim.opt.smartindent = true
vim.opt.cursorline = true

-- -- 24bit色有効
 vim.opt.termguicolors = true

-- -- システムクリップボードとの連携
 vim.opt.clipboard = 'unnamedplus'

-- -- バックアップ・スワップファイル無効化
vim.opt.backup = false
vim.opt.writebackup = false
vim.opt.swapfile = false
--
-- -- undoファイルを有効化
vim.opt.undofile = true


-- キーバインド

-- セミコロンをコロンとして動作させる
vim.keymap.set('n', ';', ':')

-- コロンをセミコロンとして動作させる
vim.keymap.set('n', ':', ';')

-- NERDTree のキーマッピング
vim.keymap.set('n', ';t', ':NERDTreeToggle<CR>', { noremap = true, silent = true })

-- NERDTree の設定
vim.g.NERDTreeMapOpenSplit = "-"

-- nvim + lua
local fn = vim.fn
local jetpackfile = fn.stdpath('data') .. '/site/pack/jetpack/opt/vim-jetpack/plugin/jetpack.vim'
local jetpackurl = 'https://raw.githubusercontent.com/tani/vim-jetpack/master/plugin/jetpack.vim'
if fn.filereadable(jetpackfile) == 0 then
  fn.system('curl -fsSLo ' .. jetpackfile .. ' --create-dirs ' .. jetpackurl)
end


-- Jetpack 関連の設定
vim.cmd('packadd vim-jetpack')
require('jetpack.paq'){
  {'preservim/nerdtree'},
  {'nvim-telescope/telescope.nvim'},
  {'nvim-lualine/lualine.nvim'},
  {'nvim-treesitter/nvim-treesitter'},
  {'gruvbox-community/gruvbox'},
  {'vim-airline/vim-airline'},
  {'vim-airline/vim-airline-themes'},
  {'mattn/gist-vim'},
  {'mattn/webapi-vim'},
}

-- 例: Telescopeのキーバインド設定
vim.api.nvim_set_keymap('n', '<leader>ff', ':Telescope find_files<CR>', { noremap = true, silent = true })

-- Pluginの自動インストール
local jetpack = require('jetpack')
for _, name in ipairs(jetpack.names()) do
  if not jetpack.tap(name) then
    jetpack.sync()
    break
  end
end
