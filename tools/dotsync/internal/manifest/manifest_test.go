package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/a2chub/dotsync/internal/config"
)

func TestRunFullScan(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tmpDir := t.TempDir()

	dotfilesDir := filepath.Join(tmpDir, "dotfiles")
	homeDir := filepath.Join(tmpDir, "home")
	configDir := filepath.Join(homeDir, ".config")

	// ディレクトリ構造を作成
	for _, dir := range []string{
		filepath.Join(dotfilesDir, "config", "nvim"),
		filepath.Join(dotfilesDir, "scripts"),
		filepath.Join(dotfilesDir, "packages"),
		configDir,
	} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("ディレクトリ作成に失敗: %v", err)
		}
	}

	// dotfiles/config/nvim/init.lua を作成
	nvimInit := filepath.Join(dotfilesDir, "config", "nvim", "init.lua")
	if err := os.WriteFile(nvimInit, []byte("-- nvim config"), 0o644); err != nil {
		t.Fatalf("ファイル作成に失敗: %v", err)
	}

	// link-dotfiles.sh を作成（create_link 1行）
	// テスト用に絶対パスで記述（$HOME は実際のホームディレクトリに展開されるため）
	scriptPath := filepath.Join(dotfilesDir, "scripts", "link-dotfiles.sh")
	scriptContent := fmt.Sprintf(`#!/bin/bash
create_link "$DOTFILES_DIR/config/nvim" "%s"
`, filepath.Join(configDir, "nvim"))
	if err := os.WriteFile(scriptPath, []byte(scriptContent), 0o755); err != nil {
		t.Fatalf("スクリプト作成に失敗: %v", err)
	}

	// ~/.config/nvim をシンボリックリンクとして作成（管理対象）
	nvimLink := filepath.Join(configDir, "nvim")
	nvimSource := filepath.Join(dotfilesDir, "config", "nvim")
	if err := os.Symlink(nvimSource, nvimLink); err != nil {
		t.Fatalf("シンボリックリンク作成に失敗: %v", err)
	}

	// ~/.config/unmanaged-app を作成（非管理対象ディレクトリ）
	unmanagedDir := filepath.Join(configDir, "unmanaged-app")
	if err := os.MkdirAll(unmanagedDir, 0o755); err != nil {
		t.Fatalf("ディレクトリ作成に失敗: %v", err)
	}
	// 非管理対象ディレクトリ内にファイルを追加
	if err := os.WriteFile(filepath.Join(unmanagedDir, "config.toml"), []byte("key = \"value\""), 0o644); err != nil {
		t.Fatalf("ファイル作成に失敗: %v", err)
	}

	// Brewfile を作成
	brewfilePath := filepath.Join(dotfilesDir, "packages", "Brewfile")
	brewfileContent := `brew "git"
brew "jq"
`
	if err := os.WriteFile(brewfilePath, []byte(brewfileContent), 0o644); err != nil {
		t.Fatalf("Brewfile作成に失敗: %v", err)
	}

	// Config を作成（テスト用に watchDirs を直接指定）
	cfg := &config.Config{
		Dotfiles: config.DotfilesConfig{
			Repo:       dotfilesDir,
			ConfigDir:  "config",
			LinkScript: "scripts/link-dotfiles.sh",
			Brewfile:   "packages/Brewfile",
		},
		Scan: config.ScanConfig{
			WatchDirs: []string{configDir},
			Ignore:    []string{},
		},
		Secrets: config.SecretsConfig{
			EntropyThreshold: 4.5,
		},
	}

	// フルスキャンを実行
	result, err := RunFullScan(cfg)
	if err != nil {
		t.Fatalf("RunFullScan がエラーを返しました: %v", err)
	}

	// 基本的な検証
	if result.DotfilesDir != dotfilesDir {
		t.Errorf("DotfilesDir: got %q, want %q", result.DotfilesDir, dotfilesDir)
	}
	if result.HomeDir == "" {
		t.Error("HomeDir が空です")
	}
	if result.ScannedAt.IsZero() {
		t.Error("ScannedAt がゼロ値です")
	}

	// コンフィグアイテムを検証（managed も含まれるようになった）
	var configItems []ScanItem
	var brewItems []ScanItem
	for _, item := range result.Items {
		switch item.Kind {
		case KindConfig:
			configItems = append(configItems, item)
		case KindBrew:
			brewItems = append(brewItems, item)
		}
	}

	// managed (nvim) と unmanaged-app の2つが検出されるはず
	if len(configItems) < 1 {
		t.Fatalf("コンフィグアイテム数: got %d, want >= 1", len(configItems))
	}

	// unmanaged-app を探す
	var unmanagedItem *ScanItem
	for i, item := range configItems {
		if item.Name == "unmanaged-app" {
			unmanagedItem = &configItems[i]
			break
		}
	}
	if unmanagedItem == nil {
		t.Fatal("unmanaged-app が検出されませんでした")
	}

	if unmanagedItem.Status != "unmanaged" {
		t.Errorf("Status: got %q, want %q", unmanagedItem.Status, "unmanaged")
	}
	if unmanagedItem.Kind != KindConfig {
		t.Errorf("Kind: got %d, want %d (KindConfig)", unmanagedItem.Kind, KindConfig)
	}

	// ファイル情報の検証
	if len(unmanagedItem.Files) != 1 {
		t.Fatalf("Files 数: got %d, want 1", len(unmanagedItem.Files))
	}
	if unmanagedItem.Files[0].RelPath != "config.toml" {
		t.Errorf("Files[0].RelPath: got %q, want %q", unmanagedItem.Files[0].RelPath, "config.toml")
	}

	// ID が連番であることを検証
	for i, item := range result.Items {
		expectedID := i + 1
		if item.ID != expectedID {
			t.Errorf("Items[%d].ID: got %d, want %d", i, item.ID, expectedID)
		}
	}

	// brew 関連: brew が利用できない環境ではエラーとして記録され、brew アイテムは空になる
	// brew が利用可能な場合は missing-from-brewfile アイテムが含まれる可能性がある
	for _, item := range brewItems {
		if item.Status != "missing-from-brewfile" {
			t.Errorf("Brew item status: got %q, want %q", item.Status, "missing-from-brewfile")
		}
		if item.Kind != KindBrew {
			t.Errorf("Brew item kind: got %d, want %d (KindBrew)", item.Kind, KindBrew)
		}
	}

	// brew が利用できない場合、エラーに記録されていることを確認
	hasBrew := true
	for _, scanErr := range result.Errors {
		if scanErr.Context == "brew leaves の実行" {
			hasBrew = false
			break
		}
	}
	if !hasBrew {
		t.Log("brew が利用できない環境のため、brew leaves のエラーを確認")
		if len(brewItems) != 0 {
			t.Errorf("brew が利用できない場合、brew アイテムは空であるべき: got %d", len(brewItems))
		}
	}
}
