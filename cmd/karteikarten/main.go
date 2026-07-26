package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Dennis-H1/cards/internal/api"
	"github.com/Dennis-H1/cards/internal/auth"
	"github.com/Dennis-H1/cards/internal/db"
	"github.com/Dennis-H1/cards/internal/mcp"
	"github.com/Dennis-H1/cards/internal/service"
	"github.com/Dennis-H1/cards/internal/store"
)

func main() {
	port := envOr("PORT", "8080")
	dbPath := envOr("DB_PATH", "./data/karteikarten.db")

	// The MCP server is reachable from the public internet (behind
	// Cloudflare Zero Trust) -- refuse to start without a credential rather
	// than silently exposing an open write endpoint.
	apiKey := os.Getenv("MCP_API_KEY")
	if apiKey == "" {
		log.Fatal("MCP_API_KEY must be set")
	}

	// The REST API/frontend are reachable from the public internet too -- a
	// single-account login gate keeps them from being wide open, the same
	// way MCP_API_KEY gates the MCP endpoint.
	authCfg := auth.Config{
		Username:      os.Getenv("AUTH_USERNAME"),
		Password:      os.Getenv("AUTH_PASSWORD"),
		SessionSecret: []byte(os.Getenv("SESSION_SECRET")),
	}
	if authCfg.Username == "" || authCfg.Password == "" || len(authCfg.SessionSecret) == 0 {
		log.Fatal("AUTH_USERNAME, AUTH_PASSWORD, and SESSION_SECRET must be set")
	}

	sqlDB, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer sqlDB.Close()

	svc := service.New(store.New(sqlDB))

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcp.NewHTTPHandler(svc, apiKey))
	mux.HandleFunc("POST /api/login", auth.LoginHandler(authCfg))
	mux.HandleFunc("POST /api/logout", auth.LogoutHandler(authCfg))
	mux.Handle("/", api.NewRouter(svc, auth.RequireAuth(authCfg)))

	log.Printf("listening on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
