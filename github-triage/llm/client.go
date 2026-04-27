package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github-triage/models"
)

type Client interface {
	Analyze(text string) (*models.Analysis, error)
}

type OpenAIClient struct {
	apiKey string
	model  string
}

func NewOpenAIClient() *OpenAIClient {
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini"
	}

	return &OpenAIClient{
		apiKey: os.Getenv("OPENAI_API_KEY"),
		model:  model,
	}
}

type openAIRequest struct {
	Model    string        `json:"model"`
	Messages []messageItem `json:"messages"`
}

type messageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message messageItem `json:"message"`
	} `json:"choices"`
}

func (c *OpenAIClient) Analyze(text string) (*models.Analysis, error) {

	log.Println("[LLM] Starting analysis")

	prompt := fmt.Sprintf(`
You are a GitHub issue triage assistant.

Return STRICT JSON only:
{
  "category": "bug | feature | docs | other",
  "summary": "short summary",
  "action_items": ["list of actions"]
}

Issue:
%s
`, text)

	reqBody := openAIRequest{
		Model: c.model,
		Messages: []messageItem{
			{Role: "user", Content: prompt},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Println("[LLM][ERROR] Failed to marshal request:", err)
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		log.Println("[LLM][ERROR] Failed to create request:", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[LLM][ERROR] Request failed:", err)
		return nil, err
	}
	defer resp.Body.Close()

	log.Printf("[LLM] Response received (status=%d, duration=%s)\n",
		resp.StatusCode, time.Since(start))

	var result openAIResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Println("[LLM][ERROR] Failed to decode response:", err)
		return nil, err
	}

	if len(result.Choices) == 0 {
		log.Println("[LLM][ERROR] Empty response from model")
		return nil, fmt.Errorf("empty response from LLM")
	}

	content := result.Choices[0].Message.Content
	log.Println("[LLM] Raw response:", content)

	// Clean JSON
	content = extractJSON(content)

	var analysis models.Analysis
	err = json.Unmarshal([]byte(content), &analysis)
	if err != nil {
		log.Println("[LLM][ERROR] JSON parse failed:", err)
		log.Println("[LLM][ERROR] Cleaned content:", content)
		return nil, fmt.Errorf("failed to parse LLM response: %v", err)
	}

	log.Printf("[LLM] Parsed result: category=%s\n", analysis.Category)

	return &analysis, nil
}

// Extract JSON safely
func extractJSON(input string) string {
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start == -1 || end == -1 {
		return input
	}
	return input[start : end+1]
}