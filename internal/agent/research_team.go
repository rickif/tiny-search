package agent

import (
	"context"

	"github.com/tmc/langchaingo/llms/openai"
)

var _ Node = (*ResearchTeam)(nil)

type ResearchTeam struct {
	llm *openai.LLM
}

func NewResearchTeam(llm *openai.LLM) *ResearchTeam {
	return &ResearchTeam{
		llm: llm,
	}
}

func (planner *ResearchTeam) Execute(ctx context.Context, state *AgentState) (nextStep string, output string, err error) {

	if state.CurrentPlan == nil {
		return StepPlanner, "", nil
	}

	var step *Step
	for i := range state.CurrentPlan.Steps {
		if state.CurrentPlan.Steps[i].ExecutionResult == "" {
			step = &state.CurrentPlan.Steps[i]
			break
		}
	}

	if step == nil {
		return StepPlanner, "", nil
	}

	switch step.StepType {
	case StepTypeReasearch:
		return StepResearcher, "", nil
	case StepTypeProcessing:
		return StepCoder, "", nil
	default:
		return StepPlanner, "", nil
	}
}
