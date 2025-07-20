package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/rickif/tiny-research/internal/agent"
	"github.com/rickif/tiny-research/internal/config"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, _ := a.Value.Any().(*slog.Source)
				if source != nil {
					source.File = filepath.Base(source.File)
				}
			}
			return a
		}})))

	config, err := config.LoadConfig()
	if err != nil {
		slog.Error("load config", "error", err)
		return
	}
	agent, err := agent.NewAgent(config)
	if err != nil {
		slog.Error("new agent", "error", err)
		return
	}

	result, err := agent.Research(context.Background(), "What's the weather like in Chengdu today?")
	if err != nil {
		slog.Error("research", "error", err)
		return
	}

	fmt.Println(result)
}
