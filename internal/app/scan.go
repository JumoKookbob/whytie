package app

import (
	"io"

	"github.com/JumoKookbob/whytie/internal/formatter"
	"github.com/JumoKookbob/whytie/internal/reasoning"
	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/storage"
	"github.com/JumoKookbob/whytie/internal/syncer"
)

func Scan(w io.Writer, store storage.Store, path string) error {
	comments, err := scanner.ScanDir(path)
	if err != nil {
		return err
	}

	if _, err := syncer.Sync(store, comments); err != nil {
		return err
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
