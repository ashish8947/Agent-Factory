package loader

import (
	"os"
	"path/filepath"
)

func LoadFiles(dir string) (map[string]string, error) {
	docs := make(map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		docs[path] = string(content)
		return nil
	})

	return docs, err
}
