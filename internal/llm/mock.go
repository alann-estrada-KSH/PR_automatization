package llm

import "context"

// MockClient returns a static response, useful for tests and --dry-run.
type MockClient struct {
	Response string
}

func (m *MockClient) Name() string { return "mock" }

func (m *MockClient) Generate(_ context.Context, _ string) (string, error) {
	if m.Response != "" {
		return m.Response, nil
	}
	return `## 📌 Resumen del cambio

Este es un PR de prueba generado por el cliente mock de prgen.
Los cambios incluyen mejoras en la arquitectura y refactorización del código principal.
Se optimizaron las consultas a la base de datos para mejorar el rendimiento.
Se implementaron nuevos endpoints en el controlador principal.
Se actualizaron las dependencias del proyecto a sus versiones más recientes.

## 🔍 ¿Qué problema soluciona?

Resuelve el problema de rendimiento en la generación de reportes grandes.
Reduce el tiempo de respuesta de la API en un 40%.

## 🚀 ¿Cómo probarlo?

1. Clona el repositorio y cambia a esta rama.
2. Ejecuta las migraciones con ` + "`php artisan migrate`" + `.
3. Prueba el endpoint ` + "`GET /api/reports`" + `.

## ⚠️ Consideraciones adicionales

Ninguna.
`, nil
}
