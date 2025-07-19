package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	LLMModel   string
	LLMBaseURL string
	LLMToken   string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load env file: %w", err)
	}

	return &Config{
		LLMModel:   os.Getenv("LLM_MODEL"),
		LLMBaseURL: os.Getenv("LLM_BASE_URL"),
		LLMToken:   os.Getenv("LLM_TOKEN"),
	}, nil
}
