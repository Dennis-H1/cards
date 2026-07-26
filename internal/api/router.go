// Package api is the REST HTTP layer, wired directly to the service layer.
package api

import (
	"net/http"

	"github.com/Dennis-H1/cards/internal/service"
)

// NewRouter builds the REST API mux. authMiddleware wraps every /api/...
// route so the login/logout endpoints (mounted separately by the caller) and
// /healthz stay reachable without a session.
func NewRouter(svc *service.Service, authMiddleware func(http.HandlerFunc) http.HandlerFunc) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handleHealthz)

	mux.HandleFunc("POST /api/cards", authMiddleware(handleCreateCard(svc)))
	mux.HandleFunc("GET /api/cards/due", authMiddleware(handleDueCards(svc)))
	mux.HandleFunc("GET /api/cards/search", authMiddleware(handleSearchCards(svc)))
	mux.HandleFunc("GET /api/cards/{id}", authMiddleware(handleGetCard(svc)))
	mux.HandleFunc("PATCH /api/cards/{id}", authMiddleware(handleUpdateCard(svc)))
	mux.HandleFunc("POST /api/cards/{id}/grade", authMiddleware(handleGradeCard(svc)))

	mux.HandleFunc("GET /api/tags", authMiddleware(handleListTags(svc)))
	mux.HandleFunc("GET /api/tags/{name}/overview", authMiddleware(handleTagOverview(svc)))

	mux.HandleFunc("GET /api/activity", authMiddleware(handleListActivity(svc)))
	mux.HandleFunc("POST /api/activity/seen", authMiddleware(handleMarkActivitySeen(svc)))

	return mux
}
