package storage

import (
	"fmt"
	"time"

	"github.com/JumoKookbob/whytie/internal/history"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

// ListRecentHistory returns events in reverse insertion order.
func (s *SQLiteStore) ListRecentHistory(limit int) ([]history.Event, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("history limit must be positive: %d", limit)
	}

	rows, err := s.db.Query(`
		SELECT
			memory_id,
			event_type,
			kind,
			text,
			path,
			line,
			commit_hash,
			occurred_at
		FROM history_events
		ORDER BY id DESC
		LIMIT ?
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []history.Event

	for rows.Next() {
		var event history.Event
		var eventType string
		var kind string
		var occurredAt string

		if err := rows.Scan(
			&event.MemoryID,
			&eventType,
			&kind,
			&event.Text,
			&event.Path,
			&event.Line,
			&event.CommitHash,
			&occurredAt,
		); err != nil {
			return nil, err
		}

		event.Type = history.EventType(eventType)
		event.Kind = syntax.Kind(kind)

		event.OccurredAt, err = time.Parse(time.RFC3339Nano, occurredAt)
		if err != nil {
			return nil, fmt.Errorf(
				"parse history timestamp for memory %s: %w",
				event.MemoryID,
				err,
			)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
