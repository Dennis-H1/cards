// Package api is the REST HTTP layer, wired directly to the service layer.
package api

import (
	"net/http"

	"github.com/Dennis-H1/cards/internal/service"
)

func NewRouter(svc *service.Service) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handleHealthz)

	mux.HandleFunc("POST /api/cards", handleCreateCard(svc))
	mux.HandleFunc("GET /api/cards/due", handleDueCards(svc))
	mux.HandleFunc("GET /api/cards/search", handleSearchCards(svc))
	mux.HandleFunc("GET /api/cards/{id}", handleGetCard(svc))
	mux.HandleFunc("PATCH /api/cards/{id}", handleUpdateCard(svc))
	mux.HandleFunc("POST /api/cards/{id}/grade", handleGradeCard(svc))

	mux.HandleFunc("GET /api/tags", handleListTags(svc))
	mux.HandleFunc("GET /api/tags/{name}/overview", handleTagOverview(svc))

	mux.HandleFunc("GET /api/activity", handleListActivity(svc))
	mux.HandleFunc("POST /api/activity/seen", handleMarkActivitySeen(svc))

	return mux
}
