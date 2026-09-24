package app

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JumoKookbob/whytie/internal/history"
	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestRunListOpensRepositoryAndWritesMemories(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	source := `package example

//+ SQLite를 사용한다
//< local-first에 적합하기 때문에
`

	sourcePath := filepath.Join(root, "example.go")

	if err := os.WriteFile(sourcePath, []byte(source), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	store, err := OpenStore(root)
	if err != nil {
		t.Fatalf("OpenStore() error = %v", err)
	}

	var scanOutput bytes.Buffer

	if err := Scan(&scanOutput, store, root); err != nil {
		store.Close()
		t.Fatalf("Scan() error = %v", err)
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	var output bytes.Buffer

	if err := RunList(&output, root); err != nil {
		t.Fatalf("RunList() error = %v", err)
	}

	want := "" +
		"📄 example.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"    3  ✅ DECISION  SQLite를 사용한다\n" +
		"    4  └─ 💡 REASON    local-first에 적합하기 때문에\n"

	if output.String() != want {
		t.Errorf(
			"RunList() output =\n%q\nwant:\n%q",
			output.String(),
			want,
		)
	}
}

func TestOpenStoreOpensRepositoryDatabase(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	store, err := OpenStore(root)
	if err != nil {
		t.Fatalf("OpenStore() error = %v", err)
	}
	defer store.Close()

	path := filepath.Join(root, ".whytie", "whytie.db")

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("database file not created: %v", err)
	}
}

func TestRunScanOpensRepositoryScansAndPersists(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	source := `package example

//+ SQLite를 사용한다
//< local-first에 적합하기 때문에
`

	sourcePath := filepath.Join(root, "example.go")

	if err := os.WriteFile(sourcePath, []byte(source), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var output bytes.Buffer

	if err := RunScan(&output, root, root); err != nil {
		t.Fatalf("RunScan() error = %v", err)
	}

	want := "" +
		"📄 example.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"    3  ✅ DECISION  SQLite를 사용한다\n" +
		"    4  └─ 💡 REASON    local-first에 적합하기 때문에\n"

	if output.String() != want {
		t.Errorf(
			"RunScan() output =\n%q\nwant:\n%q",
			output.String(),
			want,
		)
	}

	store, err := OpenStore(root)
	if err != nil {
		t.Fatalf("OpenStore() error = %v", err)
	}
	defer store.Close()

	memories, err := store.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(memories) != 2 {
		t.Fatalf("persisted memories = %d, want 2", len(memories))
	}
}

func TestRunListFromNestedDirectoryFindsRepositoryRoot(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	source := `package example

//+ SQLite를 사용한다
`

	sourcePath := filepath.Join(root, "example.go")
	if err := os.WriteFile(sourcePath, []byte(source), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var scanOutput bytes.Buffer
	if err := RunScan(&scanOutput, root, root); err != nil {
		t.Fatalf("RunScan() error = %v", err)
	}

	nested := filepath.Join(root, "internal", "example")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	var output bytes.Buffer

	if err := RunListFrom(&output, nested); err != nil {
		t.Fatalf("RunListFrom() error = %v", err)
	}

	want := "" +
		"📄 example.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"    3  ✅ DECISION  SQLite를 사용한다\n"

	if output.String() != want {
		t.Errorf(
			"RunListFrom() output =\n%q\nwant:\n%q",
			output.String(),
			want,
		)
	}
}

func TestRunScanFromNestedDirectoryFindsRepositoryRoot(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	source := `package example

//+ SQLite를 사용한다
//< local-first에 적합하기 때문에
`

	sourcePath := filepath.Join(root, "example.go")
	if err := os.WriteFile(sourcePath, []byte(source), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	nested := filepath.Join(root, "internal", "example")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}

	var output bytes.Buffer

	if err := RunScanFrom(&output, nested, root); err != nil {
		t.Fatalf("RunScanFrom() error = %v", err)
	}

	scanWant := "" +
		"📄 example.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"    3  ✅ DECISION  SQLite를 사용한다\n" +
		"    4  └─ 💡 REASON    local-first에 적합하기 때문에\n"

	if output.String() != scanWant {
		t.Errorf(
			"RunScanFrom() output =\n%q\nwant:\n%q",
			output.String(),
			scanWant,
		)
	}

	var listOutput bytes.Buffer

	if err := RunListFrom(&listOutput, nested); err != nil {
		t.Fatalf("RunListFrom() error = %v", err)
	}

	listWant := "" +
		"📄 example.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"    3  ✅ DECISION  SQLite를 사용한다\n" +
		"    4  └─ 💡 REASON    local-first에 적합하기 때문에\n"

	if listOutput.String() != listWant {
		t.Errorf(
			"persisted output =\n%q\nwant:\n%q",
			listOutput.String(),
			listWant,
		)
	}
}

func TestRunHelp(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "no arguments", args: nil},
		{name: "help", args: []string{"help"}},
		{name: "long flag", args: []string{"--help"}},
		{name: "short flag", args: []string{"-h"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			exitCode := Run(tt.args, &stdout, &stderr, ".")

			if exitCode != 0 {
				t.Errorf("Run() exit code = %d, want 0", exitCode)
			}

			wants := []string{
				"WhyTie",
				"Usage:",
				"whytie <command> [arguments]",
				"Commands:",
				"init",
				"scan <path>",
				"list",
				"why <memory-id|file:line>",
				"history <memory-id|file:line>",
				"resume",
				"version",
				"help",
				"Annotations:",
				"// ?  question",
				"// +  decision",
				"// -  rejected",
				"// x  failed",
				"// <  reason",
				"// !  important",
			}

			for _, want := range wants {
				if !strings.Contains(stdout.String(), want) {
					t.Errorf(
						"Run() stdout does not contain %q\nstdout:\n%s",
						want,
						stdout.String(),
					)
				}
			}

			if stderr.Len() != 0 {
				t.Errorf("stderr = %q, want empty", stderr.String())
			}
		})
	}
}

func TestRunRejectsUnknownCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"wat"},
		&stdout,
		&stderr,
		".",
	)

	if exitCode != 1 {
		t.Errorf("Run() exit code = %d, want 1", exitCode)
	}

	if stdout.String() != "" {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}

	want := "unknown command: wat\n"

	if stderr.String() != want {
		t.Errorf(
			"stderr = %q, want %q",
			stderr.String(),
			want,
		)
	}
}

