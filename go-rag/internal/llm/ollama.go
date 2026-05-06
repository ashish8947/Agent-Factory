package llm

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Generate(prompt string) string {

	reqBody := map[string]any{
		"model":  "tinyllama",
		"prompt": prompt,
		"stream": false,
	}

	data, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		"http://localhost:11434/api/generate",
		"application/json",
		bytes.NewBuffer(data),
	)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)

	return result["response"].(string)
}