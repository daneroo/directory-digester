package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	// "github.com/daneroo/directory-digester/go/logsetup"
	"github.com/daneroo/directory-digester/go/logsetup"
)

type FileInfo struct {
	Name    string `json:"name"`
	ModTime string `json:"mod_time"`
	Mode    string `json:"mode"`
	Sha256  string `json:"sha256"`
}

func main() {
	logsetup.SetupFormat()

	// Define the directory to walk recursively
	// root := "path/to/root/directory"
	root := "/Users/daniel/Downloads"
	if len(os.Args) > 1 {
		root = os.Args[len(os.Args)-1]
	}

	log.Printf("directory-digester root:%s\n", root) // TODO(daneroo): add version,buildDate
	// log.

	// Call the filepath.Walk function to recursively walk the directory tree
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// If the path is a directory, print the directory name
		if info.IsDir() {
			log.Println("Directory:", path)
			return nil
		}

		// Create a FileInfo struct for the file
		fileInfo, err := fileInfo(path, info)
		if err != nil {
			return err
		}

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
		log.Fatal(err)
	}
}

// given info as a os.FileInfo return a FileInfo struct, or error if unable to open file or calculate sha256
func fileInfo(path string, info os.FileInfo) (*FileInfo, error) {
	fileInfo := &FileInfo{
		Name:    info.Name(),
		ModTime: info.ModTime().String(),
		Mode:    info.Mode().String(),
	}
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Calculate the sha256 digest of the file
	digester := sha256.New()
	if _, err := io.Copy(digester, file); err != nil {
		return nil, err
	}
	fileInfo.Sha256 = hex.EncodeToString(digester.Sum(nil))

	return fileInfo, nil
}
