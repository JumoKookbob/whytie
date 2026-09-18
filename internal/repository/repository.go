package repository

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrAlreadyInitialized = errors.New("WhyTie repository already initialized")
	ErrNotRepository      = errors.New("not inside a WhyTie repository")
)

func Init(root string) error {
	dir := filepath.Join(root, ".whytie")

	info, err := os.Stat(dir)
	if err == nil {
		if info.IsDir() {
			return ErrAlreadyInitialized
		}

		return errors.New(".whytie exists and is not a directory")
	}

	if !os.IsNotExist(err) {
		return err
	}

	return os.Mkdir(dir, 0755)
}

func FindRoot(start string) (string, error) {
	current, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		whytieDir := filepath.Join(current, ".whytie")

		info, err := os.Stat(whytieDir)
		if err == nil && info.IsDir() {
			return current, nil
		}

		if err != nil && !os.IsNotExist(err) {
			return "", err
		}

		parent := filepath.Dir(current)

		if parent == current {
			return "", ErrNotRepository
		}

		current = parent
	}
}
