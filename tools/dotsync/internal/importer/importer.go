package importer

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// ActionType はインポート時に実行するアクションの種別を表す。
type ActionType int

const (
	ActionCopyDir ActionType = iota
	ActionCopyFile
	ActionAppendScript
	ActionAppendBrewfile
	ActionRepairDrift
)

// Action はインポートプラン内の個別アクションを表す。
type Action struct {
	Type        ActionType
	Description string
	Source      string
	Dest        string
}

// ImportPlan はインポート対象のアイテムと実行するアクション群をまとめたものを表す。
type ImportPlan struct {
	ItemName string
	ItemKind string // "config" or "brew"
	Actions  []Action
}

// ImportResult はインポート実行結果を表す。
type ImportResult struct {
	Plan   *ImportPlan
	Errors []error
}

// PlanConfigImport は設定のインポートプランを作成する。
// sourcePath がディレクトリならディレクトリコピー、ファイルならファイルコピーを計画する。
func PlanConfigImport(name, sourcePath, dotfilesDir, scriptPath string) *ImportPlan {
	dest := filepath.Join(dotfilesDir, "config", name)

	// ソースがファイルかディレクトリかで分岐
	info, err := os.Stat(sourcePath)
	actionType := ActionCopyDir
	if err == nil && !info.IsDir() {
		actionType = ActionCopyFile
	}

	return &ImportPlan{
		ItemName: name,
		ItemKind: "config",
		Actions: []Action{
			{
				Type:        actionType,
				Description: fmt.Sprintf("Copy %s to %s", sourcePath, dest),
				Source:      sourcePath,
				Dest:        dest,
			},
			{
				Type:        ActionAppendScript,
				Description: fmt.Sprintf("Append create_link entry for %s to %s", name, scriptPath),
				Source:      name,
				Dest:        scriptPath,
			},
		},
	}
}

// PlanDriftRepair は drifted エントリの修復プランを作成する。
// 1. 実ディレクトリの内容を dotfiles リポにコピー（既存ファイルはスキップ）
// 2. 実ディレクトリを削除し、dotfiles リポへのシンボリックリンクを作成
func PlanDriftRepair(name, sourcePath, dotfilesDir, scriptPath string) *ImportPlan {
	dest := filepath.Join(dotfilesDir, "config", name)

	// ソースがファイルかディレクトリかで分岐
	info, err := os.Stat(sourcePath)
	copyAction := ActionCopyDir
	if err == nil && !info.IsDir() {
		copyAction = ActionCopyFile
	}

	return &ImportPlan{
		ItemName: name,
		ItemKind: "drift",
		Actions: []Action{
			{
				Type:        copyAction,
				Description: fmt.Sprintf("Copy %s to %s (merge, skip existing)", sourcePath, dest),
				Source:      sourcePath,
				Dest:        dest,
			},
			{
				Type:        ActionRepairDrift,
				Description: fmt.Sprintf("Remove real dir %s and create symlink to %s", sourcePath, dest),
				Source:      dest,
				Dest:        sourcePath,
			},
		},
	}
}

// PlanBrewImport は Brewfile へのパッケージ追加プランを作成する。
func PlanBrewImport(name, brewfilePath string) *ImportPlan {
	return &ImportPlan{
		ItemName: name,
		ItemKind: "brew",
		Actions: []Action{
			{
				Type:        ActionAppendBrewfile,
				Description: fmt.Sprintf("Append brew \"%s\" to %s", name, brewfilePath),
				Source:      name,
				Dest:        brewfilePath,
			},
		},
	}
}

// ExecuteImport はプラン内の全アクションを実行する。各アクションは独立して実行され、エラーは収集される。
func ExecuteImport(plan *ImportPlan) *ImportResult {
	result := &ImportResult{Plan: plan}

	for _, action := range plan.Actions {
		var err error
		switch action.Type {
		case ActionCopyDir:
			err = copyDir(action.Source, action.Dest)
		case ActionCopyFile:
			// 宛先の親ディレクトリを作成
			if mkErr := os.MkdirAll(filepath.Dir(action.Dest), 0755); mkErr != nil {
				err = fmt.Errorf("create parent dir: %w", mkErr)
			} else {
				err = copyFile(action.Source, action.Dest)
			}
		case ActionAppendScript:
			err = AppendLinkEntry(action.Dest, action.Source)
		case ActionAppendBrewfile:
			err = AppendBrewEntry(action.Dest, action.Source)
		case ActionRepairDrift:
			err = repairDrift(action.Source, action.Dest)
		default:
			err = fmt.Errorf("unknown action type: %d", action.Type)
		}
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("%s: %w", action.Description, err))
		}
	}

	return result
}

// copyDir はディレクトリを再帰的にコピーする。パーミッションを保持する。
// dst が既に存在する場合は新しいファイルのみをマージする（既存ファイルはスキップ）。
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("source stat: %w", err)
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory: %s", src)
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("create destination: %w", err)
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 宛先にシンボリックリンクが既にある場合はスキップ（管理済み）
		if linfo, lerr := os.Lstat(dstPath); lerr == nil {
			if linfo.Mode()&os.ModeSymlink != 0 {
				return nil // シンボリックリンクは触らない
			}
			// 通常ファイルが既にある場合もスキップ
			return nil
		}

		return copyFile(path, dstPath)
	})
}

// repairDrift は実ディレクトリを削除し、dotfiles リポへのシンボリックリンクを作成する。
func repairDrift(dotfilesDest, systemPath string) error {
	// dotfiles 側にコピー済みであることを確認
	if _, err := os.Stat(dotfilesDest); err != nil {
		return fmt.Errorf("dotfiles destination does not exist: %w", err)
	}

	// 実ディレクトリを削除
	if err := os.RemoveAll(systemPath); err != nil {
		return fmt.Errorf("remove real directory: %w", err)
	}

	// 親ディレクトリを作成（必要に応じて）
	if err := os.MkdirAll(filepath.Dir(systemPath), 0755); err != nil {
		return fmt.Errorf("create parent directory: %w", err)
	}

	// シンボリックリンクを作成
	if err := os.Symlink(dotfilesDest, systemPath); err != nil {
		return fmt.Errorf("create symlink: %w", err)
	}

	return nil
}

// copyFile は単一ファイルをコピーする。パーミッションを保持する。
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
