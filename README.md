# prgen — AI-powered PR Generator

> Generate detailed Pull Request descriptions from git commits using your preferred LLM.

**Providers:** Ollama (local) · Groq · OpenAI · OpenRouter

---

## Quick Start

```bash
# 1. Clone
git clone https://github.com/ksh/prgen
cd prgen

# 2. Install (adds prgen to your PATH automatically)
bash scripts/install.sh          # macOS / Linux
# PowerShell -ExecutionPolicy Bypass -File scripts\install.ps1  # Windows

# 3. Use (from any git repo)
prgen                            # generates PR from last commit
```

---

## Usage

```
prgen [generate] [flags]

Flags:
  -c, --commits int             Number of commits to analyze (default 1)
  -n, --notes string            Additional instructions (inline)
  -f, --notes-file string       Read additional instructions from file
  -i, --interactive-notes       Enter multiline notes (end with 'END')
  -p, --provider string         LLM provider (ollama|openai|groq|openrouter|mock)
  -m, --model string            Model override
      --no-clipboard            Do not copy output to clipboard
      --dry-run                 Skip LLM call (for testing)
      --dump-prompt             Print prompt and exit
      --debug                   Enable debug output

Commands:
  prgen generate                Main generate flow (default)
  prgen version                 Show version + build date
  prgen update                  Pull latest from git (safe, with confirmation)
  prgen config                  Show active configuration
```

---

## Configuration

Config is loaded in this priority order (highest wins):

1. `config.yaml` in current dir or binary dir
2. `~/.prgen/config.yaml` ← your personal overrides
3. Environment variables
4. CLI flags

### Key environment variables

| Variable | Description |
|---|---|
| `PRGEN_PROVIDER` | `ollama` \| `openai` \| `groq` \| `openrouter` |
| `PRGEN_MODEL` | Model name (e.g. `llama3.1`, `llama-3.1-70b-versatile`) |
| `PRGEN_API_KEY` | API key (also accepts `GROQ_API_KEY`, `OPENAI_API_KEY`) |
| `PRGEN_API_BASE_URL` | Custom base URL (auto-set for groq/openai/openrouter) |
| `PRGEN_OLLAMA_URL` | Ollama URL (default: `http://localhost:11434`) |

### Use Groq (recommended for speed)

```bash
export PRGEN_PROVIDER=groq
export GROQ_API_KEY=gsk_xxxx
export PRGEN_MODEL=llama-3.1-70b-versatile
prgen
```

### Use Ollama (local, no API key)

```bash
ollama pull llama3.1
prgen  # provider defaults to ollama
```

---

## Customizing the Prompt

The prompt system has two layers:

1. **`prompts/base.md`** — base prompt, versioned with the tool. Supports `{{.Branch}}`, `{{.Stats}}`, `{{.Logs}}`, `{{.ProjectType}}` placeholders.
2. **`~/.prgen/extra_prompt.md`** — your personal/team additions. Not tracked by git. Injected after the base prompt.

---

## Update

```bash
prgen update            # fetches, shows new commits, confirms before pulling
```

After pulling, reinstall:
```bash
bash scripts/install.sh
```

---

## Building manually

```bash
VERSION=$(cat VERSION)
go build \
  -ldflags "-s -w -X github.com/ksh/prgen/internal/version.Version=$VERSION" \
  -o prgen \
  ./cmd/prgen
```

---

## License

MIT — KSH
