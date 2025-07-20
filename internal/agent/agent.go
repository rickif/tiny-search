package agent

import (
	"context"

	"github.com/rickif/tiny-research/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Agent struct {
	llm *openai.LLM
}

func NewAgent(config config.Config) (*Agent, error) {
	llm, err := openai.New(openai.WithBaseURL(config.LLMBaseURL), openai.WithModel(config.LLMModel), openai.WithToken(config.LLMToken))
	if err != nil {
		return nil, err
	}
	return &Agent{
		llm: llm,
	}, nil
}

func (agent *Agent) Research(ctx context.Context, query string) (string, error) {
	coordinator := NewCoordinator(agent.llm)
	update, err := coordinator.Coordinate(&AgentState{
		Messages: []llms.MessageContent{
			{
				Role:  llms.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{llms.TextContent{Text: query}},
			},
		},
	})
	if err != nil {
		return "", err
	}

	return update.Output, nil
}
