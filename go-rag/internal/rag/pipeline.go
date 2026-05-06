package rag

import (
	"github.com/ashish8947/Agent-Factory/go-rag/internal/llm"
	"github.com/ashish8947/Agent-Factory/go-rag/internal/vectorstore"
)

func Ask(question string) string {

	vec := Embed(question)

	context := vectorstore.Search(vec)

	prompt := BuildPrompt(context, question)

	answer := llm.Generate(prompt)

	return answer
}
