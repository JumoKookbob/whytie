package app

import (
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

	for i, block := range blocks {
		if i > 0 {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}

		formatter.WriteBlock(w, block)
	}

	return nil
}
