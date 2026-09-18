package formatter

import (
	"fmt"
	"io"

	"github.com/JumoKookbob/whytie/internal/reasoning"
)

func WriteBlock(w io.Writer, block reasoning.Block) {
	items := reasoning.AttachReasons(block)

	for i, item := range items {
		if i > 0 {
			fmt.Fprintln(w)
		}

		comment := item.Comment

		fmt.Fprintf(
			w,
			"%s: %s  %s:%d\n",
			comment.Kind,
			comment.Text,
			comment.RelativePath,
			comment.Line,
		)

		for _, reason := range item.Reasons {
			fmt.Fprintf(
				w,
				"└─ reason: %s  %s:%d\n",
				reason.Text,
				reason.RelativePath,
				reason.Line,
			)
		}
	}
}
