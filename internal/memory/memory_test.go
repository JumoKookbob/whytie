package memory

import (
	"testing"

	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestFromSourceCommentPreservesInitialLocation(t *testing.T) {
	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite",
		RelativePath: "internal/store/db.go",
		Line:         21,
	}

	got := FromSourceComment(source)

	if got.Kind != syntax.Decision {
		t.Errorf("Kind = %q, want %q", got.Kind, syntax.Decision)
	}

	if got.Text != "SQLite" {
		t.Errorf("Text = %q, want %q", got.Text, "SQLite")
	}

	if got.CreatedPath != "internal/store/db.go" {
		t.Errorf("CreatedPath = %q", got.CreatedPath)
	}

	if got.CreatedLine != 21 {
		t.Errorf("CreatedLine = %d, want 21", got.CreatedLine)
	}

	if got.CurrentPath != "internal/store/db.go" {
		t.Errorf("CurrentPath = %q", got.CurrentPath)
	}

	if got.CurrentLine != 21 {
		t.Errorf("CurrentLine = %d, want 21", got.CurrentLine)
	}
}

func TestNewAssignsStableIdentity(t *testing.T) {
	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite",
		RelativePath: "internal/store/db.go",
		Line:         21,
	}

	got, err := New(source)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if got.ID == "" {
		t.Fatal("New() ID is empty")
	}

	if got.Kind != syntax.Decision {
		t.Errorf("Kind = %q, want %q", got.Kind, syntax.Decision)
	}

	if got.CreatedPath != "internal/store/db.go" {
		t.Errorf("CreatedPath = %q", got.CreatedPath)
	}

	if got.CurrentPath != "internal/store/db.go" {
		t.Errorf("CurrentPath = %q", got.CurrentPath)
	}
}

func TestNewAssignsDifferentIDs(t *testing.T) {
	source := scanner.SourceComment{
		Kind:         syntax.Decision,
		Text:         "SQLite",
		RelativePath: "internal/store/db.go",
		Line:         21,
	}

	first, err := New(source)
	if err != nil {
		t.Fatalf("first New() error = %v", err)
	}

	second, err := New(source)
	if err != nil {
		t.Fatalf("second New() error = %v", err)
	}

	if first.ID == second.ID {
		t.Fatalf("New() generated duplicate ID %q", first.ID)
	}
}

func TestToSourceCommentUsesCurrentLocation(t *testing.T) {
	m := Memory{
		ID:          "memory-1",
		Kind:        syntax.Decision,
		Text:        "SQLite",
		CreatedPath: "internal/store/db.go",
		CreatedLine: 21,
		CurrentPath: "internal/store/db.go",
		CurrentLine: 40,
	}

	got := ToSourceComment(m)

	if got.Kind != syntax.Decision {
		t.Errorf("Kind = %q, want %q", got.Kind, syntax.Decision)
	}

	if got.Text != "SQLite" {
		t.Errorf("Text = %q, want %q", got.Text, "SQLite")
	}

	if got.RelativePath != "internal/store/db.go" {
		t.Errorf(
			"RelativePath = %q, want %q",
			got.RelativePath,
			"internal/store/db.go",
		)
	}

	if got.Line != 40 {
		t.Errorf("Line = %d, want 40", got.Line)
	}
}
