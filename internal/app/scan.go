package app

import (
	"io"

	"github.com/JumoKookbob/whytie/internal/formatter"
	"github.com/JumoKookbob/whytie/internal/gitinfo"
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

	// History records the repository state in which a reasoning event
	// was observed. This is intentionally the current HEAD commit rather
	// than git blame provenance for an individual source line.
	commitHash := ""

	if info, err := gitinfo.Current(path); err == nil {
		commitHash = info.CommitHash
	}

	provenance := func(_ string, _ int) string {
		return commitHash
	}

	if _, err := syncer.SyncWithProvenance(
		store,
		comments,
		provenance,
	); err != nil {
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
