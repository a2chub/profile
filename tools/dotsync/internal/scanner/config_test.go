package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseLinkScript(t *testing.T) {
	tmpDir := t.TempDir()

	script := `#!/bin/bash
# コメント行はスキップされる
DOTFILES_DIR="$1"

create_link "$DOTFILES_DIR/.zshrc" "$HOME/.zshrc"
create_link "$DOTFILES_DIR/config/nvim" "$HOME/.config/nvim"
create_link "$DOTFILES_DIR/config/starship.toml" "$HOME/.config/starship.toml"
`
	scriptPath := filepath.Join(tmpDir, "link-dotfiles.sh")
	if err := os.WriteFile(scriptPath, []byte(script), 0644); err != nil {
		t.Fatal(err)
	}

	dotfilesDir := "/home/user/dotfiles"
	homeDir := "/home/user"

	links, errs := ParseLinkScript(scriptPath, dotfilesDir, homeDir)

	if len(errs) > 0 {
		t.Errorf("予期しないエラー: %v", errs)
	}

	if len(links) != 3 {
		t.Fatalf("リンク数が期待値と異なる: got %d, want 3", len(links))
	}

	tests := []struct {
		wantSrc  string
		wantDest string
	}{
		{"/home/user/dotfiles/.zshrc", "/home/user/.zshrc"},
		{"/home/user/dotfiles/config/nvim", "/home/user/.config/nvim"},
		{"/home/user/dotfiles/config/starship.toml", "/home/user/.config/starship.toml"},
	}

	for i, tt := range tests {
		if links[i].Source != tt.wantSrc {
			t.Errorf("links[%d].Source = %q, want %q", i, links[i].Source, tt.wantSrc)
		}
		if links[i].Dest != tt.wantDest {
			t.Errorf("links[%d].Dest = %q, want %q", i, links[i].Dest, tt.wantDest)
		}
		if links[i].Line == 0 {
			t.Errorf("links[%d].Line が 0 になっている", i)
		}
	}
}

func TestScanConfigDir_Managed(t *testing.T) {
	tmpDir := t.TempDir()
	dotfilesDir := filepath.Join(tmpDir, "dotfiles")
	configDir := filepath.Join(tmpDir, "config")

	// dotfiles 側のディレクトリを作成
	nvimSrc := filepath.Join(dotfilesDir, "config", "nvim")
	if err := os.MkdirAll(nvimSrc, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	// シンボリックリンクを作成
	nvimDest := filepath.Join(configDir, "nvim")
	if err := os.Symlink(nvimSrc, nvimDest); err != nil {
		t.Fatal(err)
	}

	managed := []ManagedLink{
		{Source: nvimSrc, Dest: nvimDest, Line: 10},
	}

	entries, err := ScanConfigDir(configDir, managed, dotfilesDir, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf("エントリ数が期待値と異なる: got %d, want 1", len(entries))
	}

	if entries[0].Status != StatusManaged {
		t.Errorf("Status = %d, want StatusManaged (%d)", entries[0].Status, StatusManaged)
	}
	if entries[0].Name != "nvim" {
		t.Errorf("Name = %q, want %q", entries[0].Name, "nvim")
	}
	if entries[0].ExpectedLink == nil {
		t.Error("ExpectedLink が nil")
	}
}

func TestScanConfigDir_Unmanaged(t *testing.T) {
	tmpDir := t.TempDir()
	dotfilesDir := filepath.Join(tmpDir, "dotfiles")
	configDir := filepath.Join(tmpDir, "config")

	if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 管理対象外のディレクトリを作成
	unmanagedDir := filepath.Join(configDir, "someapp")
	if err := os.MkdirAll(unmanagedDir, 0755); err != nil {
		t.Fatal(err)
	}

	entries, err := ScanConfigDir(configDir, nil, dotfilesDir, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf("エントリ数が期待値と異なる: got %d, want 1", len(entries))
	}

	if entries[0].Status != StatusUnmanaged {
		t.Errorf("Status = %d, want StatusUnmanaged (%d)", entries[0].Status, StatusUnmanaged)
	}
	if entries[0].Name != "someapp" {
		t.Errorf("Name = %q, want %q", entries[0].Name, "someapp")
	}
}

func TestScanConfigDir_Drifted(t *testing.T) {
	tmpDir := t.TempDir()
	dotfilesDir := filepath.Join(tmpDir, "dotfiles")
	configDir := filepath.Join(tmpDir, "config")

	if err := os.MkdirAll(dotfilesDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatal(err)
	}

	// 実ディレクトリとして存在（シンボリックリンクではない）
	driftedDir := filepath.Join(configDir, "nvim")
	if err := os.MkdirAll(driftedDir, 0755); err != nil {
		t.Fatal(err)
	}

	// managed マップにはこのパスが登録されている
	managed := []ManagedLink{
		{
			Source: filepath.Join(dotfilesDir, "config", "nvim"),
			Dest:   driftedDir,
			Line:   5,
		},
	}

	entries, err := ScanConfigDir(configDir, managed, dotfilesDir, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf("エントリ数が期待値と異なる: got %d, want 1", len(entries))
	}

	if entries[0].Status != StatusDrifted {
		t.Errorf("Status = %d, want StatusDrifted (%d)", entries[0].Status, StatusDrifted)
	}
	if entries[0].ExpectedLink == nil {
		t.Error("ExpectedLink が nil")
	}
}

func TestIsSymlinkIntoDotfiles(t *testing.T) {
	tmpDir := t.TempDir()
	dotfilesDir := filepath.Join(tmpDir, "dotfiles")
	target := filepath.Join(dotfilesDir, "config", "nvim")

	// ターゲットディレクトリを作成
	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatal(err)
	}

	// シンボリックリンクを作成
	linkPath := filepath.Join(tmpDir, "nvim-link")
	if err := os.Symlink(target, linkPath); err != nil {
		t.Fatal(err)
	}

	isManaged, resolvedTarget, err := isSymlinkIntoDotfiles(linkPath, dotfilesDir)
	if err != nil {
		t.Fatal(err)
	}

	if !isManaged {
		t.Error("isManaged = false, want true")
	}
	if resolvedTarget != target {
		t.Errorf("target = %q, want %q", resolvedTarget, target)
	}

	// dotfiles 外を指すシンボリックリンク
	otherTarget := filepath.Join(tmpDir, "other")
	if err := os.MkdirAll(otherTarget, 0755); err != nil {
		t.Fatal(err)
	}
	otherLink := filepath.Join(tmpDir, "other-link")
	if err := os.Symlink(otherTarget, otherLink); err != nil {
		t.Fatal(err)
	}

	isManaged, _, err = isSymlinkIntoDotfiles(otherLink, dotfilesDir)
	if err != nil {
		t.Fatal(err)
	}
	if isManaged {
		t.Error("isManaged = true, want false (dotfiles 外を指している)")
	}
}
