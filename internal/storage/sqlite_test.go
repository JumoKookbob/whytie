package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/repository"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestOpenSQLiteCreatesDatabaseFile(t *testing.T) {
	root := t.TempDir()

	if err := repository.Init(root); err != nil {
		t.Fatalf("repository.Init() error = %v", err)
	}

	store, err := OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	dbPath := filepath.Join(root, ".whytie", "whytie.db")

	info, err := os.Stat(dbPath)
	if err != nil {
		t.Fatalf("database file not created: %v", err)
	}

	if info.IsDir() {
		t.Fatalf("%q is a directory, want database file", dbPath)
	}
}

func TestOpenSQLiteCreatesMemoriesTable(t *testing.T) {
	root := t.TempDir()

	if err := repository.Init(root); err != nil {
		t.Fatalf("repository.Init() error = %v", err)
	}

	store, err := OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	var name string

	err = store.db.QueryRow(`
		SELECT name
		FROM sqlite_master
		WHERE type = 'table'
		  AND name = 'memories'
	`).Scan(&name)

	if err != nil {
		t.Fatalf("memories table not found: %v", err)
	}

	if name != "memories" {
		t.Fatalf("table name = %q, want %q", name, "memories")
	}
}

func TestSQLiteStoreSaveMemory(t *testing.T) {
	root := t.TempDir()

	if err := repository.Init(root); err != nil {
		t.Fatalf("repository.Init() error = %v", err)
	}

	store, err := OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	input := memory.Memory{
		ID:          "memory-1",
		Kind:        syntax.Decision,
		Text:        "SQLite",
		CreatedPath: "internal/store/db.go",
		CreatedLine: 21,
		CurrentPath: "internal/store/db.go",
		CurrentLine: 21,
	}

	if err := store.Save(input); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	var (
		id          string
		kind        string
		text        string
		createdPath string
		createdLine int
		currentPath string
		currentLine int
	)

	err = store.db.QueryRow(`
		SELECT
			id,
			kind,
			text,
			created_path,
			created_line,
			current_path,
			current_line
		FROM memories
		WHERE id = ?
	`, input.ID).Scan(
		&id,
		&kind,
		&text,
		&createdPath,
		&createdLine,
		&currentPath,
		&currentLine,
	)
	if err != nil {
		t.Fatalf("query saved memory: %v", err)
	}

	if id != input.ID {
		t.Errorf("id = %q, want %q", id, input.ID)
	}

	if kind != string(input.Kind) {
		t.Errorf("kind = %q, want %q", kind, input.Kind)
	}

	if text != input.Text {
		t.Errorf("text = %q, want %q", text, input.Text)
	}

	if createdPath != input.CreatedPath || createdLine != input.CreatedLine {
		t.Errorf(
			"created location = %s:%d, want %s:%d",
			createdPath,
			createdLine,
			input.CreatedPath,
			input.CreatedLine,
		)
	}

	if currentPath != input.CurrentPath || currentLine != input.CurrentLine {
		t.Errorf(
			"current location = %s:%d, want %s:%d",
			currentPath,
			currentLine,
			input.CurrentPath,
			input.CurrentLine,
		)
	}
}

func TestSQLiteStoreGetMemoryAfterReopen(t *testing.T) {
	root := t.TempDir()

	if err := repository.Init(root); err != nil {
		t.Fatalf("repository.Init() error = %v", err)
	}

	store, err := OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}

	input := memory.Memory{
		ID:          "memory-1",
		Kind:        syntax.Decision,
		Text:        "SQLite",
		CreatedPath: "internal/store/db.go",
		CreatedLine: 21,
		CurrentPath: "internal/store/db.go",
		CurrentLine: 21,
	}

	if err := store.Save(input); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	store, err = OpenSQLite(root)
	if err != nil {
		t.Fatalf("reopen OpenSQLite() error = %v", err)
	}
	defer store.Close()

	got, err := store.Get(input.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got != input {
		t.Errorf("Get() = %#v, want %#v", got, input)
	}
}

func TestSQLiteStoreListMemories(t *testing.T) {
	root := t.TempDir()

	if err := repository.Init(root); err != nil {
		t.Fatalf("repository.Init() error = %v", err)
	}

	store, err := OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	first := memory.Memory{
		ID:          "memory-1",
		Kind:        syntax.Question,
		Text:        "어떤 DB를 쓸까?",
		CreatedPath: "internal/store/db.go",
		CreatedLine: 20,
		CurrentPath: "internal/store/db.go",
		CurrentLine: 20,
	}

	second := memory.Memory{
		ID:          "memory-2",
		Kind:        syntax.Decision,
		Text:        "SQLite",
		CreatedPath: "internal/store/db.go",
		CreatedLine: 21,
		CurrentPath: "internal/store/db.go",
		CurrentLine: 21,
	}

	if err := store.Save(first); err != nil {
		t.Fatalf("Save(first) error = %v", err)
	}

	if err := store.Save(second); err != nil {
		t.Fatalf("Save(second) error = %v", err)
	}

	got, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("List() returned %d memories, want 2", len(got))
	}

	if got[0] != first {
		t.Errorf("got[0] = %#v, want %#v", got[0], first)
	}

	if got[1] != second {
		t.Errorf("got[1] = %#v, want %#v", got[1], second)
	}
}

func TestSQLiteStoreSaveUpdatesExistingMemory(t *testing.T) {
	root := t.TempDir()

	whytieDir := filepath.Join(root, ".whytie")
	if err := os.MkdirAll(whytieDir, 0755); err != nil {
		t.Fatalf("create .whytie: %v", err)
	}

	store, err := OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	original := memory.Memory{
		ID:          "memory-1",
		Kind:        syntax.Decision,
		Text:        "SQLite",
		CreatedPath: "internal/store/db.go",
		CreatedLine: 21,
		CurrentPath: "internal/store/db.go",
		CurrentLine: 21,
	}

	if err := store.Save(original); err != nil {
		t.Fatalf("first Save() error = %v", err)
	}

	moved := original
	moved.CurrentLine = 40

	if err := store.Save(moved); err != nil {
		t.Fatalf("second Save() error = %v", err)
	}

	got, err := store.Get("memory-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if got.CreatedLine != 21 {
		t.Errorf("CreatedLine = %d, want 21", got.CreatedLine)
	}

	if got.CurrentLine != 40 {
		t.Errorf("CurrentLine = %d, want 40", got.CurrentLine)
	}
}
