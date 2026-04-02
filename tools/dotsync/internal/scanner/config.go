package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ConfigStatus はコンフィグエントリの管理状態を表す
type ConfigStatus int

const (
	StatusManaged          ConfigStatus = iota // dotfiles からのシンボリックリンクで管理されている
	StatusPartiallyManaged                     // 一部のファイルのみ管理されている
	StatusDrifted                              // 管理対象だが実ディレクトリになっている
	StatusUnmanaged                            // dotfiles で管理されていない
)

// ManagedLink は link-dotfiles.sh で定義された1つのシンボリックリンクを表す
type ManagedLink struct {
	Source string // dotfiles リポジトリ内のパス
	Dest   string // システム上の配置先パス
	Line   int    // スクリプト内の行番号
}

// ConfigEntry は ~/.config/ 配下の1つのエントリを表す
type ConfigEntry struct {
	Name           string        // エントリ名（例: "nvim"）
	Path           string        // フルパス（例: "/Users/foo/.config/nvim"）
	Status         ConfigStatus  // 管理状態
	ManagedFiles   []string      // dotfiles で管理されているファイル
	UnmanagedFiles []string      // 管理されていないファイル
	ExpectedLink   *ManagedLink  // 管理対象の場合、期待されるリンク情報
}

// createLinkRegex は create_link 呼び出しにマッチする正規表現
var createLinkRegex = regexp.MustCompile(`create_link\s+"([^"]+)"\s+"([^"]+)"`)

// ParseLinkScript は link-dotfiles.sh を解析し、ManagedLink のリストを返す。
// パースできない行はエラーとして収集するが、パース可能な行はすべて返す。
func ParseLinkScript(scriptPath, dotfilesDir, homeDir string) ([]ManagedLink, []error) {
	f, err := os.Open(scriptPath)
	if err != nil {
		return nil, []error{fmt.Errorf("スクリプトを開けません: %w", err)}
	}
	defer f.Close()

	var links []ManagedLink
	var errs []error

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// コメント行をスキップ
		if strings.HasPrefix(line, "#") {
			continue
		}

		// create_link を含まない行はスキップ
		if !strings.Contains(line, "create_link") {
			continue
		}

		matches := createLinkRegex.FindStringSubmatch(line)
		if matches == nil {
			errs = append(errs, fmt.Errorf("行 %d: create_link のパースに失敗: %s", lineNum, line))
			continue
		}

		src := replaceVars(matches[1], dotfilesDir, homeDir)
		dest := replaceVars(matches[2], dotfilesDir, homeDir)

		links = append(links, ManagedLink{
			Source: src,
			Dest:   dest,
			Line:   lineNum,
		})
	}

	if err := scanner.Err(); err != nil {
		errs = append(errs, fmt.Errorf("スクリプトの読み取りエラー: %w", err))
	}

	return links, errs
}

// replaceVars は $DOTFILES_DIR と $HOME をそれぞれ実際のパスに置換する
func replaceVars(s, dotfilesDir, homeDir string) string {
	s = strings.ReplaceAll(s, "$DOTFILES_DIR", dotfilesDir)
	s = strings.ReplaceAll(s, "$HOME", homeDir)
	return s
}

