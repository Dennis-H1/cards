package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Dennis-H1/cards/internal/model"
)

var ErrCardNotFound = errors.New("card not found")

func (s *Store) CreateCard(ctx context.Context, front, back string, source *string, tagNames []string, now time.Time) (*model.Card, error) {
	var cardID int64
	err := s.withTx(ctx, func(tx *sql.Tx) error {
		res, err := tx.ExecContext(ctx, `INSERT INTO cards (front, back, source, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
			front, back, source, now, now)
		if err != nil {
			return err
		}
		cardID, err = res.LastInsertId()
		if err != nil {
			return err
		}

		tagIDs, err := upsertTagIDs(ctx, tx, tagNames)
		if err != nil {
			return err
		}
		for _, tagID := range tagIDs {
			if _, err := tx.ExecContext(ctx, `INSERT INTO card_tags (card_id, tag_id) VALUES (?, ?)`, cardID, tagID); err != nil {
				return err
			}
		}

		if _, err := tx.ExecContext(ctx, `INSERT INTO reviews (card_id, due_at) VALUES (?, ?)`, cardID, now); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, `INSERT INTO events (type, card_id, created_at) VALUES (?, ?, ?)`,
			model.EventCardCreated, cardID, now); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.GetCard(ctx, cardID)
}

func (s *Store) UpdateCard(ctx context.Context, id int64, front, back *string, tagNames []string, now time.Time) (*model.Card, error) {
	err := s.withTx(ctx, func(tx *sql.Tx) error {
		if front != nil || back != nil {
			res, err := tx.ExecContext(ctx, `
				UPDATE cards SET
					front = COALESCE(?, front),
					back = COALESCE(?, back),
					updated_at = ?
				WHERE id = ?`, front, back, now, id)
			if err != nil {
				return err
			}
			if n, err := res.RowsAffected(); err != nil {
				return err
			} else if n == 0 {
				return ErrCardNotFound
			}
		}

		if tagNames != nil {
			if _, err := tx.ExecContext(ctx, `DELETE FROM card_tags WHERE card_id = ?`, id); err != nil {
				return err
			}
			tagIDs, err := upsertTagIDs(ctx, tx, tagNames)
			if err != nil {
				return err
			}
			for _, tagID := range tagIDs {
				if _, err := tx.ExecContext(ctx, `INSERT INTO card_tags (card_id, tag_id) VALUES (?, ?)`, id, tagID); err != nil {
					return err
				}
			}
		}

		_, err := tx.ExecContext(ctx, `INSERT INTO events (type, card_id, created_at) VALUES (?, ?, ?)`,
			model.EventCardEdited, id, now)
		return err
	})
	if err != nil {
		return nil, err
	}
	return s.GetCard(ctx, id)
}

func (s *Store) GetCard(ctx context.Context, id int64) (*model.Card, error) {
	var c model.Card
	var source sql.NullString
	err := s.db.QueryRowContext(ctx, `SELECT id, front, back, source, created_at, updated_at FROM cards WHERE id = ?`, id).
		Scan(&c.ID, &c.Front, &c.Back, &source, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrCardNotFound
	}
	if err != nil {
		return nil, err
	}
	if source.Valid {
		c.Source = &source.String
	}
	tags, err := loadCardTags(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	c.Tags = tags
	return &c, nil
}

// SearchCards does a substring search over front/back, most recently
// created first.
func (s *Store) SearchCards(ctx context.Context, query string, limit int) ([]model.Card, error) {
	like := "%" + query + "%"
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, front, back, source, created_at, updated_at
		FROM cards
		WHERE front LIKE ? OR back LIKE ?
		ORDER BY created_at DESC
		LIMIT ?`, like, like, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCards(ctx, s.db, rows)
}

func (s *Store) DueCards(ctx context.Context, now time.Time, limit int) ([]model.Card, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id, c.front, c.back, c.source, c.created_at, c.updated_at
		FROM cards c
		JOIN reviews r ON r.card_id = c.id
		WHERE r.due_at <= ?
		ORDER BY r.due_at ASC
		LIMIT ?`, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCards(ctx, s.db, rows)
}

func scanCards(ctx context.Context, q querier, rows *sql.Rows) ([]model.Card, error) {
	var cards []model.Card
	for rows.Next() {
		var c model.Card
		var source sql.NullString
		if err := rows.Scan(&c.ID, &c.Front, &c.Back, &source, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if source.Valid {
			c.Source = &source.String
		}
		cards = append(cards, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range cards {
		tags, err := loadCardTags(ctx, q, cards[i].ID)
		if err != nil {
			return nil, err
		}
		cards[i].Tags = tags
	}
	return cards, nil
}

func loadCardTags(ctx context.Context, q querier, cardID int64) ([]string, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT t.name FROM tags t
		JOIN card_tags ct ON ct.tag_id = t.id
		WHERE ct.card_id = ?
		ORDER BY t.name`, cardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		names = append(names, name)
	}
	return names, rows.Err()
}
