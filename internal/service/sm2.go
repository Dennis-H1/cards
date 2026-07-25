package service

import (
	"math"
	"time"

	"github.com/Dennis-H1/cards/internal/model"
)

var gradeScores = map[model.Grade]float64{
	model.GradeAgain: 0,
	model.GradeHard:  3,
	model.GradeGood:  4,
	model.GradeEasy:  5,
}

func applySM2(r model.Review, grade model.Grade, now time.Time) model.Review {
	score := gradeScores[grade]

	if grade == model.GradeAgain {
		r.Repetitions = 0
		r.IntervalDays = 1
	} else {
		switch r.Repetitions {
		case 0:
			r.IntervalDays = 1
		case 1:
			r.IntervalDays = 6
		default:
			r.IntervalDays = int(math.Round(float64(r.IntervalDays) * r.EaseFactor))
		}
		r.Repetitions++
	}

	r.EaseFactor = math.Max(1.3, r.EaseFactor+(0.1-(5-score)*(0.08+(5-score)*0.02)))
	r.DueAt = now.AddDate(0, 0, r.IntervalDays)
	r.ReviewCount++
	lastReviewed := now
	r.LastReviewedAt = &lastReviewed

	return r
}
