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

func GenerateJSON(ctx context.Context, llm llms.Model, prompt string, result any, fields []string, maxRetries int, options ...llms.CallOption) error {
	var retErr error
	for i := 0; i < maxRetries; i++ {
		resp, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt, append(options, llms.WithJSONMode())...)
		if err != nil {
			slog.Error("generate from single prompt", "error", err)
			retErr = err
			continue
		}

		resp = util.FixJSON(resp)

		if err := json.Unmarshal([]byte(resp), result); err != nil {
			slog.Error("unmarshal json", "error", err)
			retErr = err
			continue
		}

		validator := validator.New()
		if err := validator.StructPartial(result, fields...); err != nil {
			slog.Error("validate json", "error", err)
			retErr = err
			continue
		}
	}
	return fmt.Errorf("generate json: %w", retErr)
}
