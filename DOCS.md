# prgen — Documentación

**prgen** es un generador de Pull Requests impulsado por IA. Analiza los commits de git recientes y genera documentación técnica completa usando el proveedor LLM que elijas.

---

## Instalación

### Windows
```powershell
PowerShell -ExecutionPolicy Bypass -File scripts\install.ps1
```

### macOS / Linux
```bash
bash scripts/install.sh
```

El instalador:
- Compila el binario con la versión embebida
- Lo agrega al `PATH` automáticamente (permanente)
- Copia `config.yaml` a `~/.prgen/config.yaml` como tu config personal

---

## Uso básico

```bash
# Desde cualquier repositorio git:
prgen                        # PR del último commit
prgen generate --commits 3   # PR usando los últimos 3 commits
```

### Salida

El PR se guarda en:
```
~/KSH/Projects/<nombre-repo> - PR/<DD-MM-YYYY>/PR_<hash>.md
```

También se copia al portapapeles automáticamente.

---

## Comandos

### `prgen generate` (alias: `prgen`, `prgen gen`)

Genera el PR. Es el comando por defecto.

| Flag | Corto | Descripción |
|------|-------|-------------|
| `--commits` | `-c` | Número de commits a analizar (default: 1) |
| `--notes` | `-n` | Instrucciones adicionales en línea |
| `--notes-file` | `-f` | Leer instrucciones desde un archivo |
| `--interactive-notes` | `-i` | Entrada multilinea (termina con `END`) |
| `--provider` | `-p` | Override del proveedor LLM |
| `--model` | `-m` | Override del modelo |
| `--no-clipboard` | | No copiar al portapapeles |
| `--dry-run` | | Salta la llamada al LLM (para probar el flujo) |
| `--dump-prompt` | | Imprime el prompt y sale (para depuración) |
| `--debug` | | Muestra información de depuración |
| `--tasks` | `-t` | Task IDs a referenciar en el PR (separados por coma, ej: `"TK-123,TK-456"`) |

**Ejemplos:**
```bash
# PR de los últimos 5 commits
prgen generate --commits 5

# Agregar notas inline
prgen generate --notes "Esto cierra el ticket #123"

# Leer notas de un archivo
prgen generate --notes-file ~/notas-pr.md

# Notas multilinea interactivas
prgen generate --interactive-notes

# Usar Groq en lugar de Ollama
prgen generate --provider groq --model llama-3.1-70b-versatile

# Ver el prompt sin llamar a la IA
prgen generate --dump-prompt

# Probar el flujo sin gastar tokens
prgen generate --dry-run
```

---

### `prgen version`

Muestra la versión instalada y la fecha de compilación.

```
prgen v0.1.0 (built 2026-02-19)
```

---

### `prgen config`

Muestra la configuración activa (valores combinados de `config.yaml` + variables de entorno + flags).

```bash
prgen config
```

Útil para verificar qué proveedor está usando, si la API key está configurada, etc.

---

### `prgen update`

Actualiza prgen desde el repositorio git de forma segura.

```bash
prgen update                       # rama main por defecto
prgen update --branch develop      # usar otra rama
prgen update --remote upstream     # usar otro remote
```

El comando:
1. Verifica que el directorio esté limpio (sin cambios sin commitear)
2. Muestra los commits nuevos disponibles
3. Pide confirmación antes de hacer pull
4. Sugiere reinstalar después de actualizar

Tras actualizar:
```bash
bash scripts/install.sh       # macOS/Linux
# o
PowerShell -ExecutionPolicy Bypass -File scripts\install.ps1   # Windows
```

---

## Configuración

### Prioridad (de menor a mayor)

```
Defaults → config.yaml → ~/.prgen/config.yaml → Variables de entorno → Flags CLI
```

### Archivo `config.yaml`

```yaml
provider: ollama          # ollama | openai | groq | openrouter
model: llama3.1
ollama_url: http://localhost:11434
api_key: ""               # usa la variable de entorno en su lugar
api_base_url: ""          # auto-configurado para proveedores conocidos

prompts:
  base: prompts/base.md           # prompt base (versionado con la herramienta)
  extra: ~/.prgen/extra_prompt.md # instrucciones de tu equipo (opcional)

output:
  save_path: ~/KSH/Projects       # dónde se guardan los PRs
  copy_to_clipboard: true

debug: false
```

Tu config personal va en `~/.prgen/config.yaml` — no se versiona con git.

### Variables de entorno

