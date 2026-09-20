package scanner_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestScanFileGoRawStringClosesAndReopensOnSameLine(t *testing.T) {
	source := strings.Join([]string{
		"package sample",
		"var sample = `first",
		"` + \"`\" + `",
		"// ? fake question inside string",
		"`",
		"// + real decision outside string",
		"",
	}, "\n")

	path := filepath.Join(t.TempDir(), "sample.go")
	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}

	comments, err := scanner.ScanFile(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(comments) != 1 {
		t.Fatalf("expected one real comment, got %#v", comments)
	}

	got := comments[0]
	if got.Kind != syntax.Decision ||
		got.Text != "real decision outside string" ||
		got.Line != 6 {
		t.Fatalf("wrong comment detected: %#v", got)
	}
}
