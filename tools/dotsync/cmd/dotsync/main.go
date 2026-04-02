package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"

	"github.com/a2chub/dotsync/internal/config"
	"github.com/a2chub/dotsync/internal/tui"
)

func main() {
	dotfilesFlag := flag.String("dotfiles", "~/dotfiles", "dotfiles リポジトリのパス")
	dryRun := flag.Bool("dry-run", false, "変更を行わずに表示のみ")
	flag.Parse()

	// ~ を展開
	dotfilesDir := expandHome(*dotfilesFlag)

	// dotfiles ディレクトリの存在確認
	if _, err := os.Stat(dotfilesDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: dotfiles directory not found: %s\n", dotfilesDir)
		os.Exit(1)
	}

	// 設定ファイルを読み込み
	cfg, err := config.LoadConfig(dotfilesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load config: %v\n", err)
		os.Exit(1)
	}

	// TUI モデル作成・実行
	m := tui.NewModel(cfg, *dryRun)
	p := tea.NewProgram(m)

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// 終了後にコミットメッセージを表示
	if fm, ok := finalModel.(tui.Model); ok {
		if msg := fm.CommitMsg(); msg != "" {
			fmt.Println(msg)
		}
	}
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
