package mcp

import (
	"context"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Dennis-H1/cards/internal/model"
	"github.com/Dennis-H1/cards/internal/service"
)

func registerTools(s *mcpsdk.Server, svc *service.Service) {
	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        "create_card",
		Description: "Persist a new flashcard, creating any tags that don't exist yet. Creates a default (immediately due) Review row and logs a card_created event.",
	}, createCardHandler(svc))

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        "update_card",
		Description: "Update an existing flashcard's front, back, and/or tags. Omit a field to leave it unchanged. Logs a card_edited event.",
	}, updateCardHandler(svc))

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        "list_tags",
		Description: "List all tags with their card counts.",
	}, listTagsHandler(svc))

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        "get_tag_overview",
		Description: "Get the stored synthesized summary and the cards under a tag.",
	}, getTagOverviewHandler(svc))

	mcpsdk.AddTool(s, &mcpsdk.Tool{
		Name:        "search_cards",
		Description: "Text search over card front/back. Use this to check for near-duplicates before creating a new card.",
	}, searchCardsHandler(svc))
}

type createCardInput struct {
	Front  string   `json:"front" jsonschema:"card front content, as markdown"`
	Back   string   `json:"back" jsonschema:"card back content, as markdown"`
	Tags   []string `json:"tags,omitempty" jsonschema:"tag names to attach to the card"`
	Source string   `json:"source,omitempty" jsonschema:"optional source URL or short free-text context"`
}

func createCardHandler(svc *service.Service) mcpsdk.ToolHandlerFor[createCardInput, model.Card] {
	return func(ctx context.Context, _ *mcpsdk.CallToolRequest, in createCardInput) (*mcpsdk.CallToolResult, model.Card, error) {
		var source *string
		if in.Source != "" {
			source = &in.Source
		}
		card, err := svc.CreateCard(ctx, in.Front, in.Back, in.Tags, source)
		if err != nil {
			return nil, model.Card{}, err
		}
		return nil, *card, nil
	}
}

type updateCardInput struct {
	CardID int64    `json:"card_id" jsonschema:"id of the card to update"`
	Front  *string  `json:"front,omitempty" jsonschema:"new front content; omit to leave unchanged"`
	Back   *string  `json:"back,omitempty" jsonschema:"new back content; omit to leave unchanged"`
	Tags   []string `json:"tags,omitempty" jsonschema:"replacement tag list; omit to leave unchanged"`
}

func updateCardHandler(svc *service.Service) mcpsdk.ToolHandlerFor[updateCardInput, model.Card] {
	return func(ctx context.Context, _ *mcpsdk.CallToolRequest, in updateCardInput) (*mcpsdk.CallToolResult, model.Card, error) {
		card, err := svc.UpdateCard(ctx, in.CardID, in.Front, in.Back, in.Tags)
		if err != nil {
			return nil, model.Card{}, err
		}
		return nil, *card, nil
	}
}

type listTagsInput struct{}

// tagsResult wraps the tag list in an object, as required by the MCP spec's
// output schema (a tool's structured output must be a JSON object).
type tagsResult struct {
	Tags []model.Tag `json:"tags"`
}

func listTagsHandler(svc *service.Service) mcpsdk.ToolHandlerFor[listTagsInput, tagsResult] {
	return func(ctx context.Context, _ *mcpsdk.CallToolRequest, _ listTagsInput) (*mcpsdk.CallToolResult, tagsResult, error) {
		tags, err := svc.ListTags(ctx)
		return nil, tagsResult{Tags: tags}, err
	}
}

type getTagOverviewInput struct {
	Tag string `json:"tag" jsonschema:"the tag name"`
}

func getTagOverviewHandler(svc *service.Service) mcpsdk.ToolHandlerFor[getTagOverviewInput, service.TagOverview] {
	return func(ctx context.Context, _ *mcpsdk.CallToolRequest, in getTagOverviewInput) (*mcpsdk.CallToolResult, service.TagOverview, error) {
		overview, err := svc.GetTagOverview(ctx, in.Tag)
		if err != nil {
			return nil, service.TagOverview{}, err
		}
		return nil, *overview, nil
	}
}

type searchCardsInput struct {
	Query string `json:"query" jsonschema:"free-text search over card front/back"`
}

type cardsResult struct {
	Cards []model.Card `json:"cards"`
}

func searchCardsHandler(svc *service.Service) mcpsdk.ToolHandlerFor[searchCardsInput, cardsResult] {
	return func(ctx context.Context, _ *mcpsdk.CallToolRequest, in searchCardsInput) (*mcpsdk.CallToolResult, cardsResult, error) {
		cards, err := svc.SearchCards(ctx, in.Query)
		return nil, cardsResult{Cards: cards}, err
	}
}
