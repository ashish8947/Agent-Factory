package similarity

import (
	"encoding/json"
	"os"
)

const cacheFile = "embedding_cache.json"

func SaveCache() {
	file, _ := os.Create(cacheFile)
	defer file.Close()
	json.NewEncoder(file).Encode(cache)
}

func LoadCache() {
	file, err := os.Open(cacheFile)
	if err != nil {
		return
	}
	defer file.Close()
	json.NewDecoder(file).Decode(&cache)
}
