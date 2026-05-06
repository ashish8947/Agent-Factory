package rag

import "strings"

// ChunkText splits text into overlapping word-based chunks.
// This improves embedding + retrieval quality compared to character splitting.
func ChunkText(text string, size int, overlap int) []string {

	if text == "" {
		return nil
	}

	if size <= 0 {
		size = 200
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= size {
		overlap = size / 5
	}

	text = normalize(text)

	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var chunks []string

	step := size - overlap
	if step <= 0 {
		step = size
	}

	for start := 0; start < len(words); start += step {

		end := start + size
		if end > len(words) {
			end = len(words)
		}

		chunk := strings.Join(words[start:end], " ")
		chunks = append(chunks, chunk)

		if end == len(words) {
			break
		}
	}

	return chunks
}

// normalize cleans text for better embedding quality.
func normalize(text string) string {
	text = strings.TrimSpace(text)

	// normalize whitespace
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return text
}
