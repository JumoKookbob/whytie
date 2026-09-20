package history

import (
	"testing"

	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestEventRepresentsReasoningCreation(t *testing.T) {
	event := Event{
		MemoryID:   "memory-1",
		Type:       EventCreated,
		Path:       "internal/storage/sqlite.go",
		Line:       42,
		CommitHash: "abc123",
	}

	if event.MemoryID != "memory-1" {
		t.Errorf("MemoryID = %q, want %q", event.MemoryID, "memory-1")
	}

	if event.Type != EventCreated {
		t.Errorf("Type = %q, want %q", event.Type, EventCreated)
	}

	if event.Path != "internal/storage/sqlite.go" {
		t.Errorf(
			"Path = %q, want %q",
			event.Path,
			"internal/storage/sqlite.go",
		)
	}

	if event.Line != 42 {
		t.Errorf("Line = %d, want %d", event.Line, 42)
	}

	if event.CommitHash != "abc123" {
		t.Errorf(
			"CommitHash = %q, want %q",
			event.CommitHash,
			"abc123",
		)
	}
}

func TestEventTypesAreStable(t *testing.T) {
	tests := []struct {
		name string
		got  EventType
		want EventType
	}{
		{
			name: "created",
			got:  EventCreated,
			want: "created",
		},
		{
			name: "moved",
			got:  EventMoved,
			want: "moved",
		},
		{
			name: "changed",
			got:  EventChanged,
			want: "changed",
		},
		{
			name: "deleted",
			got:  EventDeleted,
			want: "deleted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf(
					"event type = %q, want %q",
					tt.got,
					tt.want,
				)
			}
		})
	}
}

func TestDetectReasoningChanges(t *testing.T) {
	base := memory.Memory{
		ID:          "memory-1",
		Kind:        syntax.Decision,
		Text:        "SQLite를 사용한다",
		CreatedPath: "main.go",
		CreatedLine: 10,
		CurrentPath: "main.go",
		CurrentLine: 10,
	}

	tests := []struct {
		name     string
		previous memory.Memory
		current  memory.Memory
		wantType EventType
		want     bool
	}{
		{
			name:     "unchanged",
			previous: base,
			current:  base,
			wantType: "",
			want:     false,
		},
		{
			name:     "moved to another file",
			previous: base,
			current: func() memory.Memory {
				m := base
				m.CurrentPath = "internal/storage/sqlite.go"
				m.CurrentLine = 42
				return m
			}(),
			wantType: EventMoved,
			want:     true,
		},
		{
			name:     "moved within same file",
			previous: base,
			current: func() memory.Memory {
				m := base
				m.CurrentLine = 20
				return m
			}(),
			wantType: EventMoved,
			want:     true,
		},
		{
			name:     "text changed",
			previous: base,
			current: func() memory.Memory {
				m := base
				m.Text = "SQLite WAL 모드를 사용한다"
				return m
			}(),
			wantType: EventChanged,
			want:     true,
		},
		{
			name:     "kind changed",
			previous: base,
			current: func() memory.Memory {
				m := base
				m.Kind = syntax.Reason
				return m
			}(),
			wantType: EventChanged,
			want:     true,
		},
		{
			name:     "different identity",
			previous: base,
			current: func() memory.Memory {
				m := base
				m.ID = "memory-2"
				return m
			}(),
			wantType: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, got := Detect(tt.previous, tt.current)

			if got != tt.want {
				t.Fatalf(
					"Detect() changed = %v, want %v",
					got,
					tt.want,
				)
			}

			if gotType != tt.wantType {
				t.Errorf(
					"Detect() type = %q, want %q",
					gotType,
					tt.wantType,
				)
			}
		})
	}
}

func TestEventPreservesReasoningSnapshot(t *testing.T) {
	event := Event{
		MemoryID: "decision-1",
		Type:     EventDeleted,
		Kind:     syntax.Decision,
		Text:     "SQLite WAL 모드를 사용한다",
		Path:     "database.go",
		Line:     4,
	}

	if event.Kind != syntax.Decision {
		t.Errorf("Kind = %q, want %q", event.Kind, syntax.Decision)
	}

	if event.Text != "SQLite WAL 모드를 사용한다" {
		t.Errorf(
			"Text = %q, want %q",
			event.Text,
			"SQLite WAL 모드를 사용한다",
		)
	}
}
