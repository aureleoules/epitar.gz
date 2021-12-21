package archive

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
)

func listIndexableFiles(rootDir string, accept []string) []string {
	var files []string
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !checkExtension(accept, path) {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files
}

func getExtension(path string) string {
	ext := filepath.Ext(path)
	if ext == "" {
		return ""
	}
	return ext[1:]
}

func checkExtension(accept []string, path string) bool {
	ext := getExtension(path)
	if ext == "" {
		return false
	}
	for _, a := range accept {
		if a == ext {
			return true
		}
	}
	return false
}

func getFileID(data []byte) string {
	h := sha1.Sum(data)
	return hex.EncodeToString(h[:])
}
