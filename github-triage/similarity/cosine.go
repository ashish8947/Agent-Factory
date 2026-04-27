package similarity

import "math"

func CosineSimilarity(a, b []float64) float64 {
	var dot, magA, magB float64

	for i := 0; i < len(a); i++ {
		dot += a[i] * b[i]
		magA += a[i] * a[i]
		magB += b[i] * b[i]
	}

	return dot / (math.Sqrt(magA) * math.Sqrt(magB))
}