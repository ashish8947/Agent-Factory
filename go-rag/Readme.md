# 🚀 Go-RAG — Local RAG System in Go

A **Retrieval-Augmented Generation (RAG) system** built in Go, fully running locally with no external API dependencies.

It combines:
- 🧠 Local LLM (TinyLlama via Ollama)
- 🔎 Vector search (Qdrant)
- 📦 Embeddings (nomic-embed-text)
- ⚡ Custom Go RAG pipeline
- 🌐 HTTP API layer

---

# 🏗️ System Overview

Go-RAG is designed as a **modular, production-ready backend system** for building AI applications.

It supports:
- Document ingestion
- Chunk-based embedding
- Vector search retrieval
- Context-aware LLM generation

---

# 🧩 Architecture

```text
Client
   ↓
HTTP API (Go)
   ↓
Handler Layer
   ↓
RAG Pipeline
   ↓
Retriever (Vector Search)
   ↓
Qdrant (Vector Database)
   ↓
Embedding Model (Ollama)
   ↓
Context Builder
   ↓
TinyLlama (LLM)
   ↓
Response
```

# 📁 Project Structure

cmd/server/              → Application entrypoint
internal/api/            → HTTP handlers + routing
internal/rag/            → Chunking, embeddings, pipeline, retrieval
internal/vectorstore/    → Qdrant client implementation
internal/llm/            → Ollama LLM client

# 🚀 Prerequisites

1. Install Go (1.21+)
```
go version
```
https://go.dev/dl/

2. Install Ollama

```
ollama --version
```
https://ollama.com

3. Pull Required Models

```
ollama pull tinyllama
ollama pull nomic-embed-text
```
4. Run Qdrant
```
docker run -p 6333:6333 qdrant/qdrant
```
Verify:
curl http://localhost:6333/collections

# ▶️ Running the System

```
go run ./cmd/server/main.go
Server starts at:
http://localhost:8080
```
# 📡 API Reference

🔹 Health Check
GET /health
Response
{  "status": "ok"}

🔹 Ingest Document
POST /ingest
Request
{  "doc_id": "doc1",  "text": "Go is a programming language created at Google..."}
Response
{  "status": "ingested"}

🔹 Ask Question (RAG Query)
POST /ask
Request
{  "question": "What is Go?"}
Response
{  "answer": "Go is a programming language created at Google..."}

# 🔄 Data Flow

**Ingestion Flow**

Text → Chunking → Embedding → Qdrant Storage

**Query Flow**

Question → Embedding → Vector Search → Context → LLM → Answer

# 🧱 Design Principles


🔹 Fully local-first (no cloud dependency)

🔹 Modular architecture

🔹 Replaceable components (LLM / DB / embeddings)

🔹 Clean Go idioms

🔹 Scalable backend design



# 📊 Performance Characteristics

Low latency retrieval via Qdrant

Lightweight LLM inference via TinyLlama

Efficient chunk-based embeddings

Stateless API layer


# ⭐ Status

✔ RAG pipeline working

✔ Qdrant integration working

✔ Ollama integration working

✔ API server running

✔ End-to-end retrieval + generation
