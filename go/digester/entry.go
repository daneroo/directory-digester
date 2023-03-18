package digester

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"
)

// 	Name    string `json:"name"`
// 	ModTime string `json:"mod_time"`
// 	Mode    string `json:"mode"`
// 	Sha256  string `json:"sha256"`

type DigestInfo struct {
	Path    string      `json:"path"`
	Size    int64       `json:"size"`
	ModTime time.Time   `json:"mod_time"`
	Mode    os.FileMode `json:"mode"`
	Sha256  string      `json:"sha256"`
}

func Entry(path string, fileInfo os.FileInfo) (DigestInfo, error) {

	if fileInfo.IsDir() {
		return DigestInfo{
			Path: path,
			// Will be overridden by sum of children
			Size:    fileInfo.Size(),
			ModTime: fileInfo.ModTime().UTC(),
			Mode:    fileInfo.Mode(),
			// Sha256: // not yet calculated - digest of children
		}, nil
	}

	// Otherwise we have a regular file

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return DigestInfo{}, err
	}
	defer file.Close()

	// Calculate the sha256 digest of the file
	digester := sha256.New()
	if _, err := io.Copy(digester, file); err != nil {
		return DigestInfo{}, err
	}

	return DigestInfo{
		Path:    path,
		Size:    fileInfo.Size(),
		ModTime: fileInfo.ModTime().UTC(),
		Mode:    fileInfo.Mode(),
		// same as hex.EncodeToString(sha[:])
		Sha256: fmt.Sprintf("%x", digester.Sum(nil)),
	}, nil
}
