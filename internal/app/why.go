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

	if target.Kind == syntax.Question {
		decision, ok := FindDecision(memories, id)
		if !ok {
			return nil
		}

		fmt.Fprintf(
			w,
			"└─ %s: %s  %s:%d\n",
			decision.Kind,
			decision.Text,
			decision.CurrentPath,
			decision.CurrentLine,
		)

		for _, reason := range FindReasons(memories, decision.ID) {
			fmt.Fprintf(
				w,
				"   └─ %s: %s  %s:%d\n",
				reason.Kind,
				reason.Text,
				reason.CurrentPath,
				reason.CurrentLine,
			)
		}

		return nil
	}

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

func FindDecision(memories []memory.Memory, id string) (memory.Memory, bool) {
	for i, m := range memories {
		if m.ID != id {
			continue
		}

		if m.Kind != syntax.Question {
			return memory.Memory{}, false
		}

		if i+1 >= len(memories) {
			return memory.Memory{}, false
		}

		candidate := memories[i+1]

		if candidate.CurrentPath != m.CurrentPath {
			return memory.Memory{}, false
		}

		if candidate.CurrentLine != m.CurrentLine+1 {
			return memory.Memory{}, false
		}

		if candidate.Kind != syntax.Decision {
			return memory.Memory{}, false
		}

		return candidate, true
	}

	return memory.Memory{}, false
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
