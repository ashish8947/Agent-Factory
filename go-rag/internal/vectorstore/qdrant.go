package vectorstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"

	"github.com/ashish8947/Agent-Factory/go-rag/internal/rag"
)

type QdrantClient struct {
	BaseURL    string
	Collection string
	Client     *http.Client
}

func NewQdrantClient(baseURL, collection string) *QdrantClient {
	return &QdrantClient{
		BaseURL:    baseURL,
		Collection: collection,
		Client:     &http.Client{},
	}
}

//
// -----------------------------
// COLLECTION INIT
// -----------------------------
//

func (q *QdrantClient) EnsureCollection(ctx context.Context, vectorSize int) error {

	// Try creating collection
	payload := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     vectorSize,
			"distance": "Cosine",
		},
	}

	b, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(
		ctx,
		"PUT",
		fmt.Sprintf("%s/collections/%s", q.BaseURL, q.Collection),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := q.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	// ✅ IMPORTANT: treat "already exists" as OK
	if res.StatusCode == 200 || res.StatusCode == 201 {
		fmt.Println("✅ Collection ready:", q.Collection)
		return nil
	}

	if res.StatusCode == 409 {
		fmt.Println("⚠️ Collection already exists:", q.Collection)
		return nil
	}

	if res.StatusCode >= 300 {
		return fmt.Errorf("failed to create collection: %s", string(body))
	}

	return nil
}

//
// -----------------------------
// UPSERT
// -----------------------------
//

type qdrantUpsertRequest struct {
	Points []qdrantPointUpsert `json:"points"`
}

type qdrantPointUpsert struct {
	ID      uint64                 `json:"id"`
	Vector  []float64              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

func (q *QdrantClient) Upsert(
	ctx context.Context,
	docID string,
	chunkIndex int,
	text string,
	vector []float64,
) error {

	reqBody := qdrantUpsertRequest{
		Points: []qdrantPointUpsert{
			{
				ID:     safeID(docID, chunkIndex),
				Vector: vector,
				Payload: map[string]interface{}{
					"text":   text,
					"doc_id": docID,
					"chunk":  chunkIndex,
				},
			},
		},
	}

	b, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("%s/collections/%s/points", q.BaseURL, q.Collection)

	req, err := http.NewRequestWithContext(
		ctx,
		"PUT",
		url,
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := q.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)

	if res.StatusCode >= 300 {
		return fmt.Errorf("qdrant upsert failed: %s", string(respBody))
	}

	return nil
}

//
// -----------------------------
// SIMPLE HASH (for stable IDs)
// -----------------------------
//

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}
func (q *QdrantClient) Search(
	ctx context.Context,
	vector []float64,
	topK int,
) ([]rag.Document, error) {

	reqBody := map[string]interface{}{
		"vector":       vector,
		"limit":        topK,
		"with_payload": true,
	}

	b, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("%s/collections/%s/points/search", q.BaseURL, q.Collection)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := q.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("search failed: %s", string(body))
	}

	// parse response
	var result struct {
		Result []struct {
			ID      interface{}            `json:"id"`
			Score   float64                `json:"score"`
			Payload map[string]interface{} `json:"payload"`
		} `json:"result"`
	}

	_ = json.Unmarshal(body, &result)

	var docs []rag.Document

	for _, r := range result.Result {
		text, _ := r.Payload["text"].(string)

		docs = append(docs, rag.Document{
			ID:      fmt.Sprintf("%v", r.ID),
			Content: text,
			Score:   r.Score,
			Meta:    r.Payload,
		})
	}

	return docs, nil
}
func safeID(docID string, chunkIndex int) uint64 {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%s-%d", docID, chunkIndex)))
	return h.Sum64()
}
