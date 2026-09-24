package formatter

import (
	"fmt"
	"io"

	"github.com/JumoKookbob/whytie/internal/reasoning"
	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

const (
	FileSeparator  = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	BlockSeparator = "────────────────────────────────────────────────────"
)

// WriteFileHeader writes the shared WhyTie file heading.
func WriteFileHeader(w io.Writer, path string) {
	fmt.Fprintf(w, "📄 %s\n", path)
	fmt.Fprintln(w, FileSeparator)
	fmt.Fprintln(w)
}

// WriteBlock writes one contiguous reasoning block.
//
// Separators between different blocks are intentionally handled by the
// caller because a Block already represents adjacent reasoning comments.
func WriteBlock(w io.Writer, block reasoning.Block) {
	items := reasoning.AttachReasons(block)
	if len(items) == 0 {
		return
	}

	WriteFileHeader(w, items[0].Comment.RelativePath)

	WriteBlockBody(w, block)
}

// WriteBlockBody writes a reasoning block without a file header.
// This is useful when several blocks from the same file are displayed.
func WriteBlockBody(w io.Writer, block reasoning.Block) {
	items := reasoning.AttachReasons(block)

	for i, item := range items {
		if i > 0 {
			fmt.Fprintln(w)
		}

		WriteComment(w, item.Comment, false)

		for _, reason := range item.Reasons {
			WriteComment(w, reason, true)
		}
	}
}

// WriteComment writes one reasoning annotation using the shared WhyTie UI.
func WriteComment(
	w io.Writer,
	comment scanner.SourceComment,
	child bool,
) {
	icon, label := kindStyle(comment.Kind)

	if child {
		fmt.Fprintf(
			w,
			"  %3d  └─ %s %-9s %s\n",
			comment.Line,
			icon,
			label,
			comment.Text,
		)
		return
	}

	fmt.Fprintf(
		w,
		"  %3d  %s %-9s %s\n",
		comment.Line,
		icon,
		label,
		comment.Text,
	)
}

func kindStyle(kind syntax.Kind) (icon string, label string) {
	switch kind {
	case syntax.Question:
		return "❓", "QUESTION"
	case syntax.Decision:
		return "✅", "DECISION"
	case syntax.Rejected:
		return "⛔", "REJECTED"
	case syntax.Failed:
		return "❌", "FAILED"
	case syntax.Reason:
		return "💡", "REASON"
	case syntax.Important:
		return "❗", "IMPORTANT"
	default:
		return "•", string(kind)
	}
}
