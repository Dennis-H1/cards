package model

import "time"

type Card struct {
	ID        int64
	Front     string
	Back      string
	Source    *string
	CreatedAt time.Time
	UpdatedAt time.Time
	Tags      []string
}

type Tag struct {
	ID                int64
	Name              string
	Overview          *string
	OverviewUpdatedAt *time.Time
	CardCount         int
}

type Review struct {
	CardID         int64
	EaseFactor     float64
	IntervalDays   int
	Repetitions    int
	DueAt          time.Time
	ReviewCount    int
	LastReviewedAt *time.Time
}

type EventType string

const (
	EventCardCreated EventType = "card_created"
	EventCardEdited  EventType = "card_edited"
)

type Event struct {
	ID        int64
	Type      EventType
	CardID    int64
	CreatedAt time.Time
	SeenAt    *time.Time
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
