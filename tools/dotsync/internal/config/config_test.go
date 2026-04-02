package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_Default(t *testing.T) {
	// 設定ファイルが存在しない場合、デフォルト値が返される
	tmpDir := t.TempDir()

	cfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig がエラーを返しました: %v", err)
	}

	if cfg.Dotfiles.Repo != tmpDir {
		t.Errorf("Repo: got %q, want %q", cfg.Dotfiles.Repo, tmpDir)
	}
	if cfg.Dotfiles.ConfigDir != "config" {
		t.Errorf("ConfigDir: got %q, want %q", cfg.Dotfiles.ConfigDir, "config")
	}
	if cfg.Dotfiles.LinkScript != "scripts/link-dotfiles.sh" {
		t.Errorf("LinkScript: got %q, want %q", cfg.Dotfiles.LinkScript, "scripts/link-dotfiles.sh")
	}
	if cfg.Dotfiles.Brewfile != "packages/Brewfile" {
		t.Errorf("Brewfile: got %q, want %q", cfg.Dotfiles.Brewfile, "packages/Brewfile")
	}
	if len(cfg.Scan.WatchDirs) != 1 || cfg.Scan.WatchDirs[0] != "~/.config" {
		t.Errorf("WatchDirs: got %v, want [~/.config]", cfg.Scan.WatchDirs)
	}
	if len(cfg.Scan.Ignore) == 0 {
		t.Error("Ignore: デフォルトの ignore パターンが空です")
	}
	if cfg.Secrets.EntropyThreshold != 4.5 {
		t.Errorf("EntropyThreshold: got %f, want 4.5", cfg.Secrets.EntropyThreshold)
	}
	if len(cfg.Secrets.Patterns) == 0 {
		t.Error("Patterns: デフォルトのパターンが空です")
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	tmpDir := t.TempDir()

	// .dotsync.toml を作成
	configContent := `
[dotfiles]
config_dir = "my-config"
link_script = "my-scripts/link.sh"
brewfile = "my-packages/Brewfile"

[scan]
watch_dirs = ["~/.config", "~/.local/share"]
ignore = ["*.tmp"]

[secrets]
entropy_threshold = 5.0
`
	configPath := filepath.Join(tmpDir, ".dotsync.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("設定ファイル作成に失敗: %v", err)
	}

	cfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig がエラーを返しました: %v", err)
	}

	// ファイルから読み込まれた値を確認
	if cfg.Dotfiles.ConfigDir != "my-config" {
		t.Errorf("ConfigDir: got %q, want %q", cfg.Dotfiles.ConfigDir, "my-config")
	}
	if cfg.Dotfiles.LinkScript != "my-scripts/link.sh" {
		t.Errorf("LinkScript: got %q, want %q", cfg.Dotfiles.LinkScript, "my-scripts/link.sh")
	}
	if cfg.Dotfiles.Brewfile != "my-packages/Brewfile" {
		t.Errorf("Brewfile: got %q, want %q", cfg.Dotfiles.Brewfile, "my-packages/Brewfile")
	}
	if len(cfg.Scan.WatchDirs) != 2 {
		t.Errorf("WatchDirs: got %v, want 2 entries", cfg.Scan.WatchDirs)
	}
	if len(cfg.Scan.Ignore) != 1 || cfg.Scan.Ignore[0] != "*.tmp" {
		t.Errorf("Ignore: got %v, want [*.tmp]", cfg.Scan.Ignore)
	}
	if cfg.Secrets.EntropyThreshold != 5.0 {
		t.Errorf("EntropyThreshold: got %f, want 5.0", cfg.Secrets.EntropyThreshold)
	}

	// repo はファイルに設定されていないのでデフォルト（dotfilesDir）が使われる
	if cfg.Dotfiles.Repo != tmpDir {
		t.Errorf("Repo: got %q, want %q", cfg.Dotfiles.Repo, tmpDir)
	}
}

func TestLoadConfig_PartialOverride(t *testing.T) {
	tmpDir := t.TempDir()

	// ignore のみ設定
	configContent := `
[scan]
ignore = ["*.cache", "node_modules/"]
`
	configPath := filepath.Join(tmpDir, ".dotsync.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("設定ファイル作成に失敗: %v", err)
	}

	cfg, err := LoadConfig(tmpDir)
	if err != nil {
		t.Fatalf("LoadConfig がエラーを返しました: %v", err)
	}

	// ignore はファイルの値が使われる
	if len(cfg.Scan.Ignore) != 2 {
		t.Errorf("Ignore: got %v, want 2 entries", cfg.Scan.Ignore)
	}

	// 他のフィールドはデフォルト値が保持される
	if cfg.Dotfiles.ConfigDir != "config" {
		t.Errorf("ConfigDir: got %q, want %q (default)", cfg.Dotfiles.ConfigDir, "config")
	}
	if cfg.Dotfiles.LinkScript != "scripts/link-dotfiles.sh" {
		t.Errorf("LinkScript: got %q, want %q (default)", cfg.Dotfiles.LinkScript, "scripts/link-dotfiles.sh")
	}
	if cfg.Secrets.EntropyThreshold != 4.5 {
		t.Errorf("EntropyThreshold: got %f, want 4.5 (default)", cfg.Secrets.EntropyThreshold)
	}
	if len(cfg.Scan.WatchDirs) != 1 || cfg.Scan.WatchDirs[0] != "~/.config" {
		t.Errorf("WatchDirs: got %v, want [~/.config] (default)", cfg.Scan.WatchDirs)
	}
}
