package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alann-estrada-KSH/ai-pr-generator/internal/llm"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/ui"
)

func newBranchCmd() *cobra.Command {
	var (
		provider string
		model    string
	)

	cmd := &cobra.Command{
		Use:   "branch [description]",
		Short: "Suggest branch names from a description",
		Long: `Given a description of what you're going to work on, suggests several
branch names following common naming conventions:

  feature/<ticket?>-<short-desc>
  fix/<ticket?>-<short-desc>
  hotfix/<ticket?>-<short-desc>
  refactor/<scope>-<short-desc>
  chore/<scope>-<short-desc>

Examples:
  prgen branch "agregar login con google oauth"
  prgen branch "corregir error 500 en endpoint de pagos"
  prgen branch "refactorizar el módulo de permisos"`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustLoadConfig()
			if provider != "" {
				cfg.Provider = provider
			}
			if model != "" {
				cfg.Model = model
			}

			description := strings.Join(args, " ")

			prompt := buildBranchPrompt(description)

			prog := ui.New()
			llmClient, err := llm.NewClient(cfg)
			if err != nil {
				return fmt.Errorf("inicializando LLM: %w", err)
			}

			prog.Start(fmt.Sprintf("Generando sugerencias con %s...", llmClient.Name()))
			prog.Update("Pensando en nombres...", 50)

			result, err := llmClient.Generate(context.Background(), prompt)
			if err != nil {
				prog.Stop("Error en LLM", false)
				return fmt.Errorf("generación LLM:\n%w", err)
			}
			prog.Stop("Sugerencias listas", true)

			fmt.Println()
			fmt.Println("┌─ Sugerencias de rama ──────────────────────────────────────")
			fmt.Println(strings.TrimSpace(result))
			fmt.Println("└────────────────────────────────────────────────────────────")
			fmt.Println()
			fmt.Println("Para crear la rama elegida:")
			fmt.Println("  git checkout -b <nombre-elegido>")
			fmt.Println()

			return nil
		},
	}

	cmd.Flags().StringVarP(&provider, "provider", "p", "", "LLM provider override")
	cmd.Flags().StringVarP(&model, "model", "m", "", "Model override")

	return cmd
}

func buildBranchPrompt(description string) string {
	return fmt.Sprintf(`You are an expert in Git and development team conventions.
Your task is to suggest 5 branch names for the following task.

TASK DESCRIPTION:
%s

RULES:
1. Respond ONLY with the list of branch names. No explanations or extra text.
2. Branch names must be in ENGLISH. Use kebab-case (lowercase, hyphens). No spaces or extra slashes.
3. Valid prefixes (choose the most accurate):
   feature/   → new feature
   fix/        → bug fix
   hotfix/     → urgent production fix
   refactor/   → code change without new functionality
   chore/      → maintenance, deps, scripts, config
   docs/       → documentation only
   test/       → only tests
   ci/         → pipelines or CI/CD configuration
4. If the description mentions a ticket number (e.g., TK-123, JIRA-456), include it after the prefix.
5. Maximum 50 characters per branch name.
6. Order from most to least descriptive.
7. Next to each option show: → git checkout -b <nombre>

RESPONSE FORMAT (no markdown, plain text only):
  feature/description-in-english          → git checkout -b feature/description-in-english
  fix/description-in-english              → git checkout -b fix/description-in-english
  ...
`, description)
}
