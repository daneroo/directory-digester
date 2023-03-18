package digester

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Directory(dirPath string) ([]byte, error) {
	var files []DigestInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// Process only regular files
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			hash := sha256.New()
			if _, err := io.Copy(hash, file); err != nil {
				return err
			}

			files = append(files, DigestInfo{
				Name:    info.Name(),
				ModTime: info.ModTime().UTC(),
				Mode:    info.Mode(),
				Sha256:  fmt.Sprintf("%x", hash.Sum(nil)),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(files)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
