package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Dennis-H1/cards/internal/model"
)

var ErrReviewNotFound = errors.New("review not found")

func (s *Store) GetReview(ctx context.Context, cardID int64) (*model.Review, error) {
	var r model.Review
	var lastReviewedAt sql.NullTime
	err := s.db.QueryRowContext(ctx, `
		SELECT card_id, ease_factor, interval_days, repetitions, due_at, review_count, last_reviewed_at
		FROM reviews WHERE card_id = ?`, cardID).
		Scan(&r.CardID, &r.EaseFactor, &r.IntervalDays, &r.Repetitions, &r.DueAt, &r.ReviewCount, &lastReviewedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrReviewNotFound
	}
	if err != nil {
		return nil, err
	}
	if lastReviewedAt.Valid {
		r.LastReviewedAt = &lastReviewedAt.Time
	}
	return &r, nil
}

func (s *Store) UpdateReview(ctx context.Context, r *model.Review) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE reviews SET
			ease_factor = ?,
			interval_days = ?,
			repetitions = ?,
			due_at = ?,
			review_count = ?,
			last_reviewed_at = ?
		WHERE card_id = ?`,
		r.EaseFactor, r.IntervalDays, r.Repetitions, r.DueAt, r.ReviewCount, r.LastReviewedAt, r.CardID)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrReviewNotFound
	}
	return nil
}