func TestRunInitCreatesRepository(t *testing.T) {
	root := t.TempDir()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"init"},
		&stdout,
		&stderr,
		root,
	)

	if exitCode != 0 {
		t.Errorf("Run() exit code = %d, want 0", exitCode)
	}

	want := "Initialized WhyTie repository in .whytie\n"

	if stdout.String() != want {
		t.Errorf(
			"stdout = %q, want %q",
			stdout.String(),
			want,
		)
	}

	if stderr.String() != "" {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}

	path := filepath.Join(root, ".whytie")

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat(%q) error = %v", path, err)
	}

	if !info.IsDir() {
		t.Fatalf("%q is not a directory", path)
	}
}

func TestRunListCommand(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	source := `package example

//+ SQLite를 사용한다
//< local-first에 적합하기 때문에
`

	sourcePath := filepath.Join(root, "example.go")
	if err := os.WriteFile(sourcePath, []byte(source), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var scanOutput bytes.Buffer
	if err := RunScanFrom(&scanOutput, root, root); err != nil {
		t.Fatalf("RunScanFrom() error = %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"list"},
		&stdout,
		&stderr,
		root,
	)

	if exitCode != 0 {
		t.Errorf("Run() exit code = %d, want 0", exitCode)
	}

	want := "" +
		"📄 example.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"    3  ✅ DECISION  SQLite를 사용한다\n" +
		"    4  └─ 💡 REASON    local-first에 적합하기 때문에\n"

	if stdout.String() != want {
		t.Errorf(
			"stdout =\n%q\nwant:\n%q",
			stdout.String(),
			want,
		)
	}

	if stderr.String() != "" {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunScanRequiresPath(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"scan"},
		&stdout,
		&stderr,
		".",
	)

	if exitCode != 1 {
		t.Errorf("Run() exit code = %d, want 1", exitCode)
	}

	if stdout.String() != "" {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}

	want := "usage: whytie scan <path>\n"

	if stderr.String() != want {
		t.Errorf(
			"stderr = %q, want %q",
			stderr.String(),
			want,
		)
	}
}

func TestRunWhyRequiresTarget(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"why"},
		&stdout,
		&stderr,
		".",
	)

	if exitCode != 1 {
		t.Errorf("Run() exit code = %d, want 1", exitCode)
	}

	want := "usage: whytie why <memory-id|file:line>\n"

	if stderr.String() != want {
		t.Errorf("stderr = %q, want %q", stderr.String(), want)
	}
}

