package agent

import "github.com/tmc/langchaingo/llms"

type AgentState struct {
	Messages []llms.MessageContent
}
