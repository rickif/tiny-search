package agent

import "github.com/tmc/langchaingo/llms"

const (
	StepTypeReasearch  = "research"
	StepTypeProcessing = "processing"
)

type Step struct {
	NeedWebSearch   bool   `json:"need_web_search"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	StepType        string `json:"step_type"`
	ExecutionResult string `json:"_"`
}

type Plan struct {
	HasEnoughContext bool   `json:"has_enough_context"`
	Thought          string `json:"thought"`
	Title            string `json:"title"`
	Steps            []Step `json:"steps"`
}

type AgentState struct {
	Messages       []llms.MessageContent
	LastPlan       *Plan
	CurrentPlan    *Plan
	PlanIterations int
}

const (
	StepEnd          = "__end__"
	StepPlanner      = "__planner__"
	StepResearchTeam = "__research_team__"
	StepReporter     = "__reporter__"
	StepResearcher   = "__researcher__"
	StepCoder        = "__coder__"
)
