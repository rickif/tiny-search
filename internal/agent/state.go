package agent

type AgentState struct {
	Contexts    []string
	Knowledges  []Knowledge
	BadAttempts []BadAttempt
	ActionState ActionState
}

type Knowledge struct {
	Question  string
	Answer    string
	Reference string
}

type BadAttempt struct {
	Question     string
	Answer       string
	RejectReason string
	ActionRecap  string
	ActionBlame  string
	Improvement  string
}

type ActionState struct {
	AllURLs []string
}
