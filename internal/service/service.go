package service

import (
	"context"
	"time"

	"github.com/Dennis-H1/cards/internal/model"
	"github.com/Dennis-H1/cards/internal/store"
)

const (
	defaultSearchLimit = 50
	defaultDueLimit    = 100
)

type Service struct {
	store *store.Store
	now   func() time.Time
}

func New(st *store.Store) *Service {
	return &Service{store: st, now: time.Now}
}

func (s *Service) CreateCard(ctx context.Context, front, back string, tags []string, source *string) (*model.Card, error) {
	return s.store.CreateCard(ctx, front, back, source, tags, s.now())
}

// UpdateCard updates only the fields that are non-nil; tags is replaced
// wholesale when non-nil.
func (s *Service) UpdateCard(ctx context.Context, id int64, front, back *string, tags []string) (*model.Card, error) {
	return s.store.UpdateCard(ctx, id, front, back, tags, s.now())
}

func (s *Service) GetCard(ctx context.Context, id int64) (*model.Card, error) {
	return s.store.GetCard(ctx, id)
}

func (s *Service) SearchCards(ctx context.Context, query string) ([]model.Card, error) {
	return s.store.SearchCards(ctx, query, defaultSearchLimit)
}

func (s *Service) ListTags(ctx context.Context) ([]model.Tag, error) {
	return s.store.ListTags(ctx)
}

type TagOverview struct {
	Tag   model.Tag    `json:"tag"`
	Cards []model.Card `json:"cards"`
}

func (s *Service) GetTagOverview(ctx context.Context, tagName string) (*TagOverview, error) {
	tag, err := s.store.GetTagByName(ctx, tagName)
	if err != nil {
		return nil, err
	}
	cards, err := s.store.CardsByTag(ctx, tagName)
	if err != nil {
		return nil, err
	}
	return &TagOverview{Tag: *tag, Cards: cards}, nil
}

func (s *Service) SetTagOverview(ctx context.Context, tagName, overview string) error {
	tag, err := s.store.GetTagByName(ctx, tagName)
	if err != nil {
		return err
	}
	return s.store.SetTagOverview(ctx, tag.ID, overview, s.now())
}

func (s *Service) DueQueue(ctx context.Context) ([]model.DueCard, error) {
	return s.store.DueCards(ctx, s.now(), defaultDueLimit)
}

func (s *Service) GradeCard(ctx context.Context, cardID int64, grade model.Grade) (*model.Review, error) {
	review, err := s.store.GetReview(ctx, cardID)
	if err != nil {
		return nil, err
	}
	updated := applySM2(*review, grade, s.now())
	if err := s.store.UpdateReview(ctx, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *Service) ListActivity(ctx context.Context, limit int) ([]model.Event, error) {
	return s.store.ListEvents(ctx, limit)
}

func (s *Service) UnseenActivity(ctx context.Context) ([]model.Event, error) {
	return s.store.ListUnseenEvents(ctx)
}

func (s *Service) MarkActivitySeen(ctx context.Context) error {
	return s.store.MarkAllEventsSeen(ctx, s.now())
}
