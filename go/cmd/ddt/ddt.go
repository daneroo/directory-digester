package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/daneroo/directory-digester/go/digester"
	"github.com/daneroo/directory-digester/go/logsetup"
)

type DirInfo struct {
	Info        digester.DigestInfo
	Children    []digester.DigestInfo `json:"children,omitempty"`
	NumChildren int
}

func main() {
	logsetup.SetupFormat()

	// Define the directory to walk recursively
	root := "/Users/daniel/Downloads"
	if len(os.Args) > 1 {
		root = os.Args[len(os.Args)-1]
	}
	log.Printf("directory-digester root:%s\n", root) // TODO(daneroo): add version,buildDate

	// Call the filepath.Walk function to recursively walk the directory tree
	var dirStack []DirInfo
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		// This happens if the walker was unable to access the file or directory
		if err != nil {
			return err
		}

		digestInfo, err := digester.Entry(path, info)
		if err != nil {
			return err
		}

		// If the path is a directory, add it to the stack
		if info.IsDir() {

			children, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}

			// add the DirInfo to the stack, but only if we have children
			// otherwise, we can just skip it, but wrap up the DigestInfo with Size and Sha256
			if len(children) == 0 {
				digestInfo.Size = 0
				digestInfo.Sha256 = fmt.Sprintf("%x", sha256.Sum256([]byte{}))
				log.Printf("0Directory: %s (%d)", path, len(children))
			} else {
				dirStack = append(dirStack, DirInfo{
					Info:        digestInfo,
					NumChildren: len(children),
				})
				log.Printf(">Directory: %s (%d)", path, len(children))
			}
		} else {
			// Now we know it is a file
			log.Printf("=File:      %s", digestInfo.Path)
		}

		// add ourselves to our parent
		if len(dirStack) > 0 {
			parent := &dirStack[len(dirStack)-1]
			parent.Children = append(parent.Children, digestInfo)
		}

		// unwind the stack?
		// for len(dirStack) > 0 {
		// 	// pop
		// 	log.Printf("pop? |dirStack|:      (%d)", len(dirStack))
		// 	parent := &dirStack[len(dirStack)-1]
		// 	if len(parent.Children) == parent.NumChildren {
		// 		log.Printf("<Directory: %s (%d =?= %d)", parent.Info.Path, len(parent.Children), parent.NumChildren)
		// 		dirStack = dirStack[:len(dirStack)-1]
		// 	} else {
		// 		break
		// 	}
		// }

		return nil
	})

	// Check for any errors while walking the directory tree
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Current implementation: the stack contains all directories with children

	// Calculate the sha256 digest of the FileInfo JSON structures for each directory's children
	for i := len(dirStack) - 1; i >= 0; i-- {
		dir := &dirStack[i]

		// Sort the list of children by name - should already be sorted, but just in case
		sort.SliceStable(dir.Children, func(i, j int) bool {
			return dir.Children[i].Path < dir.Children[j].Path
		})

		// Sum th sizes of the children
		// Encode the list of children as JSON and calculate its sha256 digest
		childrenJson, err := json.Marshal(dir.Children)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// fmt.Println("CH", string(childrenJson))
		dir.Info.Sha256 = fmt.Sprintf("%x", sha256.Sum256(childrenJson))

		// Encode the DirInfo struct as JSON and print it
		dirInfoJson, err := json.Marshal(dir.Info)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(string(dirInfoJson))
	}
}
