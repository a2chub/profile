-- .vimrcからエンコーディング設定を読み込む（EUC-JP, Shift_JIS対応）
local vimrc = vim.fn.expand('~/.vimrc')
if vim.fn.filereadable(vimrc) == 1 then
  vim.cmd('source ' .. vimrc)
end

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

-- NeoTree の設定
vim.keymap.set('n', '<C-n>', ':Neotree toggle<CR>')

-- ペイン移動 (Ctrl + h/j/k/l)
vim.keymap.set('n', '<C-h>', '<C-w>h', { noremap = true, silent = true, desc = 'Move to left pane' })
vim.keymap.set('n', '<C-j>', '<C-w>j', { noremap = true, silent = true, desc = 'Move to lower pane' })
vim.keymap.set('n', '<C-k>', '<C-w>k', { noremap = true, silent = true, desc = 'Move to upper pane' })
vim.keymap.set('n', '<C-l>', '<C-w>l', { noremap = true, silent = true, desc = 'Move to right pane' })

-- ターミナル設定
vim.keymap.set('n', '<Leader>t', ':belowright split | terminal<CR>', { noremap = true, silent = true, desc = 'Open terminal at bottom' })
vim.keymap.set('t', '<Esc>', '<C-\\><C-n>', { noremap = true, silent = true, desc = 'Exit terminal mode' })

-- 保存時の自動整形（特定ファイルタイプのみ）
local auto_format_filetypes = {
  'javascript', 'typescript', 'javascriptreact', 'typescriptreact',
  'html', 'css', 'json', 'yaml',
  'python', 'ruby', 'lua',
}
vim.api.nvim_create_augroup('AutoFormat', { clear = true })
vim.api.nvim_create_autocmd('BufWritePre', {
  group = 'AutoFormat',
  pattern = '*',
  callback = function()
    if vim.tbl_contains(auto_format_filetypes, vim.bo.filetype) then
      -- 行末空白を削除
      local save_cursor = vim.fn.getpos('.')
      vim.cmd([[%s/\s\+$//ge]])
      vim.fn.setpos('.', save_cursor)
    end
  end,
})


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
  {'tani/vim-jetpack', opt = true},
  {'preservim/nerdtree'},
  {'nvim-telescope/telescope.nvim'},
  {'nvim-treesitter/nvim-treesitter'},
  {'gruvbox-community/gruvbox'},
  {'vim-airline/vim-airline'},
  {'vim-airline/vim-airline-themes'},
  {'mattn/gist-vim'},
  {'mattn/webapi-vim'},
  {'motemen/hatena-vim'},
  {"nvim-neo-tree/neo-tree.nvim"},
  {"nvim-lua/plenary.nvim"},
  {"nvim-tree/nvim-web-devicons"},
  {"MunifTanjim/nui.nvim"},
}

-- 例: Telescopeのキーバインド設定
vim.api.nvim_set_keymap('n', '<leader>ff', ':Telescope find_files<CR>', { noremap = true, silent = true })

-- Pluginの自動インストール (jetpack.paq経由で管理)


-- neo-tree のセットアップ
require("neo-tree").setup({
  -- 最後のウィンドウを閉じたときに neo-tree も閉じる
  close_if_last_window = true,
  -- ポップアップウィンドウの枠線スタイル
  popup_border_style = "rounded",
  -- Git のステータス表示を有効化
  enable_git_status = true,
  -- 診断情報 (エラー、警告など) の表示を有効化
  enable_diagnostics = true,
  -- 各コンポーネントのデフォルト設定
  default_component_configs = {
    icon = {
      folder_closed = "",
      folder_open = "",
      folder_empty = "",
    },
    git_status = {
      symbols = {
        added     = "✚",
        modified  = "",
        deleted   = "✖",
        renamed   = "",
        untracked = "",
        ignored   = "",
        unstaged  = "",
        staged    = "✓",
        conflict  = "",
      },
    },
  },
  -- neo-tree ウィンドウの設定
  window = {
    position = "left",  -- 表示位置（left/right/float）
    width = 40,         -- ウィンドウ幅
    mapping_options = {
      noremap = true,
      nowait = true,
    },
    mappings = {
      ["<space>"] = "toggle_node",
      ["<2-LeftMouse>"] = "open",  -- マウスダブルクリックで開く
      ["l"]  = "open",
      ["<cr>"] = "open",
      ["o"]  = "open",
      ["O"]  = "open_split",
      ["i"]  = "open_vsplit",
      ["t"]  = "open_tabnew",
      ["<esc>"] = "revert_preview",
      ["P"]  = { "toggle_preview", config = { use_float = true } },
      ["C"]  = "close_node",
      ["a"]  = { "add", config = { show_path = "relative" } },
      ["A"]  = "add_directory",
      ["d"]  = "delete",
      ["r"]  = "rename",
      ["y"]  = "copy_to_clipboard",
      ["x"]  = "cut_to_clipboard",
      ["p"]  = "paste_from_clipboard",
      ["q"]  = "close_window",
      ["R"]  = "refresh",
    },
  },
  -- ファイルシステム関連の設定
  filesystem = {
    filtered_items = {
      visible = false,      -- 非表示にする項目はリストに表示しない
      hide_dotfiles = true, -- ドットファイルを隠す
      hide_gitignored = true, -- Git に無視されているファイルも非表示にする
    },
  },
  -- バッファ関連の設定（オプション）
  buffers = {
    follow_current_file = true,  -- 現在編集中のファイルに追従する
  },
})




