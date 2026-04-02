package importer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPlanDriftRepair(t *testing.T) {
	plan := PlanDriftRepair("zellij", "/home/user/.config/zellij", "/home/user/dotfiles", "/home/user/dotfiles/scripts/link-dotfiles.sh")

	if plan.ItemName != "zellij" {
		t.Errorf("ItemName: got %q, want %q", plan.ItemName, "zellij")
	}
	if plan.ItemKind != "drift" {
		t.Errorf("ItemKind: got %q, want %q", plan.ItemKind, "drift")
	}
	if len(plan.Actions) != 2 {
		t.Fatalf("Actions count: got %d, want 2", len(plan.Actions))
	}

	// 1番目のアクション: CopyDir
	a0 := plan.Actions[0]
	if a0.Type != ActionCopyDir {
		t.Errorf("Action[0].Type: got %d, want %d (ActionCopyDir)", a0.Type, ActionCopyDir)
	}
	if a0.Source != "/home/user/.config/zellij" {
		t.Errorf("Action[0].Source: got %q, want %q", a0.Source, "/home/user/.config/zellij")
	}
	expectedDest := filepath.Join("/home/user/dotfiles", "config", "zellij")
	if a0.Dest != expectedDest {
		t.Errorf("Action[0].Dest: got %q, want %q", a0.Dest, expectedDest)
	}

	// 2番目のアクション: RepairDrift
	a1 := plan.Actions[1]
	if a1.Type != ActionRepairDrift {
		t.Errorf("Action[1].Type: got %d, want %d (ActionRepairDrift)", a1.Type, ActionRepairDrift)
	}
	// RepairDrift: Source=dotfiles側, Dest=システム側
	if a1.Source != expectedDest {
		t.Errorf("Action[1].Source: got %q, want %q", a1.Source, expectedDest)
	}
	if a1.Dest != "/home/user/.config/zellij" {
		t.Errorf("Action[1].Dest: got %q, want %q", a1.Dest, "/home/user/.config/zellij")
	}
}

func TestRepairDrift(t *testing.T) {
	// テンポラリディレクトリをセットアップ
	tmpDir := t.TempDir()

	// dotfiles 側のディレクトリを作成（コピー済みを模擬）
	dotfilesDest := filepath.Join(tmpDir, "dotfiles", "config", "zellij")
	if err := os.MkdirAll(dotfilesDest, 0755); err != nil {
		t.Fatal(err)
	}
	// dotfiles 側にファイルを配置
	if err := os.WriteFile(filepath.Join(dotfilesDest, "config.kdl"), []byte("test config"), 0644); err != nil {
		t.Fatal(err)
	}

	// システム側の実ディレクトリを作成（drifted 状態を模擬）
	systemPath := filepath.Join(tmpDir, "home", ".config", "zellij")
	if err := os.MkdirAll(systemPath, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(systemPath, "config.kdl"), []byte("old config"), 0644); err != nil {
		t.Fatal(err)
	}

	// repairDrift を実行
	if err := repairDrift(dotfilesDest, systemPath); err != nil {
		t.Fatalf("repairDrift failed: %v", err)
	}

	// systemPath がシンボリックリンクになっていることを確認
	linkInfo, err := os.Lstat(systemPath)
	if err != nil {
		t.Fatalf("Lstat failed: %v", err)
	}
	if linkInfo.Mode()&os.ModeSymlink == 0 {
		t.Errorf("systemPath is not a symlink: mode=%v", linkInfo.Mode())
	}

	// シンボリックリンクの参照先が正しいことを確認
	target, err := os.Readlink(systemPath)
	if err != nil {
		t.Fatalf("Readlink failed: %v", err)
	}
	if target != dotfilesDest {
		t.Errorf("symlink target: got %q, want %q", target, dotfilesDest)
	}

	// シンボリックリンク経由でファイルにアクセスできることを確認
	content, err := os.ReadFile(filepath.Join(systemPath, "config.kdl"))
	if err != nil {
		t.Fatalf("ReadFile via symlink failed: %v", err)
	}
	if string(content) != "test config" {
		t.Errorf("file content via symlink: got %q, want %q", string(content), "test config")
	}
}

func TestRepairDrift_MissingDotfilesDest(t *testing.T) {
	tmpDir := t.TempDir()
	dotfilesDest := filepath.Join(tmpDir, "nonexistent")
	systemPath := filepath.Join(tmpDir, "system")

	if err := os.MkdirAll(systemPath, 0755); err != nil {
		t.Fatal(err)
	}

	err := repairDrift(dotfilesDest, systemPath)
	if err == nil {
		t.Error("expected error when dotfiles dest does not exist, got nil")
	}
}
