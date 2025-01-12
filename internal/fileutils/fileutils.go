package fileutils

import (
	"os"
	"path/filepath"
	"strings"
)

func FindFiles(root string) ([]string, error) {
	var result []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && isCodeFile(path) {
			result = append(result, path)
		}
		return nil
	})
	return result, err
}

func isCodeFile(path string) bool {
	extensions := []string{".go", ".py", ".js"}
	lowerPath := strings.ToLower(path)
	for _, ext := range extensions {
		if strings.HasSuffix(lowerPath, ext) {
			return true
		}
	}
	return false
}
