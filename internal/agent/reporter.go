package agent

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

var _ Node = (*Reporter)(nil)

type Reporter struct {
	llm *openai.LLM
}

func NewReporter(llm *openai.LLM) *Reporter {
	return &Reporter{
		llm: llm,
	}
}

func (reporter *Reporter) Execute(ctx context.Context, state *AgentState) (nextStep string, output string, err error) {
	slog.Info("reporter starts")

	content, err := os.ReadFile("./internal/prompts/reporter.md")
	if err != nil {
		slog.Error("read planner prompt file", "error", err)
		return "", "", err
	}
	promptTemplate, err := prompts.NewPromptTemplate(string(content), []string{"current_time", "max_step_num"}).Format(map[string]any{
		"current_time": time.Now().Format(time.RFC3339),
	})
	if err != nil {
		slog.Error("format planner prompt", "error", err)
		return "", "", err
	}

	var messages []llms.MessageContent
	messages = append(messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeSystem,
		Parts: []llms.ContentPart{llms.TextContent{Text: promptTemplate}},
	})

	messages = append(messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextContent{Text: fmt.Sprintf("# Research Requirements\n\n## Task\n\n%s\n\n## Description\n\n%s", state.CurrentPlan.Title, state.CurrentPlan.Thought)}},
	})

	messages = append(messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeSystem,
		Parts: []llms.ContentPart{llms.TextContent{Text: "IMPORTANT: Structure your report according to the format in the prompt. Remember to include:\n\n1. Key Points - A bulleted list of the most important findings\n2. Overview - A brief introduction to the topic\n3. Detailed Analysis - Organized into logical sections\n4. Survey Note (optional) - For more comprehensive reports\n5. Key Citations - List all references at the end\n\nFor citations, DO NOT include inline citations in the text. Instead, place all citations in the 'Key Citations' section at the end using the format: `- [Source Title](URL)`. Include an empty line between each citation for better readability.\n\nPRIORITIZE USING MARKDOWN TABLES for data presentation and comparison. Use tables whenever presenting comparative data, statistics, features, or options. Structure tables with clear headers and aligned columns. Example table format:\n\n| Feature | Description | Pros | Cons |\n|---------|-------------|------|------|\n| Feature 1 | Description 1 | Pros 1 | Cons 1 |\n| Feature 2 | Description 2 | Pros 2 | Cons 2 |"}},
	})

	messages = append(messages, state.Messages...)

	resp, err := reporter.llm.GenerateContent(ctx, messages)
	if err != nil {
		slog.Error("generate plan", "error", err)
		return "", "", err
	}

	slog.Info("reporter ends")
	return StepEnd, resp.Choices[0].Content, nil
}
