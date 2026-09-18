package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesWhyTieRepository(t *testing.T) {
	dir := t.TempDir()

	if err := Init(dir); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	path := filepath.Join(dir, ".whytie")

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat(%q) error = %v", path, err)
	}

	if !info.IsDir() {
		t.Fatalf("%q is not a directory", path)
	}
}
