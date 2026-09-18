package syncer

import (
	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/reconcile"
	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/storage"
)

func Sync(
	store storage.Store,
	sources []scanner.SourceComment,
) ([]memory.Memory, error) {
	existing, err := store.List()
	if err != nil {
		return nil, err
	}

	reconciled, err := reconcile.ReconcileAll(sources, existing)
	if err != nil {
		return nil, err
	}

	for _, m := range reconciled {
		if err := store.Save(m); err != nil {
			return nil, err
		}
	}

	return reconciled, nil
}
