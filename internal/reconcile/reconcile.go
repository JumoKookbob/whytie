package reconcile

import (
	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/scanner"
)

func MatchExact(
	source scanner.SourceComment,
	existing []memory.Memory,
) (memory.Memory, bool) {
	for _, candidate := range existing {
		if candidate.Kind != source.Kind {
			continue
		}

		if candidate.Text != source.Text {
			continue
		}

		if candidate.CurrentPath != source.RelativePath {
			continue
		}

		if candidate.CurrentLine != source.Line {
			continue
		}

		return candidate, true
	}

	return memory.Memory{}, false
}

func MatchEdited(
	source scanner.SourceComment,
	existing []memory.Memory,
) (memory.Memory, bool) {
	var match memory.Memory
	count := 0

	for _, candidate := range existing {
		if candidate.Kind != source.Kind {
			continue
		}

		if candidate.CurrentPath != source.RelativePath {
			continue
		}

		if candidate.CurrentLine != source.Line {
			continue
		}

		if candidate.Text == source.Text {
			continue
		}

		match = candidate
		count++

		if count > 1 {
			return memory.Memory{}, false
		}
	}

	if count != 1 {
		return memory.Memory{}, false
	}

	return match, true
}

func MatchMoved(
	source scanner.SourceComment,
	existing []memory.Memory,
) (memory.Memory, bool) {
	var match memory.Memory
	count := 0

	for _, candidate := range existing {
		if candidate.Kind != source.Kind {
			continue
		}

		if candidate.Text != source.Text {
			continue
		}

		if candidate.CurrentPath != source.RelativePath {
			continue
		}

		if candidate.CurrentLine == source.Line {
			continue
		}

		match = candidate
		count++

		if count > 1 {
			return memory.Memory{}, false
		}
	}

	if count != 1 {
		return memory.Memory{}, false
	}

	return match, true
}

func UpdateLocation(
	existing memory.Memory,
	source scanner.SourceComment,
) memory.Memory {
	existing.CurrentPath = source.RelativePath
	existing.CurrentLine = source.Line

	return existing
}

func UpdateContent(
	existing memory.Memory,
	source scanner.SourceComment,
) memory.Memory {
	existing.Text = source.Text
	existing.CurrentPath = source.RelativePath
	existing.CurrentLine = source.Line

	return existing
}

func Reconcile(
	source scanner.SourceComment,
	existing []memory.Memory,
) (memory.Memory, error) {
	if matched, ok := MatchExact(source, existing); ok {
		return matched, nil
	}

	if matched, ok := MatchEdited(source, existing); ok {
		return UpdateContent(matched, source), nil
	}

	if matched, ok := MatchMoved(source, existing); ok {
		return UpdateLocation(matched, source), nil
	}

	return memory.New(source)
}

func ReconcileAll(
	sources []scanner.SourceComment,
	existing []memory.Memory,
) ([]memory.Memory, error) {
	result := make([]memory.Memory, 0, len(sources))

	known := make([]memory.Memory, len(existing))
	copy(known, existing)

	for _, source := range sources {
		reconciled, err := Reconcile(source, known)
		if err != nil {
			return nil, err
		}

		result = append(result, reconciled)

		found := false

		for i := range known {
			if known[i].ID == reconciled.ID {
				known[i] = reconciled
				found = true
				break
			}
		}

		if !found {
			known = append(known, reconciled)
		}
	}

	return result, nil
}
