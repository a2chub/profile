package scanner

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// BrewEntry はBrewfileの1エントリを表す
type BrewEntry struct {
	Name string
	Kind string // brew, cask, tap
	Line int
}

// BrewDiff はBrewfileとbrew leavesの差分を表す
type BrewDiff struct {
	Missing []BrewEntry // brew leavesにあるがBrewfileにない
	Extra   []BrewEntry // Brewfileにあるがbrew leavesにない
}

var brewLineRe = regexp.MustCompile(`^(brew|cask|tap)\s+"([^"]+)"`)

// ParseBrewfile はBrewfileをパースしてエントリのスライスを返す
func ParseBrewfile(path string) ([]BrewEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Brewfileを開けません: %w", err)
	}
	defer f.Close()

	var entries []BrewEntry
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 空行とコメント行をスキップ
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		matches := brewLineRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		entries = append(entries, BrewEntry{
			Name: matches[2],
			Kind: matches[1],
			Line: lineNum,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Brewfileの読み込みエラー: %w", err)
	}

	return entries, nil
}

// RunBrewLeaves はbrew leavesを実行してパッケージ一覧を返す
func RunBrewLeaves() ([]string, error) {
	brewPath, err := exec.LookPath("brew")
	if err != nil {
		return nil, fmt.Errorf("brewが見つかりません: %w", err)
	}

	cmd := exec.Command(brewPath, "leaves")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("brew leavesの実行に失敗: %w", err)
	}

	return ParseLeaves(string(out)), nil
}

// ParseLeaves はbrew leavesの出力をパースしてスライスを返す
func ParseLeaves(output string) []string {
	var result []string
	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

// shortName はタップ付きの名前から短い名前（最後のパスコンポーネント）を取得する
func shortName(name string) string {
	parts := strings.Split(name, "/")
	return parts[len(parts)-1]
}

// DiffBrew はBrewfileエントリとbrew leavesの差分を計算する
// brewエントリのみを比較対象とする
func DiffBrew(brewfileEntries []BrewEntry, leaves []string) BrewDiff {
	// Brewfileのbrewエントリからルックアップマップを構築
	// フルネームと短い名前の両方を登録
	brewfileNames := make(map[string]bool)
	var brewEntries []BrewEntry
	for _, e := range brewfileEntries {
		if e.Kind != "brew" {
			continue
		}
		brewEntries = append(brewEntries, e)
		brewfileNames[e.Name] = true
		short := shortName(e.Name)
		if short != e.Name {
			brewfileNames[short] = true
		}
	}

	// leavesからルックアップマップを構築
	leavesSet := make(map[string]bool)
	for _, l := range leaves {
		leavesSet[l] = true
		short := shortName(l)
		if short != l {
			leavesSet[short] = true
		}
	}

	// Missing: leavesにあるがBrewfileにない
	var missing []BrewEntry
	for _, l := range leaves {
		short := shortName(l)
		if !brewfileNames[l] && !brewfileNames[short] {
			missing = append(missing, BrewEntry{Name: l, Kind: "brew"})
		}
	}

	// Extra: Brewfileにあるがleavesにない（情報提供用）
	var extra []BrewEntry
	for _, e := range brewEntries {
		short := shortName(e.Name)
		if !leavesSet[e.Name] && !leavesSet[short] {
			extra = append(extra, e)
		}
	}

	return BrewDiff{
		Missing: missing,
		Extra:   extra,
	}
}
