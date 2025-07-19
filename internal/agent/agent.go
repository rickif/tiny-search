package agent

import "context"

type ResearchQuery struct {
	Query       string
	TokenBudget int
	BadAttempts int
}

type Agent struct {
}

func (a *Agent) Research(ctx context.Context, query ResearchQuery) (string, error) {
	return "", nil
}
