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

var coordinatorTool = llms.Tool{
	Type: "function",
	Function: &llms.FunctionDefinition{
		Name:        "handoff_to_planner",
		Description: "Handoff to planner agent to do plan",
		Parameters:  map[string]any{},
	},
}

var _ Node = (*Coordinator)(nil)

type Coordinator struct {
	llm *openai.LLM
}

func NewCoordinator(llm *openai.LLM) *Coordinator {
	return &Coordinator{
		llm: llm,
	}
}

func (coord *Coordinator) Execute(ctx context.Context, state *AgentState) (nextStep string, output string, err error) {
	// load prompts
	content, err := os.ReadFile("./internal/prompts/coordinator.md")
	if err != nil {
		slog.Info("read coordinator prompt file", "error", err)
		return "", "", err
	}
	promptTemplate, err := prompts.NewPromptTemplate(string(content), []string{"current_time", "locale"}).Format(map[string]any{
		"current_time": time.Now().Format(time.RFC3339),
		"locale":       state.Locale,
	})
	if err != nil {
		return "", "", err
	}

	var messages []llms.MessageContent
	messages = append(messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeSystem,
		Parts: []llms.ContentPart{llms.TextContent{Text: promptTemplate}},
	})

	messages = append(messages, state.Messages...)

	resp, err := coord.llm.GenerateContent(context.Background(), messages, llms.WithTools([]llms.Tool{coordinatorTool}))
	if err != nil {
		return "", "", err
	}

	if len(resp.Choices) == 0 {
		return "", "", fmt.Errorf("empty response")
	}

	if len(resp.Choices[0].ToolCalls) > 0 {
		return StepPlanner, "", nil
	}

	return StepEnd, resp.Choices[0].Content, nil
}
