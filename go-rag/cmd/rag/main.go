package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ashish8947/Agent-Factory/go-rag/internal/llm"
	"github.com/ashish8947/Agent-Factory/go-rag/internal/loader"
	"github.com/ashish8947/Agent-Factory/go-rag/internal/rag"
	"github.com/ashish8947/Agent-Factory/go-rag/internal/vectorstore"
)

func main() {

	ctx := context.Background()

	// ---------------------------------------------------
	// 1. VECTOR STORE (QDRANT)
	// ---------------------------------------------------
	store := vectorstore.NewQdrantClient(
		"http://localhost:6333",
		"rag_collection",
	)

	// Auto-create collection (IMPORTANT)
	err := store.EnsureCollection(ctx, 768)
	if err != nil {
		log.Fatal("failed to init qdrant:", err)
	}

	// ---------------------------------------------------
	// 2. EMBEDDER (nomic-embed-text)
	// ---------------------------------------------------
	embedder := rag.NewEmbedder(
		"http://localhost:11434",
		"nomic-embed-text",
	)

	// ---------------------------------------------------
	// 3. RETRIEVER
	// ---------------------------------------------------
	retriever := rag.NewRetriever(embedder, store)

	// ---------------------------------------------------
	// 4. LLM (TinyLlama)
	// ---------------------------------------------------
	llmClient := llm.NewOllamaClient(
		"http://localhost:11434",
		"tinyllama",
	)

	// ---------------------------------------------------
	// 5. PIPELINE
	// ---------------------------------------------------
	pipeline := rag.NewPipeline(retriever, llmClient)

	// ---------------------------------------------------
	// 6. INGESTOR
	// ---------------------------------------------------
	ingestor := rag.NewIngestor(embedder, store)

	sampleText := `
Go is an open-source programming language created at Google.
It is designed for simplicity, concurrency, and performance.
It is widely used for backend systems, cloud services, and microservices.
`

	err = ingestor.Ingest(
		ctx,
		"doc1",
		sampleText,
		200, // chunk size
		30,  // overlap
	)

	if err != nil {
		log.Fatal("ingest failed:", err)
	}

	fmt.Println("✅ Ingestion completed")
	docs, err := loader.LoadFiles("./data")
	if err != nil {
		log.Fatal(err)
	}

	for path, content := range docs {
		err := ingestor.Ingest(
			ctx,
			path,
			content,
			200,
			30,
		)
		if err != nil {
			log.Fatal("ingest failed:", err)
		}
	}

	fmt.Println("✅ All files ingested")
	// ---------------------------------------------------
	// 7. QUERY
	// ---------------------------------------------------
	query := "What is Go programming language?"

	answer, err := pipeline.Ask(ctx, query)
	if err != nil {
		log.Fatal("query failed:", err)
	}

	fmt.Println("\n====================")
	fmt.Println("ANSWER:")
	fmt.Println("====================")
	fmt.Println(answer)
}
