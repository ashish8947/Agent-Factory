package vectorstore

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const qdrantURL = "http://localhost:6333"

func Search(vector []float32) []string {

	reqBody := map[string]any{
		"vector": vector,
		"limit":  3,
	}

	data, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		qdrantURL+"/collections/docs/points/search",
		"application/json",
		bytes.NewBuffer(data),
	)

	if err != nil {
		return []string{}
	}

	defer resp.Body.Close()

	return []string{
		"sample context 1",
		"sample context 2",
	}
}