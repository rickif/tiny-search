package agent

import (
	"context"
	"log/slog"

	"github.com/rickif/tiny-research/internal/config"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type Node interface {
	Execute(ctx context.Context, state *AgentState) (nextStep string, output string, err error)
}

type Agent struct {
	llm    *openai.LLM
	config *config.Config
}

func NewAgent(config config.Config) (*Agent, error) {
	llm, err := openai.New(openai.WithBaseURL(config.LLMBaseURL), openai.WithModel(config.LLMModel), openai.WithToken(config.LLMToken))
	if err != nil {
		return nil, err
	}
	return &Agent{
		llm:    llm,
		config: &config,
	}, nil
}

func (wf *Agent) Research(ctx context.Context, query string) (string, error) {
	state := AgentState{
		Messages: []llms.MessageContent{
			{
				Role:  llms.ChatMessageTypeHuman,
				Parts: []llms.ContentPart{llms.TextContent{Text: query}},
			},
		},
		Locale: "zh-CN",
	}
	coordinator := NewCoordinator(wf.llm)
	planner := NewPlanner(wf.llm, 3, 3)
	researchTeam := NewResearchTeam(wf.llm)
	researcher := NewResearcher(wf.llm, wf.config.TavilyKey)
	coder := NewCoder(wf.llm)

	nextStep, output, err := coordinator.Execute(ctx, &state)
	if err != nil {
		slog.Error("coordinate", "error", err)
		return "", err
	}

	for {
		switch nextStep {
		case StepPlanner:
			nextStep, output, err = planner.Execute(ctx, &state)
		case StepResearchTeam:
			nextStep, output, err = researchTeam.Execute(ctx, &state)
		case StepResearcher:
			nextStep, output, err = researcher.Execute(ctx, &state)
		case StepCoder:
			nextStep, output, err = coder.Execute(ctx, &state)
		case StepEnd:
			return output, nil
		default:
			slog.Error("unknown step", "step", nextStep)
			return "", err
		}
		if err != nil {
			slog.Error("execute", "error", err)
			return "", err
		}
	}
}
