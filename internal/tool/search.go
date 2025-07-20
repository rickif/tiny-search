package tool

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/strrl/tavily-go/pkg/tavily"
	"github.com/tmc/langchaingo/llms"
)

var SearchTool = llms.Tool{
	Type: "function",
	Function: &llms.FunctionDefinition{
		Name:        "tavily_search",
		Description: "Tool that queries the Tavily Search API and gets back json",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "search query to look up.",
				},
			},
			"required": []string{"query"},
		},
	},
}

type TavilySearchTool struct {
	key string
}

func NewTavilySearchTool(key string) *TavilySearchTool {
	return &TavilySearchTool{key: key}
}

func (t *TavilySearchTool) Search(ctx context.Context, query string) (string, error) {
	client := tavily.NewClient(t.key)
	resp, err := client.Search(ctx, query)
	if err != nil {
		slog.Error("tavily search", "error", err)
		return "", err
	}

	b, _ := json.Marshal(resp)

	return string(b), nil
}
