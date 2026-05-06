package rag

func ChunkText(text string, size int, overlap int) []string {
	var chunks []string
	runes := []rune(text)

	start := 0

	for start < len(runes) {
		end := start + size
		if end > len(runes) {
			end = len(runes)
		}

		chunks = append(chunks, string(runes[start:end]))
		start += size - overlap
	}

	return chunks
}