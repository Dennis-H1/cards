package mcp

import (
	"context"
	_ "embed"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed design-spec.json
var designSpecJSON string

const designSpecURI = "resource://karteikarten/design-spec"

func registerResources(s *mcpsdk.Server) {
	s.AddResource(&mcpsdk.Resource{
		URI:         designSpecURI,
		Name:        "design-spec",
		Description: "Canonical design tokens (colors, fonts, motion) and reference card markup. This is the single source of truth for card visual design -- fetch it before rendering any in-chat card preview so it stays pixel-consistent with the real app.",
		MIMEType:    "application/json",
	}, designSpecHandler)
}

func designSpecHandler(context.Context, *mcpsdk.ReadResourceRequest) (*mcpsdk.ReadResourceResult, error) {
	return &mcpsdk.ReadResourceResult{
		Contents: []*mcpsdk.ResourceContents{
			{
				URI:      designSpecURI,
				MIMEType: "application/json",
				Text:     designSpecJSON,
			},
		},
	}, nil
}
