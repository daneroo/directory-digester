package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

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

		// Print the file name
		fmt.Println("File:", path)

		// Print the last modification time
		fmt.Println("Last modified:", info.ModTime())

		// Print the permissions mode
		fmt.Println("Permissions mode:", info.Mode().String(), "(", info.Mode(), ")")

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
		digest := hex.EncodeToString(hasher.Sum(nil))
		fmt.Println("SHA256 digest:", digest)

		return nil
	})

	// Check for any errors while walking the directory tree
	if err != nil {
		fmt.Println("Error:", err)
	}
}
