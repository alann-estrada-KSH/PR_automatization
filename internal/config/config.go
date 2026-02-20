package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all runtime configuration for prgen.
type Config struct {
	// LLM provider: "ollama" | "openai" | "groq" | "openrouter" | "mock"
	Provider string `mapstructure:"provider"`
	Model    string `mapstructure:"model"`

	// Ollama
	OllamaURL string `mapstructure:"ollama_url"`

	// OpenAI-compatible (Groq, OpenRouter, OpenAI, …)
	APIKey     string `mapstructure:"api_key"`
	APIBaseURL string `mapstructure:"api_base_url"`

	// Prompts
	Prompts PromptConfig `mapstructure:"prompts"`

	// Output
	Output OutputConfig `mapstructure:"output"`

	// Debug / dry-run
	Debug  bool `mapstructure:"debug"`
}

type PromptConfig struct {
	Base  string `mapstructure:"base"`
	Extra string `mapstructure:"extra"`
}

type OutputConfig struct {
	SavePath        string `mapstructure:"save_path"`
	CopyToClipboard bool   `mapstructure:"copy_to_clipboard"`
}

// Load reads config from (in order of priority, highest last wins):
//  1. Built-in defaults
//  2. config.yaml next to the binary (or current dir)
//  3. ~/.prgen/config.yaml
//  4. Environment variables
//  5. CLI flags (applied by the caller after Load)
func Load() (*Config, error) {
	v := viper.New()

	// ── Defaults ──────────────────────────────────────────────────────────
	v.SetDefault("provider", "ollama")
	v.SetDefault("model", "llama3.1")
	v.SetDefault("ollama_url", "http://localhost:11434")
	v.SetDefault("api_key", "")
	v.SetDefault("api_base_url", "")
	v.SetDefault("prompts.base", "prompts/base.md")
	v.SetDefault("prompts.extra", filepath.Join(userHome(), ".prgen", "extra_prompt.md"))
	v.SetDefault("output.save_path", filepath.Join(userHome(), "KSH", "Projects"))
	v.SetDefault("output.copy_to_clipboard", true)
	v.SetDefault("debug", false)

	// ── Config files ──────────────────────────────────────────────────────
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	// Search order: current dir → binary dir → ~/.prgen
	v.AddConfigPath(".")
	v.AddConfigPath(binaryDir())
	v.AddConfigPath(filepath.Join(userHome(), ".prgen"))

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// No config file found — use defaults + env
	}

	// ── Environment variables ─────────────────────────────────────────────
	v.SetEnvPrefix("PRGEN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Explicit env mappings for common API key conventions
	bindEnv(v, "api_key", "PRGEN_API_KEY", "OPENAI_API_KEY", "GROQ_API_KEY")
	bindEnv(v, "api_base_url", "PRGEN_API_BASE_URL")
	bindEnv(v, "ollama_url", "PRGEN_OLLAMA_URL")
	bindEnv(v, "provider", "PRGEN_PROVIDER")
	bindEnv(v, "model", "PRGEN_MODEL")

	cfg := &Config{}
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}

	// Post-process: expand ~ in paths
	cfg.Prompts.Base = expandPath(cfg.Prompts.Base)
	cfg.Prompts.Extra = expandPath(cfg.Prompts.Extra)
	cfg.Output.SavePath = expandPath(cfg.Output.SavePath)

	return cfg, nil
}

// ApplyProviderDefaults sets sane API base URL defaults for known providers.
func ApplyProviderDefaults(cfg *Config) {
	if cfg.APIBaseURL != "" {
		return
	}
	switch strings.ToLower(cfg.Provider) {
	case "groq":
		cfg.APIBaseURL = "https://api.groq.com/openai/v1"
	case "openai":
		cfg.APIBaseURL = "https://api.openai.com/v1"
	case "openrouter":
		cfg.APIBaseURL = "https://openrouter.ai/api/v1"
	}
}

// ── helpers ──────────────────────────────────────────────────────────────────

func userHome() string {
	h, _ := os.UserHomeDir()
	return h
}

func binaryDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

func expandPath(p string) string {
	if strings.HasPrefix(p, "~/") || p == "~" {
		return filepath.Join(userHome(), p[1:])
	}
	return p
}

// bindEnv tries each env key in order and binds the first one that is set.
func bindEnv(v *viper.Viper, key string, envKeys ...string) {
	for _, ek := range envKeys {
		if val := os.Getenv(ek); val != "" {
			_ = v.BindEnv(key, ek)
			return
		}
	}
	// Bind the first one anyway so AutomaticEnv picks it up later
	if len(envKeys) > 0 {
		_ = v.BindEnv(key, envKeys[0])
	}
}
