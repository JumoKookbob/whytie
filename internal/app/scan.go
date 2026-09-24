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

	currentFile := ""

	for _, block := range blocks {
		items := reasoning.AttachReasons(block)
		if len(items) == 0 {
			continue
		}

		path := items[0].Comment.RelativePath

		if path != currentFile {
			if currentFile != "" {
				if _, err := io.WriteString(w, "\n"); err != nil {
					return err
				}
			}

			formatter.WriteFileHeader(w, path)
			currentFile = path
		} else {
			if _, err := io.WriteString(
				w,
				"\n"+formatter.BlockSeparator+"\n\n",
			); err != nil {
				return err
			}
		}

		formatter.WriteBlockBody(w, block)
	}

	return nil
}
