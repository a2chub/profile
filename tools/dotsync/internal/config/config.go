package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/a2chub/dotsync/internal/scanner"
)

// Config はアプリケーション全体の設定を表す
type Config struct {
	Dotfiles DotfilesConfig `toml:"dotfiles"`
	Scan     ScanConfig     `toml:"scan"`
	Secrets  SecretsConfig  `toml:"secrets"`
}

// DotfilesConfig は dotfiles リポジトリに関する設定
type DotfilesConfig struct {
	Repo       string `toml:"repo"`
	ConfigDir  string `toml:"config_dir"`
	LinkScript string `toml:"link_script"`
	Brewfile   string `toml:"brewfile"`
}

// ScanConfig はスキャン対象に関する設定
type ScanConfig struct {
	WatchDirs []string `toml:"watch_dirs"`
	Ignore    []string `toml:"ignore"`
}

// SecretsConfig はシークレット検出に関する設定
type SecretsConfig struct {
	Patterns         []string `toml:"patterns"`
	EntropyThreshold float64  `toml:"entropy_threshold"`
}

// DefaultConfig はデフォルト設定を返す
func DefaultConfig(dotfilesDir string) *Config {
	return &Config{
		Dotfiles: DotfilesConfig{
			Repo:       dotfilesDir,
			ConfigDir:  "config",
			LinkScript: "scripts/link-dotfiles.sh",
			Brewfile:   "packages/Brewfile",
		},
		Scan: ScanConfig{
			WatchDirs: []string{"~/.config"},
			Ignore:    []string{"*.cache", "*.log", "Code/", "google-chrome/", "Slack/", "configstore/", "gcloud/"},
		},
		Secrets: SecretsConfig{
			Patterns:         append([]string{}, scanner.DangerousPatterns...),
			EntropyThreshold: 4.5,
		},
	}
}

// LoadConfig は設定ファイルを検索し、デフォルト値とマージした Config を返す。
// 検索順序: {dotfilesDir}/.dotsync.toml -> ~/.config/dotsync/config.toml -> デフォルト値のみ
func LoadConfig(dotfilesDir string) (*Config, error) {
	cfg := DefaultConfig(dotfilesDir)

	// 設定ファイルの候補パスを構築
	candidates := []string{
		filepath.Join(dotfilesDir, ".dotsync.toml"),
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		candidates = append(candidates, filepath.Join(homeDir, ".config", "dotsync", "config.toml"))
	}

	// 最初に見つかったファイルを読み込む
	for _, path := range candidates {
		if _, err := os.Stat(path); err != nil {
			continue
		}

		if _, err := toml.DecodeFile(path, cfg); err != nil {
			return nil, err
		}

		// repo が空文字列ならデフォルトに戻す
		if cfg.Dotfiles.Repo == "" {
			cfg.Dotfiles.Repo = dotfilesDir
		}

		break
	}

	return cfg, nil
}

// ScriptPath は link-dotfiles.sh の絶対パスを返す
func (c *Config) ScriptPath() string {
	return filepath.Join(c.Dotfiles.Repo, c.Dotfiles.LinkScript)
}

// BrewfilePath は Brewfile の絶対パスを返す
func (c *Config) BrewfilePath() string {
	return filepath.Join(c.Dotfiles.Repo, c.Dotfiles.Brewfile)
}
