package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/a2chub/dotsync/internal/config"
	"github.com/a2chub/dotsync/internal/scanner"
)

// ItemKind はスキャンアイテムの種別を表す
type ItemKind int

const (
	KindConfig ItemKind = iota
	KindBrew
)

// ScanItem はスキャン結果の1アイテムを表す
type ScanItem struct {
	ID             int
	Kind           ItemKind
	Name           string
	Path           string // ファイルシステムパス（config）またはパッケージ名（brew）
	Status         string // "unmanaged", "drifted", "partially-managed", "missing-from-brewfile"
	Files          []FileInfo
	SecretWarnings []scanner.SecretWarning
	Selected       bool
}

// FileInfo はファイルの情報を表す
type FileInfo struct {
	RelPath   string
	IsSymlink bool
	SymTarget string
	Size      int64
}

// ScanError は非致命的なスキャンエラーを表す
type ScanError struct {
	Context string
	Err     error
}

// ScanResult はフルスキャンの結果を表す
type ScanResult struct {
	Items       []ScanItem
	Errors      []ScanError
	DotfilesDir string
	HomeDir     string
	ScannedAt   time.Time
}

// statusString は ConfigStatus を文字列に変換する
func statusString(s scanner.ConfigStatus) string {
	switch s {
	case scanner.StatusManaged:
		return "managed"
	case scanner.StatusUnmanaged:
		return "unmanaged"
	case scanner.StatusDrifted:
		return "drifted"
	case scanner.StatusPartiallyManaged:
		return "partially-managed"
	default:
		return "unmanaged"
	}
}

// collectFileInfo はディレクトリ配下（1階層）のファイル情報を収集する
func collectFileInfo(dirPath string) []FileInfo {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil
	}

	var files []FileInfo
	for _, entry := range entries {
		if entry.Name() == ".DS_Store" {
			continue
		}

		fullPath := filepath.Join(dirPath, entry.Name())
		fi, err := os.Lstat(fullPath)
		if err != nil {
			continue
		}

		info := FileInfo{
			RelPath: entry.Name(),
			Size:    fi.Size(),
		}

		if fi.Mode()&os.ModeSymlink != 0 {
			info.IsSymlink = true
			target, err := os.Readlink(fullPath)
			if err == nil {
				info.SymTarget = target
			}
		}

		files = append(files, info)
	}

	return files
}

// expandHome は ~ をホームディレクトリに展開する
func expandHome(path, homeDir string) string {
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

// RunFullScan はすべてのスキャナを実行し、統合された ScanResult を返す
func RunFullScan(cfg *config.Config) (*ScanResult, error) {
	dotfilesDir := cfg.Dotfiles.Repo
	scriptPath := cfg.ScriptPath()
	brewfilePath := cfg.BrewfilePath()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("ホームディレクトリの取得に失敗: %w", err)
	}

	result := &ScanResult{
		DotfilesDir: dotfilesDir,
		HomeDir:     homeDir,
		ScannedAt:   time.Now(),
	}

	nextID := 1

	// 1. link-dotfiles.sh をパースして管理対象リンクを取得
	managed, parseErrs := scanner.ParseLinkScript(scriptPath, dotfilesDir, homeDir)
	for _, e := range parseErrs {
		result.Errors = append(result.Errors, ScanError{
			Context: "link-dotfiles.sh のパース",
			Err:     e,
		})
	}

	// 2. WatchDirs をスキャン
	var allConfigEntries []scanner.ConfigEntry
	for _, watchDir := range cfg.Scan.WatchDirs {
		expandedDir := expandHome(watchDir, homeDir)
		configEntries, err := scanner.ScanConfigDir(expandedDir, managed, dotfilesDir, cfg.Scan.Ignore)
		if err != nil {
			result.Errors = append(result.Errors, ScanError{
				Context: fmt.Sprintf("%s ディレクトリのスキャン", watchDir),
				Err:     err,
			})
			continue
		}
		allConfigEntries = append(allConfigEntries, configEntries...)
	}

	// 3. 全エントリのファイル情報を収集（管理済みは表示のみ、それ以外はシークレットスキャンも実行）
	for _, entry := range allConfigEntries {
		status := statusString(entry.Status)

		var secretWarnings []scanner.SecretWarning
		if entry.Status != scanner.StatusManaged {
			warnings, scanErr := scanner.ScanForSecrets(entry.Path, cfg.Secrets.EntropyThreshold)
			if scanErr != nil {
				result.Errors = append(result.Errors, ScanError{
					Context: fmt.Sprintf("シークレットスキャン: %s", entry.Name),
					Err:     scanErr,
				})
			} else {
				secretWarnings = warnings
			}
		}

		files := collectFileInfo(entry.Path)

		item := ScanItem{
			ID:             nextID,
			Kind:           KindConfig,
			Name:           entry.Name,
			Path:           entry.Path,
			Status:         status,
			Files:          files,
			SecretWarnings: secretWarnings,
		}
		result.Items = append(result.Items, item)
		nextID++
	}

	// 4. Brewfile をパースし、brew leaves を取得
	brewEntries, err := scanner.ParseBrewfile(brewfilePath)
	if err != nil {
		result.Errors = append(result.Errors, ScanError{
			Context: "Brewfile のパース",
			Err:     err,
		})
	}

	leaves, err := scanner.RunBrewLeaves()
	if err != nil {
		result.Errors = append(result.Errors, ScanError{
			Context: "brew leaves の実行",
			Err:     err,
		})
	}

	// 5. Brew の差分を計算（Brewfile パースと leaves 取得の両方が成功した場合のみ）
	if brewEntries != nil && leaves != nil {
		diff := scanner.DiffBrew(brewEntries, leaves)

		// Missing: leaves にあるが Brewfile にないパッケージ
		for _, missing := range diff.Missing {
			item := ScanItem{
				ID:     nextID,
				Kind:   KindBrew,
				Name:   missing.Name,
				Path:   missing.Name,
				Status: "missing-from-brewfile",
			}
			result.Items = append(result.Items, item)
			nextID++
		}
	}

	return result, nil
}
