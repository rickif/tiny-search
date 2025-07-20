package tool

import (
	"context"
	"log/slog"
	"os/exec"

	"github.com/tmc/langchaingo/llms"
)

var PythonTool = llms.Tool{
	Type: "function",
	Function: &llms.FunctionDefinition{
		Name:        "python-executor",
		Description: "Use this to execute python code and do data analysis or calculation. If you want to see the output of a value, you should print it out with `print(...)`. This is visible to the user.",
		Parameters: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"code": map[string]any{
					"type":        "string",
					"description": "The python code to execute to do further analysis or calculation.",
				},
			},
			"required": []string{"code"},
		},
	},
}

func Python(ctx context.Context, code string) (string, error) {
	output, err := exec.Command("python", "-c", code).Output()
	if err != nil {
		slog.Error("execute python command", "error", err, "command", code)
		return "", err
	}
	return string(output), nil
}
