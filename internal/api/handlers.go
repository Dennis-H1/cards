package api

import (
	"net/http"
	"strconv"

	"github.com/Dennis-H1/cards/internal/model"
	"github.com/Dennis-H1/cards/internal/service"
)

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func handleCreateCard(svc *service.Service) http.HandlerFunc {
	type request struct {
		Front  string   `json:"front"`
		Back   string   `json:"back"`
		Tags   []string `json:"tags"`
		Source *string  `json:"source"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := readJSON(r, &req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if req.Front == "" || req.Back == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "front and back are required"})
			return
		}
		card, err := svc.CreateCard(r.Context(), req.Front, req.Back, req.Tags, req.Source)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, card)
	}
}

func handleGetCard(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid card id"})
			return
		}
		card, err := svc.GetCard(r.Context(), id)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, card)
	}
}

func handleUpdateCard(svc *service.Service) http.HandlerFunc {
	type request struct {
		Front *string  `json:"front"`
		Back  *string  `json:"back"`
		Tags  []string `json:"tags"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid card id"})
			return
		}
		var req request
		if err := readJSON(r, &req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		card, err := svc.UpdateCard(r.Context(), id, req.Front, req.Back, req.Tags)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, card)
	}
}

func handleDueCards(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cards, err := svc.DueQueue(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, cards)
	}
}

func handleSearchCards(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cards, err := svc.SearchCards(r.Context(), r.URL.Query().Get("q"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, cards)
	}
}

func handleGradeCard(svc *service.Service) http.HandlerFunc {
	type request struct {
		Grade model.Grade `json:"grade"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid card id"})
			return
		}
		var req request
		if err := readJSON(r, &req); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		review, err := svc.GradeCard(r.Context(), id, req.Grade)
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, review)
	}
}

func handleListTags(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := svc.ListTags(r.Context())
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, tags)
	}
}

func handleTagOverview(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		overview, err := svc.GetTagOverview(r.Context(), r.PathValue("name"))
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, overview)
	}
}

func handleListActivity(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			events []model.Event
			err    error
		)
		if r.URL.Query().Get("unseen") == "true" {
			events, err = svc.UnseenActivity(r.Context())
		} else {
			events, err = svc.ListActivity(r.Context(), 100)
		}
		if err != nil {
			writeError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, events)
	}
}

func handleMarkActivitySeen(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.MarkActivitySeen(r.Context()); err != nil {
			writeError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
