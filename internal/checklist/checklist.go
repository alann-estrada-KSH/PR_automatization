package checklist

import (
	"strings"

	"github.com/ksh/prgen/internal/detect"
)

// ── Technical checklists ──────────────────────────────────────────────────────

var templateLaravel = []string{
	"- [ ] Nuevo endpoint en el controlador `PermissionController` (o lógica backend)",
	"- [ ] Modificación de la base de datos (nueva migración)",
	"- [ ] Actualización de pruebas unitarias e integración",
}

var templatePython = []string{
	"- [ ] Cambios en lógica principal (.py)",
	"- [ ] Modificación de dependencias (requirements/pip)",
	"- [ ] Actualización de tests (pytest)",
}

var templateDolibarr = []string{
	"- [ ] Cambios en descriptores de módulo o SQL",
	"- [ ] Modificación de lógica PHP/Core",
	"- [ ] Cambios en interfaz (CSS/JS)",
}

var templateGo = []string{
	"- [ ] Cambios en lógica principal (.go)",
	"- [ ] Modificación de dependencias (go.mod / go.sum)",
	"- [ ] Actualización de tests (go test)",
}

var templateNode = []string{
	"- [ ] Cambios en lógica principal (.js/.ts)",
	"- [ ] Modificación de dependencias (package.json)",
	"- [ ] Actualización de tests",
}

// ── Merge checklists ──────────────────────────────────────────────────────────

var mergeTemplates = map[detect.ProjectType]string{
	detect.Laravel: `## ✅ Checklist antes de hacer merge
- [x] Código probado localmente
- [ ] Pruebas unitarias pasan (` + "`php artisan test`" + `)
- [ ] Pruebas de integración pasan
- [ ] Revisado por al menos 1 desarrollador`,

	detect.Python: `## ✅ Checklist antes de hacer merge
- [x] Código probado localmente
- [ ] Pruebas unitarias pasan (` + "`pytest`" + `)
- [ ] Linter verificado (` + "`flake8` / `black`" + `)
- [ ] Revisado por al menos 1 desarrollador`,

	detect.Dolibarr: `## ✅ Checklist antes de hacer merge
- [x] Código probado localmente
- [ ] Módulo activado y verificado en entorno de pruebas
- [ ] Scripts SQL ejecutados sin errores
- [ ] Revisado por al menos 1 desarrollador`,

	detect.Go: `## ✅ Checklist antes de hacer merge
- [x] Código probado localmente
- [ ] Tests pasan (` + "`go test ./...`" + `)
- [ ] Linter verificado (` + "`golangci-lint`" + `)
- [ ] Revisado por al menos 1 desarrollador`,

	detect.Node: `## ✅ Checklist antes de hacer merge
- [x] Código probado localmente
- [ ] Tests pasan (` + "`npm test`" + `)
- [ ] Linter verificado (` + "`eslint`" + `)
- [ ] Revisado por al menos 1 desarrollador`,
}

var mergeGeneric = `## ✅ Checklist antes de hacer merge
- [x] Código probado localmente
- [ ] Pruebas manuales completadas
- [ ] Documentación actualizada
- [ ] Revisado por al menos 1 desarrollador`

// Technical builds the technical checklist for the given project, with
// items auto-checked based on the diff stats.
func Technical(pt detect.ProjectType, stats string) string {
	lower := strings.ToLower(stats)
	var items []string

	switch pt {
	case detect.Laravel:
		items = cloneSlice(templateLaravel)
		if containsAny(lower, "controller", "route", "api.php", "web.php", "trait", "service", "request") {
			items[0] = checked(items[0])
		}
		if containsAny(lower, "migration", "schema", "model", "database", "pivot") {
			items[1] = checked(items[1])
		}
		if strings.Contains(lower, "test") || strings.Contains(lower, "phpunit") {
			items[2] = checked(items[2])
		}

	case detect.Python:
		items = cloneSlice(templatePython)
		if strings.Contains(lower, ".py") {
			items[0] = checked(items[0])
		}
		if strings.Contains(lower, "requirements") || strings.Contains(lower, ".toml") {
			items[1] = checked(items[1])
		}
		if strings.Contains(lower, "test") {
			items[2] = checked(items[2])
		}

	case detect.Dolibarr:
		items = cloneSlice(templateDolibarr)
		if strings.Contains(lower, "sql") || strings.Contains(lower, "descriptor") {
			items[0] = checked(items[0])
		}
		if strings.Contains(lower, ".php") {
			items[1] = checked(items[1])
		}
		if strings.Contains(lower, ".css") || strings.Contains(lower, ".js") {
			items[2] = checked(items[2])
		}

	case detect.Go:
		items = cloneSlice(templateGo)
		if strings.Contains(lower, ".go") {
			items[0] = checked(items[0])
		}
		if strings.Contains(lower, "go.mod") || strings.Contains(lower, "go.sum") {
			items[1] = checked(items[1])
		}
		if strings.Contains(lower, "_test.go") {
			items[2] = checked(items[2])
		}

	case detect.Node:
		items = cloneSlice(templateNode)
		if containsAny(lower, ".js", ".ts", ".jsx", ".tsx") {
			items[0] = checked(items[0])
		}
		if strings.Contains(lower, "package.json") {
			items[1] = checked(items[1])
		}
		if strings.Contains(lower, "test") || strings.Contains(lower, "spec") {
			items[2] = checked(items[2])
		}

	default:
		items = []string{"- [ ] Revisión manual de cambios genéricos"}
	}

	return strings.Join(items, "\n")
}

// Merge returns the merge checklist for the given project type.
func Merge(pt detect.ProjectType) string {
	if tpl, ok := mergeTemplates[pt]; ok {
		return tpl
	}
	return mergeGeneric
}

// ── helpers ───────────────────────────────────────────────────────────────────

func cloneSlice(s []string) []string {
	out := make([]string, len(s))
	copy(out, s)
	return out
}

func checked(s string) string {
	return strings.Replace(s, "[ ]", "[x]", 1)
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
