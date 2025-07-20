package agent

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/rickif/tiny-research/internal/llm"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

var _ Node = (*Planner)(nil)

type Planner struct {
	llm           *openai.LLM
	maxIterations int
	maxStepNum    int
}

func NewPlanner(llm *openai.LLM, maxIterations int, maxStepNum int) *Planner {
	return &Planner{
		llm:           llm,
		maxIterations: maxIterations,
		maxStepNum:    maxStepNum,
	}
}

func (planner *Planner) Execute(ctx context.Context, state *AgentState) (nextStep string, output string, err error) {
	slog.Info("planner starts")
	if state.PlanIterations >= planner.maxIterations {
		return StepReporter, "", nil
	}

	content, err := os.ReadFile("./internal/prompts/planner.md")
	if err != nil {
		slog.Error("read planner prompt file", "error", err)
		return "", "", err
	}
	promptTemplate, err := prompts.NewPromptTemplate(string(content), []string{"current_time", "max_step_num"}).Format(map[string]any{
		"current_time": time.Now().Format(time.RFC3339),
		"max_step_num": planner.maxStepNum,
		"locale":       state.Locale,
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

	messages = append(messages, state.Messages...)

	var plan Plan
	output, err = llm.GenerateJSON(ctx, planner.llm, messages, &plan, 3)
	if err != nil {
		slog.Error("generate plan", "error", err)
		return "", "", err
	}

	nextStep = StepResearchTeam
	if plan.HasEnoughContext {
		nextStep = StepReporter
	}

	state.Messages = append(state.Messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextContent{Text: output}},
	})
	state.LastPlan = state.CurrentPlan
	state.CurrentPlan = &plan
	state.PlanIterations += 1

	slog.Info("planner ends", "plan_iterations", state.PlanIterations)

	return nextStep, output, nil
}
