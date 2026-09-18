package reconcile

import (
	"testing"

	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestMatchExactExistingMemory(t *testing.T) {
	existing := []memory.Memory{
		{
			ID:          "memory-1",
			Kind:        syntax.Decision,
			Text:        "SQLite",
			CreatedPath: "internal/store/db.go",
			CreatedLine: 21,
			CurrentPath: "internal/store/db.go",
			CurrentLine: 21,
		},
	}

	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite",
		RelativePath: "internal/store/db.go",
		Line:         21,
	}

	got, ok := MatchExact(source, existing)

	if !ok {
		t.Fatal("MatchExact() did not find existing memory")
	}

	if got.ID != "memory-1" {
		t.Errorf("MatchExact() ID = %q, want %q", got.ID, "memory-1")
	}
}

func TestMatchExactRejectsChangedText(t *testing.T) {
	existing := []memory.Memory{
		{
			ID:          "memory-1",
			Kind:        syntax.Decision,
			Text:        "SQLite",
			CreatedPath: "internal/store/db.go",
			CreatedLine: 21,
			CurrentPath: "internal/store/db.go",
			CurrentLine: 21,
		},
	}

	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "PostgreSQL",
		RelativePath: "internal/store/db.go",
		Line:         21,
	}

	_, ok := MatchExact(source, existing)

	if ok {
		t.Fatal("MatchExact() matched memory with changed text")
	}
}

func TestMatchExactRejectsDifferentLocation(t *testing.T) {
	existing := []memory.Memory{
		{
			ID:          "memory-1",
			Kind:        syntax.Decision,
			Text:        "SQLite",
			CreatedPath: "internal/store/db.go",
			CreatedLine: 21,
			CurrentPath: "internal/store/db.go",
			CurrentLine: 21,
		},
	}

	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite",
		RelativePath: "internal/store/db.go",
		Line:         40,
	}

	_, ok := MatchExact(source, existing)

	if ok {
		t.Fatal("MatchExact() matched memory at different location")
	}
}

func TestMatchMovedWithinSameFile(t *testing.T) {
	existing := []memory.Memory{
		{
			ID:          "memory-1",
			Kind:        syntax.Decision,
			Text:        "SQLite",
			CreatedPath: "internal/store/db.go",
			CreatedLine: 21,
			CurrentPath: "internal/store/db.go",
			CurrentLine: 21,
		},
	}

	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite",
		RelativePath: "internal/store/db.go",
		Line:         40,
	}

	got, ok := MatchMoved(source, existing)

	if !ok {
		t.Fatal("MatchMoved() did not find moved memory")
	}

	if got.ID != "memory-1" {
		t.Errorf("MatchMoved() ID = %q, want %q", got.ID, "memory-1")
	}
}

func TestMatchMovedRejectsAmbiguousCandidates(t *testing.T) {
	existing := []memory.Memory{
		{
			ID:          "memory-1",
			Kind:        syntax.Decision,
			Text:        "SQLite",
			CreatedPath: "internal/store/db.go",
			CreatedLine: 20,
			CurrentPath: "internal/store/db.go",
			CurrentLine: 20,
		},
		{
			ID:          "memory-2",
			Kind:        syntax.Decision,
			Text:        "SQLite",
			CreatedPath: "internal/store/db.go",
			CreatedLine: 80,
			CurrentPath: "internal/store/db.go",
			CurrentLine: 80,
		},
	}

	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite",
		RelativePath: "internal/store/db.go",
		Line:         40,
	}

	_, ok := MatchMoved(source, existing)

	if ok {
		t.Fatal("MatchMoved() matched ambiguous candidates")
	}
}

func TestUpdateLocationPreservesCreatedLocation(t *testing.T) {
	existing := memory.Memory{
		ID:          "memory-1",
		Kind:        syntax.Decision,
		Text:        "SQLite",
		CreatedPath: "internal/store/db.go",
		CreatedLine: 21,
		CurrentPath: "internal/store/db.go",
		CurrentLine: 21,
	}

	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite",
		RelativePath: "internal/store/db.go",
		Line:         40,
	}

	got := UpdateLocation(existing, source)

	if got.ID != existing.ID {
		t.Errorf("ID = %q, want %q", got.ID, existing.ID)
	}

	if got.CreatedPath != "internal/store/db.go" {
		t.Errorf(
			"CreatedPath = %q, want %q",
			got.CreatedPath,
			"internal/store/db.go",
		)
	}

	if got.CreatedLine != 21 {
		t.Errorf("CreatedLine = %d, want %d", got.CreatedLine, 21)
	}

	if got.CurrentPath != "internal/store/db.go" {
		t.Errorf(
			"CurrentPath = %q, want %q",
			got.CurrentPath,
			"internal/store/db.go",
		)
	}

	if got.CurrentLine != 40 {
		t.Errorf("CurrentLine = %d, want %d", got.CurrentLine, 40)
	}
}

