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

func newCommitCmd() *cobra.Command {
	var (
		apply    bool
		provider string
		model    string
		dryRun   bool
	)

	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Generate a Conventional Commit message from staged changes",
		Long: `Analyzes the currently staged changes (git diff --cached) and generates
a commit message following the Conventional Commits specification.

Types used: feat | fix | refactor | chore | docs | style | test | perf | ci | build

Example:
  git add .
  prgen commit          # preview the suggested message
  prgen commit --apply  # generate and run git commit directly`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustLoadConfig()
			if provider != "" {
				cfg.Provider = provider
			}
			if model != "" {
				cfg.Model = model
			}

			// ── Check staged changes ──────────────────────────────────────
			if !git.HasStagedChanges() {
				ui.Warnf("No hay cambios staged. Usa 'git add <archivos>' primero.")
				return nil
			}

			prog := ui.New()

			// ── Collect staged diff ───────────────────────────────────────
			prog.Start("Leyendo cambios staged...")
			prog.Update("Leyendo cambios staged...", 20)
			stats := git.StagedStat()
			diff := git.FilteredStagedDiff(cfg.Diff.Ignore)

			// Truncate diff if it exceeds max_chars
			if len(diff) > cfg.Diff.MaxChars {
				diff = diff[:cfg.Diff.MaxChars] + "\n\n...(diff truncado por límite de configuración)..."
			}
			prog.Stop("Cambios staged leídos", true)

			// ── Load commit prompt ────────────────────────────────────────
			promptText := buildCommitPrompt(cfg.Prompts.Commit, cfg.Provider, stats, diff)

			if dryRun {
				ui.Warnf("Dry-run: prompt que se enviaría al LLM:\n\n%s", promptText)
				return nil
			}

			// ── LLM call ─────────────────────────────────────────────────
			llmClient, err := llm.NewClient(cfg)
			if err != nil {
				return fmt.Errorf("inicializando LLM: %w", err)
			}

			prog.Start(fmt.Sprintf("Generando mensaje con %s...", llmClient.Name()))
			prog.Update(fmt.Sprintf("Generando mensaje con %s...", llmClient.Name()), 50)

			suggestion, err := llmClient.Generate(context.Background(), promptText)
			if err != nil {
				prog.Stop("Error en LLM", false)
				return fmt.Errorf("generación LLM:\n%w", err)
			}
			prog.Stop("Mensaje generado", true)

			// Clean up the suggestion (trim whitespace, extra backticks)
			suggestion = strings.TrimSpace(suggestion)
			suggestion = strings.Trim(suggestion, "`")
			suggestion = strings.TrimSpace(suggestion)

			// ── Show suggestion ───────────────────────────────────────────
			fmt.Println()
			fmt.Println("┌─ Mensaje sugerido ─────────────────────────────────────────")
			fmt.Println(suggestion)
			fmt.Println("└────────────────────────────────────────────────────────────")
			fmt.Println()

			// ── Apply ─────────────────────────────────────────────────────
			if apply {
				if err := runGitCommit(suggestion); err != nil {
					return fmt.Errorf("git commit: %w", err)
				}
				ui.Successf("Commit creado.")
			} else {
				fmt.Println("Para aplicarlo:")
				fmt.Printf("  git commit -m %q\n", strings.Split(suggestion, "\n")[0])
				fmt.Println("\nO ejecuta:")
				fmt.Println("  prgen commit --apply")
			}
			return nil
		},
	}

	cmd.Flags().BoolVar(&apply, "apply", false, "Run git commit with the generated message")
	cmd.Flags().StringVarP(&provider, "provider", "p", "", "LLM provider override")
	cmd.Flags().StringVarP(&model, "model", "m", "", "Model override")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print prompt and exit without calling LLM")

	return cmd
}

func buildCommitPrompt(promptPath string, provider string, stats, diff string) string {
	// Try to load from file
	data, err := os.ReadFile(promptPath)
	if err != nil {
		// Embedded fallback
		return fmt.Sprintf(`Genera un mensaje de commit en Conventional Commits para estos cambios staged.

Archivos: %s

Diff:
%s

Responde ÚNICAMENTE con el mensaje de commit. Sin explicaciones.`, stats, diff)
	}

	tmpl := string(data)
	tmpl = strings.ReplaceAll(tmpl, "{{.ProjectType}}", provider)
	tmpl = strings.ReplaceAll(tmpl, "{{.Stats}}", stats)
	tmpl = strings.ReplaceAll(tmpl, "{{.Diff}}", diff)
	return tmpl
}

func runGitCommit(message string) error {
	// Use -m for single-line, or write to temp file for multi-line
	lines := strings.Split(message, "\n")
	if len(lines) == 1 {
		_, err := git.Run("commit", "-m", message)
		return err
	}

	// Multi-line commit via temp file
	tmp, err := os.CreateTemp("", "prgen-commit-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.WriteString(message); err != nil {
		tmp.Close()
		return err
	}
	tmp.Close()

	_, err = git.Run("commit", "-F", tmp.Name())
	return err
}
