package app

import (
	"bytes"
	"database/sql"
	"testing"

	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestFindReasonsReturnsReasonsFollowingDecision(t *testing.T) {
	memories := []memory.Memory{
		{
			ID:          "decision-1",
			Kind:        syntax.Decision,
			Text:        "SQLite를 사용한다",
			CurrentPath: "internal/store/db.go",
			CurrentLine: 20,
		},
		{
			ID:          "reason-1",
			Kind:        syntax.Reason,
			Text:        "local-first에 적합하기 때문에",
			CurrentPath: "internal/store/db.go",
			CurrentLine: 21,
		},
	}

	reasons := FindReasons(memories, "decision-1")

	if len(reasons) != 1 {
		t.Fatalf("FindReasons() returned %d reasons, want 1", len(reasons))
	}

	if reasons[0].ID != "reason-1" {
		t.Errorf(
			"reason ID = %q, want %q",
			reasons[0].ID,
			"reason-1",
		)
	}

	if reasons[0].Text != "local-first에 적합하기 때문에" {
		t.Errorf(
			"reason text = %q, want %q",
			reasons[0].Text,
			"local-first에 적합하기 때문에",
		)
	}
}

func TestFindReasonsReturnsMultipleFollowingReasons(t *testing.T) {
	memories := []memory.Memory{
		{
			ID:          "decision-1",
			Kind:        syntax.Decision,
			Text:        "SQLite를 사용한다",
			CurrentPath: "internal/store/db.go",
			CurrentLine: 20,
		},
		{
			ID:          "reason-1",
			Kind:        syntax.Reason,
			Text:        "local-first에 적합하기 때문에",
			CurrentPath: "internal/store/db.go",
			CurrentLine: 21,
		},
		{
			ID:          "reason-2",
			Kind:        syntax.Reason,
			Text:        "별도 서버가 필요 없기 때문에",
			CurrentPath: "internal/store/db.go",
			CurrentLine: 22,
		},
	}

	reasons := FindReasons(memories, "decision-1")

	if len(reasons) != 2 {
		t.Fatalf("FindReasons() returned %d reasons, want 2", len(reasons))
	}

	if reasons[0].ID != "reason-1" {
		t.Errorf("first reason ID = %q, want %q", reasons[0].ID, "reason-1")
	}

	if reasons[1].ID != "reason-2" {
		t.Errorf("second reason ID = %q, want %q", reasons[1].ID, "reason-2")
	}
}

func TestWhyWritesMemoryAndReasons(t *testing.T) {
	store := &fakeWhyStore{
		memories: []memory.Memory{
			{
				ID:          "decision-1",
				Kind:        syntax.Decision,
				Text:        "SQLite를 사용한다",
				CurrentPath: "example.go",
				CurrentLine: 3,
			},
			{
				ID:          "reason-1",
				Kind:        syntax.Reason,
				Text:        "local-first에 적합하기 때문에",
				CurrentPath: "example.go",
				CurrentLine: 4,
			},
		},
	}

	var stdout bytes.Buffer

	err := Why(&stdout, store, "decision-1")
	if err != nil {
		t.Fatalf("Why() error = %v", err)
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
}

type fakeWhyStore struct {
	memories []memory.Memory
}

func (s *fakeWhyStore) Save(memory.Memory) error {
	return nil
}

func (s *fakeWhyStore) Delete(id string) error {
	return nil
}

func (s *fakeWhyStore) Get(id string) (memory.Memory, error) {
	for _, m := range s.memories {
		if m.ID == id {
			return m, nil
		}
	}

	return memory.Memory{}, sql.ErrNoRows
}

func (s *fakeWhyStore) List() ([]memory.Memory, error) {
	return s.memories, nil
}

func TestFindMemoryAtReturnsMemoryAtCurrentLocation(t *testing.T) {
	memories := []memory.Memory{
		{
			ID:          "decision-1",
			Kind:        syntax.Decision,
			Text:        "SQLite를 사용한다",
			CurrentPath: "example.go",
			CurrentLine: 10,
		},
		{
			ID:          "decision-2",
			Kind:        syntax.Decision,
			Text:        "JSON을 사용한다",
			CurrentPath: "other.go",
			CurrentLine: 20,
		},
	}

	got, ok := FindMemoryAt(memories, "example.go", 10)

	if !ok {
		t.Fatal("FindMemoryAt() ok = false, want true")
	}

	if got.ID != "decision-1" {
		t.Errorf("FindMemoryAt() ID = %q, want %q", got.ID, "decision-1")
	}

	if got.Text != "SQLite를 사용한다" {
		t.Errorf(
			"FindMemoryAt() Text = %q, want %q",
			got.Text,
			"SQLite를 사용한다",
		)
	}
}

func TestParseLocation(t *testing.T) {
	path, line, err := ParseLocation("internal/store/db.go:42")
	if err != nil {
		t.Fatalf("ParseLocation() error = %v", err)
	}

	if path != "internal/store/db.go" {
		t.Errorf("path = %q, want %q", path, "internal/store/db.go")
	}

	if line != 42 {
		t.Errorf("line = %d, want 42", line)
	}
}

func TestParseLocationSupportsWindowsPath(t *testing.T) {
	path, line, err := ParseLocation(`C:\project\internal\store\db.go:42`)
	if err != nil {
		t.Fatalf("ParseLocation() error = %v", err)
	}

	wantPath := `C:\project\internal\store\db.go`

	if path != wantPath {
		t.Errorf("path = %q, want %q", path, wantPath)
	}

	if line != 42 {
		t.Errorf("line = %d, want 42", line)
	}
}

func TestWhyAtWritesMemoryAndReasonsFromLocation(t *testing.T) {
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
		},
	}

	var stdout bytes.Buffer

	err := WhyAt(&stdout, store, "example.go:10")
	if err != nil {
		t.Fatalf("WhyAt() error = %v", err)
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
}

func TestParseLocationRejectsEmptyPath(t *testing.T) {
	_, _, err := ParseLocation(":42")
	if err == nil {
		t.Fatal("ParseLocation() error = nil, want error")
	}
}

func TestParseLocationRejectsInvalidLocations(t *testing.T) {
	tests := []string{
		"",
		"example.go",
		":42",
		"example.go:",
		"example.go:abc",
		"example.go:0",
		"example.go:-1",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, _, err := ParseLocation(input)
			if err == nil {
				t.Fatalf("ParseLocation(%q) error = nil, want error", input)
			}
		})
	}
}

