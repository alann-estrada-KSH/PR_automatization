# prgen — Documentación Completa

**prgen** es un conjunto de herramientas de IA para el flujo de trabajo Git: genera Pull Requests detallados, revisa código, crea commits semánticos y sugiere nombres de rama — todo desde la terminal, con el LLM que elijas.

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

## Comandos

| Comando | Descripción |
|---|---|
| [`prgen generate`](#prgen-generate) | Genera descripción de PR desde commits recientes |
| [`prgen commit`](#prgen-commit) | Genera mensaje de commit CDE desde staged changes |
| [`prgen review`](#prgen-review) | Revisión de código por IA |
| [`prgen branch`](#prgen-branch) | Sugiere nombres de rama |
| [`prgen version`](#prgen-version) | Muestra versión instalada |
| [`prgen update`](#prgen-update) | Actualiza desde git |
| [`prgen config`](#prgen-config) | Muestra configuración activa |

---

## `prgen generate`

Analiza N commits recientes y genera una descripción completa del PR usando el LLM configurado. El resultado se guarda en disco y se copia al portapapeles.

### Uso básico

```bash
prgen                                              # PR del último commit
prgen generate --commits 3                         # últimos 3 commits

# Comparar ramas arbitrarias
prgen generate --from develop                      # todo lo que difiere de develop
prgen generate --from develop --to feature/oauth   # rango específico entre dos ramas

# Incluir referencias de tareas
prgen generate --tasks "TK-123,TK-456"
# → agrega sección "## 🗂️ Referencias de tareas" al PR

# Instrucciones adicionales al LLM
prgen generate --notes "este PR cierra el sprint 14, enfocarse en la capa de pagos"
prgen generate --notes-file ~/contexto-pr.md
prgen generate --interactive-notes    # entrada multilinea, termina con 'END'
```

### Archivo de salida

```
~/KSH/Projects/<nombre-repo> - PR/<DD-MM-YYYY>/PR_<hash>.md
```

### Flags completos

| Flag | Corto | Descripción | Default |
|---|---|---|---|
| `--commits` | `-c` | Commits a analizar | `1` |
| `--from` | | Rama/ref base para comparar | |
| `--to` | | Rama/ref destino | `HEAD` |
| `--tasks` | `-t` | Task IDs separados por coma | |
| `--notes` | `-n` | Instrucciones adicionales inline | |
| `--notes-file` | `-f` | Instrucciones desde archivo | |
| `--interactive-notes` | `-i` | Notas multilinea | |
| `--provider` | `-p` | Override del proveedor LLM | |
| `--model` | `-m` | Override del modelo | |
| `--no-clipboard` | | No copiar al portapapeles | |
| `--dry-run` | | Salta la llamada al LLM | |
| `--dump-prompt` | | Imprime el prompt y sale | |
| `--debug` | | Modo debug | |

---

## `prgen commit`

Analiza los cambios staged (`git diff --cached`) y genera un mensaje de commit siguiendo el **Método CDE** (Document as you code) + especificación [Conventional Commits](https://www.conventionalcommits.org/).

### ¿Qué es el Método CDE?

CDE significa que el mensaje de commit **es documentación técnica viva** del sistema. No es un log de cambios — es el registro de *por qué* el sistema evolucionó de una forma particular.

Un commit CDE tiene dos partes:

1. **Asunto** — *qué* en imperativo presente, máx. 72 chars
2. **Cuerpo** — *por qué* se tomó esa decisión técnica (2-4 líneas)

### Uso

```bash
git add .
prgen commit           # muestra el mensaje sugerido
prgen commit --apply   # genera y ejecuta git commit directamente
```

### Tipos de commit

| Tipo | Cuándo usarlo |
|---|---|
| `feat` | Nueva funcionalidad visible |
| `fix` | Corrección de bug con comportamiento incorrecto |
| `refactor` | Cambio de código sin nuevas features ni bug fix |
| `chore` | Mantenimiento: deps, config, scripts, archivos generados |
| `docs` | Solo documentación |
| `style` | Formato, espacios (sin cambio de lógica) |
| `test` | Añadir o corregir tests |
| `perf` | Mejora de rendimiento medible |
| `ci` | Configuración de CI/CD |
| `build` | Sistema de compilación |

### Ejemplos de commits CDE

**Nueva funcionalidad con contexto de decisión:**
```
feat(auth): implementar rotación de refresh tokens en TokenService

Los tokens de acceso estáticos representaban una ventana de exposición
ilimitada si eran comprometidos. Con la rotación, cada uso de un refresh
token invalida el anterior, limitando el daño a una sola sesión.
TokenService.refreshToken() ahora retorna un par {access, refresh} nuevo.
```

**Corrección de bug con causa raíz:**
```
fix(payments): corregir cálculo de IVA en órdenes con descuento

OrderCalculator.computeTax() aplicaba el IVA sobre el total bruto antes
de descontar promociones, resultando en cobros mayores al esperado.
Se cambió el orden de operaciones: descuento → base imponible → IVA.
```

**Refactor con justificación de arquitectura:**
```
refactor(api): extraer validación de permisos a PermissionMiddleware

La lógica de permisos estaba duplicada en 7 controladores distintos.
Centralizar en un middleware elimina la inconsistencia y facilita auditar
qué rutas tienen restricciones sin revisar cada controlador.
```

**Chore simple (sin cuerpo):**
```
chore(deps): actualizar dependencias de seguridad (febrero 2026)
```

**Breaking change con footer:**
```
feat(api): requerir Authorization header en todos los endpoints

Esta versión elimina el acceso anónimo a rutas públicas deprecadas
desde v1.2. Los clientes deben incluir Bearer token en cada request.

BREAKING CHANGE: rutas /api/public/* ya no aceptan requests sin auth.
```

### Banderas

| Flag | Descripción |
|---|---|
| `--apply` | Ejecuta `git commit` con el mensaje generado |
| `--provider` / `-p` | Override del proveedor LLM |
| `--model` / `-m` | Override del modelo |
| `--dry-run` | Imprime el prompt sin llamar al LLM |

---

## `prgen review`

Revisa el diff y produce un reporte técnico estructurado. Útil para auto-revisión antes de abrir el PR, o para ayudar al reviewer con contexto.

### Uso

```bash
prgen review                       # revisa el último commit
prgen review --commits 3           # últimos 3 commits
prgen review --from develop        # todo lo diferente de develop
prgen review --from develop --to feature/auth   # rango específico
```

### Estructura del reporte

```markdown
## 🐛 Posibles bugs y casos no manejados
## 🔒 Problemas de seguridad
## ⚠️ Manejo de errores
## 🔧 Sugerencias de refactor
## ✅ Resumen general
```

Cada hallazgo tiene un nivel de severidad: **🔴 Crítico | 🟡 Advertencia | 🟢 Sugerencia**

### Banderas

| Flag | Corto | Descripción |
|---|---|---|
| `--commits` | `-c` | Commits a revisar (default: 1) |
| `--from` | | Rama/ref base |
| `--to` | | Rama/ref destino (default: HEAD) |
| `--provider` | `-p` | Override del proveedor LLM |
| `--model` | `-m` | Override del modelo |
| `--dry-run` | | Imprime el prompt y sale |

---

## `prgen branch`

Dado lo que necesitas hacer, sugiere 5 nombres de rama con los prefijos correctos y en kebab-case.

### Uso

```bash
prgen branch "agregar login con google oauth"
prgen branch "corregir error 500 en endpoint de pagos"
prgen branch "refactorizar el módulo de permisos para soportar RBAC"
```

### Ejemplo de salida

```
┌─ Sugerencias de rama ──────────────────────────────────────
  feature/google-oauth-login           → git checkout -b feature/google-oauth-login
  feature/oauth-google-integration     → git checkout -b feature/oauth-google-integration
  feature/auth-google-login            → git checkout -b feature/auth-google-login
  feat/google-social-login             → git checkout -b feat/google-social-login
  feature/social-auth-google           → git checkout -b feature/social-auth-google
└────────────────────────────────────────────────────────────
```

Si la descripción menciona un ticket (TK-123, JIRA-456), se incluye en el nombre:
```bash
prgen branch "TK-456 corregir error en módulo de pagos"
# → fix/TK-456-error-modulo-pagos
```

---

## `prgen version`

```bash
prgen version
# prgen v0.2.0 (built 2026-02-19)
```

---

## `prgen config`

Muestra la configuración activa combinando `config.yaml` + variables de entorno + flags:

```bash
prgen config
```

Útil para verificar qué proveedor está activo, si la API key está configurada, etc.

---

## `prgen update`

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
4. Sugiere reinstalar después

```bash
# Después de actualizar:
bash scripts/install.sh                                           # macOS/Linux
PowerShell -ExecutionPolicy Bypass -File scripts\install.ps1     # Windows
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
api_key: ""               # usar variable de entorno
api_base_url: ""          # auto-configurado para proveedores conocidos

prompts:
  base:   prompts/base.md           # prompt para generate (versionado)
  commit: prompts/commit.md         # prompt CDE para commit
  review: prompts/review.md         # prompt para review
  extra:  ~/.prgen/extra_prompt.md  # instrucciones de tu equipo (opcional)

output:
  save_path: ~/KSH/Projects         # dónde se guardan los PRs
  copy_to_clipboard: true

diff:
  max_chars: 20000        # chars máximos del diff enviado al LLM
  ignore:                 # archivos excluidos del diff automáticamente
    - "package-lock.json"
    - "composer.lock"
    - "yarn.lock"
    - "go.sum"
    - "*.min.js"
    - "*.min.css"
    # agrega los tuyos:
    # - "vendor/**"
    # - "*.generated.*"

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
ollama pull llama3.1
prgen   # ollama es el default
```

### Groq (cloud, muy rápido — recomendado)

```bash
export PRGEN_PROVIDER=groq
export GROQ_API_KEY=gsk_xxxxxxxxxxxx
export PRGEN_MODEL=llama-3.1-70b-versatile

prgen
```

O en `~/.prgen/config.yaml`:
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

### OpenRouter

```bash
export PRGEN_PROVIDER=openrouter
export PRGEN_API_KEY=sk-or-xxxxxxxxxxxx
export PRGEN_MODEL=meta-llama/llama-3.1-70b-instruct
prgen
```

---

## Personalizar prompts

### Prompts disponibles

| Archivo | Comando | Placeholders |
|---|---|---|
| `prompts/base.md` | `generate` | `{{.ProjectType}}`, `{{.Branch}}`, `{{.Stats}}`, `{{.Logs}}`, `{{.Diff}}` |
| `prompts/commit.md` | `commit` | `{{.ProjectType}}`, `{{.Stats}}`, `{{.Diff}}` |
| `prompts/review.md` | `review` | `{{.ProjectType}}`, `{{.Branch}}`, `{{.Stats}}`, `{{.Diff}}` |

### Prompt extra del equipo (`~/.prgen/extra_prompt.md`)

Instrucciones adicionales que se inyectan al final del prompt base en `generate`. No se versiona con git.

```bash
cat > ~/.prgen/extra_prompt.md << 'EOF'
- Esta empresa usa Jira. Menciona el ticket si aparece en el commit o la rama.
- Los PRs van hacia la rama develop, nunca a main.
- Siempre mencionar si el cambio requiere ejecutar migraciones.
EOF
```

---

## Tipos de proyecto detectados

prgen detecta automáticamente el tipo de proyecto y adapta los checklists y el tono del prompt:

| Proyecto | Archivo detectado | Checklist |
|---|---|---|
| Laravel | `artisan` | Controladores, migraciones, tests PHPUnit |
| Python | `requirements.txt` / `pyproject.toml` | `.py`, dependencias, pytest |
| Go | `go.mod` | `.go`, `go.mod`, `_test.go` |
| Node.js | `package.json` | `.js/.ts`, `package.json`, tests |
| Dolibarr | `main.inc.php` | SQL, `.php`, CSS/JS |
| Genérico | (ninguno anterior) | Revisión manual |

---

## Flujo de generación de PR

```
1. git log + git diff (filtrado de archivos noise)
       ↓
2. Construir prompt (base.md + extra_prompt.md + notas + diff real)
       ↓
3. LLM genera el contenido (basándose en el diff, no en el commit message)
       ↓
4. Limpieza (cleaner): normalización de headers, emojis, bullets
       ↓
5. Inyección de secciones:
   - ## 🗂️ Referencias de tareas (si --tasks)
   - ## 📝 Instrucciones adicionales (si --notes)
   - ## 🛠️ Cambios realizados (checklist técnico)
   - ## ✅ Checklist antes de hacer merge
       ↓
6. Guardar en ~/KSH/Projects/<repo> - PR/<fecha>/PR_<hash>.md
7. Copiar al portapapeles (PowerShell Set-Clipboard en Windows)
```

---

## Resolución de problemas

### `ollama: connection refused`
Ollama no está corriendo.
```bash
ollama serve   # en otra terminal
```

### `Modelo no encontrado (404)`
El modelo no está descargado.
```bash
ollama pull llama3.1
```

### `API key inválida (401)`
```bash
echo $GROQ_API_KEY   # verificar que esté seteada
prgen config         # ver qué api_key tiene configurada
```

### Los emojis del PR salen como caracteres extraños al pegar
Estaba usando `clip.exe` que no maneja UTF-8. Desde v0.1.1 se usa `PowerShell Set-Clipboard`. Rebuilda: `PowerShell -ExecutionPolicy Bypass -File scripts\install.ps1`.

### El PR es genérico y no refleja los cambios reales
- Asegúrate de que el binario sea ≥ v0.1.1 (usa `prgen version`)
- Verifica con `prgen generate --dump-prompt` que el diff aparezca en el prompt
- Si el diff está vacío, puede que todos los archivos estén en `diff.ignore`

### `prgen commit` no encuentra cambios staged
```bash
git status         # ver qué archivos están modificados
git add <archivo>  # stagear los cambios
prgen commit
```

---

## Build manual

```bash
VERSION=$(cat VERSION)
go build \
  -ldflags "-s -w \
    -X github.com/ksh/prgen/internal/version.Version=$VERSION \
    -X github.com/ksh/prgen/internal/version.BuildDate=$(date +%Y-%m-%d)" \
  -o prgen \
  ./cmd/prgen
```
