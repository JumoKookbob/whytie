package repository

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestInitCreatesWhyTieDirectory(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	path := filepath.Join(root, ".whytie")

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat .whytie: %v", err)
	}

	if !info.IsDir() {
		t.Fatalf(".whytie is not a directory")
	}
}

func TestInitRejectsAlreadyInitializedRepository(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("first Init() error = %v", err)
	}

	err := Init(root)

	if !errors.Is(err, ErrAlreadyInitialized) {
		t.Fatalf(
			"second Init() error = %v, want %v",
			err,
			ErrAlreadyInitialized,
		)
	}
}

func TestFindRootFromNestedDirectory(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	nested := filepath.Join(root, "internal", "store", "sqlite")

	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("create nested directory: %v", err)
	}

	got, err := FindRoot(nested)
	if err != nil {
		t.Fatalf("FindRoot() error = %v", err)
	}

	if got != root {
		t.Errorf("FindRoot() = %q, want %q", got, root)
	}
}

func TestFindRootRejectsDirectoryOutsideRepository(t *testing.T) {
	root := t.TempDir()

	nested := filepath.Join(root, "some", "nested", "directory")

	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("create nested directory: %v", err)
	}

	_, err := FindRoot(nested)

	if !errors.Is(err, ErrNotRepository) {
		t.Fatalf(
			"FindRoot() error = %v, want %v",
			err,
			ErrNotRepository,
		)
	}
}

func TestInitCreatesInternalGitIgnore(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	path := filepath.Join(root, ".whytie", ".gitignore")

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}

	want := "*\n!.gitignore\n"

	if string(content) != want {
		t.Errorf(
			".whytie/.gitignore = %q, want %q",
			string(content),
			want,
		)
	}
}
