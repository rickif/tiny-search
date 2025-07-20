package tool

import (
	"context"
	"log/slog"
	"os/exec"

	"github.com/tmc/langchaingo/llms"
)

var BashTool = llms.Tool{
	Type: "function",
	Function: &llms.FunctionDefinition{
		Name:        "bash",
		Description: "Use this to execute bash command and do necessary operations.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"cmd": map[string]any{
					"type":        "string",
					"description": "The bash command to execute.",
				},
			},
			"required": []string{"cmd"},
		},
	},
}

func Bash(ctx context.Context, cmd string) (string, error) {
	output, err := exec.Command(cmd).Output()
	if err != nil {
		slog.Error("execute bash command", "error", err, "command", cmd)
		return "", err
	}
	return string(output), nil
}
