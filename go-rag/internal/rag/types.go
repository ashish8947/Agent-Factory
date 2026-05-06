package rag

import "context"

// Document is the core unit of RAG system
type Document struct {
	ID      string
	Content string
	Score   float64
	Meta    map[string]any
}

// Embedder defines embedding contract
type Embedder interface {
	Embed(ctx context.Context, text string) ([]float64, error)
}

// VectorStore defines storage + retrieval contract
type VectorStore interface {
	Upsert(ctx context.Context, docID string, chunkIndex int, text string, vector []float64) error
	Search(ctx context.Context, vector []float64, topK int) ([]Document, error)
}