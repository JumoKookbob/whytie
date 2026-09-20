package syncer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/JumoKookbob/whytie/internal/history"
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

func TestSyncRecordsCreatedHistory(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatalf("MkdirAll(.whytie) error = %v", err)
	}

	store, err := storage.OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	sources := []scanner.SourceComment{
		{
			Kind:         syntax.Decision,
			Text:         "SQLite를 사용한다",
			RelativePath: "internal/storage/db.go",
			Line:         20,
		},
	}

	result, err := Sync(store, sources)
	if err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("Sync() returned %d memories, want 1", len(result))
	}

	events, err := store.ListHistory(result[0].ID)
	if err != nil {
		t.Fatalf("ListHistory() error = %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("history contains %d events, want 1", len(events))
	}

	event := events[0]

	if event.MemoryID != result[0].ID {
		t.Errorf("MemoryID = %q, want %q", event.MemoryID, result[0].ID)
	}

	if event.Type != "created" {
		t.Errorf("Type = %q, want %q", event.Type, "created")
	}
}

func TestSyncDoesNotDuplicateHistoryWhenUnchanged(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatalf("MkdirAll(.whytie) error = %v", err)
	}

	store, err := storage.OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	sources := []scanner.SourceComment{
		{
			Kind:         syntax.Decision,
			Text:         "SQLite를 사용한다",
			RelativePath: "internal/storage/db.go",
			Line:         20,
		},
	}

	first, err := Sync(store, sources)
	if err != nil {
		t.Fatalf("first Sync() error = %v", err)
	}

	if len(first) != 1 {
		t.Fatalf("first Sync() returned %d memories, want 1", len(first))
	}

	memoryID := first[0].ID

	for i := 0; i < 2; i++ {
		result, err := Sync(store, sources)
		if err != nil {
			t.Fatalf("repeated Sync() %d error = %v", i+1, err)
		}

		if len(result) != 1 {
			t.Fatalf(
				"repeated Sync() %d returned %d memories, want 1",
				i+1,
				len(result),
			)
		}

		if result[0].ID != memoryID {
			t.Fatalf(
				"repeated Sync() %d ID = %q, want %q",
				i+1,
				result[0].ID,
				memoryID,
			)
		}
	}

	events, err := store.ListHistory(memoryID)
	if err != nil {
		t.Fatalf("ListHistory() error = %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("history contains %d events, want exactly 1", len(events))
	}

	if events[0].Type != "created" {
		t.Errorf("event type = %q, want %q", events[0].Type, "created")
	}
}

func TestSyncRecordsMovedHistory(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatalf("MkdirAll(.whytie) error = %v", err)
	}

	store, err := storage.OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	firstSources := []scanner.SourceComment{
		{
			Kind:         syntax.Decision,
			Text:         "SQLite를 사용한다",
			RelativePath: "old/storage.go",
			Line:         20,
		},
	}

	first, err := Sync(store, firstSources)
	if err != nil {
		t.Fatalf("first Sync() error = %v", err)
	}

	if len(first) != 1 {
		t.Fatalf("first Sync() returned %d memories, want 1", len(first))
	}

	memoryID := first[0].ID

	secondSources := []scanner.SourceComment{
		{
			Kind:         syntax.Decision,
			Text:         "SQLite를 사용한다",
			RelativePath: "internal/storage/db.go",
			Line:         42,
		},
	}

	second, err := Sync(store, secondSources)
	if err != nil {
		t.Fatalf("second Sync() error = %v", err)
	}

	if len(second) != 1 {
		t.Fatalf("second Sync() returned %d memories, want 1", len(second))
	}

	if second[0].ID != memoryID {
		t.Fatalf(
			"ID after move = %q, want preserved ID %q",
			second[0].ID,
			memoryID,
		)
	}

	events, err := store.ListHistory(memoryID)
	if err != nil {
		t.Fatalf("ListHistory() error = %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("history contains %d events, want 2", len(events))
	}

	if events[0].Type != "created" {
		t.Errorf("first event type = %q, want %q", events[0].Type, "created")
	}

	if events[1].Type != "moved" {
		t.Errorf("second event type = %q, want %q", events[1].Type, "moved")
	}

	if events[1].Path != "internal/storage/db.go" {
		t.Errorf(
			"moved event path = %q, want %q",
			events[1].Path,
			"internal/storage/db.go",
		)
	}

	if events[1].Line != 42 {
		t.Errorf("moved event line = %d, want 42", events[1].Line)
	}
}

