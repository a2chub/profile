package tui

import (
	"fmt"
	"strings"

	"github.com/a2chub/dotsync/internal/manifest"
	"github.com/a2chub/dotsync/internal/scanner"
)

// renderDetailView はアイテムの詳細画面を描画する
func renderDetailView(m Model, item manifest.ScanItem, width int) string {
	var b strings.Builder

	// アイテム名とステータス
	b.WriteString(titleStyle.Render(fmt.Sprintf("Detail: %s", item.Name)))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("  Kind:   %s\n", kindLabel(item.Kind)))
	b.WriteString(fmt.Sprintf("  Status: %s\n", item.Status))
	b.WriteString(fmt.Sprintf("  Path:   %s\n", item.Path))

	// ファイル一覧
	if len(item.Files) > 0 {
		b.WriteString("\n  Files:\n")
		for _, f := range item.Files {
			sizeStr := formatSize(f.Size)
			symInfo := ""
			if f.IsSymlink {
				symInfo = fmt.Sprintf(" -> %s", f.SymTarget)
			}
			b.WriteString(fmt.Sprintf("    %s  %s%s\n", sizeStr, f.RelPath, symInfo))
		}
	}

	// シークレット警告
	if len(item.SecretWarnings) > 0 {
		b.WriteString("\n  Secret Warnings:\n")
		for _, w := range item.SecretWarnings {
			style := warningStyle
			icon := "⚠"
			if w.Risk == scanner.RiskHigh {
				style = dangerStyle
				icon = "🚨"
			}
			b.WriteString(style.Render(fmt.Sprintf("    %s %s (risk: %s, entropy: %.2f)",
				icon, w.Reason, riskLabel(w.Risk), w.Entropy)))
			b.WriteString("\n")
			b.WriteString(style.Render(fmt.Sprintf("      file: %s", w.FilePath)))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  esc: back"))

	return b.String()
}

func kindLabel(k manifest.ItemKind) string {
	switch k {
	case manifest.KindConfig:
		return "Config"
	case manifest.KindBrew:
		return "Brew"
	default:
		return "Unknown"
	}
}

func riskLabel(r scanner.SecretRisk) string {
	switch r {
	case scanner.RiskLow:
		return "low"
	case scanner.RiskMedium:
		return "medium"
	case scanner.RiskHigh:
		return "high"
	default:
		return "none"
	}
}

func formatSize(size int64) string {
	switch {
	case size >= 1024*1024:
		return fmt.Sprintf("%5.1fM", float64(size)/(1024*1024))
	case size >= 1024:
		return fmt.Sprintf("%5.1fK", float64(size)/1024)
	default:
		return fmt.Sprintf("%5dB", size)
	}
}
