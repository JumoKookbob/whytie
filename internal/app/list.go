package app

import (
	"fmt"
	"io"

	"github.com/JumoKookbob/whytie/internal/formatter"
	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/reasoning"
	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/storage"
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

			formatter.WriteFileHeader(w, path)
			currentFile = path
		} else {
			if _, err := fmt.Fprintf(
				w,
				"\n%s\n\n",
				formatter.BlockSeparator,
			); err != nil {
				return err
			}
		}

		formatter.WriteBlockBody(w, block)
	}

	return nil
}
