package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/daneroo/directory-digester/go/digester"
	"github.com/daneroo/directory-digester/go/logsetup"
)

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
			digests, err := digester.Directory(path)
			if err != nil {
				return err
			}
			jsonBytes, err := json.Marshal(digests)
			if err != nil {
				return nil
			}

			fmt.Println(string(jsonBytes))
			return nil
		}

		// Create a FileInfo struct for the file
		digestInfo, err := digester.File(path, info)
		if err != nil {
			return err
		}

		// Encode the FileInfo struct as JSON and print it
		digestJson, err := json.Marshal(digestInfo)
		if err != nil {
			return err
		}
		fmt.Println(string(digestJson))

		return nil
	})

	// Check for any errors while walking the directory tree
	if err != nil {
		log.Fatal(err)
	}
}
