package app

import (
	"errors"
	"fmt"
	"io"
	"strings"

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

func RunHistory(w io.Writer, root string, target string) error {
	store, err := OpenStore(root)
	if err != nil {
		return err
	}
	defer store.Close()

	if strings.Contains(target, ":") {
		if _, _, err := ParseLocation(target); err != nil {
			return err
		}

		return HistoryAt(w, store, store, target)
	}

	return HistoryByID(w, store, store, target)
}

func RunHistoryFrom(w io.Writer, start string, target string) error {
	root, err := repository.FindRoot(start)
	if err != nil {
		return err
	}

	return RunHistory(w, root, target)
}

func Run(args []string, stdout io.Writer, stderr io.Writer, start string) int {
	if len(args) == 0 {
		fmt.Fprintln(stdout, "WhyTie")
		return 0
	}

	switch args[0] {
	case "version":
		fmt.Fprintln(stdout, "WhyTie v1.0.0")
		return 0

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

	case "resume":
		if len(args) != 1 {
			fmt.Fprintln(stderr, "usage: whytie resume")
			return 1
		}

		root, err := repository.FindRoot(start)
		if err != nil {
			fmt.Fprintf(stderr, "resume failed: %v\n", err)
			return 1
		}

		store, err := OpenStore(root)
		if err != nil {
			fmt.Fprintf(stderr, "resume failed: %v\n", err)
			return 1
		}
		defer store.Close()

		if err := Resume(stdout, store); err != nil {
			fmt.Fprintf(stderr, "resume failed: %v\n", err)
			return 1
		}

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

	case "history":
		if len(args) < 2 {
			fmt.Fprintln(stderr, "usage: whytie history <memory-id|file:line>")
			return 1
		}

		if err := RunHistoryFrom(stdout, start, args[1]); err != nil {
			fmt.Fprintf(stderr, "history failed: %v\n", err)
			return 1
		}

		return 0

	case "why":
		if len(args) < 2 {
			fmt.Fprintln(stderr, "usage: whytie why <memory-id|file:line>")
			return 1
		}

		root, err := repository.FindRoot(start)
		if err != nil {
			fmt.Fprintf(stderr, "why failed: %v\n", err)
			return 1
		}

		store, err := OpenStore(root)
		if err != nil {
			fmt.Fprintf(stderr, "why failed: %v\n", err)
			return 1
		}
		defer store.Close()

		target := args[1]

		if strings.Contains(target, ":") {
			if _, _, parseErr := ParseLocation(target); parseErr != nil {
				fmt.Fprintf(stderr, "why failed: %v\n", parseErr)
				return 1
			}

			err = WhyAt(stdout, store, target)
		} else {
			err = Why(stdout, store, target)
		}
		if err != nil {
			fmt.Fprintf(stderr, "why failed: %v\n", err)
			return 1
		}

		return 0

	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		return 1
	}
}
