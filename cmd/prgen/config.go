package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Show the active configuration",
		Long:  `Prints the merged configuration (defaults → config.yaml → env vars → CLI flags).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustLoadConfig()
			fmt.Printf(`
 ⚙️  prgen — active configuration
 ─────────────────────────────────────────
  Provider:           %s
  Model:              %s
  Ollama URL:         %s
  API Base URL:       %s
  API Key:            %s

  Prompt base:        %s
  Prompt extra:       %s

  Output save path:   %s
  Copy to clipboard:  %v
  Debug:              %v
 ─────────────────────────────────────────

 Tip: override any value via:
   • ~/.prgen/config.yaml
   • Environment variables (PRGEN_PROVIDER, PRGEN_API_KEY, …)
   • CLI flags (--provider, --model, …)
`,
				cfg.Provider,
				cfg.Model,
				cfg.OllamaURL,
				cfg.APIBaseURL,
				maskKey(cfg.APIKey),
				cfg.Prompts.Base,
				cfg.Prompts.Extra,
				cfg.Output.SavePath,
				cfg.Output.CopyToClipboard,
				cfg.Debug,
			)
			return nil
		},
	}
}

func maskKey(key string) string {
	if key == "" {
		return "(not set)"
	}
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}