| Variable | Equivalente en config |
|---|---|
| `PRGEN_PROVIDER` | `provider` |
| `PRGEN_MODEL` | `model` |
| `PRGEN_API_KEY` | `api_key` |
| `PRGEN_API_BASE_URL` | `api_base_url` |
| `PRGEN_OLLAMA_URL` | `ollama_url` |
| `GROQ_API_KEY` | `api_key` (cuando provider=groq) |
| `OPENAI_API_KEY` | `api_key` (cuando provider=openai) |

---

## Proveedores LLM

### Ollama (local, sin costo)

```bash
# 1. Instalar y descargar el modelo
ollama pull llama3.1

# 2. Usar prgen (ollama es el default)
prgen
```

El modelo corre localmente. No requiere API key ni internet.

### Groq (cloud, muy rápido)

```bash
export PRGEN_PROVIDER=groq
export GROQ_API_KEY=gsk_xxxxxxxxxxxx
export PRGEN_MODEL=llama-3.1-70b-versatile   # o llama-3.3-70b-versatile

prgen
```

O en tu `~/.prgen/config.yaml`:
```yaml
provider: groq
model: llama-3.1-70b-versatile
api_key: gsk_xxxxxxxxxxxx
```

### OpenAI

```bash
export PRGEN_PROVIDER=openai
export OPENAI_API_KEY=sk-xxxxxxxxxxxx
export PRGEN_MODEL=gpt-4o-mini

prgen
```

### OpenRouter (acceso a múltiples modelos)

```bash
export PRGEN_PROVIDER=openrouter
export PRGEN_API_KEY=sk-or-xxxxxxxxxxxx
export PRGEN_MODEL=meta-llama/llama-3.1-70b-instruct

prgen
```

---

## Personalizar el prompt

### Prompt base (`prompts/base.md`)

Define la estructura del PR. Usa los siguientes placeholders:

| Placeholder | Contenido |
|---|---|
| `{{.ProjectType}}` | Tipo de proyecto detectado |
| `{{.Branch}}` | Nombre de la rama actual |
| `{{.Stats}}` | Output de `git diff --stat` |
| `{{.Logs}}` | Output de `git log --pretty` |

Puedes editarlo directamente, pero se recomienda no cambiar los títulos de las secciones para no romper el parser.

### Prompt extra del equipo (`~/.prgen/extra_prompt.md`)

Para instrucciones adicionales de tu equipo sin modificar el prompt base:

```bash
# Crea el archivo (no se versiona)
cat > ~/.prgen/extra_prompt.md << 'EOF'
- Esta empresa usa Jira. Menciona el ticket si aparece en el commit.
- Los PRs van a la rama develop, no a main.
- El equipo usa convención de commits Conventional Commits.
EOF
```

---

## Tipos de proyecto detectados

prgen detecta automáticamente el tipo de proyecto y adapta los checklists:

| Proyecto | Archivo detectado | Checklist |
|---|---|---|
| Laravel | `artisan` | Controladores, migraciones, tests PHPUnit |
| Python | `requirements.txt` / `pyproject.toml` | `.py`, dependencias, pytest |
| Go | `go.mod` | `.go`, `go.mod`, `_test.go` |
| Node.js | `package.json` | `.js/.ts`, `package.json`, tests |
| Dolibarr | `main.inc.php` | SQL, `.php`, CSS/JS |
| Genérico | (ninguno anterior) | Revisión manual |

---

## Flujo de un PR generado

```
1. git log + git diff --stat
       ↓
2. Construir prompt (base.md + extra_prompt.md + notas)
       ↓
3. LLM genera el contenido
       ↓
4. Limpieza (cleaner): normalización de headers, emojis, bullets
       ↓
5. Inyección de secciones:
   - ## 📝 Instrucciones adicionales (si hay --notes)
   - ## 🛠️ Cambios realizados (checklist técnico auto-marcado)
   - ## ✅ Checklist antes de hacer merge (por tipo de proyecto)
       ↓
6. Guardar en ~/KSH/Projects/<repo> - PR/<fecha>/PR_<hash>.md
7. Copiar al portapapeles
```

---

## Resolución de problemas

### "ollama: connection refused"
Ollama no está corriendo. Ejecuta `ollama serve` en otra terminal.

### "Modelo no encontrado (404)"
El modelo no está descargado. Ejecuta `ollama pull <modelo>`.

### "API key inválida (401)"
Verifica tu variable de entorno: `echo $GROQ_API_KEY` o `prgen config`.

### Los headers del PR salen mal
Usa `--dump-prompt` para ver el prompt enviado al LLM, y `--debug` para ver la respuesta raw.

### Los caracteres especiales salen garbled
El problema suele ser el encoding del terminal o del editor. Los archivos se generan en UTF-8 estándar.
