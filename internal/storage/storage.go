package storage

import "github.com/JumoKookbob/whytie/internal/memory"

type Store interface {
	Save(memory.Memory) error
	Get(id string) (memory.Memory, error)
	List() ([]memory.Memory, error)
	Delete(id string) error
}
