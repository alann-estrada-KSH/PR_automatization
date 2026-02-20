# Changelog

All notable changes to **prgen** will be documented in this file.

Format based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) — versioning follows [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

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