func TestReconcilePreservesExactMemory(t *testing.T) {
	existing := []memory.Memory{
		{
			ID:          "memory-1",
			Kind:        syntax.Decision,
			Text:        "SQLite",
			CreatedPath: "internal/store/db.go",
			CreatedLine: 21,
			CurrentPath: "internal/store/db.go",
			CurrentLine: 21,
		},
	}

	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite",
		RelativePath: "internal/store/db.go",
		Line:         21,
	}

	got, err := Reconcile(source, existing)
	if err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}

	if got != existing[0] {
		t.Errorf("Reconcile() = %#v, want %#v", got, existing[0])
	}
}

func TestReconcileUpdatesMovedMemoryLocation(t *testing.T) {
	existing := []memory.Memory{
		{
			ID:          "memory-1",
			Kind:        syntax.Decision,
			Text:        "SQLite",
			CreatedPath: "internal/store/db.go",
			CreatedLine: 21,
			CurrentPath: "internal/store/db.go",
			CurrentLine: 21,
		},
	}

	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite",
		RelativePath: "internal/store/db.go",
		Line:         40,
	}

	got, err := Reconcile(source, existing)
	if err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}

	if got.ID != "memory-1" {
		t.Errorf("ID = %q, want %q", got.ID, "memory-1")
	}

	if got.CreatedPath != "internal/store/db.go" {
		t.Errorf(
			"CreatedPath = %q, want %q",
			got.CreatedPath,
			"internal/store/db.go",
		)
	}

	if got.CreatedLine != 21 {
		t.Errorf("CreatedLine = %d, want %d", got.CreatedLine, 21)
	}

	if got.CurrentPath != "internal/store/db.go" {
		t.Errorf(
			"CurrentPath = %q, want %q",
			got.CurrentPath,
			"internal/store/db.go",
		)
	}

	if got.CurrentLine != 40 {
		t.Errorf("CurrentLine = %d, want %d", got.CurrentLine, 40)
	}
}

func TestReconcileCreatesNewMemory(t *testing.T) {
	source := scanner.SourceComment{
		Kind:         syntax.Question,
		Text:         "어떤 DB를 쓸까?",
		RelativePath: "internal/store/db.go",
		Line:         20,
	}

	got, err := Reconcile(source, nil)
	if err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}

	if got.ID == "" {
		t.Error("Reconcile() created memory with empty ID")
	}

	if got.Kind != source.Kind {
		t.Errorf("Kind = %q, want %q", got.Kind, source.Kind)
	}

	if got.Text != source.Text {
		t.Errorf("Text = %q, want %q", got.Text, source.Text)
	}

	if got.CreatedPath != source.RelativePath {
		t.Errorf(
			"CreatedPath = %q, want %q",
			got.CreatedPath,
			source.RelativePath,
		)
	}

	if got.CreatedLine != source.Line {
		t.Errorf("CreatedLine = %d, want %d", got.CreatedLine, source.Line)
	}

	if got.CurrentPath != source.RelativePath {
		t.Errorf(
			"CurrentPath = %q, want %q",
			got.CurrentPath,
			source.RelativePath,
		)
	}

	if got.CurrentLine != source.Line {
		t.Errorf("CurrentLine = %d, want %d", got.CurrentLine, source.Line)
	}
}

