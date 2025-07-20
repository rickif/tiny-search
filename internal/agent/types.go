package agent

const (
	StepEnd     = "__end__"
	StepPlanner = "__planner__"
)

type AgentUpdate struct {
	NextStep string
	Output   string
}
