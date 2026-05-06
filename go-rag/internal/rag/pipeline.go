package rag

import (
	"context"
	"fmt"
	"strings"
)

// -----------------------------
// Interfaces (decoupled design)
// -----------------------------

type LLM interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type Pipeline struct {
	Retriever *Retriever
	LLM       LLM
}

// -----------------------------
// Constructor
// -----------------------------

func NewPipeline(r *Retriever, llm LLM) *Pipeline {
	return &Pipeline{
		Retriever: r,
		LLM:       llm,
	}
}

// -----------------------------
// Public Entry Point (CORE)
// -----------------------------

func (p *Pipeline) Ask(ctx context.Context, query string) (string, error) {

	// 1. Retrieve relevant documents
	docs, err := p.Retriever.Retrieve(ctx, query)
	if err != nil {
		return "", fmt.Errorf("retrieval failed: %w", err)
	}

	// 2. Build context from docs
	contextStr := p.Retriever.BuildContext(docs)

	// 3. Build final prompt
	prompt := p.buildPrompt(query, contextStr)

	// 4. Call local LLM (Ollama)
	response, err := p.LLM.Generate(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("llm generation failed: %w", err)
	}

	return response, nil
}

// -----------------------------
// Prompt Builder (IMPORTANT)
// -----------------------------

func (p *Pipeline) buildPrompt(query, context string) string {

	var sb strings.Builder

	sb.WriteString("You are a helpful assistant.\n")
	sb.WriteString("Use ONLY the provided context to answer.\n")
	sb.WriteString("If the answer is not in the context, say: 'I don't know'.\n\n")

	sb.WriteString("### Context:\n")
	sb.WriteString(context)
	sb.WriteString("\n\n")

	sb.WriteString("### Question:\n")
	sb.WriteString(query)
	sb.WriteString("\n\n")

	sb.WriteString("### Answer:")

	return sb.String()
}