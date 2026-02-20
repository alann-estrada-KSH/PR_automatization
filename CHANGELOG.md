# Changelog

All notable changes to **prgen** will be documented in this file.

Format based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) — versioning follows [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

## [0.2.0] - 2026-02-19

### Added
- `prgen commit` — genera mensaje de commit Conventional Commits desde staged changes
  - `--apply` ejecuta el `git commit` directamente
- `prgen review` — revisión de código por IA (bugs, seguridad, error handling, refactor)
  - soporta `--commits N`, `--from`, `--to`
- `prgen branch <descripción>` — sugiere 5 nombres de rama con convenciones (feature/, fix/, etc.)
- `prgen generate --from <base> --to <target>` — genera PR comparando ramas arbitrarias
- Filtrado inteligente del diff: excluye `package-lock.json`, `composer.lock`, `go.sum`, `*.min.*`
- `DiffConfig` en config: `diff.max_chars` (default 20000) y `diff.ignore` (lista configurable en YAML)
- Prompts editables: `prompts/commit.md` y `prompts/review.md`

### Changed
- Límite del diff aumentado de 6000 a 20000 chars (configurable)
- `prgen generate` ahora usa `FilteredDiff` — excluye archivos de ruido automáticamente

## [0.1.1] - 2026-02-19

### Added
- `git diff` real incluido en el prompt como fuente de verdad para el LLM
- Regla explícita: ignorar commits vagos y analizar el diff directamente
- Guía estructurada de 5 párrafos en la sección "Resumen del cambio"
- `git.Diff(n)` function en el paquete `internal/git`
- Campo `Diff` en `prompt.Context` con truncado automático a 6000 chars

### Fixed
- Clipboard en Windows ahora usa PowerShell `Set-Clipboard` (soluciona encoding UTF-8/emojis)
- Header regex en cleaner ahora detecta emojis entre `##` y el texto del título
- Git commands forzados a UTF-8 en Windows (`-c i18n.logOutputEncoding=UTF-8`)
- Import paths incorrectos (`alann-estrada-KSH` → `ksh/prgen`) en varios archivos

### Changed
- `base.md` reescrito: incluye diff, reglas más estrictas, menciona funciones concretas
- Terminal interactivo con spinner + barra de progreso + tiempo transcurrido
- Errores de API incluyen cuerpo de respuesta y sugerencia de solución

## [0.1.0] - 2026-02-19

### Added
- Initial rewrite from Python to Go
- Multi-provider LLM support: **Ollama** (local), **Groq**, **OpenAI-compatible** (OpenRouter, etc.)
- Provider selectable via `.env` / `config.yaml` / `--provider` flag
- Prompt system: immutable `prompts/base.md` + optional `~/.prgen/extra_prompt.md`
- Context providers: manual task input, file-based notes
- CLI commands: `generate`, `version`, `update`, `config`
- Flags: `--commits`, `--notes`, `--notes-file`, `--interactive-notes`, `--no-clipboard`, `--dry-run`, `--dump-prompt`, `--debug`
- Centralized config via `config.yaml` + environment variable overrides
- Safe `update` command (shows diff, asks confirmation before pulling)
- Cross-platform install scripts: `scripts/install.sh` (macOS/Linux) and `scripts/install.ps1` (Windows)
- Automatic PATH configuration on install for Windows, macOS, and Linux
- Technical checklists (Laravel, Python, Dolibarr, generic) and merge templates
- Output cleaning and header normalization for LLM responses
- PR saved to `~/KSH/Projects/<repo> - PR/<date>/PR_<hash>.md`
- Clipboard copy support (Windows, macOS, Linux)
