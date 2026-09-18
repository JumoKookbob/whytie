package app

import "github.com/JumoKookbob/whytie/internal/repository"

func Init(path string) error {
	return repository.Init(path)
}
