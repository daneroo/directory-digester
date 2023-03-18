package digester

import (
	"crypto/sha256"
	"encoding/json"
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
	Name    string      `json:"name"`
	ModTime time.Time   `json:"mod_time"`
	Mode    os.FileMode `json:"mode"`
	Sha256  string      `json:"sha256"`
}

func File(path string, fileInfo os.FileInfo) (DigestInfo, error) {
	// Get the file info
	// fileInfo, err := os.Stat(path)
	// if err != nil {
	// 	return FileInfo{}, err
	// }

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
		Name:    path,
		ModTime: fileInfo.ModTime().UTC(),
		Mode:    fileInfo.Mode(),
		// same as hex.EncodeToString(sha[:])
		Sha256: fmt.Sprintf("%x", digester.Sum(nil)),
	}, nil
}

func EncodeJSON(fileInfo DigestInfo) ([]byte, error) {
	jsonBytes, err := json.MarshalIndent(fileInfo, "", "  ")
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}