func TestRunWhyCommand(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	store, err := OpenStore(root)
	if err != nil {
		t.Fatalf("OpenStore() error = %v", err)
	}

	decision := memory.Memory{
		ID:          "decision-1",
		Kind:        syntax.Decision,
		Text:        "SQLite를 사용한다",
		CreatedPath: "example.go",
		CreatedLine: 3,
		CurrentPath: "example.go",
		CurrentLine: 3,
	}

	reason := memory.Memory{
		ID:          "reason-1",
		Kind:        syntax.Reason,
		Text:        "local-first에 적합하기 때문에",
		CreatedPath: "example.go",
		CreatedLine: 4,
		CurrentPath: "example.go",
		CurrentLine: 4,
	}

	if err := store.Save(decision); err != nil {
		store.Close()
		t.Fatalf("Save(decision) error = %v", err)
	}

	if err := store.Save(reason); err != nil {
		store.Close()
		t.Fatalf("Save(reason) error = %v", err)
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"why", "decision-1"},
		&stdout,
		&stderr,
		root,
	)

	if exitCode != 0 {
		t.Errorf("Run() exit code = %d, want 0", exitCode)
	}

	want := "" +
		"📄 example.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"    3  ✅ DECISION  SQLite를 사용한다\n" +
		"    4  └─ 💡 REASON    local-first에 적합하기 때문에\n"

	if stdout.String() != want {
		t.Errorf(
			"stdout =\n%q\nwant:\n%q",
			stdout.String(),
			want,
		)
	}

	if stderr.String() != "" {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunWhyCommandAcceptsLocation(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	store, err := OpenStore(root)
	if err != nil {
		t.Fatalf("OpenStore() error = %v", err)
	}

	decision := memory.Memory{
		ID:          "decision-location",
		Kind:        syntax.Decision,
		Text:        "SQLite를 사용한다",
		CreatedPath: "example.go",
		CreatedLine: 10,
		CurrentPath: "example.go",
		CurrentLine: 10,
	}

	reason := memory.Memory{
		ID:          "reason-location",
		Kind:        syntax.Reason,
		Text:        "local-first에 적합하기 때문에",
		CreatedPath: "example.go",
		CreatedLine: 11,
		CurrentPath: "example.go",
		CurrentLine: 11,
	}

	if err := store.Save(decision); err != nil {
		store.Close()
		t.Fatalf("Save(decision) error = %v", err)
	}

	if err := store.Save(reason); err != nil {
		store.Close()
		t.Fatalf("Save(reason) error = %v", err)
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"why", "example.go:10"},
		&stdout,
		&stderr,
		root,
	)

	if exitCode != 0 {
		t.Errorf("Run() exit code = %d, want 0", exitCode)
	}

	want := "" +
		"📄 example.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"   10  ✅ DECISION  SQLite를 사용한다\n" +
		"   11  └─ 💡 REASON    local-first에 적합하기 때문에\n"

	if stdout.String() != want {
		t.Errorf(
			"stdout =\n%q\nwant:\n%q",
			stdout.String(),
			want,
		)
	}

	if stderr.String() != "" {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunWhyRejectsInvalidLocation(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"why", "example.go:abc"},
		&stdout,
		&stderr,
		root,
	)

	if exitCode != 1 {
		t.Errorf("Run() exit code = %d, want 1", exitCode)
	}

	want := "why failed: invalid line number: abc\n"

	if stderr.String() != want {
		t.Errorf("stderr = %q, want %q", stderr.String(), want)
	}

	if stdout.String() != "" {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
}

func TestRunWhyRejectsMissingMemoryID(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"why", "does-not-exist"},
		&stdout,
		&stderr,
		root,
	)

	if exitCode != 1 {
		t.Errorf("Run() exit code = %d, want 1", exitCode)
	}

	if stdout.String() != "" {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}

	if stderr.String() == "" {
		t.Fatal("stderr = empty, want why failure")
	}
}

func TestWhyOnReasonWritesOnlyReason(t *testing.T) {
	store := &fakeWhyStore{
		memories: []memory.Memory{
			{
				ID:          "decision-1",
				Kind:        syntax.Decision,
				Text:        "SQLite를 사용한다",
				CurrentPath: "example.go",
				CurrentLine: 10,
			},
			{
				ID:          "reason-1",
				Kind:        syntax.Reason,
				Text:        "local-first에 적합하기 때문에",
				CurrentPath: "example.go",
				CurrentLine: 11,
			},
			{
				ID:          "reason-2",
				Kind:        syntax.Reason,
				Text:        "설정이 단순하기 때문에",
				CurrentPath: "example.go",
				CurrentLine: 12,
			},
		},
	}

	var stdout bytes.Buffer

	err := Why(&stdout, store, "reason-1")
	if err != nil {
		t.Fatalf("Why() error = %v", err)
	}

	want := "" +
		"📄 example.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"   11  💡 REASON    local-first에 적합하기 때문에\n"

	if stdout.String() != want {
		t.Errorf(
			"stdout =\n%q\nwant:\n%q",
			stdout.String(),
			want,
		)
	}
}

