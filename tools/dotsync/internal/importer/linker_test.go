package importer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const sampleLinkScript = `#!/bin/bash
set -e

DOTFILES_DIR="$1"

create_link() {
    local src="$1"
    local dest="$2"
    ln -s "$src" "$dest"
}

echo "Creating symlinks..."

create_link "$DOTFILES_DIR/.zshrc" "$HOME/.zshrc"
create_link "$DOTFILES_DIR/config/nvim" "$HOME/.config/nvim"

print_success "All symlinks created"
`

func TestAppendLinkEntry(t *testing.T) {
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "link-dotfiles.sh")

	if err := os.WriteFile(scriptPath, []byte(sampleLinkScript), 0755); err != nil {
		t.Fatal(err)
	}

	if err := AppendLinkEntry(scriptPath, "wezterm"); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatal(err)
	}

	result := string(content)

	// エントリが存在することを確認
	expectedLine := GenerateLinkLine("wezterm")
	if !strings.Contains(result, expectedLine) {
		t.Errorf("expected link line not found in output:\n%s", result)
	}

	// print_success の前に挿入されていることを確認
	lines := strings.Split(result, "\n")
	linkIdx := -1
	printIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "wezterm") {
			linkIdx = i
		}
		if strings.Contains(line, "print_success") {
			printIdx = i
		}
	}

	if linkIdx < 0 {
		t.Fatal("link entry not found")
	}
	if printIdx < 0 {
		t.Fatal("print_success not found")
	}
	if linkIdx >= printIdx {
		t.Errorf("link entry (line %d) should appear before print_success (line %d)", linkIdx, printIdx)
	}
}

func TestAppendLinkEntry_Idempotent(t *testing.T) {
	dir := t.TempDir()
	scriptPath := filepath.Join(dir, "link-dotfiles.sh")

	if err := os.WriteFile(scriptPath, []byte(sampleLinkScript), 0755); err != nil {
		t.Fatal(err)
	}

	// 同じエントリを2回追加
	if err := AppendLinkEntry(scriptPath, "alacritty"); err != nil {
		t.Fatal(err)
	}
	if err := AppendLinkEntry(scriptPath, "alacritty"); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatal(err)
	}

	// エントリが1行だけ存在することを確認
	lines := strings.Split(string(content), "\n")
	count := 0
	for _, line := range lines {
		if strings.Contains(line, "config/alacritty") {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected 1 line containing alacritty entry, got %d", count)
	}
}

func TestGenerateLinkLine(t *testing.T) {
	got := GenerateLinkLine("wezterm")
	expected := `create_link "$DOTFILES_DIR/config/wezterm" "$HOME/.config/wezterm"`
	if got != expected {
		t.Errorf("GenerateLinkLine(\"wezterm\")\n  got:  %s\n  want: %s", got, expected)
	}
}
