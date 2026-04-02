package scanner

import (
	"math"
	"os"
	"path/filepath"
)

// SecretRisk はシークレット検出のリスクレベルを表す
type SecretRisk int

const (
	RiskNone   SecretRisk = iota
	RiskLow
	RiskMedium
	RiskHigh
)

// SecretWarning はシークレット検出の警告を表す
type SecretWarning struct {
	FilePath string
	Risk     SecretRisk
	Reason   string
	Entropy  float64
}

// DangerousPatterns は危険なファイル名のグロブパターン一覧
var DangerousPatterns = []string{
	"*.pem",
	"*.key",
	"*.p12",
	"*.pfx",
	"*.keystore",
	".env",
	".env.*",
	"credentials.json",
	"secrets.json",
	"token.json",
	".netrc",
	".npmrc",
	"id_rsa",
	"id_ed25519",
	"id_ecdsa",
	"*.secret",
	"*.secrets",
}

// エントロピー閾値
const entropyThreshold = 4.5

// ScanForSecrets はディレクトリを走査し、シークレットの可能性があるファイルを検出する。
// threshold <= 0 の場合はデフォルト閾値（4.5）を使用する。
func ScanForSecrets(dirPath string, threshold float64) ([]SecretWarning, error) {
	if threshold <= 0 {
		threshold = entropyThreshold
	}

	var warnings []SecretWarning

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// シンボリックリンクをスキップ
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		// ディレクトリはスキップ
		if info.IsDir() {
			return nil
		}

		// 空ファイル、16バイト未満、1MB超をスキップ
		size := info.Size()
		if size == 0 || size < 16 || size > 1024*1024 {
			return nil
		}

		filename := filepath.Base(path)
		patternMatch := matchesDangerousPattern(filename)

		// ファイル内容を読み取りバイナリチェックとエントロピー計算
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if IsBinary(data) {
			// バイナリファイルでもパターンマッチがあれば警告
			if patternMatch {
				warnings = append(warnings, SecretWarning{
					FilePath: path,
					Risk:     RiskMedium,
					Reason:   "ファイル名が危険なパターンに一致（バイナリファイル）",
					Entropy:  0.0,
				})
			}
			return nil
		}

		entropy := ShannonEntropy(data)
		highEntropy := entropy > threshold

		switch {
		case patternMatch && highEntropy:
			warnings = append(warnings, SecretWarning{
				FilePath: path,
				Risk:     RiskHigh,
				Reason:   "ファイル名が危険なパターンに一致し、高エントロピー",
				Entropy:  entropy,
			})
		case patternMatch:
			warnings = append(warnings, SecretWarning{
				FilePath: path,
				Risk:     RiskMedium,
				Reason:   "ファイル名が危険なパターンに一致",
				Entropy:  entropy,
			})
		case highEntropy:
			warnings = append(warnings, SecretWarning{
				FilePath: path,
				Risk:     RiskMedium,
				Reason:   "高エントロピー（シークレットの可能性）",
				Entropy:  entropy,
			})
		}

		return nil
	})

	return warnings, err
}

// ShannonEntropy はバイト列のシャノンエントロピーをビット単位で計算する（0.0〜8.0）
func ShannonEntropy(data []byte) float64 {
	if len(data) == 0 {
		return 0.0
	}

	var freq [256]int
	for _, b := range data {
		freq[b]++
	}

	total := float64(len(data))
	entropy := 0.0

	for _, count := range freq {
		if count == 0 {
			continue
		}
		p := float64(count) / total
		entropy -= p * math.Log2(p)
	}

	return entropy
}

// IsBinary は先頭8KBにヌルバイトが含まれるかを判定する
func IsBinary(data []byte) bool {
	limit := 8192
	if len(data) < limit {
		limit = len(data)
	}
	for i := 0; i < limit; i++ {
		if data[i] == 0 {
			return true
		}
	}
	return false
}

// matchesDangerousPattern はファイル名が危険なパターンに一致するか判定する
func matchesDangerousPattern(filename string) bool {
	for _, pattern := range DangerousPatterns {
		matched, err := filepath.Match(pattern, filename)
		if err == nil && matched {
			return true
		}
	}
	return false
}
