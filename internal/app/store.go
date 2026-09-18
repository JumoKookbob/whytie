package app

import (
	"github.com/JumoKookbob/whytie/internal/storage"
)

func OpenStore(root string) (*storage.SQLiteStore, error) {
	return storage.OpenSQLite(root)
}
