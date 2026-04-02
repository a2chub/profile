package importer

import (
	"fmt"
	"strings"
)

// GenerateCommitMessage は ImportPlan のリストから conventional commit メッセージを生成する。
// 単一の config: "feat: add {name} config to dotfiles"
// 単一の brew: "feat: add {name} to Brewfile"
// 複数: "feat: import {name1}, {name2}, ... into dotfiles" + 本文に詳細
func GenerateCommitMessage(plans []*ImportPlan) string {
	if len(plans) == 0 {
		return ""
	}

	if len(plans) == 1 {
		p := plans[0]
		switch p.ItemKind {
		case "config":
			return fmt.Sprintf("feat: add %s config to dotfiles", p.ItemName)
		case "brew":
			return fmt.Sprintf("feat: add %s to Brewfile", p.ItemName)
		case "drift":
			return fmt.Sprintf("fix: repair drifted %s symlink", p.ItemName)
		default:
			return fmt.Sprintf("feat: add %s to dotfiles", p.ItemName)
		}
	}

	// 複数アイテムの場合
	names := make([]string, 0, len(plans))
	for _, p := range plans {
		names = append(names, p.ItemName)
	}

	subject := fmt.Sprintf("feat: import %s into dotfiles", strings.Join(names, ", "))

	// 本文に種別ごとの詳細を記載
	var configs, brews, drifts []string
	for _, p := range plans {
		switch p.ItemKind {
		case "config":
			configs = append(configs, p.ItemName)
		case "brew":
			brews = append(brews, p.ItemName)
		case "drift":
			drifts = append(drifts, p.ItemName)
		}
	}

	var body strings.Builder
	if len(configs) > 0 {
		body.WriteString("\nConfigs: ")
		body.WriteString(strings.Join(configs, ", "))
	}
	if len(brews) > 0 {
		body.WriteString("\nBrewfile: ")
		body.WriteString(strings.Join(brews, ", "))
	}
	if len(drifts) > 0 {
		body.WriteString("\nRepaired: ")
		body.WriteString(strings.Join(drifts, ", "))
	}

	return subject + "\n" + body.String()
}
