package agent

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

var coordinatorTools = []llms.Tool{
	{
		Type: "function",
		Function: &llms.FunctionDefinition{
			Name:        "handoff_to_planner",
			Description: "Handoff to planner agent to do plan",
			Parameters:  map[string]any{},
		},
	},
}

type Coordinator struct {
	llm *openai.LLM
}

func NewCoordinator(llm *openai.LLM) *Coordinator {
	return &Coordinator{
		llm: llm,
	}
}

func (coord *Coordinator) Coordinate(state *AgentState) (*AgentUpdate, error) {
	// load prompts
	content, err := os.ReadFile("./internal/prompts/coordinator.md")
	if err != nil {
		return nil, err
	}
	promptTemplate, err := prompts.NewPromptTemplate(string(content), []string{"CURRENT_TIME"}).Format(map[string]any{
		"CURRENT_TIME": time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return nil, err
	}

	var msgs []llms.MessageContent
	msgs = append(msgs, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextContent{Text: promptTemplate}},
	})

	msgs = append(msgs, state.Messages...)

	resp, err := coord.llm.GenerateContent(context.Background(), msgs, llms.WithTools(coordinatorTools))
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("empty response")
	}

	if len(resp.Choices[0].ToolCalls) > 0 {
		return &AgentUpdate{
			NextStep: StepPlanner,
		}, nil
	}

	return &AgentUpdate{
		NextStep: StepEnd,
		Output:   resp.Choices[0].Content,
	}, nil
}