func TestRunVersionCommand(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"version"},
		&stdout,
		&stderr,
		".",
	)

	if exitCode != 0 {
		t.Fatalf("exitCode = %d, want 0", exitCode)
	}

	want := "WhyTie v1.2.1\n"
	if stdout.String() != want {
		t.Errorf("stdout = %q, want %q", stdout.String(), want)
	}

	if stderr.String() != "" {
		t.Errorf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunHistoryCommand(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	store, err := OpenStore(root)
	if err != nil {
		t.Fatalf("OpenStore() error = %v", err)
	}

	target := memory.Memory{
		ID:          "history-decision",
		Kind:        syntax.Decision,
		Text:        "SQLite를 사용한다",
		CreatedPath: "example.go",
		CreatedLine: 10,
		CurrentPath: "internal/storage/db.go",
		CurrentLine: 24,
	}

	if err := store.Save(target); err != nil {
		store.Close()
		t.Fatalf("Save() error = %v", err)
	}

	if err := store.SaveHistory(history.Event{
		MemoryID:   target.ID,
		Type:       history.EventCreated,
		Path:       "example.go",
		Line:       10,
		CommitHash: "1234567890abcdef",
	}); err != nil {
		store.Close()
		t.Fatalf("SaveHistory(created) error = %v", err)
	}

	if err := store.SaveHistory(history.Event{
		MemoryID:   target.ID,
		Type:       history.EventMoved,
		Path:       "internal/storage/db.go",
		Line:       24,
		CommitHash: "abcdef1234567890",
	}); err != nil {
		store.Close()
		t.Fatalf("SaveHistory(moved) error = %v", err)
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"history", target.ID},
		&stdout,
		&stderr,
		root,
	)

	if exitCode != 0 {
		t.Fatalf(
			"Run() exit code = %d, want 0\nstderr:\n%s",
			exitCode,
			stderr.String(),
		)
	}

	wants := []string{
		"DECISION  SQLite를 사용한다",
		"internal/storage/db.go",
		"24  ✅ DECISION",
		"History",
		"CREATED",
		"example.go:10",
		"commit 1234567",
		"MOVED",
		"commit abcdef1",
	}

	for _, want := range wants {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf(
				"Run() stdout does not contain %q\nstdout:\n%s",
				want,
				stdout.String(),
			)
		}
	}

	if stderr.Len() != 0 {
		t.Errorf("Run() stderr = %q, want empty", stderr.String())
	}
}

func TestRunHistoryCommandAcceptsLocation(t *testing.T) {
	root := t.TempDir()

	if err := Init(root); err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	store, err := OpenStore(root)
	if err != nil {
		t.Fatalf("OpenStore() error = %v", err)
	}

	target := memory.Memory{
		ID:          "history-location",
		Kind:        syntax.Decision,
		Text:        "SQLite를 사용한다",
		CreatedPath: "example.go",
		CreatedLine: 10,
		CurrentPath: "example.go",
		CurrentLine: 10,
	}

	if err := store.Save(target); err != nil {
		store.Close()
		t.Fatalf("Save() error = %v", err)
	}

	if err := store.SaveHistory(history.Event{
		MemoryID:   target.ID,
		Type:       history.EventCreated,
		Path:       "example.go",
		Line:       10,
		CommitHash: "1234567890abcdef",
	}); err != nil {
		store.Close()
		t.Fatalf("SaveHistory() error = %v", err)
	}

	if err := store.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := Run(
		[]string{"history", "example.go:10"},
		&stdout,
		&stderr,
		root,
	)

	if exitCode != 0 {
		t.Fatalf(
			"Run() exit code = %d, want 0\nstderr:\n%s",
			exitCode,
			stderr.String(),
		)
	}

	wants := []string{
		"DECISION  SQLite를 사용한다",
		"example.go",
		"10  ✅ DECISION",
		"History",
		"CREATED",
		"commit 1234567",
	}

	for _, want := range wants {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf(
				"Run() stdout does not contain %q\nstdout:\n%s",
				want,
				stdout.String(),
			)
		}
	}

	if stderr.Len() != 0 {
		t.Errorf("Run() stderr = %q, want empty", stderr.String())
	}
}