func TestWhyAtReturnsErrorWhenLocationHasNoMemory(t *testing.T) {
	store := &fakeWhyStore{
		memories: []memory.Memory{},
	}

	var stdout bytes.Buffer

	err := WhyAt(&stdout, store, "missing.go:42")
	if err == nil {
		t.Fatal("WhyAt() error = nil, want error")
	}

	want := "no memory at missing.go:42"

	if err.Error() != want {
		t.Errorf("WhyAt() error = %q, want %q", err.Error(), want)
	}

	if stdout.String() != "" {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
}

func TestWhyReturnsErrorWhenMemoryDoesNotExist(t *testing.T) {
	store := &fakeWhyStore{
		memories: []memory.Memory{},
	}

	var stdout bytes.Buffer

	err := Why(&stdout, store, "does-not-exist")
	if err == nil {
		t.Fatal("Why() error = nil, want error")
	}

	if stdout.String() != "" {
		t.Errorf("stdout = %q, want empty", stdout.String())
	}
}

func TestFindReasonsReturnsNoneWhenTargetIsReason(t *testing.T) {
	memories := []memory.Memory{
		{
			ID:          "reason-1",
			Kind:        syntax.Reason,
			Text:        "첫 번째 이유",
			CurrentPath: "example.go",
			CurrentLine: 10,
		},
		{
			ID:          "reason-2",
			Kind:        syntax.Reason,
			Text:        "두 번째 이유",
			CurrentPath: "example.go",
			CurrentLine: 11,
		},
	}

	reasons := FindReasons(memories, "reason-1")

	if len(reasons) != 0 {
		t.Fatalf(
			"FindReasons() returned %d reasons, want 0",
			len(reasons),
		)
	}
}

func TestFindReasonsStopsAtLineGap(t *testing.T) {
	memories := []memory.Memory{
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
			CurrentLine: 13,
		},
	}

	reasons := FindReasons(memories, "decision-1")

	if len(reasons) != 1 {
		t.Fatalf(
			"FindReasons() returned %d reasons, want 1",
			len(reasons),
		)
	}

	if reasons[0].ID != "reason-1" {
		t.Errorf(
			"reason ID = %q, want %q",
			reasons[0].ID,
			"reason-1",
		)
	}
}

