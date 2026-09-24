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
				Text:         "Which database should we use?",
				RelativePath: "internal/store/db.go",
				Line:         18,
			},
			{
				Kind:         syntax.Decision,
				Text:         "Use SQLite",
				RelativePath: "internal/store/db.go",
				Line:         19,
			},
			{
				Kind:         syntax.Reason,
				Text:         "It fits the local-first design",
				RelativePath: "internal/store/db.go",
				Line:         20,
			},
		},
	}

	var buf bytes.Buffer

	WriteBlock(&buf, block)

	want := "" +
		"📄 internal/store/db.go\n" +
		"━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n" +
		"\n" +
		"   18  ❓ QUESTION  Which database should we use?\n" +
		"\n" +
		"   19  ✅ DECISION  Use SQLite\n" +
		"   20  └─ 💡 REASON    It fits the local-first design\n"

	if got := buf.String(); got != want {
		t.Errorf(
			"WriteBlock() output mismatch\n\ngot:\n%s\nwant:\n%s",
			got,
			want,
		)
	}
}

func TestWriteComment(t *testing.T) {
	comment := scanner.SourceComment{
		Kind:         syntax.Important,
		Text:         "Do not remove this compatibility path",
		RelativePath: "compat.go",
		Line:         42,
	}

	var buf bytes.Buffer

	WriteComment(&buf, comment, false)

	want := "   42  ❗ IMPORTANT Do not remove this compatibility path\n"

	if got := buf.String(); got != want {
		t.Errorf(
			"WriteComment() output mismatch\n\ngot:\n%q\nwant:\n%q",
			got,
			want,
		)
	}
}
