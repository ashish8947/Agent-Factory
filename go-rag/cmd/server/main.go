package main

import (
	"context"
	"log"

	"github.com/ashish8947/Agent-Factory/go-rag/internal/api"
	"github.com/ashish8947/Agent-Factory/go-rag/internal/llm"
	"github.com/ashish8947/Agent-Factory/go-rag/internal/rag"
	"github.com/ashish8947/Agent-Factory/go-rag/internal/vectorstore"
)

func main() {

	ctx := context.Background()

	store := vectorstore.NewQdrantClient(
		"http://localhost:6333",
		"rag_collection",
	)

	if err := store.EnsureCollection(ctx, 768); err != nil {
		log.Fatal(err)
	}

	embedder := rag.NewEmbedder("http://localhost:11434", "nomic-embed-text")

	retriever := rag.NewRetriever(embedder, store)

	llmClient := llm.NewOllamaClient("http://localhost:11434", "tinyllama")

	pipeline := rag.NewPipeline(retriever, llmClient)
	ingestor := rag.NewIngestor(embedder, store)
	handler := api.NewHandler(pipeline, ingestor)

	api.StartServer(handler)
}
