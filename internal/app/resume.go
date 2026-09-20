package app

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/JumoKookbob/whytie/internal/formatter"
	"github.com/JumoKookbob/whytie/internal/history"
	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/reasoning"
	"github.com/JumoKookbob/whytie/internal/scanner"
)

type ResumeStore interface {
	List() ([]memory.Memory, error)
	ListRecentHistory(limit int) ([]history.Event, error)
}

func Resume(w io.Writer, store ResumeStore) error {
	events, err := store.ListRecentHistory(5)
	if err != nil {
		return fmt.Errorf("load recent history: %w", err)
	}

	memories, err := store.List()
	if err != nil {
		return fmt.Errorf("load memories: %w", err)
	}

	var out bytes.Buffer
	fmt.Fprintln(&out, "WhyTie Resume")
	fmt.Fprintln(&out, "Based on saved scans; source files have not been rescanned.")

	if len(events) == 0 {
		if len(memories) == 0 {
			fmt.Fprintln(&out, "\nNo saved reasoning yet.")
		} else {
			fmt.Fprintln(&out, "\nSaved reasoning exists, but no history is available.")
			fmt.Fprintln(&out, "View it with: whytie list")
		}
		fmt.Fprintln(&out, "Capture source changes with: whytie scan .")

		_, err := w.Write(out.Bytes())
		return err
	}

	fmt.Fprintln(&out, "\nRecent recorded events (newest saved first):")

	for _, event := range events {
		fmt.Fprintf(
			&out,
			"  %s  %s:%d\n",
			event.Type,
			event.Path,
			event.Line,
		)

		if event.Text != "" {
			fmt.Fprintf(&out, "    %s: %s\n", event.Kind, event.Text)
		} else {
			fmt.Fprintln(&out, "    No text snapshot available.")
		}

		if !event.OccurredAt.IsZero() {
			fmt.Fprintf(
				&out,
				"    recorded: %s\n",
				event.OccurredAt.UTC().Format("2006-01-02 15:04:05 UTC"),
			)
		}

		if event.CommitHash != "" {
			fmt.Fprintf(
				&out,
				"    commit: %s\n",
				shortCommitHash(event.CommitHash),
			)
		}
	}

	// Sort a copy so grouping does not mutate the store's slice.
	ordered := append([]memory.Memory(nil), memories...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].CurrentPath != ordered[j].CurrentPath {
			return ordered[i].CurrentPath < ordered[j].CurrentPath
		}
		return ordered[i].CurrentLine < ordered[j].CurrentLine
	})

	byID := make(map[string]memory.Memory, len(ordered))
	comments := make([]scanner.SourceComment, 0, len(ordered))

	for _, m := range ordered {
		byID[m.ID] = m
		comments = append(comments, memory.ToSourceComment(m))
	}

	blocks := reasoning.Group(comments)
	shown := make(map[int]bool)
	wroteContext := false

	for _, event := range events {
		// Deleted snapshots must not be presented as current context.
		if event.Type == history.EventDeleted {
			continue
		}

		current, ok := byID[event.MemoryID]
		if !ok {
			continue
		}

		for index, block := range blocks {
			if shown[index] {
				continue
			}

			matched := false
			for _, item := range block.Items {
				if item.RelativePath == current.CurrentPath &&
					item.Line == current.CurrentLine {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}

			if !wroteContext {
				fmt.Fprintln(&out, "\nRelated context from saved scans:")
				wroteContext = true
			}

			fmt.Fprintln(&out)
			formatter.WriteBlock(&out, block)
			inspectID := current.ID
			first := block.Items[0]

			for _, m := range ordered {
				if m.CurrentPath == first.RelativePath &&
					m.CurrentLine == first.Line {
					inspectID = m.ID
					break
				}
			}

			fmt.Fprintf(&out, "  Inspect: whytie why %s\n", inspectID)
			shown[index] = true
			break
		}
	}

	fmt.Fprintln(&out, "\nCapture source changes with: whytie scan .")

	_, err = w.Write(out.Bytes())
	return err
}
