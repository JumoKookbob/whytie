package storage_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/JumoKookbob/whytie/internal/history"
	"github.com/JumoKookbob/whytie/internal/storage"
)

func TestListRecentHistory(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".whytie"), 0755); err != nil {
		t.Fatal(err)
	}

	store, err := storage.OpenSQLite(root)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	empty, err := store.ListRecentHistory(5)
	if err != nil {
		t.Fatal(err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected no events, got %d", len(empty))
	}

	// 같은 시각이어도 저장된 순서로 최신 기록을 구분한다.
	at := time.Date(2026, 9, 20, 1, 0, 0, 123, time.UTC)
	events := []history.Event{
		{
			MemoryID:   "first",
			Type:       history.EventCreated,
			Kind:       "decision",
			Text:       "Use SQLite",
			Path:       "main.go",
			Line:       3,
			OccurredAt: at,
		},
		{
			MemoryID:   "second",
			Type:       history.EventChanged,
			Kind:       "question",
			Text:       "How should shutdown work?",
			Path:       "shutdown.go",
			Line:       8,
			OccurredAt: at,
		},
		{
			MemoryID:   "first",
			Type:       history.EventDeleted,
			Kind:       "decision",
			Text:       "Use SQLite",
			Path:       "main.go",
			Line:       3,
			CommitHash: "abcdef1234567890",
			OccurredAt: at,
		},
	}

	for _, event := range events {
		if err := store.SaveHistory(event); err != nil {
			t.Fatal(err)
		}
	}

	got, err := store.ListRecentHistory(2)
	if err != nil {
		t.Fatal(err)
	}

	want := []history.Event{events[2], events[1]}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("recent history mismatch:\ngot:  %#v\nwant: %#v", got, want)
	}

	for _, limit := range []int{0, -1} {
		if _, err := store.ListRecentHistory(limit); err == nil {
			t.Errorf("limit %d: expected an error", limit)
		}
	}
}
