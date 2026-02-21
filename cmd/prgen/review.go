package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alann-estrada-KSH/ai-pr-generator/internal/git"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/llm"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/ui"
)

func newReviewCmd() *cobra.Command {
	var (
		numCommits int
		from       string
		to         string
		provider   string
		model      string
		dryRun     bool
	)

	cmd := &cobra.Command{
		Use:   "review",
		Short: "AI code review: bugs, security, error handling, refactor suggestions",
		Long: `Analyzes the diff and produces a structured code review report covering:
  - Possible bugs and unhandled edge cases
  - Security issues (SQL injection, secrets, missing auth)  
  - Error handling gaps
  - Refactor suggestions

Examples:
  prgen review                      # review last commit
  prgen review --commits 3          # review last 3 commits
  prgen review --from develop       # review changes since develop
  prgen review --from develop --to feature/auth  # specific range`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustLoadConfig()
			if provider != "" {
				cfg.Provider = provider
			}
			if model != "" {
				cfg.Model = model
			}

			prog := ui.New()

			// ── Collect diff ──────────────────────────────────────────────
			prog.Start("Leyendo cambios para revisión...")
			prog.Update("Leyendo cambios...", 15)

			var stats, diff, ref string

			if from != "" {
				effTo := to
				if effTo == "" {
					effTo = "HEAD"
				}
				stats = git.StatBetween(from, effTo)
				diff = git.FilteredDiffBetween(from, effTo, cfg.Diff.Ignore)
				ref = fmt.Sprintf("%s...%s", from, effTo)
			} else {
				stats = git.DiffStat(numCommits)
				diff = git.FilteredDiff(numCommits, cfg.Diff.Ignore)
				ref = fmt.Sprintf("últimos %d commits", numCommits)
			}

			// Truncate diff if it exceeds max_chars
			if len(diff) > cfg.Diff.MaxChars {
				diff = diff[:cfg.Diff.MaxChars] + "\n\n...(diff truncado por límite de configuración)..."
			}
			prog.Stop(fmt.Sprintf("Cambios leídos (%s)", ref), true)

			// ── Build prompt ──────────────────────────────────────────────
			promptText := buildReviewPrompt(cfg.Prompts.Review, ref, stats, diff)

			if dryRun {
				fmt.Println(promptText)
				return nil
			}

			// ── LLM call ─────────────────────────────────────────────────
			llmClient, err := llm.NewClient(cfg)
			if err != nil {
				return fmt.Errorf("inicializando LLM: %w", err)
			}

			prog.Start(fmt.Sprintf("Analizando con %s...", llmClient.Name()))
			prog.Update("Revisando código...", 40)

			review, err := llmClient.Generate(context.Background(), promptText)
			if err != nil {
				prog.Stop("Error en LLM", false)
				return fmt.Errorf("generación LLM:\n%w", err)
			}
			prog.Stop("Revisión completa", true)

			// ── Print ─────────────────────────────────────────────────────
			fmt.Println()
			fmt.Println(strings.TrimSpace(review))
			fmt.Println()

			return nil
		},
	}

	cmd.Flags().IntVarP(&numCommits, "commits", "c", 1, "Number of commits to review")
	cmd.Flags().StringVar(&from, "from", "", "Base ref/branch to compare from")
	cmd.Flags().StringVar(&to, "to", "HEAD", "Target ref/branch to compare to")
	cmd.Flags().StringVarP(&provider, "provider", "p", "", "LLM provider override")
	cmd.Flags().StringVarP(&model, "model", "m", "", "Model override")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print prompt and exit")

	return cmd
}

func buildReviewPrompt(promptPath, ref, stats, diff string) string {
	data, err := os.ReadFile(promptPath)
	if err != nil {
		return fmt.Sprintf(`Eres un senior code reviewer. Analiza el siguiente diff y reporta:
1. Posibles bugs o casos no manejados
2. Problemas de seguridad (SQL injection, secrets, auth)
3. Manejo de errores faltante
4. Sugerencias de refactor

Ref: %s
Archivos: %s
Diff:
%s`, ref, stats, diff)
	}

	tmpl := string(data)
	tmpl = strings.ReplaceAll(tmpl, "{{.ProjectType}}", "Go")
	tmpl = strings.ReplaceAll(tmpl, "{{.Branch}}", ref)
	tmpl = strings.ReplaceAll(tmpl, "{{.Stats}}", stats)
	tmpl = strings.ReplaceAll(tmpl, "{{.Diff}}", diff)
	return tmpl
}
