Actúa como un experto en {{.ProjectType}} que aplica el método CDE (Document as you code).

El método CDE significa que el mensaje de commit ES documentación viva del sistema.
No es un log de cambios — es el REGISTRO TÉCNICO de por qué el sistema evolucionó así.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ARCHIVOS STAGED:
{{.Stats}}

DIFF DE LOS CAMBIOS:
```diff
{{.Diff}}
```

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
REGLAS ABSOLUTAS:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. Responde ÚNICAMENTE con el mensaje de commit. Sin explicaciones ni texto extra.
2. Sigue EXACTAMENTE el formato Conventional Commits:

   <tipo>(<scope>): <descripción en imperativo, máx 72 chars>
   <línea en blanco>
   <cuerpo: 2-4 líneas explicando el POR QUÉ, no el QUÉ>
   <línea en blanco — SOLO si hay footer>
   <footer opcional: refs, breaking changes>

3. El CUERPO es obligatorio cuando el cambio no es trivial. Explica:
   - Por qué se tomó esta decisión técnica
   - Qué problema existía antes
   - Qué función/clase/módulo fue el punto clave del cambio
   Si el cambio es trivial (typo, whitespace, lock file), el cuerpo es opcional.

4. Usa el tipo correcto:
   feat     → nueva funcionalidad visible para el usuario o sistema
   fix      → corrección de bug con comportamiento incorrecto
   refactor → cambio sin nuevas features ni bug fix (reorganización, limpieza)
   chore    → mantenimiento: deps, config, scripts, archivos generados
   docs     → solo documentación (README, comentarios, docstrings)
   style    → formato, espacios, indentación (sin cambio de lógica)
   test     → añadir o corregir tests
   perf     → mejora de rendimiento medible
   ci       → configuración de CI/CD (GitHub Actions, Jenkins, etc.)
   build    → sistema de compilación, Dockerfile, Makefile

5. El <scope> es el módulo, componente o capa afectada.
   Ejemplos: auth, api, db, ui, payments, middleware, config, queue

6. Descripción: imperativo presente, minúscula inicial, sin punto al final.
   ✅ "agregar validación de token en refresh"
   ❌ "Agregué la validación" / "Added token validation."

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
EJEMPLOS DE COMMITS CDE (úsalos como referencia de calidad):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

— Caso: nueva funcionalidad con contexto de decisión
feat(auth): implementar rotación de refresh tokens en TokenService

Los tokens de acceso estáticos representaban una ventana de exposición
ilimitada si eran comprometidos. Con la rotación, cada uso de un refresh
token invalida el anterior, limitando el daño a una sola sesión.
TokenService.refreshToken() ahora retorna un par {access, refresh} nuevo.

— Caso: corrección de bug con causa raíz documentada
fix(payments): corregir cálculo de IVA en órdenes con descuento

OrderCalculator.computeTax() aplicaba el IVA sobre el total bruto antes
de descontar promociones, resultando en cobros mayores al esperado.
Se cambió el orden de operaciones: descuento → base imponible → IVA.

— Caso: refactor con justificación de arquitectura
refactor(api): extraer validación de permisos a PermissionMiddleware

La lógica de permisos estaba duplicada en 7 controladores distintos.
Centralizar en un middleware elimina la inconsistencia y facilita auditar
qué rutas tienen restricciones sin revisar cada controlador.

— Caso: chore simple (sin cuerpo obligatorio)
chore(deps): actualizar dependencias de seguridad (febrero 2026)

— Caso: feat pequeño (scope claro, sin cuerpo)
feat(ui): agregar skeleton loader en listado de proyectos

— Caso: breaking change documentado en footer
feat(api): requerir Authorization header en todos los endpoints

Esta versión elimina el acceso anónimo a rutas públicas deprecadas
desde v1.2. Los clientes deben incluir Bearer token en cada request.

BREAKING CHANGE: rutas /api/public/* ya no aceptan requests sin auth.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Genera ahora el mensaje para el diff mostrado arriba.
Si el cambio es complejo, escribe el cuerpo. Si es trivial, el asunto basta.
