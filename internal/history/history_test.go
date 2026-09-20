package history

import (
	"testing"
	"time"

	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestEventRepresentsReasoningCreation(t *testing.T) {
	when := time.Date(
		2026,
		time.September,
		20,
		10,
		30,
		0,
		0,
		time.UTC,
	)

	event := Event{
		MemoryID:   "memory-123",
		Type:       EventCreated,
		Path:       "internal/storage/db.go",
		Line:       42,
		CommitHash: "abc123",
		OccurredAt: when,
	}

	if event.MemoryID != "memory-123" {
		t.Fatalf("MemoryID = %q, want %q", event.MemoryID, "memory-123")
	}

	if event.Type != EventCreated {
		t.Fatalf("Type = %q, want %q", event.Type, EventCreated)
	}

	if event.Path != "internal/storage/db.go" {
		t.Fatalf(
			"Path = %q, want %q",
			event.Path,
			"internal/storage/db.go",
		)
	}

	if event.Line != 42 {
		t.Fatalf("Line = %d, want 42", event.Line)
	}

	if event.CommitHash != "abc123" {
		t.Fatalf(
			"CommitHash = %q, want %q",
			event.CommitHash,
			"abc123",
		)
	}

	if !event.OccurredAt.Equal(when) {
		t.Fatalf(
			"OccurredAt = %v, want %v",
			event.OccurredAt,
			when,
		)
	}
}

func TestEventTypesAreStable(t *testing.T) {
	tests := []struct {
		got  EventType
		want string
	}{
		{EventCreated, "created"},
		{EventMoved, "moved"},
		{EventChanged, "changed"},
	}

	for _, test := range tests {
		if string(test.got) != test.want {
			t.Fatalf(
				"EventType = %q, want %q",
				test.got,
				test.want,
			)
		}
	}
}

func TestDetectReasoningChanges(t *testing.T) {
	tests := []struct {
		name     string
		previous memory.Memory
		current  memory.Memory
		wantType EventType
		wantOK   bool
	}{
		{
			name: "unchanged",
			previous: memory.Memory{
				ID:          "a",
				Kind:        syntax.Kind("decision"),
				Text:        "SQLite를 사용한다",
				CurrentPath: "main.go",
				CurrentLine: 10,
			},
			current: memory.Memory{
				ID:          "a",
				Kind:        syntax.Kind("decision"),
				Text:        "SQLite를 사용한다",
				CurrentPath: "main.go",
				CurrentLine: 10,
			},
			wantOK: false,
		},
		{
			name: "moved to another file",
			previous: memory.Memory{
				ID:          "a",
				Kind:        syntax.Kind("decision"),
				Text:        "SQLite를 사용한다",
				CurrentPath: "main.go",
				CurrentLine: 10,
			},
			current: memory.Memory{
				ID:          "a",
				Kind:        syntax.Kind("decision"),
				Text:        "SQLite를 사용한다",
				CurrentPath: "internal/storage/db.go",
				CurrentLine: 42,
			},
			wantType: EventMoved,
			wantOK:   true,
		},
		{
			name: "moved within same file",
			previous: memory.Memory{
				ID:          "a",
				Kind:        syntax.Kind("decision"),
				Text:        "SQLite를 사용한다",
				CurrentPath: "main.go",
				CurrentLine: 10,
			},
			current: memory.Memory{
				ID:          "a",
				Kind:        syntax.Kind("decision"),
				Text:        "SQLite를 사용한다",
				CurrentPath: "main.go",
				CurrentLine: 30,
			},
			wantType: EventMoved,
			wantOK:   true,
		},
		{
			name: "text changed",
			previous: memory.Memory{
				ID:          "a",
				Kind:        syntax.Kind("decision"),
				Text:        "SQLite를 사용한다",
				CurrentPath: "main.go",
				CurrentLine: 10,
			},
			current: memory.Memory{
				ID:          "a",
				Kind:        syntax.Kind("decision"),
				Text:        "SQLite를 기본 저장소로 사용한다",
				CurrentPath: "main.go",
				CurrentLine: 10,
			},
			wantType: EventChanged,
			wantOK:   true,
		},
		{
			name: "different identity",
			previous: memory.Memory{
				ID:          "a",
				Kind:        syntax.Kind("decision"),
				Text:        "SQLite를 사용한다",
				CurrentPath: "main.go",
				CurrentLine: 10,
			},
			current: memory.Memory{
				ID:          "b",
				Kind:        syntax.Kind("decision"),
				Text:        "SQLite를 사용한다",
				CurrentPath: "main.go",
				CurrentLine: 10,
			},
			wantOK: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotType, gotOK := Detect(test.previous, test.current)

			if gotOK != test.wantOK {
				t.Fatalf(
					"Detect() ok = %v, want %v",
					gotOK,
					test.wantOK,
				)
			}

			if gotType != test.wantType {
				t.Fatalf(
					"Detect() type = %q, want %q",
					gotType,
					test.wantType,
				)
			}
		})
	}
}
