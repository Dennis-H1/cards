package model

import "time"

type Card struct {
	ID        int64     `json:"id"`
	Front     string    `json:"front"`
	Back      string    `json:"back"`
	Source    *string   `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []string  `json:"tags"`
}

type Tag struct {
	ID                int64      `json:"id"`
	Name              string     `json:"name"`
	Overview          *string    `json:"overview"`
	OverviewUpdatedAt *time.Time `json:"overview_updated_at"`
	CardCount         int        `json:"card_count"`
}

type Review struct {
	CardID         int64      `json:"card_id"`
	EaseFactor     float64    `json:"ease_factor"`
	IntervalDays   int        `json:"interval_days"`
	Repetitions    int        `json:"repetitions"`
	DueAt          time.Time  `json:"due_at"`
	ReviewCount    int        `json:"review_count"`
	LastReviewedAt *time.Time `json:"last_reviewed_at"`
}

// DueCard pairs a card with its review state -- the due queue needs the
// review fields (ease factor, interval, repetitions) so the client can
// classify new-vs-due and preview grading outcomes without a second request.
type DueCard struct {
	Card   Card   `json:"card"`
	Review Review `json:"review"`
}

type EventType string

const (
	EventCardCreated EventType = "card_created"
	EventCardEdited  EventType = "card_edited"
)

type Event struct {
	ID        int64      `json:"id"`
	Type      EventType  `json:"type"`
	CardID    int64      `json:"card_id"`
	CreatedAt time.Time  `json:"created_at"`
	SeenAt    *time.Time `json:"seen_at"`
}

// Grade is a review grading input (Again/Hard/Good/Easy), mapped to SM-2
// quality scores in the service layer.
type Grade string

const (
	GradeAgain Grade = "again"
	GradeHard  Grade = "hard"
	GradeGood  Grade = "good"
	GradeEasy  Grade = "easy"
)
