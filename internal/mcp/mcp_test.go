package mcp

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Dennis-H1/cards/internal/db"
	"github.com/Dennis-H1/cards/internal/service"
	"github.com/Dennis-H1/cards/internal/store"
)

func newTestSession(t *testing.T) *mcpsdk.ClientSession {
	t.Helper()
	sqlDB, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { sqlDB.Close() })
	svc := service.New(store.New(sqlDB))

	serverTransport, clientTransport := mcpsdk.NewInMemoryTransports()
	server := newServer(svc)

	ctx := context.Background()
	go server.Run(ctx, serverTransport)

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "test-client"}, nil)
	session, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client.Connect: %v", err)
	}
	t.Cleanup(func() { session.Close() })
	return session
}

func TestCreateCardToolAndDesignSpecResource(t *testing.T) {
	ctx := context.Background()
	session := newTestSession(t)

	res, err := session.CallTool(ctx, &mcpsdk.CallToolParams{
		Name: "create_card",
		Arguments: map[string]any{
			"front": "What does MCP stand for?",
			"back":  "Model Context Protocol.",
			"tags":  []string{"mcp"},
		},
	})
	if err != nil {
		t.Fatalf("CallTool create_card: %v", err)
	}
	if res.IsError {
		t.Fatalf("create_card returned tool error: %+v", res.Content)
	}
	structured, err := json.Marshal(res.StructuredContent)
	if err != nil {
		t.Fatalf("marshal structured content: %v", err)
	}
	var card struct {
		ID   int64    `json:"id"`
		Tags []string `json:"tags"`
	}
	if err := json.Unmarshal(structured, &card); err != nil {
		t.Fatalf("unmarshal card: %v", err)
	}
	if card.ID == 0 || len(card.Tags) != 1 {
		t.Fatalf("unexpected card in structured content: %s", structured)
	}

	listRes, err := session.CallTool(ctx, &mcpsdk.CallToolParams{Name: "list_tags"})
	if err != nil {
		t.Fatalf("CallTool list_tags: %v", err)
	}
	if listRes.IsError {
		t.Fatalf("list_tags returned tool error: %+v", listRes.Content)
	}

	resource, err := session.ReadResource(ctx, &mcpsdk.ReadResourceParams{URI: designSpecURI})
	if err != nil {
		t.Fatalf("ReadResource design-spec: %v", err)
	}
	if len(resource.Contents) != 1 || resource.Contents[0].Text == "" {
		t.Fatalf("unexpected design-spec contents: %+v", resource.Contents)
	}
}
