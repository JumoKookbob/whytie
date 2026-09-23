package app

import (
	"fmt"
	"io"
	"os"

	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/reasoning"
	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/storage"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

const (
	listFileSeparator  = "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	listBlockSeparator = "────────────────────────────────────────────────────"

	colorReset   = "\033[0m"
	colorBold    = "\033[1m"
	colorGray    = "\033[90m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
)

func List(w io.Writer, store storage.Store) error {
	memories, err := store.List()
	if err != nil {
		return err
	}

	comments := make([]scanner.SourceComment, 0, len(memories))

	for _, m := range memories {
		comments = append(comments, memory.ToSourceComment(m))
	}

	blocks := reasoning.Group(comments)
	color := supportsColor(w)

	currentFile := ""

	for _, block := range blocks {
		items := reasoning.AttachReasons(block)
		if len(items) == 0 {
			continue
		}

		path := items[0].Comment.RelativePath

		if path != currentFile {
			if currentFile != "" {
				if _, err := fmt.Fprintln(w); err != nil {
					return err
				}
			}

			if color {
				if _, err := fmt.Fprintf(
					w,
					"%s📄 %s%s\n%s%s%s\n\n",
					colorBold,
					path,
					colorReset,
					colorGray,
					listFileSeparator,
					colorReset,
				); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(
					w,
					"📄 %s\n%s\n\n",
					path,
					listFileSeparator,
				); err != nil {
					return err
				}
			}

			currentFile = path
		} else {
			if color {
				if _, err := fmt.Fprintf(
					w,
					"\n%s%s%s\n\n",
					colorGray,
					listBlockSeparator,
					colorReset,
				); err != nil {
					return err
				}
			} else {
				if _, err := fmt.Fprintf(
					w,
					"\n%s\n\n",
					listBlockSeparator,
				); err != nil {
					return err
				}
			}
		}

		for itemIndex, item := range items {
			if itemIndex > 0 {
				if _, err := fmt.Fprintln(w); err != nil {
					return err
				}
			}

			if err := writeListComment(
				w,
				item.Comment,
				false,
				color,
			); err != nil {
				return err
			}

			for _, reason := range item.Reasons {
				if err := writeListComment(
					w,
					reason,
					true,
					color,
				); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func writeListComment(
	w io.Writer,
	comment scanner.SourceComment,
	child bool,
	color bool,
) error {
	icon, label, ansiColor := listKindStyle(comment.Kind)

	if !color {
		if child {
			_, err := fmt.Fprintf(
				w,
				"  %3d  └─ %s %-9s %s\n",
				comment.Line,
				icon,
				label,
				comment.Text,
			)
			return err
		}

		_, err := fmt.Fprintf(
			w,
			"  %3d  %s %-9s %s\n",
			comment.Line,
			icon,
			label,
			comment.Text,
		)
		return err
	}

	if child {
		_, err := fmt.Fprintf(
			w,
			"  %3d  └─ %s%s %-9s%s %s\n",
			comment.Line,
			ansiColor,
			icon,
			label,
			colorReset,
			comment.Text,
		)
		return err
	}

	_, err := fmt.Fprintf(
		w,
		"  %3d  %s%s %-9s%s %s\n",
		comment.Line,
		ansiColor,
		icon,
		label,
		colorReset,
		comment.Text,
	)

	return err
}

func listKindStyle(kind syntax.Kind) (icon string, label string, color string) {
	switch kind {
	case syntax.Question:
		return "❓", "QUESTION", colorYellow

	case syntax.Decision:
		return "✅", "DECISION", colorGreen

	case syntax.Rejected:
		return "⛔", "REJECTED", colorYellow

	case syntax.Failed:
		return "❌", "FAILED", colorRed

	case syntax.Reason:
		return "💡", "REASON", colorBlue

	case syntax.Important:
		return "⚠️", "IMPORTANT", colorMagenta

	default:
		return "•", string(kind), colorReset
	}
}

func supportsColor(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}

	info, err := file.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}
