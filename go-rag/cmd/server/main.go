package main

import (
	"fmt"
	"go-rag/internal/api"
)

func main() {
	server := api.NewServer()
	fmt.Println("Server Start")
	server.Start()
	fmt.Println("Server Start")
}

// API (api/)
//    ↓
// RAG Pipeline (rag/)
//    ↓
// Embedding (rag/embed.go)
//    ↓
// Vector DB (vectorstore/qdrant.go)
//    ↓
// Context Builder (rag/prompt.go)
//    ↓
// LLM (llm/ollama.go)
//    ↓
// Response

// AI-powered smart contracts
// Fraud detection in blockchain
// Predictive analytics for crypto markets
// Decentralized AI systems