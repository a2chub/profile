package importer

import (
	"strings"
	"testing"
)

func TestGenerateCommitMessage_Single(t *testing.T) {
	tests := []struct {
		name     string
		plan     *ImportPlan
		expected string
	}{
		{
			name:     "single config",
			plan:     &ImportPlan{ItemName: "wezterm", ItemKind: "config"},
			expected: "feat: add wezterm config to dotfiles",
		},
		{
			name:     "single brew",
			plan:     &ImportPlan{ItemName: "ripgrep", ItemKind: "brew"},
			expected: "feat: add ripgrep to Brewfile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateCommitMessage([]*ImportPlan{tt.plan})
			if got != tt.expected {
				t.Errorf("GenerateCommitMessage()\n  got:  %q\n  want: %q", got, tt.expected)
			}
		})
	}
}

func TestGenerateCommitMessage_Multiple(t *testing.T) {
	plans := []*ImportPlan{
		{ItemName: "wezterm", ItemKind: "config"},
		{ItemName: "alacritty", ItemKind: "config"},
		{ItemName: "ripgrep", ItemKind: "brew"},
	}

	got := GenerateCommitMessage(plans)

	// サブジェクト行の検証
	if !strings.HasPrefix(got, "feat: import wezterm, alacritty, ripgrep into dotfiles") {
		t.Errorf("unexpected subject line:\n%s", got)
	}

	// 本文に config と brew の詳細が含まれることを確認
	if !strings.Contains(got, "Configs: wezterm, alacritty") {
		t.Errorf("expected configs detail in body:\n%s", got)
	}
	if !strings.Contains(got, "Brewfile: ripgrep") {
		t.Errorf("expected brewfile detail in body:\n%s", got)
	}
}
