package digester

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

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
		Path:    path,
		Size:    fileInfo.Size(),
		ModTime: fileInfo.ModTime().UTC(),
		Mode:    fileInfo.Mode(),
		// same as hex.EncodeToString(sha[:])
		Sha256: fmt.Sprintf("%x", digester.Sum(nil)),
	}, nil
}
