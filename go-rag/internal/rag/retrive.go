package rag

import (
	"context"
	"fmt"
	"sort"
)

// ----------------------------
// Core Types
// ----------------------------

// type Document struct {
// 	ID      string
// 	Content string
// 	Score   float64
// 	Meta    map[string]any
// }


// ----------------------------
// Retriever
// ----------------------------

type Retriever struct {
	embedder Embedder
	vectorDB VectorStore
	topK     int
	minScore float64
}

// Constructor
func NewRetriever(embedder Embedder, db VectorStore) *Retriever {
	return &Retriever{
		embedder: embedder,
		vectorDB: db,
		topK:     5,
		minScore: 0.70,
	}
}

// ----------------------------
// Main RAG Step: Retrieve
// ----------------------------

func (r *Retriever) Retrieve(ctx context.Context, query string) ([]Document, error) {

	// 1. Embed query
	queryVec, err := r.embedder.Embed(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("embedding failed: %w", err)
	}

	// 2. Vector search
	docs, err := r.vectorDB.Search(ctx, queryVec, r.topK)
	if err != nil {
		return nil, fmt.Errorf("vector search failed: %w", err)
	}

	// 3. Filter weak results
	filtered := make([]Document, 0)
	for _, d := range docs {
		if d.Score < r.minScore {
			continue
		}
		filtered = append(filtered, d)
	}

	// 4. Optional reranking (simple deterministic sort)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Score > filtered[j].Score
	})

	return filtered, nil
}

// ----------------------------
// Context Builder (VERY IMPORTANT)
// ----------------------------

func (r *Retriever) BuildContext(docs []Document) string {
	if len(docs) == 0 {
		return "No relevant context found."
	}

	context := "### Retrieved Context:\n"

	for i, d := range docs {
		context += fmt.Sprintf(
			"\n[Doc %d | Score: %.2f]\n%s\n",
			i+1,
			d.Score,
			d.Content,
		)
	}

	return context
}

// ----------------------------
// Utility: Single-call helper
// ----------------------------

func (r *Retriever) RetrieveAndBuild(ctx context.Context, query string) (string, error) {

	docs, err := r.Retrieve(ctx, query)
	if err != nil {
		return "", err
	}

	return r.BuildContext(docs), nil
}
