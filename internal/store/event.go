package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/Dennis-H1/cards/internal/model"
)

func (s *Store) ListUnseenEvents(ctx context.Context) ([]model.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, type, card_id, created_at, seen_at
		FROM events
		WHERE seen_at IS NULL
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEvents(rows)
}

func (s *Store) ListEvents(ctx context.Context, limit int) ([]model.Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, type, card_id, created_at, seen_at
		FROM events
		ORDER BY created_at DESC
		LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanEvents(rows)
}

func (s *Store) MarkAllEventsSeen(ctx context.Context, now time.Time) error {
	_, err := s.db.ExecContext(ctx, `UPDATE events SET seen_at = ? WHERE seen_at IS NULL`, now)
	return err
}

func scanEvents(rows *sql.Rows) ([]model.Event, error) {
	var events []model.Event
	for rows.Next() {
		var e model.Event
		var seenAt sql.NullTime
		if err := rows.Scan(&e.ID, &e.Type, &e.CardID, &e.CreatedAt, &seenAt); err != nil {
			return nil, err
		}
		if seenAt.Valid {
			e.SeenAt = &seenAt.Time
		}
		events = append(events, e)
	}
	return events, rows.Err()
}
