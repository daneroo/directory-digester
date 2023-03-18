package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name    string `json:"name"`
	ModTime string `json:"mod_time"`
	Mode    string `json:"mode"`
	Sha256  string `json:"sha256"`
}

func main() {
	// Define the directory to walk recursively
	// root := "path/to/root/directory"
	root := "/Users/daniel/Downloads"

	// Call the filepath.Walk function to recursively walk the directory tree
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If the path is a directory, print the directory name
		if info.IsDir() {
			fmt.Println("Directory:", path)
			return nil
		}

		// Create a FileInfo struct for the file
		fileInfo := FileInfo{
			Name:    path,
			ModTime: info.ModTime().String(),
			Mode:    info.Mode().String(),
		}

		// Open the file
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Calculate the sha256 digest of the file
		hasher := sha256.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return err
		}
		fileInfo.Sha256 = hex.EncodeToString(hasher.Sum(nil))

		// Encode the FileInfo struct as JSON and print it
		fileInfoJson, err := json.Marshal(fileInfo)
		if err != nil {
			return err
		}
		fmt.Println(string(fileInfoJson))

		return nil
	})

	// Check for any errors while walking the directory tree
	if err != nil {
		fmt.Println("Error:", err)
	}
}
