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

func (s *listStore) Delete(id string) error {
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
				Text:        "Use SQLite",
				CreatedPath: "internal/store/db.go",
				CreatedLine: 20,
				CurrentPath: "internal/store/db.go",
				CurrentLine: 40,
			},
			{
				ID:          "memory-2",
				Kind:        syntax.Reason,
				Text:        "It fits the local-first design",
				CreatedPath: "internal/store/db.go",
				CreatedLine: 21,
				CurrentPath: "internal/store/db.go",
				CurrentLine: 41,
			},
			{
				ID:          "memory-3",
				Kind:        syntax.Question,
				Text:        "How should caching work?",
				CreatedPath: "internal/store/db.go",
				CreatedLine: 50,
				CurrentPath: "internal/store/db.go",
				CurrentLine: 50,
			},
		},
	}

	var output bytes.Buffer

	err := List(&output, store)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	want := "" +
		"📄 internal/store/db.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"   40  ✅ DECISION  Use SQLite\n" +
		"   41  └─ 💡 REASON    It fits the local-first design\n" +
		"\n" +
		"────────────────────────────────────────────────────\n" +
		"\n" +
		"   50  ❓ QUESTION  How should caching work?\n"

	if output.String() != want {
		t.Errorf(
			"List() output =\n%q\nwant:\n%q",
			output.String(),
			want,
		)
	}
}
