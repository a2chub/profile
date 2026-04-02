package tui

import (
	tea "charm.land/bubbletea/v2"

	"github.com/a2chub/dotsync/internal/config"
	"github.com/a2chub/dotsync/internal/manifest"
)

// scanMsg はスキャン完了時に送信されるメッセージ
type scanMsg struct {
	result *manifest.ScanResult
	err    error
}

// runScan はフルスキャンを実行する tea.Cmd を返す
func runScan(cfg *config.Config) tea.Cmd {
	return func() tea.Msg {
		result, err := manifest.RunFullScan(cfg)
		return scanMsg{result: result, err: err}
	}
}
