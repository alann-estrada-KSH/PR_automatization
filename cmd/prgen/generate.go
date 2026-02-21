package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/alann-estrada-KSH/ai-pr-generator/internal/checklist"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/cleaner"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/clipboard"
	prcontext "github.com/alann-estrada-KSH/ai-pr-generator/internal/context"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/detect"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/git"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/llm"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/prompt"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/ui"
)

func newGenerateCmd() *cobra.Command {
	var (
		numCommits       int
		notes            string
		notesFile        string
		interactiveNotes bool
		noClipboard      bool
		provider         string
		model            string
		debug            bool
		dryRun           bool
		dumpPrompt       bool
		tasks            string // comma-separated task IDs, e.g. "TK-123,TK-456"
		fromRef          string // --from: base branch/ref for diff
		toRef            string // --to: target branch/ref for diff
	)

	cmd := &cobra.Command{
		Use:     "generate",
		Aliases: []string{"gen", "g"},
		Short:   "Generate a PR description from recent git commits",
		Long: `Analyzes the last N git commits and generates a detailed PR description
using the configured LLM provider (Ollama, Groq, OpenAI, etc.).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := mustLoadConfig()

			// CLI flag overrides
			if provider != "" {
				cfg.Provider = provider
			}
			if model != "" {
				cfg.Model = model
			}
			if debug || viper.GetBool("debug") {
				cfg.Debug = true
			}

			prog := ui.New()

			// ── Step 1: Project detection ──────────────────────────────────
			prog.Start("Detectando proyecto...")
			pt := detect.FromCurrentDir()
			prog.Stop(fmt.Sprintf("Proyecto detectado: %s", strings.ToUpper(pt.String())), true)

			// ── Step 2: Git info ───────────────────────────────────────────
			prog.Start("Leyendo git log y diff...")
			prog.Update("Leyendo git log y diff...", 10)

			var logs, stats, diff, branch string
			headHash := git.HeadHash()

			if fromRef != "" {
				// Cross-branch mode
				effTo := toRef
				if effTo == "" {
					effTo = "HEAD"
				}
				logs = git.LogBetween(fromRef, effTo)
				stats = git.StatBetween(fromRef, effTo)
				diff = git.DiffBetween(fromRef, effTo)
				branch = fmt.Sprintf("%s...%s", fromRef, effTo)
			} else {
				branch = git.Branch()
				logs = git.Log(numCommits)
				stats = git.DiffStat(numCommits)
				diff = git.FilteredDiff(numCommits, cfg.Diff.Ignore)
			}

			// Truncate diff if it exceeds max_chars
			if len(diff) > cfg.Diff.MaxChars {
				diff = diff[:cfg.Diff.MaxChars] + "\n\n...(diff truncado por límite de configuración)..."
			}
			prog.Stop(fmt.Sprintf("Rama: %s | diff: %d chars", branch, len(diff)), true)

			if cfg.Debug {
				fmt.Println("\n[debug] Branch:", branch)
				fmt.Println("[debug] Hash:", headHash)
				fmt.Println("[debug] Stats:\n", stats)
			}

			// ── Step 3: Collect extra instructions ─────────────────────────
			var extraInstructions string
			switch {
			case notes != "":
				extraInstructions = notes
			case notesFile != "":
				p := &prcontext.FileProvider{Path: notesFile}
				extraInstructions, _ = p.GetContext()
			case interactiveNotes:
				// Must stop spinner before reading from stdin
				p := &prcontext.MultilineProvider{}
				extraInstructions, _ = p.GetContext()
			}

			// ── Task references: populated from --tasks flag (───────────────
			// When your own task system is ready, you can also populate taskCtx
			// programmatically and the section will appear automatically.
			var taskCtx string
			if tasks != "" {
				var lines []string
				for _, id := range strings.Split(tasks, ",") {
					id = strings.TrimSpace(id)
					if id != "" {
						lines = append(lines, "- "+id)
					}
				}
				taskCtx = strings.Join(lines, "\n")
			}

			// ── Step 4: Build prompt ───────────────────────────────────────
			prog.Start("Construyendo prompt...")
			prog.Update("Construyendo prompt...", 20)
			builder := &prompt.Builder{
				BasePath:  cfg.Prompts.Base,
				ExtraPath: cfg.Prompts.Extra,
			}
			fullPrompt, err := builder.Build(prompt.Context{
				ProjectType:       pt,
				Branch:            branch,
				Logs:              logs,
				Stats:             stats,
				Diff:              diff,
				ExtraInstructions: extraInstructions,
			})
			if err != nil {
				prog.Stop("Error construyendo prompt", false)
				return fmt.Errorf("building prompt: %w", err)
			}
			prog.Stop("Prompt construido", true)

			if dumpPrompt {
				fmt.Println(fullPrompt)
				return nil
			}

			if cfg.Debug {
				fmt.Println("\n[debug] Prompt length:", len(fullPrompt), "chars")
				fmt.Println("[debug] Provider:", cfg.Provider, "| Model:", cfg.Model)
			}

			// ── Step 5: LLM call ───────────────────────────────────────────
			var aiContent string
			if dryRun {
				ui.Warnf("Dry-run: saltando llamada al LLM")
				aiContent = "### [DRY-RUN — sin respuesta de IA] ###"
			} else {
				llmClient, err := llm.NewClient(cfg)
				if err != nil {
					return fmt.Errorf("inicializando cliente LLM: %w", err)
				}

				prog.Start(fmt.Sprintf("Generando PR con %s...", llmClient.Name()))
				prog.Update(fmt.Sprintf("Generando PR con %s...", llmClient.Name()), 30)

				aiContent, err = llmClient.Generate(context.Background(), fullPrompt)
				if err != nil {
					prog.Stop("Error en la llamada al LLM", false)
					// The error from the LLM providers already includes the API response body
					return fmt.Errorf("generación LLM fallida:\n%w", err)
				}
				prog.Stop("Respuesta recibida del LLM", true)
			}

			// ── Step 6: Clean output ───────────────────────────────────────
			prog.Start("Limpiando y formateando respuesta...")
			prog.Update("Limpiando y formateando respuesta...", 75)
			content := cleaner.Process(aiContent)

			// ── Inject task references (active when taskCtx is non-empty) ──
			if taskCtx != "" {
				taskSection := "## 🗂️ Referencias de tareas\n" + taskCtx
				// Insert before "¿Qué problema soluciona?" to keep logical flow
				if strings.Contains(content, "## 🔍") {
					content = strings.Replace(content, "## 🔍", taskSection+"\n\n## 🔍", 1)
				} else {
					content += "\n\n" + taskSection
				}
			}

			// ── Inject extra instructions section ─────────────────────────
			if extraInstructions != "" {
				notesSection := "## 📝 Instrucciones adicionales\n" + extraInstructions
				if strings.Contains(content, "## ⚠️ Consideraciones adicionales") {
					content = strings.Replace(content,
						"## ⚠️ Consideraciones adicionales",
						notesSection+"\n\n## ⚠️ Consideraciones adicionales", 1)
				} else {
					content += "\n\n" + notesSection
				}
			}

			// ── Append technical checklist + merge template ────────────────
			finalPR := strings.TrimSpace(content)
			finalPR += "\n\n## 🛠️ Cambios realizados\n" + checklist.Technical(pt, stats)
			finalPR += "\n\n" + checklist.Merge(pt)
			prog.Stop("Formato aplicado", true)

			// ── Step 7: Save to file ───────────────────────────────────────
			prog.Start("Guardando archivo...")
			prog.Update("Guardando archivo...", 90)
			repoName := filepath.Base(mustCwd())
			prFolder := filepath.Join(
				cfg.Output.SavePath,
				repoName+" - PR",
				time.Now().Format("02-01-2006"),
			)
			if err := os.MkdirAll(prFolder, 0755); err != nil {
				prog.Stop("Error creando carpeta de salida", false)
				return fmt.Errorf("creating output folder: %w", err)
			}

			shortHash := headHash
			if len(shortHash) > 7 {
				shortHash = shortHash[:7]
			}
			outPath := filepath.Join(prFolder, "PR_"+shortHash+".md")
			if err := os.WriteFile(outPath, []byte(finalPR), 0644); err != nil {
				prog.Stop("Error guardando PR", false)
				return fmt.Errorf("saving PR: %w", err)
			}
			prog.Stop(fmt.Sprintf("PR guardado: %s", outPath), true)

			// ── Step 8: Clipboard ──────────────────────────────────────────
			if cfg.Output.CopyToClipboard && !noClipboard {
				if err := clipboard.Copy(finalPR); err == nil {
					ui.Successf("PR copiado al portapapeles")
				}
			}

			fmt.Println()
			return nil
		},
	}

	cmd.Flags().IntVarP(&numCommits, "commits", "c", 1, "Number of commits to analyze")
	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Additional instructions (inline)")
	cmd.Flags().StringVarP(&notesFile, "notes-file", "f", "", "Read additional instructions from file")
	cmd.Flags().BoolVarP(&interactiveNotes, "interactive-notes", "i", false, "Enter multiline notes (end with 'END')")
	cmd.Flags().StringVarP(&tasks, "tasks", "t", "", `Task IDs to reference in the PR (comma-separated, e.g. "TK-123,TK-456")`)
	cmd.Flags().StringVar(&fromRef, "from", "", "Base branch/ref to compare from (e.g. develop)")
	cmd.Flags().StringVar(&toRef, "to", "HEAD", "Target branch/ref (default: HEAD)")
	cmd.Flags().BoolVar(&noClipboard, "no-clipboard", false, "Do not copy to clipboard")
	cmd.Flags().StringVarP(&provider, "provider", "p", "", "LLM provider override (ollama|openai|groq|openrouter|mock)")
	cmd.Flags().StringVarP(&model, "model", "m", "", "Model override")
	cmd.Flags().BoolVar(&debug, "debug", false, "Enable debug output")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Skip LLM call (useful for testing)")
	cmd.Flags().BoolVar(&dumpPrompt, "dump-prompt", false, "Print the prompt and exit")

	return cmd
}

func mustCwd() string {
	cwd, _ := os.Getwd()
	return cwd
}
