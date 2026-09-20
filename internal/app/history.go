package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/JumoKookbob/whytie/internal/history"
	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/storage"
)

func History(w io.Writer, store storage.HistoryStore, m memory.Memory) error {
	events, err := store.ListHistory(m.ID)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		w,
		"%s: %s  %s:%d\n",
		m.Kind,
		m.Text,
		m.CurrentPath,
		m.CurrentLine,
	)

	if len(events) == 0 {
		fmt.Fprintln(w, "\nhistory: none")
		return nil
	}

	fmt.Fprintln(w, "\nhistory:")

	for _, event := range events {
		writeHistoryEvent(w, event)
	}

	return nil
}

func HistoryByID(
	w io.Writer,
	store storage.Store,
	historyStore storage.HistoryStore,
	id string,
) error {
	target, err := store.Get(id)
	if err != nil {
		return err
	}

	if target.ID == "" {
		return fmt.Errorf("memory not found: %s", id)
	}

	return History(w, historyStore, target)
}

func HistoryAt(
	w io.Writer,
	store storage.Store,
	historyStore storage.HistoryStore,
	location string,
) error {
	path, line, err := ParseLocation(location)
	if err != nil {
		return err
	}

	memories, err := store.List()
	if err != nil {
		return err
	}

	target, ok := FindMemoryAt(memories, path, line)
	if !ok {
		return fmt.Errorf("no memory at %s", location)
	}

	return History(w, historyStore, target)
}

func writeHistoryEvent(w io.Writer, event history.Event) {
	fmt.Fprintf(
		w,
		"  %-7s  %s:%d\n",
		event.Type,
		event.Path,
		event.Line,
	)

	if event.CommitHash != "" {
		fmt.Fprintf(
			w,
			"           commit %s\n",
			shortCommitHash(event.CommitHash),
		)
	}
}

func shortCommitHash(hash string) string {
	hash = strings.TrimSpace(hash)

	if len(hash) <= 7 {
		return hash
	}

	return hash[:7]
}
