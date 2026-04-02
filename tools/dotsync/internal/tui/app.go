package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"

	"github.com/a2chub/dotsync/internal/config"
	"github.com/a2chub/dotsync/internal/importer"
	"github.com/a2chub/dotsync/internal/manifest"
)

// Screen はTUI画面の状態を表す
type Screen int

const (
	ScreenScan   Screen = iota // スキャン中（スピナー表示）
	ScreenList                 // アイテム一覧
	ScreenDetail               // アイテム詳細
	ScreenImport               // インポート確認
	ScreenDone                 // 完了画面
)

// Model はBubble Teaのルートモデル
type Model struct {
	screen        Screen
	scanResult    *manifest.ScanResult
	items         []manifest.ScanItem
	cursor        int
	offset        int // リスト表示のスクロールオフセット
	selected      map[int]bool
	importing     bool
	importResults []*importer.ImportResult
	commitMsg     string
	width, height int
	spinner       spinner.Model
	err           error
	cfg           *config.Config
	dryRun        bool
}

// importDoneMsg はインポート完了時に送信されるメッセージ
type importDoneMsg struct {
	results   []*importer.ImportResult
	commitMsg string
}

// NewModel は新しいModelを作成する
func NewModel(cfg *config.Config, dryRun bool) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		screen:   ScreenScan,
		selected: make(map[int]bool),
		spinner:  s,
		cfg:      cfg,
		dryRun:   dryRun,
	}
}

// Init はBubble Teaの初期化コマンドを返す
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		runScan(m.cfg),
	)
}

