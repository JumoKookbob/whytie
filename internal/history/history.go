package history

import (
	"time"

	"github.com/JumoKookbob/whytie/internal/memory"
)

type EventType string

const (
	EventCreated EventType = "created"
	EventMoved   EventType = "moved"
	EventChanged EventType = "changed"
)

type Event struct {
	MemoryID string
	Type     EventType

	Path string
	Line int

	CommitHash string
	OccurredAt time.Time
}

func Detect(previous, current memory.Memory) (EventType, bool) {
	if previous.ID != current.ID {
		return "", false
	}

	if previous.Text != current.Text || previous.Kind != current.Kind {
		return EventChanged, true
	}

	if previous.CurrentPath != current.CurrentPath {
		return EventMoved, true
	}

	if previous.CurrentLine != current.CurrentLine {
		return EventMoved, true
	}

	return "", false
}
