package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestScanFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "example.go")

	source := `package example

// normal comment
func example() {
	//+ SQLite
	//< local-first
}
`

	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	comments, err := ScanFile(path)
	if err != nil {
		t.Fatalf("ScanFile() error = %v", err)
	}

	if len(comments) != 2 {
		t.Fatalf("ScanFile() returned %d comments, want 2", len(comments))
	}

	tests := []struct {
		kind syntax.Kind
		text string
		line int
	}{
		{
			kind: syntax.Decision,
			text: "SQLite",
			line: 5,
		},
		{
			kind: syntax.Reason,
			text: "local-first",
			line: 6,
		},
	}

	for i, want := range tests {
		got := comments[i]

		if got.Kind != want.kind {
			t.Errorf("comment %d kind = %q, want %q", i, got.Kind, want.kind)
		}

		if got.Text != want.text {
			t.Errorf("comment %d text = %q, want %q", i, got.Text, want.text)
		}

		if got.File != path {
			t.Errorf("comment %d file = %q, want %q", i, got.File, path)
		}
		if got.Line != want.line {
			t.Errorf("comment %d line = %d, want %d", i, got.Line, want.line)
		}
	}
}

func TestScanDir(t *testing.T) {
	root := t.TempDir()

	subdir := filepath.Join(root, "internal", "storage")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("create subdir: %v", err)
	}

	mainSource := `package main

//+ SQLite 사용
// normal comment
func main() {}
`

	storeSource := `package storage

//? persistence를 어떻게 처리할까?
//< local-first가 필요함
`

	ignoredSource := `//! this must not be scanned`

	if err := os.WriteFile(
		filepath.Join(root, "main.go"),
		[]byte(mainSource),
		0644,
	); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(subdir, "store.go"),
		[]byte(storeSource),
		0644,
	); err != nil {
		t.Fatalf("write store.go: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(root, "notes.txt"),
		[]byte(ignoredSource),
		0644,
	); err != nil {
		t.Fatalf("write notes.txt: %v", err)
	}

	comments, err := ScanDir(root)
	if err != nil {
		t.Fatalf("ScanDir() error = %v", err)
	}

	if len(comments) != 3 {
		t.Fatalf("ScanDir() returned %d comments, want 3", len(comments))
	}

	got := map[string]SourceComment{}

	for _, comment := range comments {
		key := string(comment.Kind) + ":" + comment.Text
		got[key] = comment
	}

	want := []string{
		"decision:SQLite 사용",
		"question:persistence를 어떻게 처리할까?",
		"reason:local-first가 필요함",
	}

	for _, key := range want {
		if _, ok := got[key]; !ok {
			t.Errorf("ScanDir() missing %q", key)
		}
	}

	decision := got["decision:SQLite 사용"]
	if decision.RelativePath != "main.go" {
		t.Errorf(
			"decision RelativePath = %q, want %q",
			decision.RelativePath,
			"main.go",
		)
	}

	question := got["question:persistence를 어떻게 처리할까?"]
	wantStorePath := filepath.Join("internal", "storage", "store.go")

	if question.RelativePath != wantStorePath {
		t.Errorf(
			"question RelativePath = %q, want %q",
			question.RelativePath,
			wantStorePath,
		)
	}

	if _, ok := got["important:this must not be scanned"]; ok {
		t.Error("ScanDir() scanned WhyTie comment from non-Go file")
	}
}

func TestScanDirSkipsInternalDirectories(t *testing.T) {
	root := t.TempDir()

	gitDir := filepath.Join(root, ".git")
	whytieDir := filepath.Join(root, ".whytie")

	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("create .git: %v", err)
	}

	if err := os.MkdirAll(whytieDir, 0755); err != nil {
		t.Fatalf("create .whytie: %v", err)
	}

	mainSource := `package main

//+ SQLite
`

	gitSource := `package ignored

//! should not scan git
`

	whytieSource := `package ignored

//? should not scan whytie
`

	if err := os.WriteFile(
		filepath.Join(root, "main.go"),
		[]byte(mainSource),
		0644,
	); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(gitDir, "ignored.go"),
		[]byte(gitSource),
		0644,
	); err != nil {
		t.Fatalf("write .git/ignored.go: %v", err)
	}

	if err := os.WriteFile(
		filepath.Join(whytieDir, "ignored.go"),
		[]byte(whytieSource),
		0644,
	); err != nil {
		t.Fatalf("write .whytie/ignored.go: %v", err)
	}

	comments, err := ScanDir(root)
	if err != nil {
		t.Fatalf("ScanDir() error = %v", err)
	}

	if len(comments) != 1 {
		t.Fatalf("ScanDir() returned %d comments, want 1", len(comments))
	}

	got := comments[0]

	if got.Kind != syntax.Decision {
		t.Errorf("kind = %q, want %q", got.Kind, syntax.Decision)
	}

	if got.Text != "SQLite" {
		t.Errorf("text = %q, want %q", got.Text, "SQLite")
	}
}
