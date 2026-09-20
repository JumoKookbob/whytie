package syncer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/storage"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

type fakeStore struct {
	memories []memory.Memory
}

func (s *fakeStore) Save(m memory.Memory) error {
	for i := range s.memories {
		if s.memories[i].ID == m.ID {
			s.memories[i] = m
			return nil
		}
	}

	s.memories = append(s.memories, m)
	return nil
}

func (s *fakeStore) Delete(id string) error {
	for i, m := range s.memories {
		if m.ID != id {
			continue
		}

		s.memories = append(s.memories[:i], s.memories[i+1:]...)
		return nil
	}

	return nil
}

func (s *fakeStore) Get(id string) (memory.Memory, error) {
	for _, m := range s.memories {
		if m.ID == id {
			return m, nil
		}
	}

	return memory.Memory{}, nil
}

func (s *fakeStore) List() ([]memory.Memory, error) {
	result := make([]memory.Memory, len(s.memories))
	copy(result, s.memories)
	return result, nil
}

func TestSyncReconcilesAndSavesMemories(t *testing.T) {
	store := &fakeStore{
		memories: []memory.Memory{
			{
				ID:          "memory-1",
				Kind:        syntax.Decision,
				Text:        "SQLite",
				CreatedPath: "internal/store/db.go",
				CreatedLine: 21,
				CurrentPath: "internal/store/db.go",
				CurrentLine: 21,
			},
		},
	}

	sources := []scanner.SourceComment{
		{
			Kind:         syntax.Decision,
			Text:         "SQLite",
			RelativePath: "internal/store/db.go",
			Line:         40,
		},
		{
			Kind:         syntax.Reason,
			Text:         "local-first에 적합함",
			RelativePath: "internal/store/db.go",
			Line:         41,
		},
	}

	got, err := Sync(store, sources)
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("Sync() returned %d memories, want 2", len(got))
	}

	if got[0].ID != "memory-1" {
		t.Errorf("moved memory ID = %q, want memory-1", got[0].ID)
	}

	if got[0].CreatedLine != 21 {
		t.Errorf("CreatedLine = %d, want 21", got[0].CreatedLine)
	}

	if got[0].CurrentLine != 40 {
		t.Errorf("CurrentLine = %d, want 40", got[0].CurrentLine)
	}

	if got[1].ID == "" {
		t.Error("new memory has empty ID")
	}

	saved, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(saved) != 2 {
		t.Fatalf("store contains %d memories, want 2", len(saved))
	}

	if saved[0].CurrentLine != 40 {
		t.Errorf("saved CurrentLine = %d, want 40", saved[0].CurrentLine)
	}
}

func TestSyncPersistsToSQLiteAcrossReopen(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatalf("create .whytie: %v", err)
	}

	store, err := storage.OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}

	sources := []scanner.SourceComment{
		{
			Kind:         syntax.Decision,
			Text:         "SQLite",
			RelativePath: "internal/store/db.go",
			Line:         21,
		},
	}

	first, err := Sync(store, sources)
	if err != nil {
		t.Fatalf("first Sync() error = %v", err)
	}

	if len(first) != 1 {
		t.Fatalf("first Sync() returned %d memories, want 1", len(first))
	}

	id := first[0].ID

	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	store, err = storage.OpenSQLite(root)
	if err != nil {
		t.Fatalf("reopen SQLite error = %v", err)
	}
	defer store.Close()

	movedSources := []scanner.SourceComment{
		{
			Kind:         syntax.Decision,
			Text:         "SQLite",
			RelativePath: "internal/store/db.go",
			Line:         40,
		},
		{
			Kind:         syntax.Reason,
			Text:         "local-first에 적합함",
			RelativePath: "internal/store/db.go",
			Line:         41,
		},
	}

	second, err := Sync(store, movedSources)
	if err != nil {
		t.Fatalf("second Sync() error = %v", err)
	}

	if len(second) != 2 {
		t.Fatalf("second Sync() returned %d memories, want 2", len(second))
	}

	if second[0].ID != id {
		t.Errorf("moved memory ID = %q, want %q", second[0].ID, id)
	}

	if second[0].CreatedLine != 21 {
		t.Errorf("CreatedLine = %d, want 21", second[0].CreatedLine)
	}

	if second[0].CurrentLine != 40 {
		t.Errorf("CurrentLine = %d, want 40", second[0].CurrentLine)
	}

	saved, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(saved) != 2 {
		t.Fatalf("SQLite contains %d memories, want 2", len(saved))
	}

	if saved[0].ID != id {
		t.Errorf("saved ID = %q, want %q", saved[0].ID, id)
	}

	if saved[0].CreatedLine != 21 {
		t.Errorf("saved CreatedLine = %d, want 21", saved[0].CreatedLine)
	}

	if saved[0].CurrentLine != 40 {
		t.Errorf("saved CurrentLine = %d, want 40", saved[0].CurrentLine)
	}
}

