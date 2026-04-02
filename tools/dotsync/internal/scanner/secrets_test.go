package scanner

import (
	"encoding/base64"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestShannonEntropy_AllZero(t *testing.T) {
	// 全て同じバイト値 → エントロピーは0.0
	data := make([]byte, 256)
	for i := range data {
		data[i] = 0
	}
	entropy := ShannonEntropy(data)
	if entropy != 0.0 {
		t.Errorf("全て同じバイトのエントロピーは0.0であるべき、結果: %f", entropy)
	}
}

func TestShannonEntropy_Uniform(t *testing.T) {
	// 256種類のバイト値が各1回 → エントロピーは8.0
	data := make([]byte, 256)
	for i := 0; i < 256; i++ {
		data[i] = byte(i)
	}
	entropy := ShannonEntropy(data)
	if math.Abs(entropy-8.0) > 0.001 {
		t.Errorf("均一分布のエントロピーは8.0であるべき、結果: %f", entropy)
	}
}

func TestShannonEntropy_EnglishText(t *testing.T) {
	// 典型的な英語テキスト → 3.5〜4.5の範囲
	data := []byte("The quick brown fox jumps over the lazy dog. This is a sample of typical English text that should have moderate entropy values when measured using Shannon entropy calculation.")
	entropy := ShannonEntropy(data)
	if entropy < 3.5 || entropy > 4.5 {
		t.Errorf("英語テキストのエントロピーは3.5〜4.5であるべき、結果: %f", entropy)
	}
}

func TestShannonEntropy_Base64(t *testing.T) {
	// Base64エンコードされたランダムデータ → 5.0以上
	// 擬似ランダムなバイト列を生成してBase64エンコード
	raw := make([]byte, 256)
	for i := range raw {
		raw[i] = byte((i*37 + 113) % 256)
	}
	encoded := []byte(base64.StdEncoding.EncodeToString(raw))
	entropy := ShannonEntropy(encoded)
	if entropy < 5.0 {
		t.Errorf("Base64エンコードデータのエントロピーは5.0以上であるべき、結果: %f", entropy)
	}
}

func TestScanForSecrets_FilenameMatch(t *testing.T) {
	// credentials.jsonという名前のファイルを検出
	dir := t.TempDir()
	filePath := filepath.Join(dir, "credentials.json")
	// 16バイト以上の内容を書き込む
	content := []byte(`{"client_id": "test-id", "client_secret": "test-secret-value"}`)
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("テストファイル作成失敗: %v", err)
	}

	warnings, err := ScanForSecrets(dir, 0)
	if err != nil {
		t.Fatalf("ScanForSecrets エラー: %v", err)
	}

	if len(warnings) == 0 {
		t.Fatal("credentials.json が検出されるべき")
	}

	found := false
	for _, w := range warnings {
		if filepath.Base(w.FilePath) == "credentials.json" {
			found = true
			if w.Risk < RiskMedium {
				t.Errorf("リスクレベルはMedium以上であるべき、結果: %d", w.Risk)
			}
		}
	}
	if !found {
		t.Error("credentials.json の警告が見つからない")
	}
}

func TestScanForSecrets_HighEntropy(t *testing.T) {
	// 高エントロピーのファイルを検出
	dir := t.TempDir()
	filePath := filepath.Join(dir, "data.txt")
	// 高エントロピーデータを生成
	raw := make([]byte, 512)
	for i := range raw {
		raw[i] = byte((i*37 + 113) % 256)
	}
	highEntropyData := []byte(base64.StdEncoding.EncodeToString(raw))
	if err := os.WriteFile(filePath, highEntropyData, 0644); err != nil {
		t.Fatalf("テストファイル作成失敗: %v", err)
	}

	warnings, err := ScanForSecrets(dir, 0)
	if err != nil {
		t.Fatalf("ScanForSecrets エラー: %v", err)
	}

	if len(warnings) == 0 {
		t.Fatal("高エントロピーファイルが検出されるべき")
	}

	found := false
	for _, w := range warnings {
		if filepath.Base(w.FilePath) == "data.txt" {
			found = true
			if w.Entropy <= entropyThreshold {
				t.Errorf("エントロピーは閾値を超えるべき、結果: %f", w.Entropy)
			}
		}
	}
	if !found {
		t.Error("高エントロピーファイルの警告が見つからない")
	}
}

func TestIsBinary(t *testing.T) {
	// ヌルバイトを含む → true
	binaryData := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x00, 0x57, 0x6f, 0x72, 0x6c, 0x64}
	if !IsBinary(binaryData) {
		t.Error("ヌルバイトを含むデータはバイナリと判定されるべき")
	}

	// クリーンなテキスト → false
	textData := []byte("Hello, World! This is clean text without null bytes.")
	if IsBinary(textData) {
		t.Error("クリーンなテキストはバイナリと判定されるべきではない")
	}
}
