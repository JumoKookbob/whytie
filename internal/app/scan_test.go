package app

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/storage"
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

func (s *scanStore) Delete(id string) error {
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
//< local-first에 적합하기 때문에`

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

	want := "" +
		"📄 example.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"    3  ✅ DECISION  SQLite를 사용한다\n" +
		"    4  └─ 💡 REASON    local-first에 적합하기 때문에\n"

	if output.String() != want {
		t.Errorf(
			"Scan() output =\n%q\nwant:\n%q",
			output.String(),
			want,
		)
	}
}

func TestScanRecordsGitProvenanceInHistory(t *testing.T) {
	dir := t.TempDir()

	runGit := func(args ...string) string {
		t.Helper()

		cmd := exec.Command("git", args...)
		cmd.Dir = dir

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf(
				"git %v error = %v\n%s",
				args,
				err,
				output,
			)
		}

		return strings.TrimSpace(string(output))
	}

	runGit("init")
	runGit("config", "user.email", "test@example.com")
	runGit("config", "user.name", "WhyTie Test")

	source := `package example

// + SQLite를 사용한다
// < local-first에 적합하기 때문에
`

	sourcePath := filepath.Join(dir, "example.go")

	if err := os.WriteFile(
		sourcePath,
		[]byte(source),
		0644,
	); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	runGit("add", "example.go")
	runGit("commit", "-m", "add reasoning")

	wantCommit := runGit("rev-parse", "HEAD")

	if err := os.MkdirAll(
		filepath.Join(dir, ".whytie"),
		0755,
	); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	store, err := storage.OpenSQLite(dir)
	if err != nil {
		t.Fatalf("OpenSQLite() error = %v", err)
	}
	defer store.Close()

	var output bytes.Buffer

	if err := Scan(&output, store, dir); err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	memories, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(memories) != 2 {
		t.Fatalf(
			"List() returned %d memories, want 2",
			len(memories),
		)
	}

	for _, m := range memories {
		events, err := store.ListHistory(m.ID)
		if err != nil {
			t.Fatalf(
				"ListHistory(%q) error = %v",
				m.ID,
				err,
			)
		}

		if len(events) != 1 {
			t.Fatalf(
				"memory %q history contains %d events, want 1",
				m.ID,
				len(events),
			)
		}

		if events[0].CommitHash != wantCommit {
			t.Fatalf(
				"memory %q CommitHash = %q, want %q",
				m.ID,
				events[0].CommitHash,
				wantCommit,
			)
		}
	}
}
