package scanner

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func testdataPath(name string) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "testdata", name)
}

func TestParseBrewfile(t *testing.T) {
	path := testdataPath("Brewfile")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("テスト用Brewfileが見つかりません: %s", path)
	}

	entries, err := ParseBrewfile(path)
	if err != nil {
		t.Fatalf("ParseBrewfileがエラーを返しました: %v", err)
	}

	// テスト用Brewfileには: tap x2, brew x10, cask x3 = 15エントリ
	if len(entries) != 15 {
		t.Errorf("エントリ数が期待値と異なります: got %d, want 15", len(entries))
		for _, e := range entries {
			t.Logf("  %s %q (line %d)", e.Kind, e.Name, e.Line)
		}
	}

	// 各種類のエントリが正しくパースされることを確認
	kindCounts := map[string]int{}
	for _, e := range entries {
		kindCounts[e.Kind]++
	}

	if kindCounts["tap"] != 2 {
		t.Errorf("tapエントリ数: got %d, want 2", kindCounts["tap"])
	}
	if kindCounts["brew"] != 10 {
		t.Errorf("brewエントリ数: got %d, want 10", kindCounts["brew"])
	}
	if kindCounts["cask"] != 3 {
		t.Errorf("caskエントリ数: got %d, want 3", kindCounts["cask"])
	}

	// 特定のエントリの内容を確認
	found := false
	for _, e := range entries {
		if e.Name == "felixkratz/formulae/borders" && e.Kind == "brew" {
			found = true
			break
		}
	}
	if !found {
		t.Error("tap付きbrewエントリ 'felixkratz/formulae/borders' が見つかりません")
	}

	// 行番号が正の値であることを確認
	for _, e := range entries {
		if e.Line <= 0 {
			t.Errorf("エントリ %q の行番号が不正です: %d", e.Name, e.Line)
		}
	}
}

func TestParseBrewfile_NotFound(t *testing.T) {
	_, err := ParseBrewfile("/nonexistent/path/Brewfile")
	if err == nil {
		t.Error("存在しないファイルに対してエラーが返されませんでした")
	}
}

func TestParseLeaves(t *testing.T) {
	output := "curl\nwget\nfzf\ngh\nhtop\n"
	result := ParseLeaves(output)

	if len(result) != 5 {
		t.Fatalf("パッケージ数: got %d, want 5", len(result))
	}

	expected := []string{"curl", "wget", "fzf", "gh", "htop"}
	for i, pkg := range expected {
		if result[i] != pkg {
			t.Errorf("result[%d]: got %q, want %q", i, result[i], pkg)
		}
	}
}

func TestParseLeaves_EmptyOutput(t *testing.T) {
	result := ParseLeaves("")
	if len(result) != 0 {
		t.Errorf("空の出力に対してパッケージが返されました: %v", result)
	}
}

func TestParseLeaves_TrailingNewlines(t *testing.T) {
	output := "curl\nwget\n\n\n"
	result := ParseLeaves(output)
	if len(result) != 2 {
		t.Errorf("パッケージ数: got %d, want 2", len(result))
	}
}

func TestDiffBrew_MissingPackages(t *testing.T) {
	brewfile := []BrewEntry{
		{Name: "curl", Kind: "brew", Line: 1},
		{Name: "wget", Kind: "brew", Line: 2},
	}
	leaves := []string{"curl", "wget", "ripgrep"}

	diff := DiffBrew(brewfile, leaves)

	if len(diff.Missing) != 1 {
		t.Fatalf("Missingの数: got %d, want 1", len(diff.Missing))
	}
	if diff.Missing[0].Name != "ripgrep" {
		t.Errorf("Missing[0].Name: got %q, want %q", diff.Missing[0].Name, "ripgrep")
	}

	if len(diff.Extra) != 0 {
		t.Errorf("Extraの数: got %d, want 0", len(diff.Extra))
	}
}

func TestDiffBrew_ExtraPackages(t *testing.T) {
	brewfile := []BrewEntry{
		{Name: "curl", Kind: "brew", Line: 1},
		{Name: "wget", Kind: "brew", Line: 2},
		{Name: "ripgrep", Kind: "brew", Line: 3},
	}
	leaves := []string{"curl", "wget"}

	diff := DiffBrew(brewfile, leaves)

	if len(diff.Missing) != 0 {
		t.Errorf("Missingの数: got %d, want 0", len(diff.Missing))
	}
	if len(diff.Extra) != 1 {
		t.Fatalf("Extraの数: got %d, want 1", len(diff.Extra))
	}
	if diff.Extra[0].Name != "ripgrep" {
		t.Errorf("Extra[0].Name: got %q, want %q", diff.Extra[0].Name, "ripgrep")
	}
}

func TestDiffBrew_TapNormalization(t *testing.T) {
	brewfile := []BrewEntry{
		{Name: "felixkratz/formulae/borders", Kind: "brew", Line: 1},
	}
	leaves := []string{"felixkratz/formulae/borders"}

	diff := DiffBrew(brewfile, leaves)

	if len(diff.Missing) != 0 {
		t.Errorf("Missingの数: got %d, want 0 (正規化マッチが機能していません)", len(diff.Missing))
	}
	if len(diff.Extra) != 0 {
		t.Errorf("Extraの数: got %d, want 0", len(diff.Extra))
	}
}

func TestDiffBrew_TapNormalization_ShortInLeaves(t *testing.T) {
	// Brewfileにはフルパス、leavesには短い名前のケース
	brewfile := []BrewEntry{
		{Name: "felixkratz/formulae/borders", Kind: "brew", Line: 1},
	}
	leaves := []string{"borders"}

	diff := DiffBrew(brewfile, leaves)

	if len(diff.Missing) != 0 {
		t.Errorf("Missingの数: got %d, want 0 (短い名前での正規化マッチが機能していません)", len(diff.Missing))
	}
	if len(diff.Extra) != 0 {
		t.Errorf("Extraの数: got %d, want 0", len(diff.Extra))
	}
}

func TestDiffBrew_IgnoresCasks(t *testing.T) {
	brewfile := []BrewEntry{
		{Name: "curl", Kind: "brew", Line: 1},
		{Name: "ghostty", Kind: "cask", Line: 2},
	}
	leaves := []string{"curl"}

	diff := DiffBrew(brewfile, leaves)

	// caskエントリはExtraに含まれないことを確認
	if len(diff.Extra) != 0 {
		t.Errorf("Extraの数: got %d, want 0 (caskがExtraに含まれています)", len(diff.Extra))
	}
}

func TestDiffBrew_NoDiff(t *testing.T) {
	brewfile := []BrewEntry{
		{Name: "curl", Kind: "brew", Line: 1},
		{Name: "wget", Kind: "brew", Line: 2},
	}
	leaves := []string{"curl", "wget"}

	diff := DiffBrew(brewfile, leaves)

	if len(diff.Missing) != 0 {
		t.Errorf("Missingの数: got %d, want 0", len(diff.Missing))
	}
	if len(diff.Extra) != 0 {
		t.Errorf("Extraの数: got %d, want 0", len(diff.Extra))
	}
}
