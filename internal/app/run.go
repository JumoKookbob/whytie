package app

import (
	"errors"
	"fmt"
	"io"

	"github.com/JumoKookbob/whytie/internal/repository"
)

func RunList(w io.Writer, root string) error {
	store, err := OpenStore(root)
	if err != nil {
		return err
	}
	defer store.Close()

	return List(w, store)
}

func RunScan(w io.Writer, root string, path string) error {
	store, err := OpenStore(root)
	if err != nil {
		return err
	}
	defer store.Close()

	return Scan(w, store, path)
}

func RunListFrom(w io.Writer, start string) error {
	root, err := repository.FindRoot(start)
	if err != nil {
		return err
	}

	return RunList(w, root)
}

func RunScanFrom(w io.Writer, start string, path string) error {
	root, err := repository.FindRoot(start)
	if err != nil {
		return err
	}

	return RunScan(w, root, path)
}

func Run(args []string, stdout io.Writer, stderr io.Writer, start string) int {
	if len(args) == 0 {
		fmt.Fprintln(stdout, "WhyTie")
		return 0
	}

	switch args[0] {
	case "init":
		if err := Init(start); err != nil {
			if errors.Is(err, repository.ErrAlreadyInitialized) {
				fmt.Fprintln(stderr, "WhyTie repository already initialized")
				return 1
			}

			fmt.Fprintf(stderr, "init failed: %v\n", err)
			return 1
		}

		fmt.Fprintln(stdout, "Initialized WhyTie repository in .whytie")

		return 0

	case "list":
		if err := RunListFrom(stdout, start); err != nil {
			fmt.Fprintf(stderr, "list failed: %v\n", err)
			return 1
		}

		return 0

	case "scan":
		if len(args) < 2 {
			fmt.Fprintln(stderr, "usage: whytie scan <path>")
			return 1
		}

		if err := RunScanFrom(stdout, start, args[1]); err != nil {
			fmt.Fprintf(stderr, "scan failed: %v\n", err)
			return 1
		}

		return 0

	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		return 1
	}
}
