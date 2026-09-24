package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/JumoKookbob/whytie/internal/history"
	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/storage"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestHistoryWritesReasoningHistory(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(
		filepath.Join(root, ".whytie"),
		0755,
	); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	store, err := storage.OpenSQLite(root)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	m := memory.Memory{
		ID:          "memory-1",
		Kind:        syntax.Kind("decision"),
		Text:        "SQLite를 사용한다",
		CreatedPath: "main.go",
		CreatedLine: 10,
		CurrentPath: "internal/storage/db.go",
		CurrentLine: 24,
	}

	if err := store.Save(m); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	events := []history.Event{
		{
			MemoryID:   m.ID,
			Type:       history.EventCreated,
			Path:       "main.go",
			Line:       10,
			CommitHash: "1234567890abcdef",
			OccurredAt: time.Now().UTC(),
		},
		{
			MemoryID:   m.ID,
			Type:       history.EventMoved,
			Path:       "internal/storage/db.go",
			Line:       24,
			CommitHash: "abcdef1234567890",
			OccurredAt: time.Now().UTC(),
		},
	}

	for _, event := range events {
		if err := store.SaveHistory(event); err != nil {
			t.Fatalf("SaveHistory() error = %v", err)
		}
	}

	var output bytes.Buffer

	if err := History(&output, store, m); err != nil {
		t.Fatalf("History() error = %v", err)
	}

	got := output.String()

	wants := []string{
		"DECISION  SQLite를 사용한다",
		"internal/storage/db.go",
		"24  ✅ DECISION",
		"History",
		"CREATED",
		"main.go:10",
		"commit 1234567",
		"MOVED",
		"commit abcdef1",
	}

	for _, want := range wants {
		if !strings.Contains(got, want) {
			t.Errorf(
				"History() output does not contain %q\noutput:\n%s",
				want,
				got,
			)
		}
	}
}
