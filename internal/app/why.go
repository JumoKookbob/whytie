package app

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/storage"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

func FindReasons(memories []memory.Memory, id string) []memory.Memory {
	for i, m := range memories {
		if m.ID != id {
			continue
		}

		if m.Kind != syntax.Decision {
			return nil
		}

		var reasons []memory.Memory

		for j := i + 1; j < len(memories); j++ {
			candidate := memories[j]

			if candidate.CurrentPath != m.CurrentPath {
				break
			}

			if candidate.CurrentLine != memories[j-1].CurrentLine+1 {
				break
			}

			if candidate.Kind != syntax.Reason {
				break
			}

			reasons = append(reasons, candidate)
		}

		return reasons
	}

	return nil
}

func FindMemoryAt(memories []memory.Memory, path string, line int) (memory.Memory, bool) {
	for _, m := range memories {
		if m.CurrentPath == path && m.CurrentLine == line {
			return m, true
		}
	}

	return memory.Memory{}, false
}

func Why(w io.Writer, store storage.Store, id string) error {
	target, err := store.Get(id)
	if err != nil {
		return err
	}

	memories, err := store.List()
	if err != nil {
		return err
	}

	fmt.Fprintf(
		w,
		"%s: %s  %s:%d\n",
		target.Kind,
		target.Text,
		target.CurrentPath,
		target.CurrentLine,
	)

	for _, reason := range FindReasons(memories, id) {
		fmt.Fprintf(
			w,
			"└─ %s: %s  %s:%d\n",
			reason.Kind,
			reason.Text,
			reason.CurrentPath,
			reason.CurrentLine,
		)
	}

	return nil
}

func ParseLocation(location string) (string, int, error) {
	index := strings.LastIndex(location, ":")
	if index == -1 {
		return "", 0, fmt.Errorf("invalid location: %s", location)
	}

	path := location[:index]
	lineText := location[index+1:]

	if path == "" || lineText == "" {
		return "", 0, fmt.Errorf("invalid location: %s", location)
	}

	line, err := strconv.Atoi(lineText)
	if err != nil {
		return "", 0, fmt.Errorf("invalid line number: %s", lineText)
	}

	if line <= 0 {
		return "", 0, fmt.Errorf("invalid line number: %d", line)
	}

	return path, line, nil
}

func WhyAt(w io.Writer, store storage.Store, location string) error {
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

	return Why(w, store, target.ID)
}
