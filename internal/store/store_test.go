package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/Dennis-H1/cards/internal/db"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	sqlDB, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })
	return New(sqlDB)
}

func TestCreateGetUpdateCard(t *testing.T) {
	ctx := context.Background()
	s := newTestStore(t)
	now := time.Now().UTC()

	card, err := s.CreateCard(ctx, "front text", "back text", nil, []string{"go", "sqlite"}, now)
	if err != nil {
		t.Fatalf("CreateCard: %v", err)
	}
	if len(card.Tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", card.Tags)
	}

	review, err := s.GetReview(ctx, card.ID)
	if err != nil {
		t.Fatalf("GetReview: %v", err)
	}
	if review.EaseFactor != 2.5 || !review.DueAt.Equal(now) {
		t.Fatalf("unexpected default review: %+v", review)
	}

	newBack := "updated back"
	updated, err := s.UpdateCard(ctx, card.ID, nil, &newBack, []string{"go"}, now.Add(time.Hour))
	if err != nil {
		t.Fatalf("UpdateCard: %v", err)
	}
	if updated.Back != newBack || len(updated.Tags) != 1 {
		t.Fatalf("unexpected update result: %+v", updated)
	}

	due, err := s.DueCards(ctx, now.Add(time.Minute), 10)
	if err != nil {
		t.Fatalf("DueCards: %v", err)
	}
	if len(due) != 1 || due[0].Card.ID != card.ID || due[0].Review.EaseFactor != 2.5 {
		t.Fatalf("expected card in due queue, got %+v", due)
	}

	events, err := s.ListUnseenEvents(ctx)
	if err != nil {
		t.Fatalf("ListUnseenEvents: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events (created+edited), got %d", len(events))
	}

	tags, err := s.ListTags(ctx)
	if err != nil {
		t.Fatalf("ListTags: %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags total (go, sqlite), got %+v", tags)
	}
}