func TestSyncRemovesDeletedMemoryFromStore(t *testing.T) {
	store := &fakeStore{
		memories: []memory.Memory{
			{
				ID:          "memory-1",
				Kind:        syntax.Decision,
				Text:        "SQLite를 사용한다",
				CreatedPath: "example.go",
				CreatedLine: 10,
				CurrentPath: "example.go",
				CurrentLine: 10,
			},
		},
	}

	sources := []scanner.SourceComment{}

	_, err := Sync(store, sources)
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	saved, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(saved) != 0 {
		t.Fatalf("store contains %d memories, want 0", len(saved))
	}
}

func TestSyncPreservesIdentityAcrossFileMove(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatalf("MkdirAll(.whytie) error = %v", err)
	}

	store, err := storage.OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	first := scanner.SourceComment{
		Kind:         syntax.Kind("decision"),
		Text:         "SQLite를 사용한다",
		File:         filepath.Join(root, "main.go"),
		RelativePath: "main.go",
		Line:         10,
	}

	firstResult, err := Sync(store, []scanner.SourceComment{first})
	if err != nil {
		t.Fatalf("first Sync() error = %v", err)
	}

	if len(firstResult) != 1 {
		t.Fatalf("first Sync() returned %d memories, want 1", len(firstResult))
	}

	originalID := firstResult[0].ID

	moved := scanner.SourceComment{
		Kind:         syntax.Kind("decision"),
		Text:         "SQLite를 사용한다",
		File:         filepath.Join(root, "internal", "storage", "db.go"),
		RelativePath: filepath.Join("internal", "storage", "db.go"),
		Line:         42,
	}

	secondResult, err := Sync(store, []scanner.SourceComment{moved})
	if err != nil {
		t.Fatalf("second Sync() error = %v", err)
	}

	if len(secondResult) != 1 {
		t.Fatalf("second Sync() returned %d memories, want 1", len(secondResult))
	}

	got := secondResult[0]

	if got.ID != originalID {
		t.Fatalf("ID = %q, want original ID %q", got.ID, originalID)
	}

	if got.CreatedPath != "main.go" {
		t.Fatalf("CreatedPath = %q, want %q", got.CreatedPath, "main.go")
	}

	if got.CreatedLine != 10 {
		t.Fatalf("CreatedLine = %d, want 10", got.CreatedLine)
	}

	wantCurrentPath := filepath.Join("internal", "storage", "db.go")

	if got.CurrentPath != wantCurrentPath {
		t.Fatalf("CurrentPath = %q, want %q", got.CurrentPath, wantCurrentPath)
	}

	if got.CurrentLine != 42 {
		t.Fatalf("CurrentLine = %d, want 42", got.CurrentLine)
	}

	stored, err := store.Get(originalID)
	if err != nil {
		t.Fatalf("store.Get() error = %v", err)
	}

	if stored.ID != originalID {
		t.Fatalf("stored ID = %q, want %q", stored.ID, originalID)
	}

	if stored.CreatedPath != "main.go" {
		t.Fatalf("stored CreatedPath = %q, want %q", stored.CreatedPath, "main.go")
	}

	if stored.CurrentPath != wantCurrentPath {
		t.Fatalf(
			"stored CurrentPath = %q, want %q",
			stored.CurrentPath,
			wantCurrentPath,
		)
	}

	if stored.CurrentLine != 42 {
		t.Fatalf("stored CurrentLine = %d, want 42", stored.CurrentLine)
	}
}
