package similarity

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"sync"
)

var cache = make(map[string][]float64)
var mu sync.RWMutex

type EmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
}

// 🔥 NEW: Batch embedding function
func GetEmbeddings(texts []string) ([][]float64, error) {

	// Check cache first
	results := make([][]float64, len(texts))
	missing := []string{}
	indexMap := map[int]int{} // map missing index → original index

	mu.RLock()
	for i, text := range texts {
		if emb, ok := cache[text]; ok {
			results[i] = emb
		} else {
			indexMap[len(missing)] = i
			missing = append(missing, text)
		}
	}
	mu.RUnlock()

	// If everything cached → return
	if len(missing) == 0 {
		return results, nil
	}

	// Call API for missing ones
	reqBody := map[string]interface{}{
		"input": missing,
		"model": "text-embedding-3-small",
	}

	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/embeddings",
		bytes.NewBuffer(bodyBytes),
	)

	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result EmbeddingResponse
	json.NewDecoder(resp.Body).Decode(&result)

	// Save to cache
	mu.Lock()
	for i, item := range result.Data {
		originalIndex := indexMap[i]
		text := texts[originalIndex]

		cache[text] = item.Embedding
		results[originalIndex] = item.Embedding
	}
	mu.Unlock()

	return results, nil
}
