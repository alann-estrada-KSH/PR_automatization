// Package llm defines the Client interface and the factory function that
// returns the right implementation based on the config.
package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/alann-estrada-KSH/ai-pr-generator/internal/config"
)

// Client is the abstraction for any LLM backend.
type Client interface {
	// Generate sends the prompt and returns the full response.
	Generate(ctx context.Context, prompt string) (string, error)
	// Name returns a human-readable provider identifier.
	Name() string
}

// NewClient returns the Client implementation selected by cfg.Provider.
func NewClient(cfg *config.Config) (Client, error) {
	config.ApplyProviderDefaults(cfg)

	switch strings.ToLower(cfg.Provider) {
	case "ollama":
		return newOllamaClient(cfg), nil
	case "openai", "groq", "openrouter":
		return newOpenAIClient(cfg)
	case "mock":
		return &MockClient{}, nil
	default:
		return nil, fmt.Errorf("unknown LLM provider %q — supported: ollama, openai, groq, openrouter, mock", cfg.Provider)
	}
}
