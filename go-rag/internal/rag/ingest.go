package rag

import (
	"context"
	"fmt"
)

type Ingestor struct {
	Embedder Embedder
	VectorDB VectorStore
}

// NewIngestor creates ingestion pipeline
func NewIngestor(e Embedder, v VectorStore) *Ingestor {
	return &Ingestor{
		Embedder: e,
		VectorDB: v,
	}
}

// Ingest takes raw document → stores into vector DB
func (i *Ingestor) Ingest(ctx context.Context, docID string, text string, chunkSize int, overlap int) error {

	// 1. Chunk text
	chunks := ChunkText(text, chunkSize, overlap)

	if len(chunks) == 0 {
		return fmt.Errorf("no chunks generated")
	}

	// 2. Process each chunk
	for idx, chunk := range chunks {

		// 3. Embed chunk
		vec, err := i.Embedder.Embed(ctx, chunk)
		if err != nil {
			return fmt.Errorf("embedding failed: %w", err)
		}
		fmt.Println("VECTOR SIZE:", len(vec))
		// 4. Store in vector DB
		err = i.VectorDB.Upsert(ctx, docID, idx, chunk, vec)
		if err != nil {
			return fmt.Errorf("upsert failed: %w", err)
		}
	}

	return nil
}
