package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/rickif/tiny-research/internal/tool"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

var _ Node = (*Coder)(nil)

type Coder struct {
	llm *openai.LLM
}

func NewCoder(llm *openai.LLM) *Coder {
	return &Coder{
		llm: llm,
	}
}

func (r *Coder) Execute(ctx context.Context, state *AgentState) (nextStep string, output string, err error) {
	slog.Info("coder starts")
	content, err := os.ReadFile("./internal/prompts/coder.md")
	if err != nil {
		slog.Error("read prompt template", "error", err)
		return "", "", err
	}

	promptTemplate, err := prompts.NewPromptTemplate(
		string(content),
		[]string{"current_time"},
	).Format(map[string]any{
		"current_time": time.Now().Format(time.RFC3339),
		"locale":       state.Locale,
	})
	if err != nil {
		slog.Error("format prompt", "error", err)
		return "", "", err
	}

	var step *Step
	for i := range state.CurrentPlan.Steps {
		if state.CurrentPlan.Steps[i].ExecutionResult == "" {
			step = &state.CurrentPlan.Steps[i]
			break
		}
	}

	var existingFindings []string
	for i, step := range state.CurrentPlan.Steps {
		if step.ExecutionResult != "" {
			existingFindings = append(existingFindings, fmt.Sprintf("## Existing Finding %d: %s\n\n<finding>%s</finding>", i+1, step.Title, step.ExecutionResult))
		}
	}

	messages := []llms.MessageContent{
		{
			Role:  llms.ChatMessageTypeSystem,
			Parts: []llms.ContentPart{llms.TextContent{Text: promptTemplate}},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextContent{Text: fmt.Sprintf("#Task\n\ntitle: %s\n\n##description:%s", step.Title, step.Description)}},
		},
		{
			Role:  llms.ChatMessageTypeHuman,
			Parts: []llms.ContentPart{llms.TextContent{Text: fmt.Sprintf("#Existing Findings\n\n%s", strings.Join(existingFindings, "\n\n"))}},
		},
	}

	for {
		resp, err := r.llm.GenerateContent(ctx, messages, llms.WithTools([]llms.Tool{tool.PythonTool}))
		if err != nil {
			slog.Error("generate content", "error", err)
			return "", "", err
		}

		for _, toolcall := range resp.Choices[0].ToolCalls {
			switch toolcall.FunctionCall.Name {
			case "python-executor":
				var args struct {
					Code string `json:"code"`
				}
				if err := json.Unmarshal([]byte(toolcall.FunctionCall.Arguments), &args); err != nil {
					slog.Error("unmarshal arguments", "error", err)
					return "", "", err
				}
				output, err = tool.Python(ctx, args.Code)
				if err != nil {
					slog.Error("python", "error", err)
					return "", "", err
				}
				message := llms.MessageContent{
					Role: llms.ChatMessageTypeTool,
					Parts: []llms.ContentPart{
						llms.ToolCallResponse{
							ToolCallID: toolcall.ID,
							Name:       toolcall.FunctionCall.Name,
							Content:    output,
						},
					},
				}
				messages = append(messages, message)
				slog.Info("coder use python")
			default:
				slog.Error("unexpected function call", "name", toolcall.FunctionCall.Name)
				return "", "", fmt.Errorf("unexpected function call: %v", toolcall.FunctionCall.Name)
			}
		}

		if len(resp.Choices[0].ToolCalls) == 0 {
			slog.Info("coder finished")
			step.ExecutionResult = resp.Choices[0].Content
			break
		}
		var toolCalls []string
		for _, toolcall := range resp.Choices[0].ToolCalls {
			toolCalls = append(toolCalls, toolcall.FunctionCall.Name)
		}
		slog.Info("coder use tools", "tools", toolCalls)
	}
	state.Messages = append(state.Messages, llms.MessageContent{
		Role:  llms.ChatMessageTypeHuman,
		Parts: []llms.ContentPart{llms.TextContent{Text: step.ExecutionResult}},
	})
	return StepResearchTeam, output, nil
}