func TestFindReasonsStopsAtDifferentFile(t *testing.T) {
	memories := []memory.Memory{
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
			CurrentPath: "other.go",
			CurrentLine: 12,
		},
	}

	reasons := FindReasons(memories, "decision-1")

	if len(reasons) != 1 {
		t.Fatalf(
			"FindReasons() returned %d reasons, want 1",
			len(reasons),
		)
	}

	if reasons[0].ID != "reason-1" {
		t.Errorf(
			"reason ID = %q, want %q",
			reasons[0].ID,
			"reason-1",
		)
	}
}

func TestFindReasonsStopsAtNonReason(t *testing.T) {
	memories := []memory.Memory{
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
			ID:          "question-1",
			Kind:        syntax.Question,
			Text:        "WAL을 사용할까?",
			CurrentPath: "example.go",
			CurrentLine: 12,
		},
		{
			ID:          "reason-2",
			Kind:        syntax.Reason,
			Text:        "성능 때문에",
			CurrentPath: "example.go",
			CurrentLine: 13,
		},
	}

	reasons := FindReasons(memories, "decision-1")

	if len(reasons) != 1 {
		t.Fatalf(
			"FindReasons() returned %d reasons, want 1",
			len(reasons),
		)
	}

	if reasons[0].ID != "reason-1" {
		t.Errorf(
			"reason ID = %q, want %q",
			reasons[0].ID,
			"reason-1",
		)
	}
}

func TestFindReasonsReturnsNoneWhenFirstFollowingMemoryIsNotReason(t *testing.T) {
	memories := []memory.Memory{
		{
			ID:          "decision-1",
			Kind:        syntax.Decision,
			Text:        "SQLite를 사용한다",
			CurrentPath: "example.go",
			CurrentLine: 10,
		},
		{
			ID:          "question-1",
			Kind:        syntax.Question,
			Text:        "WAL을 사용할까?",
			CurrentPath: "example.go",
			CurrentLine: 11,
		},
		{
			ID:          "reason-1",
			Kind:        syntax.Reason,
			Text:        "성능 때문에",
			CurrentPath: "example.go",
			CurrentLine: 12,
		},
	}

	reasons := FindReasons(memories, "decision-1")

	if len(reasons) != 0 {
		t.Fatalf(
			"FindReasons() returned %d reasons, want 0",
			len(reasons),
		)
	}
}

func TestFindReasonsReturnsNoneWhenTargetDoesNotExist(t *testing.T) {
	memories := []memory.Memory{
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
	}

	reasons := FindReasons(memories, "does-not-exist")

	if len(reasons) != 0 {
		t.Fatalf(
			"FindReasons() returned %d reasons, want 0",
			len(reasons),
		)
	}
}

func TestFindReasonsHandlesEmptyMemories(t *testing.T) {
	var memories []memory.Memory

	reasons := FindReasons(memories, "decision-1")

	if len(reasons) != 0 {
		t.Fatalf(
			"FindReasons() returned %d reasons, want 0",
			len(reasons),
		)
	}
}

func TestWhyOnQuestionWritesDecisionAndReasons(t *testing.T) {
	store := &fakeWhyStore{
		memories: []memory.Memory{
			{
				ID:          "question-1",
				Kind:        syntax.Question,
				Text:        "어떤 DB를 쓸까?",
				CurrentPath: "example.go",
				CurrentLine: 10,
			},
			{
				ID:          "decision-1",
				Kind:        syntax.Decision,
				Text:        "SQLite를 사용한다",
				CurrentPath: "example.go",
				CurrentLine: 11,
			},
			{
				ID:          "reason-1",
				Kind:        syntax.Reason,
				Text:        "local-first에 적합하기 때문에",
				CurrentPath: "example.go",
				CurrentLine: 12,
			},
			{
				ID:          "reason-2",
				Kind:        syntax.Reason,
				Text:        "별도 서버가 필요 없기 때문에",
				CurrentPath: "example.go",
				CurrentLine: 13,
			},
		},
	}

	var stdout bytes.Buffer

	err := Why(&stdout, store, "question-1")
	if err != nil {
		t.Fatalf("Why() error = %v", err)
	}

	want := "" +
		"📄 example.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"   10  ❓ QUESTION  어떤 DB를 쓸까?\n" +
		"\n" +
		"   11  ✅ DECISION  SQLite를 사용한다\n" +
		"   12  └─ 💡 REASON    local-first에 적합하기 때문에\n" +
		"   13  └─ 💡 REASON    별도 서버가 필요 없기 때문에\n"

	if stdout.String() != want {
		t.Errorf(
			"stdout =\n%q\nwant:\n%q",
			stdout.String(),
			want,
		)
	}
}
