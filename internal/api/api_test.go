package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Dennis-H1/cards/internal/db"
	"github.com/Dennis-H1/cards/internal/model"
	"github.com/Dennis-H1/cards/internal/service"
	"github.com/Dennis-H1/cards/internal/store"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	sqlDB, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })
	svc := service.New(store.New(sqlDB))
	return httptest.NewServer(NewRouter(svc))
}

func TestCardLifecycle(t *testing.T) {
	srv := newTestServer(t)
	defer srv.Close()

	createBody, _ := json.Marshal(map[string]any{
		"front": "What is SM-2?",
		"back":  "A spaced repetition scheduling algorithm.",
		"tags":  []string{"srs"},
	})
	resp, err := http.Post(srv.URL+"/api/cards", "application/json", bytes.NewReader(createBody))
	if err != nil {
		t.Fatalf("create card: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create card: status %d", resp.StatusCode)
	}
	var card model.Card
	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		t.Fatalf("decode card: %v", err)
	}
	if card.ID == 0 || len(card.Tags) != 1 {
		t.Fatalf("unexpected card: %+v", card)
	}

	getResp, err := http.Get(srv.URL + "/api/cards/" + strconv.FormatInt(card.ID, 10))
	if err != nil {
		t.Fatalf("get card: %v", err)
	}
	defer getResp.Body.Close()
	if getResp.StatusCode != http.StatusOK {
		t.Fatalf("get card: status %d", getResp.StatusCode)
	}
	var fetched model.Card
	if err := json.NewDecoder(getResp.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode fetched card: %v", err)
	}
	if fetched.ID != card.ID || fetched.Front != card.Front {
		t.Fatalf("fetched card mismatch: %+v", fetched)
	}

	dueResp, err := http.Get(srv.URL + "/api/cards/due")
	if err != nil {
		t.Fatalf("due queue: %v", err)
	}
	defer dueResp.Body.Close()
	var due []model.DueCard
	if err := json.NewDecoder(dueResp.Body).Decode(&due); err != nil {
		t.Fatalf("decode due: %v", err)
	}
	if len(due) != 1 || due[0].Card.ID != card.ID || due[0].Review.ReviewCount != 0 {
		t.Fatalf("expected card in due queue, got %+v", due)
	}

	gradeBody, _ := json.Marshal(map[string]string{"grade": string(model.GradeGood)})
	gradeReq, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/cards/"+strconv.FormatInt(card.ID, 10)+"/grade", bytes.NewReader(gradeBody))
	gradeResp, err := http.DefaultClient.Do(gradeReq)
	if err != nil {
		t.Fatalf("grade card: %v", err)
	}
	defer gradeResp.Body.Close()
	if gradeResp.StatusCode != http.StatusOK {
		t.Fatalf("grade card: status %d", gradeResp.StatusCode)
	}
	var review model.Review
	if err := json.NewDecoder(gradeResp.Body).Decode(&review); err != nil {
		t.Fatalf("decode review: %v", err)
	}
	if review.Repetitions != 1 || review.IntervalDays != 1 {
		t.Fatalf("unexpected review after grading: %+v", review)
	}

	overviewResp, err := http.Get(srv.URL + "/api/tags/srs/overview")
	if err != nil {
		t.Fatalf("tag overview: %v", err)
	}
	defer overviewResp.Body.Close()
	if overviewResp.StatusCode != http.StatusOK {
		t.Fatalf("tag overview: status %d", overviewResp.StatusCode)
	}
	var overview service.TagOverview
	if err := json.NewDecoder(overviewResp.Body).Decode(&overview); err != nil {
		t.Fatalf("decode overview: %v", err)
	}
	if len(overview.Cards) != 1 {
		t.Fatalf("expected 1 card under tag, got %+v", overview.Cards)
	}

	activityResp, err := http.Get(srv.URL + "/api/activity")
	if err != nil {
		t.Fatalf("activity: %v", err)
	}
	defer activityResp.Body.Close()
	var events []model.Event
	if err := json.NewDecoder(activityResp.Body).Decode(&events); err != nil {
		t.Fatalf("decode events: %v", err)
	}
	if len(events) != 1 || events[0].Type != model.EventCardCreated {
		t.Fatalf("unexpected activity: %+v", events)
	}
}
