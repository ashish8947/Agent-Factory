package rag

import (
	"go-rag/internal/llm"
	"go-rag/internal/vectorstore"
)

func Ask(question string) string {

	vec := Embed(question)

	context := vectorstore.Search(vec)

	prompt := BuildPrompt(context, question)

	answer := llm.Generate(prompt)

	return answer
}
