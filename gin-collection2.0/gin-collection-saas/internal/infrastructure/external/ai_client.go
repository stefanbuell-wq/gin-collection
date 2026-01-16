package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// AIClient handles AI requests (supports Ollama and Anthropic)
type AIClient struct {
	provider   string
	ollamaURL  string
	model      string
	apiKey     string // for Anthropic
	httpClient *http.Client
	enabled    bool
}

// GinSuggestion represents AI-generated gin information
type GinSuggestion struct {
	Description        string   `json:"description"`
	NoseNotes          string   `json:"nose_notes"`
	PalateNotes        string   `json:"palate_notes"`
	FinishNotes        string   `json:"finish_notes"`
	RecommendedTonics  []string `json:"recommended_tonics"`
	RecommendedGarnish []string `json:"recommended_garnish"`
	Country            string   `json:"country"`
	Region             string   `json:"region"`
	GinType            string   `json:"gin_type"`
	EstimatedPrice     float64  `json:"estimated_price"`
	ABV                float64  `json:"abv"`
}

// AIClientConfig holds configuration for the AI client
type AIClientConfig struct {
	Provider        string
	OllamaURL       string
	Model           string
	AnthropicAPIKey string
	Enabled         bool
}

// NewAIClient creates a new AI client (Ollama or Anthropic)
func NewAIClient(cfg *AIClientConfig) *AIClient {
	return &AIClient{
		provider:  cfg.Provider,
		ollamaURL: cfg.OllamaURL,
		model:     cfg.Model,
		apiKey:    cfg.AnthropicAPIKey,
		enabled:   cfg.Enabled,
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // Longer timeout for local models
		},
	}
}

// IsEnabled returns whether the AI service is enabled
func (c *AIClient) IsEnabled() bool {
	if !c.enabled {
		return false
	}
	if c.provider == "anthropic" {
		return c.apiKey != ""
	}
	// For Ollama, always enabled if config says so (local service)
	return true
}

// SuggestGinInfo generates gin information based on name and brand
func (c *AIClient) SuggestGinInfo(name, brand string) (*GinSuggestion, error) {
	if !c.IsEnabled() {
		return nil, fmt.Errorf("AI service is not enabled")
	}

	prompt := buildGinPrompt(name, brand)

	var responseText string
	var err error

	switch c.provider {
	case "ollama":
		responseText, err = c.callOllama(prompt)
	case "anthropic":
		responseText, err = c.callAnthropic(prompt)
	default:
		return nil, fmt.Errorf("unknown AI provider: %s", c.provider)
	}

	if err != nil {
		return nil, err
	}

	return parseGinSuggestion(responseText)
}

func buildGinPrompt(name, brand string) string {
	return fmt.Sprintf(`Du bist ein Gin-Experte. Generiere Informationen für folgenden Gin:

Name: %s
Marke: %s

Antworte ausschließlich im folgenden JSON-Format (keine andere Ausgabe, kein Markdown):
{"description":"Eine ausführliche Beschreibung des Gins auf Deutsch (2-3 Sätze)","nose_notes":"Aromen in der Nase, kommagetrennt auf Deutsch","palate_notes":"Geschmack am Gaumen, kommagetrennt auf Deutsch","finish_notes":"Nachklang/Abgang, kommagetrennt auf Deutsch","recommended_tonics":["Tonic 1","Tonic 2"],"recommended_garnish":["Garnish 1","Garnish 2"],"country":"Herkunftsland","region":"Region falls bekannt, sonst leer","gin_type":"London Dry, Old Tom, New Western, etc.","estimated_price":35.00,"abv":43.0}

Falls du den Gin nicht kennst, mache plausible Annahmen basierend auf dem Namen und der Marke. Gib NUR das JSON zurück, keine Erklärungen, kein Markdown.`, name, brand)
}

// Ollama API types
type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (c *AIClient) callOllama(prompt string) (string, error) {
	req := ollamaRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := strings.TrimSuffix(c.ollamaURL, "/") + "/api/generate"
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("ollama request failed (is Ollama running?): %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama error (status %d): %s", resp.StatusCode, string(body))
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to parse ollama response: %w", err)
	}

	return ollamaResp.Response, nil
}

// Anthropic API types
type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (c *AIClient) callAnthropic(prompt string) (string, error) {
	req := anthropicRequest{
		Model:     c.model,
		MaxTokens: 1024,
		Messages: []anthropicMessage{
			{Role: "user", Content: prompt},
		},
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var apiResp anthropicResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if apiResp.Error != nil {
		return "", fmt.Errorf("API error: %s", apiResp.Error.Message)
	}

	if len(apiResp.Content) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	return apiResp.Content[0].Text, nil
}

func parseGinSuggestion(responseText string) (*GinSuggestion, error) {
	// Try to extract JSON from the response (in case there's extra text)
	jsonStart := -1
	jsonEnd := -1
	braceCount := 0

	for i, char := range responseText {
		if char == '{' {
			if jsonStart == -1 {
				jsonStart = i
			}
			braceCount++
		} else if char == '}' {
			braceCount--
			if braceCount == 0 && jsonStart != -1 {
				jsonEnd = i + 1
				break
			}
		}
	}

	if jsonStart == -1 || jsonEnd == -1 {
		return nil, fmt.Errorf("no valid JSON found in response")
	}

	jsonText := responseText[jsonStart:jsonEnd]

	var suggestion GinSuggestion
	if err := json.Unmarshal([]byte(jsonText), &suggestion); err != nil {
		return nil, fmt.Errorf("failed to parse suggestion JSON: %w (response: %s)", err, jsonText)
	}

	return &suggestion, nil
}
