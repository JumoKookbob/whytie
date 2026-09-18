package app

import (
	"bytes"
	"testing"

	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

type listStore struct {
	memories []memory.Memory
}

func (s *listStore) Save(m memory.Memory) error {
	return nil
}

func (s *listStore) Get(id string) (memory.Memory, error) {
	return memory.Memory{}, nil
}

func (s *listStore) List() ([]memory.Memory, error) {
	return s.memories, nil
}

func TestListWritesReasoningStructure(t *testing.T) {
	store := &listStore{
		memories: []memory.Memory{
			{
				ID:          "memory-1",
				Kind:        syntax.Decision,
				Text:        "SQLite",
				CreatedPath: "internal/store/db.go",
				CreatedLine: 20,
				CurrentPath: "internal/store/db.go",
				CurrentLine: 40,
			},
			{
				ID:          "memory-2",
				Kind:        syntax.Reason,
				Text:        "local-first에 적합함",
				CreatedPath: "internal/store/db.go",
				CreatedLine: 21,
				CurrentPath: "internal/store/db.go",
				CurrentLine: 41,
			},
		},
	}

	var output bytes.Buffer

	err := List(&output, store)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := "decision: SQLite  internal/store/db.go:40\n" +
		"└─ reason: local-first에 적합함  internal/store/db.go:41\n"

	if output.String() != want {
		t.Errorf(
			"List() output =\n%q\nwant:\n%q",
			output.String(),
			want,
		)
	}
}
