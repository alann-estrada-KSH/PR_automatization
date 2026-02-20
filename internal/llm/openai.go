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

// openAIClient is compatible with OpenAI, Groq, OpenRouter, and any
// OpenAI-spec API — just change base_url + api_key.
type openAIClient struct {
	baseURL string
	model   string
	apiKey  string
	http    *http.Client
}

func newOpenAIClient(cfg *config.Config) (*openAIClient, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("api_key is required for provider %q — set PRGEN_API_KEY or GROQ_API_KEY env var", cfg.Provider)
	}
	return &openAIClient{
		baseURL: cfg.APIBaseURL,
		model:   cfg.Model,
		apiKey:  cfg.APIKey,
		http:    &http.Client{Timeout: 10 * time.Minute},
	}, nil
}

func (c *openAIClient) Name() string { return c.baseURL + "/" + c.model }

// ── Request / Response types ──────────────────────────────────────────────────

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type streamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

// ── Generate ──────────────────────────────────────────────────────────────────

func (c *openAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody, _ := json.Marshal(chatRequest{
		Model: c.model,
		Messages: []chatMessage{
			{Role: "user", Content: prompt},
		},
		Stream: true,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.baseURL+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("LLM request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		hint := ""
		switch resp.StatusCode {
		case 401:
			hint = " → API key inválida o no configurada (PRGEN_API_KEY / GROQ_API_KEY)"
		case 403:
			hint = " → Sin permisos para este modelo"
		case 404:
			hint = " → Modelo no encontrado. Verifica el nombre del modelo en tu config"
		case 429:
			hint = " → Rate limit alcanzado. Espera unos segundos o cambia de proveedor"
		case 503:
			hint = " → Proveedor temporalmente no disponible. Intenta de nuevo en unos minutos"
		}
		return "", fmt.Errorf("API error HTTP %d%s\nRespuesta: %s", resp.StatusCode, hint, string(raw))
	}

	var sb strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk streamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		for _, choice := range chunk.Choices {
			sb.WriteString(choice.Delta.Content)
		}
	}

	return sb.String(), scanner.Err()
}