func TestSyncRecordsChangedHistory(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatalf("MkdirAll(.whytie) error = %v", err)
	}

	store, err := storage.OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	firstSources := []scanner.SourceComment{
		{
			Kind:         syntax.Decision,
			Text:         "SQLite를 사용한다",
			RelativePath: "internal/storage/db.go",
			Line:         20,
		},
	}

	first, err := Sync(store, firstSources)
	if err != nil {
		t.Fatalf("first Sync() error = %v", err)
	}

	if len(first) != 1 {
		t.Fatalf("first Sync() returned %d memories, want 1", len(first))
	}

	memoryID := first[0].ID

	secondSources := []scanner.SourceComment{
		{
			Kind:         syntax.Decision,
			Text:         "SQLite를 기본 저장소로 사용한다",
			RelativePath: "internal/storage/db.go",
			Line:         20,
		},
	}

	second, err := Sync(store, secondSources)
	if err != nil {
		t.Fatalf("second Sync() error = %v", err)
	}

	if len(second) != 1 {
		t.Fatalf("second Sync() returned %d memories, want 1", len(second))
	}

	if second[0].ID != memoryID {
		t.Fatalf(
			"ID after change = %q, want preserved ID %q",
			second[0].ID,
			memoryID,
		)
	}

	if second[0].Text != "SQLite를 기본 저장소로 사용한다" {
		t.Errorf(
			"Text after change = %q, want %q",
			second[0].Text,
			"SQLite를 기본 저장소로 사용한다",
		)
	}

	events, err := store.ListHistory(memoryID)
	if err != nil {
		t.Fatalf("ListHistory() error = %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("history contains %d events, want 2", len(events))
	}

	if events[0].Type != "created" {
		t.Errorf(
			"first event type = %q, want %q",
			events[0].Type,
			"created",
		)
	}

	if events[1].Type != "changed" {
		t.Errorf(
			"second event type = %q, want %q",
			events[1].Type,
			"changed",
		)
	}

	if events[1].Path != "internal/storage/db.go" {
		t.Errorf(
			"changed event path = %q, want %q",
			events[1].Path,
			"internal/storage/db.go",
		)
	}

	if events[1].Line != 20 {
		t.Errorf(
			"changed event line = %d, want 20",
			events[1].Line,
		)
	}
}

func TestSyncWithProvenanceRecordsCommitHash(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	store, err := storage.OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	sources := []scanner.SourceComment{
		{
			Kind:         syntax.Kind("decision"),
			Text:         "SQLite를 사용한다",
			RelativePath: "main.go",
			Line:         10,
		},
	}

	provenance := func(path string, line int) string {
		if path != "main.go" {
			t.Fatalf(
				"provenance path = %q, want %q",
				path,
				"main.go",
			)
		}

		if line != 10 {
			t.Fatalf(
				"provenance line = %d, want 10",
				line,
			)
		}

		return "abc123"
	}

	memories, err := SyncWithProvenance(
		store,
		sources,
		provenance,
	)
	if err != nil {
		t.Fatalf(
			"SyncWithProvenance() error = %v",
			err,
		)
	}

	if len(memories) != 1 {
		t.Fatalf(
			"SyncWithProvenance() returned %d memories, want 1",
			len(memories),
		)
	}

	events, err := store.ListHistory(memories[0].ID)
	if err != nil {
		t.Fatalf("ListHistory() error = %v", err)
	}

	if len(events) != 1 {
		t.Fatalf(
			"ListHistory() returned %d events, want 1",
			len(events),
		)
	}

	if events[0].CommitHash != "abc123" {
		t.Fatalf(
			"history CommitHash = %q, want %q",
			events[0].CommitHash,
			"abc123",
		)
	}
}

func TestSyncRecordsDeletedHistory(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	store, err := storage.OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite를 사용한다",
		File:         "example.go",
		RelativePath: "example.go",
		Line:         10,
	}

	memories, err := Sync(store, []scanner.SourceComment{source})
	if err != nil {
		t.Fatalf("first Sync() error = %v", err)
	}

	if len(memories) != 1 {
		t.Fatalf("first Sync() returned %d memories, want 1", len(memories))
	}

	memoryID := memories[0].ID

	if _, err := Sync(store, nil); err != nil {
		t.Fatalf("second Sync() error = %v", err)
	}

	events, err := store.ListHistory(memoryID)
	if err != nil {
		t.Fatalf("ListHistory() error = %v", err)
	}

	if len(events) != 2 {
		t.Fatalf("history contains %d events, want 2", len(events))
	}

	if events[0].Type != history.EventCreated {
		t.Errorf(
			"first event type = %q, want %q",
			events[0].Type,
			history.EventCreated,
		)
	}

	if events[1].Type != history.EventDeleted {
		t.Errorf(
			"second event type = %q, want %q",
			events[1].Type,
			history.EventDeleted,
		)
	}

	if events[1].Path != "example.go" {
		t.Errorf(
			"deleted event path = %q, want %q",
			events[1].Path,
			"example.go",
		)
	}

	if events[1].Line != 10 {
		t.Errorf(
			"deleted event line = %d, want %d",
			events[1].Line,
			10,
		)
	}
	if events[1].Kind != syntax.Decision {
		t.Errorf(
			"deleted event kind = %q, want %q",
			events[1].Kind,
			syntax.Decision,
		)
	}

	if events[1].Text != "SQLite를 사용한다" {
		t.Errorf(
			"deleted event text = %q, want %q",
			events[1].Text,
			"SQLite를 사용한다",
		)
	}

	active, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(active) != 0 {
		t.Errorf("active memories = %d, want 0", len(active))
	}
}
