package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// -----------------------------
// Embedder implementation
// -----------------------------

type OllamaEmbedder struct {
	BaseURL string
	Model   string
	Client  *http.Client
}

// Constructor
func NewEmbedder(baseURL, model string) *OllamaEmbedder {
	return &OllamaEmbedder{
		BaseURL: baseURL,
		Model:   model,
		Client:  &http.Client{},
	}
}

// -----------------------------
// Request/Response structs
// -----------------------------

type embedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type embedResponse struct {
	Embedding []float64 `json:"embedding"`
}

// -----------------------------
// Core embedding function
// -----------------------------

func (e *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float64, error) {

	if text == "" {
		return nil, fmt.Errorf("empty text for embedding")
	}

	reqBody := embedRequest{
		Model:  e.Model,
		Prompt: text,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		e.BaseURL+"/api/embeddings",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama embed request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("embedding API error: %d", resp.StatusCode)
	}

	var result embedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode failed: %w", err)
	}

	if len(result.Embedding) == 0 {
		return nil, fmt.Errorf("empty embedding returned")
	}

	return result.Embedding, nil
}
