package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/JumoKookbob/whytie/internal/history"
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

func TestSQLiteStoreDeleteMemory(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatalf("create .whytie: %v", err)
	}

	store, err := OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	m := memory.Memory{
		ID:          "memory-delete",
		Kind:        syntax.Decision,
		Text:        "SQLite",
		CreatedPath: "example.go",
		CreatedLine: 10,
		CurrentPath: "example.go",
		CurrentLine: 10,
	}

	if err := store.Save(m); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if err := store.Delete(m.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	memories, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(memories) != 0 {
		t.Fatalf("List() returned %d memories after Delete(), want 0", len(memories))
	}
}

func TestSQLiteStoreListOrdersByCurrentSourceLocation(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatalf("create .whytie: %v", err)
	}

	store, err := OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	memories := []memory.Memory{
		{
			ID:          "reason",
			Kind:        syntax.Reason,
			Text:        "local-first",
			CreatedPath: "example.go",
			CreatedLine: 12,
			CurrentPath: "example.go",
			CurrentLine: 12,
		},
		{
			ID:          "question",
			Kind:        syntax.Question,
			Text:        "어떤 DB를 쓸까?",
			CreatedPath: "example.go",
			CreatedLine: 10,
			CurrentPath: "example.go",
			CurrentLine: 10,
		},
		{
			ID:          "decision",
			Kind:        syntax.Decision,
			Text:        "SQLite",
			CreatedPath: "example.go",
			CreatedLine: 11,
			CurrentPath: "example.go",
			CurrentLine: 11,
		},
	}

	for _, m := range memories {
		if err := store.Save(m); err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	}

	got, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("List() returned %d memories, want 3", len(got))
	}

	wantIDs := []string{"question", "decision", "reason"}

	for i, want := range wantIDs {
		if got[i].ID != want {
			t.Errorf("got[%d].ID = %q, want %q", i, got[i].ID, want)
		}
	}
}

