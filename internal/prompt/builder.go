package prompt

import (
	"fmt"
	"os"
	"strings"

	"github.com/ksh/prgen/internal/detect"
)

// Context holds all the data assembled before calling the LLM.
type Context struct {
	ProjectType       detect.ProjectType
	Branch            string
	Logs              string
	Stats             string
	ExtraInstructions string // from --notes, --notes-file, or extra_prompt.md
}

// Builder builds the final prompt from base.md + optional extra file + context.
type Builder struct {
	BasePath  string // path to prompts/base.md
	ExtraPath string // path to ~/.prgen/extra_prompt.md (optional)
}

// Build constructs the full prompt string.
func (b *Builder) Build(ctx Context) (string, error) {
	base, err := b.loadBase(ctx)
	if err != nil {
		return "", err
	}

	extra := b.loadExtra()

	all := base
	if extra != "" {
		all += "\n\n### INSTRUCCIONES ADICIONALES DEL EQUIPO\n" + extra
	}
	if ctx.ExtraInstructions != "" {
		all += "\n\n### INSTRUCCIONES ADICIONALES\n" + ctx.ExtraInstructions
	}

	return all, nil
}

func (b *Builder) loadBase(ctx Context) (string, error) {
	data, err := os.ReadFile(b.BasePath)
	if err != nil {
		// Fall back to embedded default prompt
		return defaultPrompt(ctx), nil
	}

	// Substitute placeholders in base.md
	tmpl := string(data)
	tmpl = strings.ReplaceAll(tmpl, "{{.ProjectType}}", ctx.ProjectType.String())
	tmpl = strings.ReplaceAll(tmpl, "{{.Branch}}", ctx.Branch)
	tmpl = strings.ReplaceAll(tmpl, "{{.Logs}}", ctx.Logs)
	tmpl = strings.ReplaceAll(tmpl, "{{.Stats}}", ctx.Stats)
	return tmpl, nil
}

func (b *Builder) loadExtra() string {
	if b.ExtraPath == "" {
		return ""
	}
	data, err := os.ReadFile(b.ExtraPath)
	if err != nil {
		return "" // Extra prompt is always optional
	}
	return strings.TrimSpace(string(data))
}

// defaultPrompt is used when prompts/base.md is not found.
func defaultPrompt(ctx Context) string {
	return fmt.Sprintf(`Actúa como un TECH LEAD / ARQUITECTO DE SOFTWARE experto en %s.
Tu tarea es escribir la documentación técnica de este PR.

DATOS:
- Rama: %s
- Archivos Modificados:
%s
- Mensajes de Commit:
%s

INSTRUCCIONES DE FORMATO (ESTRICTO):
1. NO escribas saludos ni introducciones.
2. NO uses subrayados ni líneas de separación (ej: "---") debajo de los títulos.
3. NO generes checkboxes, ni listas de cambios, ni checklists. Solo texto narrativo.
4. Usa listas Markdown estándar con guiones ("- Item").

ESTRUCTURA Y CONTENIDO REQUERIDO:

## Resumen del cambio
(Escribe al menos 5 párrafos detallados. Menciona nombres de archivos técnicos importantes.
Explica la arquitectura y relaciones de datos si aplica).

## ¿Qué problema soluciona?
(Enfócate en el valor técnico y de negocio).

## ¿Cómo probarlo?
1. Cambia a la rama %s.
(Lista los pasos numerados. Si incluyes código usa bloques markdown. Si hay migraciones pon el comando exacto).

## Consideraciones adicionales
(Menciona comandos extra si son necesarios: npm run build, composer install, actualizaciones de BD, permisos o riesgos de seguridad. Si no hay, pon "Ninguna").`,
		ctx.ProjectType, ctx.Branch, ctx.Stats, ctx.Logs, ctx.Branch)
}
