package formatter

import (
	"bytes"
	"testing"

	"github.com/JumoKookbob/whytie/internal/reasoning"
	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func TestWriteBlock(t *testing.T) {
	block := reasoning.Block{
		Items: []scanner.SourceComment{
			{
				Kind:         syntax.Question,
				Text:         "어떤 DB를 쓸까?",
				RelativePath: "internal/store/db.go",
				Line:         18,
			},
			{
				Kind:         syntax.Decision,
				Text:         "SQLite",
				RelativePath: "internal/store/db.go",
				Line:         19,
			},
			{
				Kind:         syntax.Reason,
				Text:         "local-first에 적합함",
				RelativePath: "internal/store/db.go",
				Line:         20,
			},
		},
	}

	var buf bytes.Buffer

	WriteBlock(&buf, block)

	want := "" +
		"question: 어떤 DB를 쓸까?  internal/store/db.go:18\n" +
		"\n" +
		"decision: SQLite  internal/store/db.go:19\n" +
		"└─ reason: local-first에 적합함  internal/store/db.go:20\n"

	if got := buf.String(); got != want {
		t.Errorf(
			"WriteBlock() output mismatch\n\ngot:\n%s\nwant:\n%s",
			got,
			want,
		)
	}
}
