package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/rickif/tiny-research/util"
	"github.com/tmc/langchaingo/llms"
)

func GenerateJSON(ctx context.Context, llm llms.Model, messages []llms.MessageContent, result any, maxRetries int, options ...llms.CallOption) (output string, err error) {
	var resp *llms.ContentResponse
	for i := 0; i < maxRetries; i++ {
		resp, err = llm.GenerateContent(ctx, messages, append(options, llms.WithJSONMode())...)
		if err != nil {
			slog.Error("generate from single prompt", "error", err)
			continue
		}
		output = resp.Choices[0].Content

		output = util.FixJSON(output)

		if err = json.Unmarshal([]byte(output), result); err != nil {
			slog.Error("unmarshal json", "error", err)
			continue
		}

		validator := validator.New()
		if err = validator.Struct(result); err != nil {
			slog.Error("validate json", "error", err)
			continue
		}
		return output, nil
	}
	return "", fmt.Errorf("generate json: %w", err)
}