func TestReconcileAll(t *testing.T) {
	existing := []memory.Memory{
		{
			ID:          "memory-1",
			Kind:        syntax.Question,
			Text:        "어떤 DB를 쓸까?",
			CreatedPath: "internal/store/db.go",
			CreatedLine: 20,
			CurrentPath: "internal/store/db.go",
			CurrentLine: 20,
		},
		{
			ID:          "memory-2",
			Kind:        syntax.Decision,
			Text:        "SQLite",
			CreatedPath: "internal/store/db.go",
			CreatedLine: 21,
			CurrentPath: "internal/store/db.go",
			CurrentLine: 21,
		},
	}

	sources := []scanner.SourceComment{
		{
			Kind:         syntax.Question,
			Text:         "어떤 DB를 쓸까?",
			RelativePath: "internal/store/db.go",
			Line:         20,
		},
		{
			Kind:         syntax.Decision,
			Text:         "SQLite",
			RelativePath: "internal/store/db.go",
			Line:         40,
		},
		{
			Kind:         syntax.Reason,
			Text:         "local-first에 적합함",
			RelativePath: "internal/store/db.go",
			Line:         41,
		},
	}

	got, err := ReconcileAll(sources, existing)
	if err != nil {
		t.Fatalf("ReconcileAll() error = %v", err)
	}

	if len(got) != 3 {
		t.Fatalf("ReconcileAll() returned %d memories, want 3", len(got))
	}

	// exact memory
	if got[0].ID != "memory-1" {
		t.Errorf("got[0].ID = %q, want %q", got[0].ID, "memory-1")
	}

	// moved memory
	if got[1].ID != "memory-2" {
		t.Errorf("got[1].ID = %q, want %q", got[1].ID, "memory-2")
	}

	if got[1].CreatedLine != 21 {
		t.Errorf("got[1].CreatedLine = %d, want 21", got[1].CreatedLine)
	}

	if got[1].CurrentLine != 40 {
		t.Errorf("got[1].CurrentLine = %d, want 40", got[1].CurrentLine)
	}

	// new memory
	if got[2].ID == "" {
		t.Error("got[2] has empty ID")
	}

	if got[2].Kind != syntax.Reason {
		t.Errorf("got[2].Kind = %q, want %q", got[2].Kind, syntax.Reason)
	}

	if got[2].CurrentLine != 41 {
		t.Errorf("got[2].CurrentLine = %d, want 41", got[2].CurrentLine)
	}
}

func TestReconcileAllUsesEarlierResults(t *testing.T) {
	sources := []scanner.SourceComment{
		{
			Kind:         syntax.Decision,
			Text:         "SQLite",
			RelativePath: "internal/store/db.go",
			Line:         20,
		},
		{
			Kind:         syntax.Decision,
			Text:         "SQLite",
			RelativePath: "internal/store/db.go",
			Line:         40,
		},
	}

	got, err := ReconcileAll(sources, nil)
	if err != nil {
		t.Fatalf("ReconcileAll() error = %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("ReconcileAll() returned %d memories, want 2", len(got))
	}

	if got[0].ID == "" {
		t.Fatal("first memory has empty ID")
	}

	if got[1].ID != got[0].ID {
		t.Errorf(
			"second memory ID = %q, want earlier ID %q",
			got[1].ID,
			got[0].ID,
		)
	}

	if got[1].CreatedLine != 20 {
		t.Errorf("CreatedLine = %d, want 20", got[1].CreatedLine)
	}

	if got[1].CurrentLine != 40 {
		t.Errorf("CurrentLine = %d, want 40", got[1].CurrentLine)
	}
}

func TestReconcileUpdatesTextAtSameLocation(t *testing.T) {
	existing := []memory.Memory{
		{
			ID:          "memory-1",
			Kind:        syntax.Decision,
			Text:        "SQLite를 사용한다",
			CreatedPath: "demo.go",
			CreatedLine: 5,
			CurrentPath: "demo.go",
			CurrentLine: 5,
		},
	}

	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite WAL 모드를 사용한다",
		RelativePath: "demo.go",
		Line:         5,
	}

	got, err := Reconcile(source, existing)
	if err != nil {
		t.Fatalf("Reconcile() error = %v", err)
	}

	if got.ID != "memory-1" {
		t.Errorf("ID = %q, want %q", got.ID, "memory-1")
	}

	if got.Text != "SQLite WAL 모드를 사용한다" {
		t.Errorf(
			"Text = %q, want %q",
			got.Text,
			"SQLite WAL 모드를 사용한다",
		)
	}

	if got.CreatedPath != "demo.go" {
		t.Errorf("CreatedPath = %q, want %q", got.CreatedPath, "demo.go")
	}

	if got.CreatedLine != 5 {
		t.Errorf("CreatedLine = %d, want 5", got.CreatedLine)
	}

	if got.CurrentPath != "demo.go" {
		t.Errorf("CurrentPath = %q, want %q", got.CurrentPath, "demo.go")
	}

	if got.CurrentLine != 5 {
		t.Errorf("CurrentLine = %d, want 5", got.CurrentLine)
	}
}