func TestSQLiteHistoryPersistsAcrossReopen(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	store, err := OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}

	baseTime := time.Date(
		2026,
		time.September,
		20,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	events := []history.Event{
		{
			MemoryID:   "memory-123",
			Type:       history.EventCreated,
			Path:       "main.go",
			Line:       10,
			CommitHash: "commit-a",
			OccurredAt: baseTime,
		},
		{
			MemoryID:   "memory-123",
			Type:       history.EventMoved,
			Path:       filepath.Join("internal", "storage", "db.go"),
			Line:       42,
			CommitHash: "commit-b",
			OccurredAt: baseTime.Add(time.Hour),
		},
		{
			MemoryID:   "memory-123",
			Type:       history.EventChanged,
			Path:       filepath.Join("internal", "storage", "db.go"),
			Line:       51,
			CommitHash: "commit-c",
			OccurredAt: baseTime.Add(2 * time.Hour),
		},
	}

	for _, event := range events {
		if err := store.SaveHistory(event); err != nil {
			store.Close()
			t.Fatalf("SaveHistory() error = %v", err)
		}
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	store, err = OpenSQLite(root)
	if err != nil {
		t.Fatalf("reopen OpenSQLite() error = %v", err)
	}
	defer store.Close()

	got, err := store.ListHistory("memory-123")
	if err != nil {
		t.Fatalf("ListHistory() error = %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("ListHistory() returned %d events, want 3", len(got))
	}

	wantTypes := []history.EventType{
		history.EventCreated,
		history.EventMoved,
		history.EventChanged,
	}

	for i := range got {
		if got[i].MemoryID != "memory-123" {
			t.Fatalf(
				"event %d MemoryID = %q, want %q",
				i,
				got[i].MemoryID,
				"memory-123",
			)
		}

		if got[i].Type != wantTypes[i] {
			t.Fatalf(
				"event %d Type = %q, want %q",
				i,
				got[i].Type,
				wantTypes[i],
			)
		}

		if got[i].CommitHash != events[i].CommitHash {
			t.Fatalf(
				"event %d CommitHash = %q, want %q",
				i,
				got[i].CommitHash,
				events[i].CommitHash,
			)
		}

		if !got[i].OccurredAt.Equal(events[i].OccurredAt) {
			t.Fatalf(
				"event %d OccurredAt = %v, want %v",
				i,
				got[i].OccurredAt,
				events[i].OccurredAt,
			)
		}
	}

	if got[0].Path != "main.go" || got[0].Line != 10 {
		t.Fatalf(
			"created location = %s:%d, want main.go:10",
			got[0].Path,
			got[0].Line,
		)
	}

	wantMovedPath := filepath.Join("internal", "storage", "db.go")

	if got[1].Path != wantMovedPath || got[1].Line != 42 {
		t.Fatalf(
			"moved location = %s:%d, want %s:42",
			got[1].Path,
			got[1].Line,
			wantMovedPath,
		)
	}

	if got[2].Path != wantMovedPath || got[2].Line != 51 {
		t.Fatalf(
			"changed location = %s:%d, want %s:51",
			got[2].Path,
			got[2].Line,
			wantMovedPath,
		)
	}
}

func TestSQLiteHistoryPersistsReasoningSnapshotAcrossReopen(t *testing.T) {
	root := t.TempDir()

	if err := os.Mkdir(
		filepath.Join(root, ".whytie"),
		0755,
	); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	store, err := OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}

	event := history.Event{
		MemoryID:   "decision-1",
		Type:       history.EventDeleted,
		Kind:       syntax.Decision,
		Text:       "SQLite WAL 모드를 사용한다",
		Path:       "database.go",
		Line:       4,
		CommitHash: "fc30570",
		OccurredAt: time.Date(
			2026,
			time.September,
			20,
			12,
			0,
			0,
			0,
			time.UTC,
		),
	}

	if err := store.SaveHistory(event); err != nil {
		store.Close()
		t.Fatalf("SaveHistory() error = %v", err)
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	store, err = OpenSQLite(root)
	if err != nil {
		t.Fatalf("reopen OpenSQLite() error = %v", err)
	}
	defer store.Close()

	events, err := store.ListHistory("decision-1")
	if err != nil {
		t.Fatalf("ListHistory() error = %v", err)
	}

	if len(events) != 1 {
		t.Fatalf("history contains %d events, want 1", len(events))
	}

	got := events[0]

	if got.Kind != syntax.Decision {
		t.Errorf(
			"Kind = %q, want %q",
			got.Kind,
			syntax.Decision,
		)
	}

	if got.Text != "SQLite WAL 모드를 사용한다" {
		t.Errorf(
			"Text = %q, want %q",
			got.Text,
			"SQLite WAL 모드를 사용한다",
		)
	}
}

func TestSQLiteFindHistoryAtReturnsDeletedReasoning(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(
		filepath.Join(root, ".whytie"),
		0755,
	); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	store, err := OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	created := history.Event{
		MemoryID:   "decision-1",
		Type:       history.EventCreated,
		Kind:       syntax.Decision,
		Text:       "SQLite를 사용한다",
		Path:       "main.go",
		Line:       4,
		CommitHash: "abc1234",
		OccurredAt: time.Now().UTC(),
	}

	deleted := history.Event{
		MemoryID:   "decision-1",
		Type:       history.EventDeleted,
		Kind:       syntax.Decision,
		Text:       "SQLite WAL 모드를 사용한다",
		Path:       "database.go",
		Line:       4,
		CommitHash: "def5678",
		OccurredAt: time.Now().UTC(),
	}

	if err := store.SaveHistory(created); err != nil {
		t.Fatalf("SaveHistory(created) error = %v", err)
	}

	if err := store.SaveHistory(deleted); err != nil {
		t.Fatalf("SaveHistory(deleted) error = %v", err)
	}

	event, ok, err := store.FindHistoryAt("database.go", 4)
	if err != nil {
		t.Fatalf("FindHistoryAt() error = %v", err)
	}

	if !ok {
		t.Fatal("FindHistoryAt() found no event, want deleted reasoning")
	}

	if event.MemoryID != "decision-1" {
		t.Errorf(
			"MemoryID = %q, want %q",
			event.MemoryID,
			"decision-1",
		)
	}

	if event.Type != history.EventDeleted {
		t.Errorf(
			"Type = %q, want %q",
			event.Type,
			history.EventDeleted,
		)
	}

	if event.Kind != syntax.Decision {
		t.Errorf(
			"Kind = %q, want %q",
			event.Kind,
			syntax.Decision,
		)
	}

	if event.Text != "SQLite WAL 모드를 사용한다" {
		t.Errorf(
			"Text = %q, want %q",
			event.Text,
			"SQLite WAL 모드를 사용한다",
		)
	}
}
