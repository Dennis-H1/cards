package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Dennis-H1/cards/internal/model"
)

var ErrTagNotFound = errors.New("tag not found")

func upsertTagIDs(ctx context.Context, q querier, names []string) ([]int64, error) {
	ids := make([]int64, 0, len(names))
	for _, name := range names {
		var id int64
		err := q.QueryRowContext(ctx, `SELECT id FROM tags WHERE name = ?`, name).Scan(&id)
		switch {
		case err == nil:
			ids = append(ids, id)
		case errors.Is(err, sql.ErrNoRows):
			res, err := q.ExecContext(ctx, `INSERT INTO tags (name) VALUES (?)`, name)
			if err != nil {
				return nil, err
			}
			id, err = res.LastInsertId()
			if err != nil {
				return nil, err
			}
			ids = append(ids, id)
		default:
			return nil, err
		}
	}
	return ids, nil
}

func (s *Store) ListTags(ctx context.Context) ([]model.Tag, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT t.id, t.name, t.overview, t.overview_updated_at, COUNT(ct.card_id)
		FROM tags t
		LEFT JOIN card_tags ct ON ct.tag_id = t.id
		GROUP BY t.id
		ORDER BY t.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var t model.Tag
		var overview sql.NullString
		var overviewUpdatedAt sql.NullTime
		if err := rows.Scan(&t.ID, &t.Name, &overview, &overviewUpdatedAt, &t.CardCount); err != nil {
			return nil, err
		}
		if overview.Valid {
			t.Overview = &overview.String
		}
		if overviewUpdatedAt.Valid {
			t.OverviewUpdatedAt = &overviewUpdatedAt.Time
		}
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

func (s *Store) GetTagByName(ctx context.Context, name string) (*model.Tag, error) {
	var t model.Tag
	var overview sql.NullString
	var overviewUpdatedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, `
		SELECT t.id, t.name, t.overview, t.overview_updated_at, COUNT(ct.card_id)
		FROM tags t
		LEFT JOIN card_tags ct ON ct.tag_id = t.id
		WHERE t.name = ?
		GROUP BY t.id`, name).Scan(&t.ID, &t.Name, &overview, &overviewUpdatedAt, &t.CardCount)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTagNotFound
	}
	if err != nil {
		return nil, err
	}
	if overview.Valid {
		t.Overview = &overview.String
	}
	if overviewUpdatedAt.Valid {
		t.OverviewUpdatedAt = &overviewUpdatedAt.Time
	}
	return &t, nil
}

func (s *Store) SetTagOverview(ctx context.Context, tagID int64, overview string, now time.Time) error {
	_, err := s.db.ExecContext(ctx, `UPDATE tags SET overview = ?, overview_updated_at = ? WHERE id = ?`, overview, now, tagID)
	return err
}

func (s *Store) CardsByTag(ctx context.Context, tagName string) ([]model.Card, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.front, c.back, c.source, c.created_at, c.updated_at
		FROM cards c
		JOIN card_tags ct ON ct.card_id = c.id
		JOIN tags t ON t.id = ct.tag_id
		WHERE t.name = ?
		ORDER BY c.created_at DESC`, tagName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCards(ctx, s.db, rows)
}
