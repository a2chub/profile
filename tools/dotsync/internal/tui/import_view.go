package tui

import (
	"fmt"
	"strings"
)

// renderImportView はインポート確認画面を描画する
func renderImportView(m Model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Import Confirmation"))
	b.WriteString("\n\n")

	b.WriteString("  The following items will be imported:\n\n")

	for i, item := range m.items {
		if !m.selected[i] {
			continue
		}
		action := "copy config + add link entry"
		if item.Status == "drifted" {
			action = "repair drift (copy + re-symlink)"
		} else if item.Kind == 1 { // KindBrew
			action = "add to Brewfile"
		}

		hasWarning := len(item.SecretWarnings) > 0
		line := fmt.Sprintf("    • %s (%s)", item.Name, action)
		if hasWarning {
			line += warningStyle.Render(" ⚠ has secret warnings")
		}
		b.WriteString(line)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	if m.dryRun {
		b.WriteString(warningStyle.Render("  [DRY RUN] No changes will be made."))
		b.WriteString("\n\n")
	}
	b.WriteString(helpStyle.Render("  y: confirm  n: cancel"))

	return b.String()
}

// renderDoneView はインポート完了画面を描画する
func renderDoneView(m Model) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Import Complete"))
	b.WriteString("\n\n")

	successCount := 0
	failCount := 0
	for _, r := range m.importResults {
		if len(r.Errors) == 0 {
			successCount++
		} else {
			failCount++
		}
	}

	b.WriteString(selectedStyle.Render(fmt.Sprintf("  Success: %d", successCount)))
	b.WriteString("\n")
	if failCount > 0 {
		b.WriteString(dangerStyle.Render(fmt.Sprintf("  Failed:  %d", failCount)))
		b.WriteString("\n\n")
		for _, r := range m.importResults {
			for _, e := range r.Errors {
				b.WriteString(dangerStyle.Render(fmt.Sprintf("    ✗ %s: %s", r.Plan.ItemName, e.Error())))
				b.WriteString("\n")
			}
		}
	}

	if m.commitMsg != "" {
		b.WriteString("\n  Suggested commit message:\n")
		b.WriteString(managedStyle.Render(fmt.Sprintf("    %s", m.commitMsg)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  q: quit"))

	return b.String()
}
