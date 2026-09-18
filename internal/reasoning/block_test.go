package reasoning

import (
	"testing"

	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestGroupAdjacentComments(t *testing.T) {
	comments := []scanner.SourceComment{
		{
			Kind:         syntax.Question,
			Text:         "persistence를 뭘로 할까?",
			File:         `C:\project\store.go`,
			RelativePath: "store.go",
			Line:         10,
		},
		{
			Kind:         syntax.Rejected,
			Text:         "PostgreSQL",
			File:         `C:\project\store.go`,
			RelativePath: "store.go",
			Line:         11,
		},
		{
			Kind:         syntax.Reason,
			Text:         "서버 운영이 필요함",
			File:         `C:\project\store.go`,
			RelativePath: "store.go",
			Line:         12,
		},
		{
			Kind:         syntax.Decision,
			Text:         "SQLite",
			File:         `C:\project\store.go`,
			RelativePath: "store.go",
			Line:         13,
		},
		{
			Kind:         syntax.Reason,
			Text:         "local-first라 적합함",
			File:         `C:\project\store.go`,
			RelativePath: "store.go",
			Line:         14,
		},
	}

	blocks := Group(comments)

	if len(blocks) != 1 {
		t.Fatalf("Group() returned %d blocks, want 1", len(blocks))
	}

	if len(blocks[0].Items) != 5 {
		t.Fatalf(
			"block contains %d items, want 5",
			len(blocks[0].Items),
		)
	}
}

func TestGroupSeparatesNonAdjacentComments(t *testing.T) {
	comments := []scanner.SourceComment{
		{
			Kind:         syntax.Question,
			Text:         "cache를 바꿀까?",
			File:         `C:\project\cache.go`,
			RelativePath: "cache.go",
			Line:         10,
		},
		{
			Kind:         syntax.Decision,
			Text:         "per-CPU cache",
			File:         `C:\project\cache.go`,
			RelativePath: "cache.go",
			Line:         15,
		},
	}

	blocks := Group(comments)

	if len(blocks) != 2 {
		t.Fatalf("Group() returned %d blocks, want 2", len(blocks))
	}

	if len(blocks[0].Items) != 1 {
		t.Errorf("first block contains %d items, want 1", len(blocks[0].Items))
	}

	if len(blocks[1].Items) != 1 {
		t.Errorf("second block contains %d items, want 1", len(blocks[1].Items))
	}
}

func TestGroupSeparatesDifferentFiles(t *testing.T) {
	comments := []scanner.SourceComment{
		{
			Kind:         syntax.Question,
			Text:         "cache를 바꿀까?",
			File:         `C:\project\cache.go`,
			RelativePath: "cache.go",
			Line:         10,
		},
		{
			Kind:         syntax.Decision,
			Text:         "per-CPU cache",
			File:         `C:\project\other.go`,
			RelativePath: "other.go",
			Line:         11,
		},
	}

	blocks := Group(comments)

	if len(blocks) != 2 {
		t.Fatalf("Group() returned %d blocks, want 2", len(blocks))
	}
}

func TestGroupEmptyComments(t *testing.T) {
	blocks := Group(nil)

	if len(blocks) != 0 {
		t.Fatalf("Group(nil) returned %d blocks, want 0", len(blocks))
	}
}

func TestAttachReasons(t *testing.T) {
	block := Block{
		Items: []scanner.SourceComment{
			{
				Kind: syntax.Failed,
				Text: "global cache",
				Line: 10,
			},
			{
				Kind: syntax.Reason,
				Text: "lock contention이 너무 큼",
				Line: 11,
			},
			{
				Kind: syntax.Decision,
				Text: "per-CPU cache",
				Line: 12,
			},
			{
				Kind: syntax.Reason,
				Text: "contention 감소",
				Line: 13,
			},
		},
	}

	items := AttachReasons(block)

	if len(items) != 2 {
		t.Fatalf("AttachReasons() returned %d items, want 2", len(items))
	}

	if items[0].Comment.Kind != syntax.Failed {
		t.Errorf("first item kind = %q, want %q", items[0].Comment.Kind, syntax.Failed)
	}

	if len(items[0].Reasons) != 1 {
		t.Fatalf("first item has %d reasons, want 1", len(items[0].Reasons))
	}

	if items[0].Reasons[0].Text != "lock contention이 너무 큼" {
		t.Errorf(
			"first reason = %q, want %q",
			items[0].Reasons[0].Text,
			"lock contention이 너무 큼",
		)
	}

	if items[1].Comment.Kind != syntax.Decision {
		t.Errorf("second item kind = %q, want %q", items[1].Comment.Kind, syntax.Decision)
	}

	if len(items[1].Reasons) != 1 {
		t.Fatalf("second item has %d reasons, want 1", len(items[1].Reasons))
	}

	if items[1].Reasons[0].Text != "contention 감소" {
		t.Errorf(
			"second reason = %q, want %q",
			items[1].Reasons[0].Text,
			"contention 감소",
		)
	}
}

func TestAttachReasonsPreservesOrphanReason(t *testing.T) {
	block := Block{
		Items: []scanner.SourceComment{
			{
				Kind: syntax.Reason,
				Text: "benchmark에서 30% 느렸음",
				Line: 10,
			},
			{
				Kind: syntax.Decision,
				Text: "SQLite",
				Line: 11,
			},
		},
	}

	items := AttachReasons(block)

	if len(items) != 2 {
		t.Fatalf("AttachReasons() returned %d items, want 2", len(items))
	}

	if items[0].Comment.Kind != syntax.Reason {
		t.Errorf(
			"first item kind = %q, want %q",
			items[0].Comment.Kind,
			syntax.Reason,
		)
	}

	if items[0].Comment.Text != "benchmark에서 30% 느렸음" {
		t.Errorf(
			"first item text = %q, want %q",
			items[0].Comment.Text,
			"benchmark에서 30% 느렸음",
		)
	}

	if items[1].Comment.Kind != syntax.Decision {
		t.Errorf(
			"second item kind = %q, want %q",
			items[1].Comment.Kind,
			syntax.Decision,
		)
	}
}

func TestAttachReasonsPreservesSourceLocation(t *testing.T) {
	block := Block{
		Items: []scanner.SourceComment{
			{
				Kind:         syntax.Decision,
				Text:         "SQLite",
				File:         `C:\project\internal\store\db.go`,
				RelativePath: "internal/store/db.go",
				Line:         21,
			},
			{
				Kind:         syntax.Reason,
				Text:         "local-first에 적합함",
				File:         `C:\project\internal\store\db.go`,
				RelativePath: "internal/store/db.go",
				Line:         22,
			},
		},
	}

	items := AttachReasons(block)

	if len(items) != 1 {
		t.Fatalf("AttachReasons() returned %d items, want 1", len(items))
	}

	decision := items[0].Comment

	if decision.RelativePath != "internal/store/db.go" {
		t.Errorf(
			"decision path = %q, want %q",
			decision.RelativePath,
			"internal/store/db.go",
		)
	}

	if decision.Line != 21 {
		t.Errorf("decision line = %d, want 21", decision.Line)
	}

	if len(items[0].Reasons) != 1 {
		t.Fatalf(
			"decision has %d reasons, want 1",
			len(items[0].Reasons),
		)
	}

	reason := items[0].Reasons[0]

	if reason.RelativePath != "internal/store/db.go" {
		t.Errorf(
			"reason path = %q, want %q",
			reason.RelativePath,
			"internal/store/db.go",
		)
	}

	if reason.Line != 22 {
		t.Errorf("reason line = %d, want 22", reason.Line)
	}
}
