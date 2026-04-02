package importer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AppendBrewEntry は Brewfile にパッケージエントリを追加する。
// 冪等性: 既に同じエントリが存在する場合は何もしない。
// アトミック書き込み: 一時ファイルに書き込んでからリネームする。
func AppendBrewEntry(brewfilePath, packageName string) error {
	content, err := os.ReadFile(brewfilePath)
	if err != nil {
		return fmt.Errorf("read brewfile: %w", err)
	}

	entry := fmt.Sprintf(`brew "%s"`, packageName)

	// 冪等性チェック: 既にエントリが存在するか確認
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == entry {
			return nil
		}
	}

	// 末尾に追加（末尾の改行を保持）
	newContent := string(content)
	if !strings.HasSuffix(newContent, "\n") {
		newContent += "\n"
	}
	newContent += entry + "\n"

	// アトミック書き込み: 一時ファイルに書き込んでからリネーム
	dir := filepath.Dir(brewfilePath)
	tmpFile, err := os.CreateTemp(dir, ".Brewfile-*.tmp")
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
	info, err := os.Stat(brewfilePath)
	if err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("stat original: %w", err)
	}
	if err := os.Chmod(tmpPath, info.Mode()); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("chmod temp file: %w", err)
	}

	if err := os.Rename(tmpPath, brewfilePath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("rename temp file: %w", err)
	}

	return nil
}
