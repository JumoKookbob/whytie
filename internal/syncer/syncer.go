package syncer

import (
	"time"

	"github.com/JumoKookbob/whytie/internal/history"
	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/reconcile"
	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/storage"
)

// ProvenanceFunc returns the Git commit hash associated with a source
// location. An empty string means provenance is unavailable.
type ProvenanceFunc func(path string, line int) string

// Sync synchronizes scanned source comments with the persistent memory
// store.
//
// It intentionally performs no Git lookup so WhyTie continues to work
// outside Git repositories.
func Sync(
	store storage.Store,
	sources []scanner.SourceComment,
) ([]memory.Memory, error) {
	return SyncWithProvenance(store, sources, nil)
}

// SyncWithProvenance synchronizes scanned source comments with the
// persistent memory store and optionally records source provenance in
// history events.
func SyncWithProvenance(
	store storage.Store,
	sources []scanner.SourceComment,
	provenance ProvenanceFunc,
) ([]memory.Memory, error) {
	existing, err := store.List()
	if err != nil {
		return nil, err
	}

	reconciled, err := reconcile.ReconcileAll(sources, existing)
	if err != nil {
		return nil, err
	}

	existingByID := make(map[string]memory.Memory, len(existing))

	for _, m := range existing {
		existingByID[m.ID] = m
	}

	activeIDs := make(map[string]struct{}, len(reconciled))

	for _, m := range reconciled {
		activeIDs[m.ID] = struct{}{}
	}

	for _, m := range existing {
		if _, ok := activeIDs[m.ID]; ok {
			continue
		}

		if err := store.Delete(m.ID); err != nil {
			return nil, err
		}
	}

	for _, m := range reconciled {
		if err := store.Save(m); err != nil {
			return nil, err
		}
	}

	historyStore, ok := store.(storage.HistoryStore)
	if !ok {
		return reconciled, nil
	}

	now := time.Now().UTC()

	for _, current := range reconciled {
		previous, existed := existingByID[current.ID]

		if !existed {
			event := history.Event{
				MemoryID: current.ID,
				Type:     history.EventCreated,
				Path:     current.CurrentPath,
				Line:     current.CurrentLine,
				CommitHash: resolveProvenance(
					provenance,
					current.CurrentPath,
					current.CurrentLine,
				),
				OccurredAt: now,
			}

			if err := historyStore.SaveHistory(event); err != nil {
				return nil, err
			}

			continue
		}

		eventType, changed := history.Detect(previous, current)
		if !changed {
			continue
		}

		event := history.Event{
			MemoryID: current.ID,
			Type:     eventType,
			Path:     current.CurrentPath,
			Line:     current.CurrentLine,
			CommitHash: resolveProvenance(
				provenance,
				current.CurrentPath,
				current.CurrentLine,
			),
			OccurredAt: now,
		}

		if err := historyStore.SaveHistory(event); err != nil {
			return nil, err
		}
	}

	return reconciled, nil
}

func resolveProvenance(
	provenance ProvenanceFunc,
	path string,
	line int,
) string {
	if provenance == nil {
		return ""
	}

	return provenance(path, line)
}
