package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/rickif/tiny-research/internal/tool"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

var _ Node = (*Researcher)(nil)

type Researcher struct {
	llm       *openai.LLM
	tavilyKey string
}

func NewResearcher(llm *openai.LLM, tavilyKey string) *Researcher {
	return &Researcher{
		llm:       llm,
		tavilyKey: tavilyKey,
	}
}

func (r *Researcher) Execute(ctx context.Context, state *AgentState) (nextStep string, output string, err error) {
	slog.Info("researcher starts")
	content, err := os.ReadFile("./internal/prompts/researcher.md")
	if err != nil {
		slog.Error("read prompt template", "error", err)
		return "", "", err
	}

	promptTemplate, err := prompts.NewPromptTemplate(
		string(content),
		[]string{"current_time", "locale"},
	).Format(map[string]any{
		"current_time": time.Now().Format(time.RFC3339),
		"locale":       state.Locale,
	})
	if err != nil {
		slog.Error("format prompt", "error", err)
		return "", "", err
	}

	var step *Step
	for i := range state.CurrentPlan.Steps {
		if state.CurrentPlan.Steps[i].ExecutionResult == "" {
			step = &state.CurrentPlan.Steps[i]
			break
		}
	}

	messages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: promptTemplate}},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextContent{Text: fmt.Sprintf("#Task\n\ntitle: %s\n\n##description:%s", step.Title, step.Description)}},
		}}

	for {
		resp, err := r.llm.GenerateContent(ctx, messages, llms.WithTools([]llms.Tool{tool.CrawlTool, tool.SearchTool}))
		if err != nil {
			slog.Error("generate content", "error", err)
			return "", "", err
		}

		for _, toolcall := range resp.Choices[0].ToolCalls {
			switch toolcall.FunctionCall.Name {
			case "crawl":
				var args struct {
					URL string `json:"url"`
				}
				if err := json.Unmarshal([]byte(toolcall.FunctionCall.Arguments), &args); err != nil {
					slog.Error("unmarshal arguments", "error", err)
					return "", "", err
				}
				output, err := tool.Crawl(ctx, args.URL)
				if err != nil {
					slog.Error("crawl", "error", err)
					return "", "", err
				}
				message := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: toolcall.ID,
							Name:       toolcall.FunctionCall.Name,
							Content:    output,
						},
					},
				}
				messages = append(messages, message)
				slog.Info("researcher use crawl", "url", args.URL)
			case "tavily_search":
				var args struct {
					Query string `json:"query"`
				}
				if err := json.Unmarshal([]byte(toolcall.FunctionCall.Arguments), &args); err != nil {
					slog.Error("unmarshal arguments", "error", err)
					return "", "", err
				}
				output, err := tool.NewTavilySearchTool(r.tavilyKey).Search(ctx, args.Query)
				if err != nil {
					slog.Error("search", "error", err)
					return "", "", err
				}
				message := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: toolcall.ID,
							Name:       toolcall.FunctionCall.Name,
							Content:    output,
						},
					},
				}
				messages = append(messages, message)
				slog.Info("researcher use tavily search", "query", args.Query)
			default:
				slog.Error("unexpected function call", "name", toolcall.FunctionCall.Name)
				return "", "", fmt.Errorf("unexpected function call: %v", toolcall.FunctionCall.Name)
			}
		}
		if len(resp.Choices[0].ToolCalls) == 0 {
			slog.Info("researcher finished")
			step.ExecutionResult = resp.Choices[0].Content
			break
		}
		var toolCalls []string
		for _, toolcall := range resp.Choices[0].ToolCalls {
			toolCalls = append(toolCalls, fmt.Sprintf("%s: %s", toolcall.FunctionCall.Name, toolcall.FunctionCall.Arguments))
		}
		slog.Info("researcher use tools", "tools", toolCalls)
	}
	state.Messages = append(state.Messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextContent{Text: step.ExecutionResult}},
	})
	return StepResearchTeam, step.ExecutionResult, nil
}
