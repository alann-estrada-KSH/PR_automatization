# prgen — AI-powered PR Generator

> Genera Pull Requests, revisa código y crea commits semánticos desde tu terminal usando IA.

**Providers:** Ollama (local) · Groq · OpenAI · OpenRouter

---

## Quick Start

```bash
# 1. Clone
git clone https://github.com/alann-estrada-KSH/ai-pr-generator/
cd ai-pr-generator

# 2. Instalar (agrega prgen a tu PATH automáticamente)
bash scripts/install.sh                                             # macOS / Linux
PowerShell -ExecutionPolicy Bypass -File scripts\install.ps1       # Windows

# 3. Usar (desde cualquier repo git)
prgen                            # genera PR del último commit
```

---

## Comandos

| Comando | Descripción |
|---|---|
| `prgen generate` | Genera descripción de PR desde commits recientes |
| `prgen commit` | Genera mensaje de commit Conventional Commits desde staged changes |
| `prgen review` | Revisión de código por IA (bugs, seguridad, refactor) |
| `prgen branch <desc>` | Sugiere nombres de rama desde una descripción |
| `prgen version` | Muestra versión instalada |
| `prgen update` | Actualiza desde git (con confirmación) |
| `prgen config` | Muestra configuración activa |

---

## `prgen generate`

```bash
prgen                                            # PR del último commit
prgen generate --commits 3                       # últimos 3 commits
prgen generate --from develop                    # todo lo que difiere de develop
prgen generate --from develop --to feature/auth  # rango específico
prgen generate --tasks "TK-123,TK-456"           # incluye task IDs en el PR
prgen generate --notes "cierra el sprint 14"     # instrucciones adicionales
prgen generate --provider groq --model llama-3.1-70b-versatile
```

**Flags:**

| Flag | Corto | Descripción |
|---|---|---|
| `--commits` | `-c` | Commits a analizar (default: 1) |
| `--from` | | Rama/ref base para comparar |
| `--to` | | Rama/ref destino (default: HEAD) |
| `--tasks` | `-t` | Task IDs separados por coma |
| `--notes` | `-n` | Instrucciones adicionales inline |
| `--notes-file` | `-f` | Instrucciones desde archivo |
| `--interactive-notes` | `-i` | Notas multilinea (termina con `END`) |
| `--provider` | `-p` | Override del proveedor LLM |
| `--model` | `-m` | Override del modelo |
| `--no-clipboard` | | No copiar al portapapeles |
| `--dry-run` | | Salta la llamada al LLM |
| `--dump-prompt` | | Imprime el prompt y sale |
| `--debug` | | Modo debug |

---

## `prgen commit`

Analiza los cambios staged (`git diff --cached`) y genera un mensaje de commit siguiendo el **Método CDE** (Document as you code) + Conventional Commits.

```bash
git add .
prgen commit           # muestra la sugerencia
prgen commit --apply   # genera y ejecuta el git commit directamente
```

**Ejemplo de salida:**
```
feat(auth): implementar rotación de refresh tokens en TokenService

Los tokens de acceso estáticos representaban una ventana de exposición
ilimitada si eran comprometidos. Con la rotación, cada uso de un refresh
token invalida el anterior, limitando el daño a una sola sesión.
```

---

## `prgen review`

Revisión de código por IA. Detecta bugs, problemas de seguridad, falta de error handling y sugiere refactors.

```bash
prgen review                      # revisa el último commit
prgen review --commits 3          # últimos 3 commits
prgen review --from develop       # todo lo diferente de develop
```

**Reporte incluye:**
- 🔴 Crítico | 🟡 Advertencia | 🟢 Sugerencia
- Bugs y casos no manejados
- Seguridad (SQL injection, secrets, auth)
- Error handling faltante
- Sugerencias de refactor

---

## `prgen branch`

Sugiere nombres de rama desde una descripción en lenguaje natural.

```bash
prgen branch "corregir error 500 en endpoint de pagos"
# → fix/error-500-endpoint-pagos       → git checkout -b fix/...
# → fix/pagos-endpoint-500             → git checkout -b fix/...
# → hotfix/pagos-500-endpoint          → git checkout -b hotfix/...
```

---

## Configuración

**Config is loaded in this priority order (highest wins):**

1. `config.yaml` in current dir or binary dir
2. `~/.prgen/config.yaml` ← your personal overrides
3. `.env` file (current dir, binary dir, or ~/.prgen/)
4. Environment variables
5. CLI flags

```yaml
provider: ollama          # ollama | openai | groq | openrouter
model: llama3.1
ollama_url: http://localhost:11434
api_key: ""

prompts:
  base:   prompts/base.md
  commit: prompts/commit.md
  review: prompts/review.md
  extra:  ~/.prgen/extra_prompt.md

output:
  save_path: ~/KSH/Projects
  copy_to_clipboard: true

diff:
  max_chars: 20000          # límite de chars del diff enviado al LLM
  ignore:                   # archivos excluidos del diff
    - "package-lock.json"
    - "composer.lock"
    - "go.sum"
    - "*.min.js"
    - "*.min.css"
```

**Variables de entorno:**

| Variable | Descripción |
|---|---|
| `PRGEN_PROVIDER` | `ollama` \| `openai` \| `groq` \| `openrouter` |
| `PRGEN_MODEL` | Nombre del modelo |
| `PRGEN_API_KEY` | API key (también acepta `GROQ_API_KEY`, `OPENAI_API_KEY`) |
| `PRGEN_OLLAMA_URL` | URL de Ollama (default: `http://localhost:11434`) |

### Usar Groq (recomendado para velocidad)

```bash
export PRGEN_PROVIDER=groq
export GROQ_API_KEY=gsk_xxxx
export PRGEN_MODEL=llama-3.1-70b-versatile
prgen
```

---

## Personalizar prompts

| Archivo | Propósito |
|---|---|
| `prompts/base.md` | Estructura del PR — usa `{{.Branch}}`, `{{.Stats}}`, `{{.Logs}}`, `{{.Diff}}` |
| `prompts/commit.md` | Guía para mensajes de commit CDE |
| `prompts/review.md` | Estructura del reporte de revisión |
| `~/.prgen/extra_prompt.md` | Instrucciones de tu equipo (no versionado) |

---

## Build manual

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
