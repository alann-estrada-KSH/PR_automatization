package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/alann-estrada-KSH/ai-pr-generator/internal/config"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/version"
)

// Execute is the single entry point called from main.
func Execute() {
	if err := newRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "prgen",
		Short: "AI-powered PR description generator",
		Long: fmt.Sprintf(`prgen %s
─────────────────────────────────────────────────────────────
Generate detailed Pull Request descriptions from git commits
using your preferred LLM (Ollama, Groq, OpenAI, OpenRouter).

Examples:
  prgen                          # generate PR from last commit
  prgen generate --commits 3     # use last 3 commits
  prgen generate --provider groq # override LLM provider
  prgen generate --dump-prompt   # print prompt and exit
  prgen generate --dry-run       # skip LLM call
  prgen version                  # show version
  prgen update                   # update from git
  prgen config                   # show active config
─────────────────────────────────────────────────────────────`, version.Version),
		// When called without a subcommand, run generate
		RunE: func(cmd *cobra.Command, args []string) error {
			genCmd := newGenerateCmd()
			genCmd.SetArgs(args)
			return genCmd.Execute()
		},
	}

	root.AddCommand(
		newGenerateCmd(),
		newCommitCmd(),
		newReviewCmd(),
		newBranchCmd(),
		newVersionCmd(),
		newUpdateCmd(),
		newConfigCmd(),
	)

	return root
}

// mustLoadConfig loads config and exits on error.
func mustLoadConfig() *config.Config {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, " ❌ Error loading config: %v\n", err)
		os.Exit(1)
	}
	return cfg
}
