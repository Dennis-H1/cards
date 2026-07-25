package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Dennis-H1/cards/internal/api"
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

	sqlDB, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer sqlDB.Close()

	svc := service.New(store.New(sqlDB))

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcp.NewHTTPHandler(svc, apiKey))
	mux.Handle("/", api.NewRouter(svc))

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
