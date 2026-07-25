package service

import (
	"testing"
	"time"

	"github.com/Dennis-H1/cards/internal/model"
)

func TestApplySM2Progression(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	r := model.Review{EaseFactor: 2.5, IntervalDays: 0, Repetitions: 0}

	r = applySM2(r, model.GradeGood, now)
	if r.IntervalDays != 1 || r.Repetitions != 1 {
		t.Fatalf("first Good: got interval=%d repetitions=%d", r.IntervalDays, r.Repetitions)
	}

	r = applySM2(r, model.GradeGood, now.AddDate(0, 0, 1))
	if r.IntervalDays != 6 || r.Repetitions != 2 {
		t.Fatalf("second Good: got interval=%d repetitions=%d", r.IntervalDays, r.Repetitions)
	}

	prevEase := r.EaseFactor
	r = applySM2(r, model.GradeGood, now.AddDate(0, 0, 7))
	wantInterval := int(float64(6) * prevEase)
	if r.IntervalDays != wantInterval || r.Repetitions != 3 {
		t.Fatalf("third Good: got interval=%d (want ~%d) repetitions=%d", r.IntervalDays, wantInterval, r.Repetitions)
	}

	r = applySM2(r, model.GradeAgain, now.AddDate(0, 0, 20))
	if r.Repetitions != 0 || r.IntervalDays != 1 {
		t.Fatalf("Again must reset: got interval=%d repetitions=%d", r.IntervalDays, r.Repetitions)
	}

	if r.EaseFactor < 1.3 {
		t.Fatalf("ease factor must never drop below 1.3, got %f", r.EaseFactor)
	}
}

func TestApplySM2EaseFactorFloor(t *testing.T) {
	now := time.Now()
	r := model.Review{EaseFactor: 1.3, IntervalDays: 1, Repetitions: 1}
	for i := 0; i < 10; i++ {
		r = applySM2(r, model.GradeAgain, now)
	}
	if r.EaseFactor < 1.3 {
		t.Fatalf("ease factor floor violated: %f", r.EaseFactor)
	}
}
