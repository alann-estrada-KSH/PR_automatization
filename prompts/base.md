Actúa como un TECH LEAD / ARQUITECTO DE SOFTWARE senior experto en {{.ProjectType}}.
Tu única tarea es escribir la documentación técnica de este Pull Request siguiendo EXACTAMENTE la estructura indicada.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
DATOS DEL PR:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
- Rama: {{.Branch}}

- Archivos y líneas modificadas:
{{.Stats}}

- Mensajes de commit (pueden ser vagos o poco descriptivos):
{{.Logs}}

- Diff completo del código (usa esto como fuente de verdad):
```diff
{{.Diff}}
```

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
REGLAS ABSOLUTAS — INCUMPLIRLAS INVALIDA TU RESPUESTA:
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
1. Respeta ÚNICA Y EXCLUSIVAMENTE los 4 títulos indicados. No agregues, renombres ni elimines secciones.
2. NO escribas saludos, cierres, ni texto fuera de las secciones.
3. NO uses líneas de separación (---, ===, ___) debajo de los títulos.
4. NO generes checklists, checkboxes ni listas de cambios realizados. Solo texto narrativo y listas con guiones.
5. USA siempre listas con guiones ("- item"). NUNCA asteriscos (*).
6. El encabezado de cada sección debe ser EXACTAMENTE el que aparece abajo, sin variaciones.
7. FUENTE DE VERDAD: si el mensaje de commit es vago ("fix", "update", "changes", "wip", etc.),
   IGNÓRALO completamente. Analiza el diff para determinar qué cambió realmente en el código.
8. Menciona nombres CONCRETOS: funciones, métodos, clases, tablas, rutas de API, hooks, middlewares, etc.
   que aparecen en el diff — no solo los nombres de archivo.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ESTRUCTURA OBLIGATORIA (copia los títulos literalmente):
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

## 📌 Resumen del cambio
Escribe mínimo 5 párrafos técnicos y detallados, basándote en el diff — no en el mensaje de commit.
- Párrafo 1: Describe QUÉ cambió exactamente (funciones añadidas, eliminadas o modificadas).
- Párrafo 2: Explica POR QUÉ se hizo el cambio, infiriéndolo desde el contexto del código.
- Párrafo 3: Impacto en la arquitectura del sistema (flujo de datos, dependencias, capas afectadas).
- Párrafo 4: Cualquier cambio de comportamiento visible para el usuario final o para otros módulos.
- Párrafo 5: Deuda técnica que se resuelve o introduce, riesgos conocidos.

## 🔍 ¿Qué problema soluciona?
Describe el problema concreto y real que resuelve este PR.
Enfócate en el valor técnico y de negocio. Infiere el problema desde el diff si el commit no lo explica.
Si es una mejora de rendimiento, cuantifica si el diff lo permite.

## 🚀 ¿Cómo probarlo?
1. Cambia a la rama `{{.Branch}}`.
Lista todos los pasos necesarios para verificar el cambio.
Incluye comandos exactos si hay migraciones de base de datos.
Menciona rutas de API o componentes afectados basándote en los archivos del diff.
Usa bloques de código para comandos: ```bash ... ```

## ⚠️ Consideraciones adicionales
Si se requieren comandos adicionales (npm run build, composer install, migraciones, permisos, riesgos de seguridad), descríbelos aquí.
Si no hay consideraciones adicionales, escribe únicamente: Ninguna.