// ScanConfigDir は configDir（通常 ~/.config/）配下のトップレベルエントリを走査し、
// 各エントリの管理状態を判定して返す。
// ignorePatterns に一致するエントリ名はスキップされる（filepath.Match で判定）。
func ScanConfigDir(configDir string, managed []ManagedLink, dotfilesDir string, ignorePatterns []string) ([]ConfigEntry, error) {
	entries, err := os.ReadDir(configDir)
	if err != nil {
		return nil, fmt.Errorf("ディレクトリの読み取りに失敗: %w", err)
	}

	// managed リンクを dest パスでインデックス化
	managedByDest := make(map[string]*ManagedLink)
	for i := range managed {
		managedByDest[managed[i].Dest] = &managed[i]
	}

	var results []ConfigEntry
	for _, entry := range entries {
		name := entry.Name()

		// .DS_Store をスキップ
		if name == ".DS_Store" {
			continue
		}

		// ignore パターンにマッチするエントリをスキップ
		if matchesIgnorePattern(name, ignorePatterns) {
			continue
		}

		fullPath := filepath.Join(configDir, name)

		ce := ConfigEntry{
			Name: name,
			Path: fullPath,
		}

		// 1. シンボリックリンクかチェック
		fi, err := os.Lstat(fullPath)
		if err != nil {
			continue
		}

		if fi.Mode()&os.ModeSymlink != 0 {
			isManaged, _, _ := isSymlinkIntoDotfiles(fullPath, dotfilesDir)
			if isManaged {
				ce.Status = StatusManaged
				if ml, ok := managedByDest[fullPath]; ok {
					ce.ExpectedLink = ml
				}
				ce.ManagedFiles = []string{fullPath}
				results = append(results, ce)
				continue
			}
		}

		// 2. 実ディレクトリで、managed マップに含まれている場合 → DRIFTED
		if ml, ok := managedByDest[fullPath]; ok {
			ce.Status = StatusDrifted
			ce.ExpectedLink = ml
			results = append(results, ce)
			continue
		}

		// 3. 実ディレクトリ内に dotfiles へのシンボリックリンクがあるか → PARTIALLY_MANAGED
		if fi.IsDir() {
			managedFiles, unmanagedFiles := scanDirForManagedLinks(fullPath, dotfilesDir)
			if len(managedFiles) > 0 {
				ce.Status = StatusPartiallyManaged
				ce.ManagedFiles = managedFiles
				ce.UnmanagedFiles = unmanagedFiles
				results = append(results, ce)
				continue
			}
		}

		// 4. それ以外 → UNMANAGED
		ce.Status = StatusUnmanaged
		results = append(results, ce)
	}

	return results, nil
}

// scanDirForManagedLinks はディレクトリ内の直下エントリを走査し、
// dotfiles へのシンボリックリンクとそれ以外を分類する
func scanDirForManagedLinks(dir, dotfilesDir string) (managed []string, unmanaged []string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil
	}

	for _, entry := range entries {
		if entry.Name() == ".DS_Store" {
			continue
		}
		fullPath := filepath.Join(dir, entry.Name())
		fi, err := os.Lstat(fullPath)
		if err != nil {
			continue
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			isM, _, _ := isSymlinkIntoDotfiles(fullPath, dotfilesDir)
			if isM {
				managed = append(managed, fullPath)
				continue
			}
		}
		unmanaged = append(unmanaged, fullPath)
	}
	return managed, unmanaged
}

// matchesIgnorePattern はエントリ名が ignore パターンのいずれかに一致するか判定する。
// パターンの末尾が "/" の場合はそれを除去して比較する。
func matchesIgnorePattern(name string, patterns []string) bool {
	for _, pattern := range patterns {
		// 末尾の "/" を除去（ディレクトリ指定のパターン対応）
		p := strings.TrimSuffix(pattern, "/")
		matched, err := filepath.Match(p, name)
		if err == nil && matched {
			return true
		}
	}
	return false
}

// isSymlinkIntoDotfiles は指定パスが dotfilesDir 配下を指すシンボリックリンクかどうかを判定する。
// 戻り値: (dotfiles 配下を指しているか, シンボリックリンクのターゲット, エラー)
func isSymlinkIntoDotfiles(path, dotfilesDir string) (bool, string, error) {
	target, err := os.Readlink(path)
	if err != nil {
		return false, "", fmt.Errorf("readlink に失敗: %w", err)
	}

	// 相対パスを絶対パスに変換
	if !filepath.IsAbs(target) {
		target = filepath.Join(filepath.Dir(path), target)
	}

	// パスを正規化
	target = filepath.Clean(target)
	dotfilesClean := filepath.Clean(dotfilesDir)

	if strings.HasPrefix(target, dotfilesClean+string(filepath.Separator)) || target == dotfilesClean {
		return true, target, nil
	}

	return false, target, nil
}
