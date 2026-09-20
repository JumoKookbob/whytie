package app_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/JumoKookbob/whytie/internal/app"
	"github.com/JumoKookbob/whytie/internal/history"
	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

type resumeTestStore struct {
	memories []memory.Memory
	events   []history.Event
}

func (s *resumeTestStore) List() ([]memory.Memory, error) {
	return s.memories, nil
}

func (s *resumeTestStore) ListRecentHistory(limit int) ([]history.Event, error) {
	events := s.events
	if len(events) > limit {
		events = events[:limit]
	}
	return events, nil
}

func TestResumeIncludesDecisionWhenReasonChanges(t *testing.T) {
	store := &resumeTestStore{
		memories: []memory.Memory{
			{
				ID:          "decision",
				Kind:        syntax.Decision,
				Text:        "Flush writes before shutdown",
				CurrentPath: "shutdown.go",
				CurrentLine: 10,
			},
			{
				ID:          "reason",
				Kind:        syntax.Reason,
				Text:        "Prevent pending write loss",
				CurrentPath: "shutdown.go",
				CurrentLine: 11,
			},
			{
				ID:          "unrelated",
				Kind:        syntax.Decision,
				Text:        "Use a blue theme",
				CurrentPath: "theme.go",
				CurrentLine: 3,
			},
		},
		events: []history.Event{
			{
				MemoryID: "reason",
				Type:     history.EventChanged,
				Kind:     syntax.Reason,
				Text:     "Prevent pending write loss",
				Path:     "shutdown.go",
				Line:     11,
			},
		},
	}

	var out bytes.Buffer
	if err := app.Resume(&out, store); err != nil {
		t.Fatal(err)
	}

	got := out.String()
	for _, want := range []string{
		"decision: Flush writes before shutdown",
		"reason: Prevent pending write loss",
		"shutdown.go:10",
		"shutdown.go:11",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in output:\n%s", want, got)
		}
	}

	if strings.Contains(got, "Use a blue theme") {
		t.Errorf("unrelated decision appeared in resume:\n%s", got)
	}
}
