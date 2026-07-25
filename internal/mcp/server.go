package mcp

import (
	"crypto/subtle"
	"net/http"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Dennis-H1/cards/internal/service"
)

func newServer(svc *service.Service) *mcpsdk.Server {
	s := mcpsdk.NewServer(&mcpsdk.Implementation{Name: "karteikarten", Version: "0.1.0"}, nil)
	registerTools(s, svc)
	registerResources(s)
	return s
}

// NewHTTPHandler serves MCP over streamable HTTP. The server is reachable
// from the public internet, so every request must carry a bearer token
// matching apiKey -- there is no open write endpoint.
func NewHTTPHandler(svc *service.Service, apiKey string) http.Handler {
	server := newServer(svc)
	mcpHandler := mcpsdk.NewStreamableHTTPHandler(func(*http.Request) *mcpsdk.Server {
		return server
	}, nil)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !validAPIKey(r, apiKey) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		mcpHandler.ServeHTTP(w, r)
	})
}

func validAPIKey(r *http.Request, apiKey string) bool {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(token), []byte(apiKey)) == 1
}