// Update はメッセージを処理して新しいモデルとコマンドを返す
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case scanMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.scanResult = msg.result
		m.items = msg.result.Items
		m.screen = ScreenList
		return m, nil

	case importDoneMsg:
		m.importing = false
		m.importResults = msg.results
		m.commitMsg = msg.commitMsg
		m.screen = ScreenDone
		return m, nil

	case spinner.TickMsg:
		if m.screen == ScreenScan {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case tea.KeyPressMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

// handleKey はキーイベントを処理する
func (m Model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	key := msg.Key()

	switch m.screen {
	case ScreenList:
		return m.handleListKey(key)
	case ScreenDetail:
		if key.Code == tea.KeyEscape {
			m.screen = ScreenList
			return m, nil
		}
	case ScreenImport:
		return m.handleImportKey(key)
	case ScreenDone:
		if key.Code == 'q' || key.Code == tea.KeyEscape {
			return m, tea.Quit
		}
	}

	// グローバルキー
	if key.Code == 'q' && m.screen != ScreenImport {
		return m, tea.Quit
	}

	return m, nil
}

// listViewHeight はリストに使える行数を返す（ヘッダー・フッター分を除く）
func (m Model) listViewHeight() int {
	// 固定行の正確なカウント:
	//   タイトル行: 1
	//   空行: 1
	//   統計行: 1
	//   空行: 1
	//   スクロール上インジケータ: 1（表示時）
	//   セクションラベル (Config:/Brew:): 最大3行（ラベル + 空行）
	//   スクロール下インジケータ: 1（表示時）
	//   エラー行: 最大2行
	//   空行: 1
	//   ステータスバー: 1
	//   ヘルプ行: 1
	//   余白: 1
	overhead := 15
	h := m.height - overhead
	if h < 3 {
		h = 3
	}
	return h
}

// calcVisibleEnd はビューポートに収まる最後のアイテムインデックス+1を返す。
// セクションラベル行も表示行数としてカウントする。
func (m Model) calcVisibleEnd(start, viewH int) (end int, sectionLines int) {
	usedLines := 0
	lastKind := manifest.ItemKind(-1)
	if start > 0 && start < len(m.items) {
		lastKind = m.items[start].Kind
	}

	for i := start; i < len(m.items); i++ {
		item := m.items[i]

		// セクションラベルの行数を加算
		if item.Kind != lastKind {
			extra := 1 // セクションラベル行
			if lastKind >= 0 {
				extra++ // セクション間の空行
			}
			if usedLines+extra+1 > viewH { // ラベル + 最低1アイテムが入らなければ終了
				return i, sectionLines
			}
			usedLines += extra
			sectionLines += extra
			lastKind = item.Kind
		}

		usedLines++ // アイテム行
		if usedLines > viewH {
			return i, sectionLines
		}
	}
	return len(m.items), sectionLines
}

// adjustOffset はカーソルがビューポート内に収まるようにオフセットを調整する
func (m *Model) adjustOffset() {
	viewH := m.listViewHeight()

	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	// カーソルが見える範囲を calcVisibleEnd で確認
	end, _ := m.calcVisibleEnd(m.offset, viewH)
	for m.cursor >= end && m.offset < len(m.items)-1 {
		m.offset++
		end, _ = m.calcVisibleEnd(m.offset, viewH)
	}
}

func (m Model) handleListKey(key tea.Key) (tea.Model, tea.Cmd) {
	switch key.Code {
	case 'j', tea.KeyDown:
		if m.cursor < len(m.items)-1 {
			m.cursor++
			m.adjustOffset()
		}
	case 'k', tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
			m.adjustOffset()
		}
	case tea.KeySpace:
		if m.cursor < len(m.items) && m.items[m.cursor].Status != "managed" {
			m.selected[m.cursor] = !m.selected[m.cursor]
		}
	case 'a':
		allSelected := m.allSelected()
		for i := range m.items {
			if m.items[i].Status != "managed" {
				m.selected[i] = !allSelected
			}
		}
	case 'i':
		if m.selectedCount() > 0 {
			m.screen = ScreenImport
		}
	case tea.KeyEnter:
		if len(m.items) > 0 {
			m.screen = ScreenDetail
		}
	case 'q':
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleImportKey(key tea.Key) (tea.Model, tea.Cmd) {
	switch key.Code {
	case 'y':
		if !m.importing {
			m.importing = true
			return m, m.runImport()
		}
	case 'n', tea.KeyEscape:
		m.screen = ScreenList
	}
	return m, nil
}

// runImport はインポートを実行する tea.Cmd を返す
func (m Model) runImport() tea.Cmd {
	// 選択されたアイテムを収集
	var selectedItems []manifest.ScanItem
	for i, item := range m.items {
		if m.selected[i] {
			selectedItems = append(selectedItems, item)
		}
	}

	dotfilesDir := m.cfg.Dotfiles.Repo
	scriptPath := m.cfg.ScriptPath()
	brewfilePath := m.cfg.BrewfilePath()
	dryRun := m.dryRun

	return func() tea.Msg {
		var plans []*importer.ImportPlan
		var results []*importer.ImportResult

		for _, item := range selectedItems {
			var plan *importer.ImportPlan
			switch item.Kind {
			case manifest.KindConfig:
				if item.Status == "drifted" {
					plan = importer.PlanDriftRepair(item.Name, item.Path, dotfilesDir, scriptPath)
				} else {
					plan = importer.PlanConfigImport(item.Name, item.Path, dotfilesDir, scriptPath)
				}
			case manifest.KindBrew:
				plan = importer.PlanBrewImport(item.Name, brewfilePath)
			}
			if plan == nil {
				continue
			}
			plans = append(plans, plan)

			if dryRun {
				results = append(results, &importer.ImportResult{Plan: plan})
			} else {
				results = append(results, importer.ExecuteImport(plan))
			}
		}

		commitMsg := importer.GenerateCommitMessage(plans)
		return importDoneMsg{results: results, commitMsg: commitMsg}
	}
}

// View はTUI画面を描画する
func (m Model) View() tea.View {
	var content string

	switch m.screen {
	case ScreenScan:
		content = m.viewScan()
	case ScreenList:
		content = m.viewList()
	case ScreenDetail:
		if m.cursor >= 0 && m.cursor < len(m.items) {
			content = renderDetailView(m, m.items[m.cursor], m.width)
		}
	case ScreenImport:
		content = renderImportView(m)
	case ScreenDone:
		content = renderDoneView(m)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

func (m Model) viewScan() string {
	return fmt.Sprintf("\n  %s Scanning dotfiles and packages...\n", m.spinner.View())
}

func (m Model) viewList() string {
	var b strings.Builder

	// ヘッダー
	b.WriteString(titleStyle.Render("dotsync — Reverse Sync TUI"))
	b.WriteString("\n\n")

	// 統計
	managedCount := 0
	unmanagedCount := 0
	brewCount := 0
	for _, item := range m.items {
		if item.Status == "managed" {
			managedCount++
		} else if item.Kind == manifest.KindBrew {
			brewCount++
		} else {
			unmanagedCount++
		}
	}
	b.WriteString(fmt.Sprintf("  Managed: %d | Unmanaged: %d | Missing Brew: %d\n\n",
		managedCount, unmanagedCount, brewCount))

	if len(m.items) == 0 {
		b.WriteString(selectedStyle.Render("  Everything is in sync!"))
		b.WriteString("\n")
	} else {
		viewH := m.listViewHeight()

		// ビューポートに収まるアイテム数を計算（セクションラベル行も考慮）
		visibleStart := m.offset
		visibleEnd, sectionLines := m.calcVisibleEnd(visibleStart, viewH)

		_ = sectionLines // セクションラベル分は viewH から差し引き済み

		// スクロールインジケータ（上）
		if visibleStart > 0 {
			b.WriteString(helpStyle.Render(fmt.Sprintf("  ↑ %d more above", visibleStart)))
			b.WriteString("\n")
		}

		// 表示範囲のアイテムをセクション分けして描画
		m.renderVisibleItems(&b, visibleStart, visibleEnd)

		// スクロールインジケータ（下）
		if visibleEnd < len(m.items) {
			b.WriteString(helpStyle.Render(fmt.Sprintf("  ↓ %d more below", len(m.items)-visibleEnd)))
			b.WriteString("\n")
		}
	}

	// スキャンエラー表示
	if m.scanResult != nil && len(m.scanResult.Errors) > 0 {
		b.WriteString("\n")
		b.WriteString(warningStyle.Render(fmt.Sprintf("  ⚠ %d scan warning(s)", len(m.scanResult.Errors))))
		b.WriteString("\n")
	}

	// ステータスバー
	b.WriteString("\n")
	selected := m.selectedCount()
	bar := fmt.Sprintf(" %d/%d selected ", selected, len(m.items))
	if m.dryRun {
		bar += "| DRY RUN "
	}
	b.WriteString(statusBarStyle.Render(bar))
	b.WriteString("\n")

	// ヘルプ
	b.WriteString(helpStyle.Render("  j/k: move  space: toggle  a: all  i: import  enter: detail  q: quit"))
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderVisibleItems(b *strings.Builder, start, end int) {
	lastKind := manifest.ItemKind(-1)
	for i := start; i < end; i++ {
		item := m.items[i]

		// セクションラベル（種別が変わったタイミングで挿入）
		if item.Kind != lastKind {
			if lastKind >= 0 {
				b.WriteString("\n")
			}
			switch item.Kind {
			case manifest.KindConfig:
				b.WriteString("  Config:\n")
			case manifest.KindBrew:
				b.WriteString("  Brew:\n")
			}
			lastKind = item.Kind
		}

		m.renderItem(b, i, item)
	}
}

func (m Model) renderItem(b *strings.Builder, i int, item manifest.ScanItem) {
	isManaged := item.Status == "managed"

	cursor := "  "
	if i == m.cursor {
		cursor = "> "
	}

	check := "○"
	if isManaged {
		check = "✓"
	} else if m.selected[i] {
		check = "●"
	}

	status := fmt.Sprintf("[%s]", item.Status)

	warning := ""
	if len(item.SecretWarnings) > 0 {
		warning = warningStyle.Render(" ⚠")
	}

	line := fmt.Sprintf("  %s%s %s %s%s", cursor, check, item.Name, managedStyle.Render(status), warning)
	if isManaged {
		line = managedStyle.Render(line)
	} else if i == m.cursor {
		line = selectedStyle.Render(line)
	}

	b.WriteString(line)
	b.WriteString("\n")
}

func (m Model) selectedCount() int {
	count := 0
	for _, v := range m.selected {
		if v {
			count++
		}
	}
	return count
}

func (m Model) allSelected() bool {
	if len(m.items) == 0 {
		return false
	}
	for i := range m.items {
		if !m.selected[i] {
			return false
		}
	}
	return true
}

// CommitMsg はプログラム終了後に表示するコミットメッセージを返す
func (m Model) CommitMsg() string {
	return m.commitMsg
}
