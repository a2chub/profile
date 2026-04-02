package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// GenerateLinkLine は指定された設定名に対する create_link コマンド文字列を返す。
func GenerateLinkLine(configName string) string {
	return fmt.Sprintf(`create_link "$DOTFILES_DIR/config/%s" "$HOME/.config/%s"`, configName, configName)
}

// AppendLinkEntry は link-dotfiles.sh に create_link エントリを追加する。
// 最後の print_success 行の直前に挿入する。
// 冪等性: 既に同じエントリが存在する場合は何もしない。
// アトミック書き込み: 一時ファイルに書き込んでからリネームする。
func AppendLinkEntry(scriptPath, configName string) error {
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("read script: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	linkLine := GenerateLinkLine(configName)

	// 冪等性チェック: 既にエントリが存在するか確認
	for _, line := range lines {
		if strings.Contains(line, fmt.Sprintf("config/%s\"", configName)) {
			return nil
		}
	}

	// 最後の print_success 行を見つけて、その直前に挿入する
	insertIdx := -1
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.Contains(lines[i], "print_success") {
			insertIdx = i
			break
		}
	}

	if insertIdx < 0 {
		// print_success が見つからない場合は末尾に追加
		lines = append(lines, linkLine)
	} else {
		// print_success の直前に空行とエントリを挿入
		newLines := make([]string, 0, len(lines)+2)
		newLines = append(newLines, lines[:insertIdx]...)
		newLines = append(newLines, linkLine)
		newLines = append(newLines, "")
		newLines = append(newLines, lines[insertIdx:]...)
		lines = newLines
	}

	newContent := strings.Join(lines, "\n")

	// アトミック書き込み: 一時ファイルに書き込んでからリネーム
	dir := filepath.Dir(scriptPath)
	tmpFile, err := os.CreateTemp(dir, ".link-dotfiles-*.tmp")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.WriteString(newContent); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("write temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("close temp file: %w", err)
	}

	// 元ファイルのパーミッションを保持
	info, err := os.Stat(scriptPath)
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("stat original: %w", err)
	}
	if err := os.Chmod(tmpPath, info.Mode()); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("chmod temp file: %w", err)
	}

	if err := os.Rename(tmpPath, scriptPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename temp file: %w", err)
	}

	return nil
}
