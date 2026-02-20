package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ksh/prgen/internal/config"
)

// ollamaClient calls the local Ollama REST API.
type ollamaClient struct {
	baseURL string
	model   string
	http    *http.Client
}

func newOllamaClient(cfg *config.Config) *ollamaClient {
	url := cfg.OllamaURL
	if url == "" {
		url = "http://localhost:11434"
	}
	return &ollamaClient{
		baseURL: url,
		model:   cfg.Model,
		http:    &http.Client{Timeout: 10 * time.Minute},
	}
}

func (c *ollamaClient) Name() string { return "ollama/" + c.model }

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaChunk struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (c *ollamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	body, _ := json.Marshal(ollamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: true,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		hint := ""
		if resp.StatusCode == 404 {
			hint = " → Modelo no encontrado. ¿Ejecutaste 'ollama pull " + c.model + "'?"
		} else if resp.StatusCode == 0 {
			hint = " → Ollama no está corriendo. Ejecuta 'ollama serve' primero"
		}
		return "", fmt.Errorf("Ollama HTTP %d%s\nRespuesta: %s", resp.StatusCode, hint, string(raw))
	}

	var sb strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var chunk ollamaChunk
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			continue
		}
		sb.WriteString(chunk.Response)
		if chunk.Done {
			break
		}
	}

	return sb.String(), scanner.Err()
}
