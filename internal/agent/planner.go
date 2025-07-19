package agent

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/rickif/tiny-research/internal/config"
	"github.com/rickif/tiny-research/internal/llm"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

type PlanOptions struct {
	Question     string
	AllowReflect bool
	AllowAnswer  bool
	AllowRead    bool
	AllowSearch  bool
	BeastMode    bool
	AgentState   AgentState
}

type Planner struct {
	llm *openai.LLM
}

func NewPlanner(config *config.Config) (*Planner, error) {
	llm, err := openai.New(openai.WithBaseURL(config.LLMBaseURL), openai.WithModel(config.LLMModel), openai.WithToken(config.LLMToken))
	if err != nil {
		return nil, fmt.Errorf("create planner: %w", err)
	}
	return &Planner{llm: llm}, nil
}

func (planner *Planner) buildPrompt(options PlanOptions) (outputPrompt string, outputFieldNames []string, err error) {
	var promptSections []string

	template, err := os.ReadFile("../prompts/planner/header.md")
	if err != nil {
		return "", nil, fmt.Errorf("read header prompt file: %w", err)
	}
	prompt, err := prompts.NewPromptTemplate(
		string(template),
		[]string{"current_time", "question"},
	).Format(map[string]any{
		"current_time": time.Now().Format("2006-01-02 15:04:05"),
		"question":     options.Question,
	})
	if err != nil {
		return "", nil, fmt.Errorf("format header prompt: %w", err)
	}
	promptSections = append(promptSections, prompt)

	if len(options.AgentState.Contexts) > 0 {
		template, err = os.ReadFile("../prompts/planner/context.md")
		if err != nil {
			return "", nil, fmt.Errorf("read context prompt file: %w", err)
		}

		contextPrompt, err := prompts.NewPromptTemplate(
			string(template),
			[]string{"context"},
		).Format(map[string]any{
			"context": strings.Join(options.AgentState.Contexts, "\n"),
		})
		if err != nil {
			return "", nil, fmt.Errorf("format context prompt: %w", err)
		}
		promptSections = append(promptSections, contextPrompt)
	}

	if len(options.AgentState.Knowledges) > 0 {
		var knowledgePrompts []string
		template, err = os.ReadFile("../prompts/planner/knowledge.md")
		if err != nil {
			return "", nil, fmt.Errorf("read knowledge prompt file: %w", err)
		}
		for _, knowledge := range options.AgentState.Knowledges {
			knowledgePrompt, _ := prompts.NewPromptTemplate(
				string(template),
				[]string{"question", "answer", "reference"},
			).Format(map[string]any{
				"question":  knowledge.Question,
				"answer":    knowledge.Answer,
				"reference": knowledge.Reference,
			})
			knowledgePrompts = append(knowledgePrompts, knowledgePrompt)
		}

		template, err = os.ReadFile("../prompts/planner/knowledges.md")
		if err != nil {
			return "", nil, fmt.Errorf("read knowledges prompt file: %w", err)
		}

		knowledgePrompt, err := prompts.NewPromptTemplate(
			string(template),
			[]string{"knowledges"},
		).Format(map[string]any{
			"knowledges": strings.Join(knowledgePrompts, "\n"),
		})
		if err != nil {
			return "", nil, fmt.Errorf("format knowledges prompt: %w", err)
		}
		promptSections = append(promptSections, knowledgePrompt)
	}

	if len(options.AgentState.BadAttempts) > 0 {
		template, err = os.ReadFile("../prompts/planner/bad-attempt.md")
		if err != nil {
			return "", nil, fmt.Errorf("read bad attempt prompt file: %w", err)
		}

		var badAttemptPrompts []string
		for _, attempt := range options.AgentState.BadAttempts {
			badAttemptPrompt, err := prompts.NewPromptTemplate(
				string(template),
				[]string{"question", "answer", "reject_reason", "recap", "blame"},
			).Format(map[string]any{
				"question":      attempt.Question,
				"answer":        attempt.Answer,
				"reject_reason": attempt.RejectReason,
				"recap":         attempt.ActionRecap,
				"blame":         attempt.ActionBlame,
			})
			if err != nil {
				return "", nil, fmt.Errorf("format bad attempt prompt: %w", err)
			}
			badAttemptPrompts = append(badAttemptPrompts, badAttemptPrompt)
		}

		var improvments []string
		for _, attempt := range options.AgentState.BadAttempts {
			improvments = append(improvments, attempt.Improvement)
		}

		template, err = os.ReadFile("../prompts/planner/bad-attempts.md")
		if err != nil {
			return "", nil, fmt.Errorf("read bad attempts prompt file: %w", err)
		}
		badAttemptsPrompt, err := prompts.NewPromptTemplate(
			string(template),
			[]string{"bad-attempts", "learned-strategy"},
		).Format(map[string]any{
			"bad_attempts":     strings.Join(badAttemptPrompts, "\n"),
			"learned_strategy": strings.Join(improvments, "\n"),
		})
		if err != nil {
			return "", nil, fmt.Errorf("format bad attempts prompt: %w", err)
		}
		promptSections = append(promptSections, badAttemptsPrompt)
	}

	var actionPromptSections []string
	if options.AllowRead && len(options.AgentState.ActionState.AllURLs) > 0 {
		template, err = os.ReadFile("../prompts/planner/action-visit.md")
		if err != nil {
			return "", nil, fmt.Errorf("read action visit prompt file: %w", err)
		}

		actionVisitPrompt, err := prompts.NewPromptTemplate(
			string(template),
			[]string{"urls"},
		).Format(map[string]any{
			"urls": options.AgentState.ActionState.AllURLs,
		})
		if err != nil {
			return "", nil, fmt.Errorf("format action visit prompt: %w", err)
		}
		actionPromptSections = append(actionPromptSections, actionVisitPrompt)
	}

	if options.AllowSearch {
		prompt, err := os.ReadFile("../prompts/planner/action-search.md")
		if err != nil {
			return "", nil, fmt.Errorf("read action search prompt file: %w", err)
		}
		actionPromptSections = append(actionPromptSections, string(prompt))
		outputFieldNames = append(outputFieldNames, "search")
	}

	if options.AllowAnswer {
		var prompt []byte
		var err error
		if options.AllowReflect {
			prompt, err = os.ReadFile("../prompts/planner/action-answer-with-reflect.md")
		} else {
			prompt, err = os.ReadFile("../prompts/planner/action-answer.md")
		}
		if err != nil {
			return "", nil, fmt.Errorf("read action answer prompt file: %w", err)
		}
		actionPromptSections = append(actionPromptSections, string(prompt))
		outputFieldNames = append(outputFieldNames, "answer")
	}

	if options.BeastMode {
		prompt, err := os.ReadFile("../prompts/planner/action-answer-beast-mode.md")
		if err != nil {
			return "", nil, fmt.Errorf("read action answer beast mode prompt file: %w", err)
		}
		actionPromptSections = append(actionPromptSections, string(prompt))
	}

	if options.AllowReflect {
		prompt, err := os.ReadFile("../prompts/planner/action-reflect.md")
		if err != nil {
			return "", nil, fmt.Errorf("read action reflect prompt file: %w", err)
		}
		actionPromptSections = append(actionPromptSections, string(prompt))
		outputFieldNames = append(outputFieldNames, "reflect")
	}

	if len(actionPromptSections) > 0 {
		template, err = os.ReadFile("../prompts/planner/action-selections.md")
		if err != nil {
			return "", nil, fmt.Errorf("read action selections prompt file: %w", err)
		}

		actionsPrompt, err := prompts.NewPromptTemplate(
			string(template),
			[]string{"action_selections"},
		).Format(map[string]any{
			"action_selections": strings.Join(actionPromptSections, "\n\n"),
		})
		if err != nil {
			return "", nil, fmt.Errorf("format action selections prompt: %w", err)
		}
		promptSections = append(promptSections, actionsPrompt)
	}

	footerPrompt, err := os.ReadFile("../prompts/planner/footer.md")
	if err != nil {
		return "", nil, fmt.Errorf("read footer prompt file: %w", err)
	}
	promptSections = append(promptSections, string(footerPrompt))

	return strings.Join(promptSections, "\n\n"), nil, nil
}

// Test only
func (planner *Planner) BuildPrompt(options PlanOptions) (string, error) {
	prompt, err := planner.buildPrompt(options)
	if err != nil {
		return "", fmt.Errorf("build prompt: %w", err)
	}
	return prompt, nil
}

type Plan struct {
}

func (planner *Planner) Plan(options PlanOptions) {
	prompt, err := planner.buildPrompt(options)
	if err != nil {
		slog.Error("build prompt", "error", err)
		return
	}

	llm.GenerateJSON(context.Background(), planner.llm, prompt, &Plan{}, []string{}, 3)
}
