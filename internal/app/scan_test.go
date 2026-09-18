package app

import (
	"bytes"
	"os"
	"testing"

	"github.com/JumoKookbob/whytie/internal/memory"
)

type scanStore struct {
	memories []memory.Memory
}

func (s *scanStore) Save(m memory.Memory) error {
	for i, existing := range s.memories {
		if existing.ID == m.ID {
			s.memories[i] = m
			return nil
		}
	}

	s.memories = append(s.memories, m)
	return nil
}

func (s *scanStore) Get(id string) (memory.Memory, error) {
	for _, m := range s.memories {
		if m.ID == id {
			return m, nil
		}
	}

	return memory.Memory{}, nil
}

func (s *scanStore) List() ([]memory.Memory, error) {
	return s.memories, nil
}

func TestScanScansSyncsAndWritesReasoning(t *testing.T) {
	dir := t.TempDir()

	source := `package example

//+ SQLite를 사용한다
//< local-first에 적합하기 때문에
`

	path := dir + "/example.go"

	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store := &scanStore{}

	var output bytes.Buffer

	err := Scan(&output, store, dir)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	if len(store.memories) != 2 {
		t.Fatalf(
			"saved memories = %d, want 2",
			len(store.memories),
		)
	}

	want := "decision: SQLite를 사용한다  example.go:3\n" +
		"└─ reason: local-first에 적합하기 때문에  example.go:4\n"

	if output.String() != want {
		t.Errorf(
			"Scan() output =\n%q\nwant:\n%q",
			output.String(),
			want,
		)
	}
}
