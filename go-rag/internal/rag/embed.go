package rag

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Embed(text string) []float32 {

	body := map[string]string{
		"model":  "nomic-embed-text",
		"prompt": text,
	}

	data, _ := json.Marshal(body)

	resp, err := http.Post(
		"http://localhost:11434/api/embeddings",
		"application/json",
		bytes.NewBuffer(data),
	)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)

	arr := result["embedding"].([]any)

	vec := make([]float32, len(arr))
	for i, v := range arr {
		vec[i] = float32(v.(float64))
	}

	return vec
}