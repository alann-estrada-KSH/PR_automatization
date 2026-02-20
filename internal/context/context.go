// Package context defines the ContextProvider interface and built-in
// implementations for gathering additional context to inject into the PR prompt.
package context

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Provider is the interface all context sources must implement.
type Provider interface {
	// Name returns a human-readable name for this provider.
	Name() string
	// GetContext returns the context string to inject into the prompt.
	GetContext() (string, error)
}

// ── ManualProvider ────────────────────────────────────────────────────────────

// ManualProvider gathers tasks interactively from the terminal.
type ManualProvider struct{}

func (m *ManualProvider) Name() string { return "manual" }

func (m *ManualProvider) GetContext() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\n ¿Cuántas tareas son? (Enter para 0): ")
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	if line == "" {
		return "", nil
	}

	n, err := strconv.Atoi(line)
	if err != nil || n <= 0 {
		return "", nil
	}

	tasks := make([]string, 0, n)
	for i := 0; i < n; i++ {
		fmt.Printf(" Tarea %d: ", i+1)
		task, _ := reader.ReadString('\n')
		tasks = append(tasks, "- "+strings.TrimSpace(task))
	}

	return strings.Join(tasks, "\n"), nil
}

// ── FileProvider ──────────────────────────────────────────────────────────────

// FileProvider loads context from a file on disk.
type FileProvider struct {
	Path string
}

func (f *FileProvider) Name() string { return "file:" + f.Path }

func (f *FileProvider) GetContext() (string, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return "", fmt.Errorf("context file %q: %w", f.Path, err)
	}
	return strings.TrimSpace(string(data)), nil
}

// ── InlineProvider ────────────────────────────────────────────────────────────

// InlineProvider holds a pre-set string (from --notes flag).
type InlineProvider struct {
	Text string
}

func (i *InlineProvider) Name() string { return "inline" }

func (i *InlineProvider) GetContext() (string, error) {
	return strings.TrimSpace(i.Text), nil
}

// ── MultilineProvider ─────────────────────────────────────────────────────────

// MultilineProvider reads from stdin until the user types "END".
type MultilineProvider struct{}

func (m *MultilineProvider) Name() string { return "multiline" }

func (m *MultilineProvider) GetContext() (string, error) {
	fmt.Println("Escribe instrucciones adicionales (termina con una línea que contenga solo 'END'):")
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "END" {
			break
		}
		lines = append(lines, line)
	}
	return strings.TrimSpace(strings.Join(lines, "\n")), nil
}
