package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alann-estrada-KSH/ai-pr-generator/internal/git"
)

func newUpdateCmd() *cobra.Command {
	var (
		remote string
		branch string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update prgen from the remote repository",
		Long: `Safely pulls the latest version from the remote Git repository.
Shows what will change and asks for confirmation before pulling.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("\n 🔄 Buscando actualizaciones para prgen...")

			// ── Safety check: working tree must be clean ───────────────────
			if !git.IsClean() {
				return fmt.Errorf("working tree has uncommitted changes — stash or commit before updating")
			}

			currentHash := git.HeadHash()
			if len(currentHash) > 8 {
				currentHash = currentHash[:8]
			}
			fmt.Printf(" 📌 Versión actual: %s\n", currentHash)

			// ── Fetch and check for new commits ────────────────────────────
			diff, err := git.FetchAndDiff(remote, branch)
			if err != nil {
				return fmt.Errorf("fetching updates: %w", err)
			}

			if strings.TrimSpace(diff) == "" {
				fmt.Println("\n ✅ Ya estás en la versión más reciente.")
				return nil
			}

			fmt.Printf("\n 📋 Commits nuevos en %s/%s:\n", remote, branch)
			for _, line := range strings.Split(diff, "\n") {
				fmt.Println("   ", line)
			}

			// ── Ask for confirmation ───────────────────────────────────────
			fmt.Print("\n ¿Deseas actualizar? (s/n): ")
			var answer string
			fmt.Scanln(&answer)
			if strings.ToLower(strings.TrimSpace(answer)) != "s" {
				fmt.Println(" ❌ Actualización cancelada.")
				return nil
			}

			// ── Pull ───────────────────────────────────────────────────────
			if err := git.Pull(remote, branch); err != nil {
				return fmt.Errorf("git pull: %w", err)
			}

			// ── Suggest rebuild ───────────────────────────────────────────
			binName := "prgen"
			if _, err := os.Stat("cmd/prgen/main.go"); err == nil {
				fmt.Printf(`
 ✅ Código actualizado.

 Para aplicar los cambios, recompila y reinstala:

   go build -ldflags "-X github.com/alann-estrada-KSH/ai-pr-generator/internal/version.Version=$(cat VERSION)" -o %s ./cmd/prgen

 O usa el script de instalación:
   scripts/install.sh   (macOS/Linux)
   scripts/install.ps1  (Windows)

`, binName)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&remote, "remote", "r", "origin", "Git remote name")
	cmd.Flags().StringVarP(&branch, "branch", "b", "main", "Remote branch to update from")

	return cmd
}
