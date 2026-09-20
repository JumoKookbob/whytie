package storage

import (
	"database/sql"
	"path/filepath"
	"time"

	"github.com/JumoKookbob/whytie/internal/history"
	"github.com/JumoKookbob/whytie/internal/memory"
	"github.com/JumoKookbob/whytie/internal/syntax"
	_ "modernc.org/sqlite"
)

type SQLiteStore struct {
	db *sql.DB
}

var _ Store = (*SQLiteStore)(nil)

func OpenSQLite(root string) (*SQLiteStore, error) {
	path := filepath.Join(root, ".whytie", "whytie.db")

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	const schema = `
	CREATE TABLE IF NOT EXISTS memories (
		id TEXT PRIMARY KEY,
		kind TEXT NOT NULL,
		text TEXT NOT NULL,
		created_path TEXT NOT NULL,
		created_line INTEGER NOT NULL,
		current_path TEXT NOT NULL,
		current_line INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS history_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		memory_id TEXT NOT NULL,
		event_type TEXT NOT NULL,
		path TEXT NOT NULL,
		line INTEGER NOT NULL,
		commit_hash TEXT NOT NULL,
		occurred_at TEXT NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_history_events_memory_id
	ON history_events(memory_id);
	`

	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return &SQLiteStore{
		db: db,
	}, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) Save(m memory.Memory) error {
	_, err := s.db.Exec(`
		INSERT INTO memories (
			id,
			kind,
			text,
			created_path,
			created_line,
			current_path,
			current_line
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
				kind = excluded.kind,
				text = excluded.text,
				current_path = excluded.current_path,
				current_line = excluded.current_line
	`,
		m.ID,
		string(m.Kind),
		m.Text,
		m.CreatedPath,
		m.CreatedLine,
		m.CurrentPath,
		m.CurrentLine,
	)

	return err
}

func (s *SQLiteStore) Delete(id string) error {
	_, err := s.db.Exec(`
		DELETE FROM memories
		WHERE id = ?
	`, id)

	return err
}

func (s *SQLiteStore) Get(id string) (memory.Memory, error) {
	var m memory.Memory

	var kind string

	err := s.db.QueryRow(`
		SELECT
			id,
			kind,
			text,
			created_path,
			created_line,
			current_path,
			current_line
		FROM memories
		WHERE id = ?
	`, id).Scan(
		&m.ID,
		&kind,
		&m.Text,
		&m.CreatedPath,
		&m.CreatedLine,
		&m.CurrentPath,
		&m.CurrentLine,
	)
	if err != nil {
		return memory.Memory{}, err
	}

	m.Kind = syntax.Kind(kind)

	return m, nil
}

func (s *SQLiteStore) List() ([]memory.Memory, error) {
	rows, err := s.db.Query(`
		SELECT
			id,
			kind,
			text,
			created_path,
			created_line,
			current_path,
			current_line
		FROM memories
		ORDER BY current_path, current_line, rowid
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memories []memory.Memory

	for rows.Next() {
		var m memory.Memory
		var kind string

		if err := rows.Scan(
			&m.ID,
			&kind,
			&m.Text,
			&m.CreatedPath,
			&m.CreatedLine,
			&m.CurrentPath,
			&m.CurrentLine,
		); err != nil {
			return nil, err
		}

		m.Kind = syntax.Kind(kind)
		memories = append(memories, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memories, nil
}

func (s *SQLiteStore) SaveHistory(event history.Event) error {
	_, err := s.db.Exec(`
		INSERT INTO history_events (
			memory_id,
			event_type,
			path,
			line,
			commit_hash,
			occurred_at
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`,
		event.MemoryID,
		string(event.Type),
		event.Path,
		event.Line,
		event.CommitHash,
		event.OccurredAt.UTC().Format(time.RFC3339Nano),
	)

	return err
}

func (s *SQLiteStore) ListHistory(memoryID string) ([]history.Event, error) {
	rows, err := s.db.Query(`
		SELECT
			memory_id,
			event_type,
			path,
			line,
			commit_hash,
			occurred_at
		FROM history_events
		WHERE memory_id = ?
		ORDER BY id
	`, memoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []history.Event

	for rows.Next() {
		var event history.Event
		var eventType string
		var occurredAt string

		if err := rows.Scan(
			&event.MemoryID,
			&eventType,
			&event.Path,
			&event.Line,
			&event.CommitHash,
			&occurredAt,
		); err != nil {
			return nil, err
		}

		event.Type = history.EventType(eventType)

		event.OccurredAt, err = time.Parse(time.RFC3339Nano, occurredAt)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
