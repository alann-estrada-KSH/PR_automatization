package cleaner

import (
	"regexp"
	"strings"
)

// ── Garbage & formatting cleanup ──────────────────────────────────────────────

var (
	// Only match pure separator lines (exactly 3+ dashes/equals, nothing else)
	reGarbage        = regexp.MustCompile(`(?m)^[-=]{3,}$`)
	reBulletFix      = regexp.MustCompile(`(?m)^\s*\*\s+`)
	reExcessNewlines = regexp.MustCompile(`\n{3,}`)

	reChangesSection   = regexp.MustCompile(`(?is)##[^\n]*Cambios realizados[^\n]*\n.*`)
	reChecklistSection = regexp.MustCompile(`(?is)##[^\n]*Checklist[^\n]*\n.*`)
)

// headerRule matches a section header line and normalises it.
// The pattern is intentionally loose between the ## and the keyword text
// so it works whether or not the LLM puts an emoji between them.
type headerRule struct {
	// re matches the FULL line (after TrimSpace)
	re  *regexp.Regexp
	// canonical is the exact replacement line we emit
	rep string
}

var headerRules = []headerRule{
	{
		// Matches: ## Resumen del cambio  /  ## 📌 Resumen ...  /  **Resumen del cambio**
		re:  regexp.MustCompile(`(?i)^#{1,3}[^a-zA-ZáéíóúÁÉÍÓÚñÑ]*resumen del cambio`),
		rep: "## 📌 Resumen del cambio",
	},
	{
		re:  regexp.MustCompile(`(?i)^#{1,3}[^a-zA-ZáéíóúÁÉÍÓÚñÑ]*qu[eé] problema soluciona`),
		rep: "## 🔍 ¿Qué problema soluciona?",
	},
	{
		re:  regexp.MustCompile(`(?i)^#{1,3}[^a-zA-ZáéíóúÁÉÍÓÚñÑ]*c[oó]mo probarlo`),
		rep: "## 🚀 ¿Cómo probarlo?",
	},
	{
		re:  regexp.MustCompile(`(?i)^#{1,3}[^a-zA-ZáéíóúÁÉÍÓÚñÑ]*consideraciones adicionales`),
		rep: "## ⚠️ Consideraciones adicionales",
	},
}

// Clean removes decorative garbage lines, fixes bullet style, and collapses
// excess blank lines in raw LLM output.
func Clean(text string) string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))

	for _, line := range lines {
		stripped := strings.TrimSpace(line)
		if reGarbage.MatchString(stripped) {
			continue
		}
		line = reBulletFix.ReplaceAllString(line, "- ")
		out = append(out, line)
	}

	text = strings.Join(out, "\n")
	text = reExcessNewlines.ReplaceAllString(text, "\n\n")
	return text
}

// RepairHeaders enforces the canonical section headers regardless of what the
// LLM produced (with or without emojis, asterisks, or varying capitalisation).
func RepairHeaders(text string) string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))

	for _, line := range lines {
		stripped := strings.TrimSpace(line)
		replaced := false
		for _, rule := range headerRules {
			if rule.re.MatchString(stripped) {
				out = append(out, "\n"+rule.rep)
				replaced = true
				break
			}
		}
		if !replaced {
			out = append(out, line)
		}
	}

	return strings.TrimSpace(strings.Join(out, "\n"))
}

// RemoveHallucinations strips sections the LLM generates despite being told
// not to (its own checklist / "cambios realizados" sections).
func RemoveHallucinations(text string) string {
	text = reChangesSection.ReplaceAllString(text, "")
	text = reChecklistSection.ReplaceAllString(text, "")
	return text
}

// Process runs the full cleaning pipeline: Clean → RepairHeaders → RemoveHallucinations.
func Process(raw string) string {
	text := Clean(raw)
	text = RepairHeaders(text)
	text = RemoveHallucinations(text)
	return strings.TrimSpace(text)
}
