package tool

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/tmc/langchaingo/llms"
)

var CrawlTool = llms.Tool{
	Type: "function",
	Function: &llms.FunctionDefinition{
		Name:        "crawl",
		Description: "Use this to crawl a url and get a readable content in markdown format.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"url": map[string]any{
					"type":        "string",
					"description": "The url to crawl.",
				},
			},
			"required": []string{"url"},
		},
	},
}

func Crawl(ctx context.Context, url string) (string, error) {
	requetURL := "https://r.jina.ai/" + url
	resp, err := http.Get(requetURL)
	if err != nil {
		slog.Error("jina ai crawler", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("jina ai crawler", "status_code", resp.StatusCode)
		return "", fmt.Errorf("crawl failed, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("jina ai crawler", "error", err)
		return "", err
	}

	return string(body), nil
}
