package rag

import "strings"

func BuildPrompt(context []string, question string) string {

	return `
You are a helpful assistant.

Use the context below:

` + strings.Join(context, "\n") + `

Question:
` + question + `

Answer clearly:
`
}