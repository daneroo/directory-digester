package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	"github.com/daneroo/directory-digester/go/digester"
	"github.com/daneroo/directory-digester/go/logsetup"
)

type DigestTree struct {
	Info     digester.DigestInfo
	Children []DigestTree
}

func getLeaf(path string, fileInfo fs.FileInfo) (DigestTree, error) {
	log.Printf("getLeaf: %s", path)
	digestInfo, err := digester.Entry(path, fileInfo)
	if err != nil {
		return DigestTree{}, err
	}

	return DigestTree{
		Info:     digestInfo,
		Children: make([]DigestTree, 0),
	}, nil
}

func getDir(path string, fileInfo fs.FileInfo) (DigestTree, error) {
	log.Printf("getDir: %s", path)
	digestInfo, err := digester.Entry(path, fileInfo)
	if err != nil {
		return DigestTree{}, err
	}

	return DigestTree{
		Info:     digestInfo,
		Children: make([]DigestTree, 0),
	}, nil
}

func getChildren(path string) ([]DigestTree, error) {
	children := make([]DigestTree, 0)
	childInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return []DigestTree{}, err
	}
	for _, childInfo := range childInfos {
		childPath := fmt.Sprintf("%s/%s", path, childInfo.Name())
		if childInfo.IsDir() {
			child, err := getDir(childPath, childInfo)
			if err != nil {
				return []DigestTree{}, err
			}
			children = append(children, child)

		} else {
			child, err := getLeaf(childPath, childInfo)
			if err != nil {
				return []DigestTree{}, err
			}
			children = append(children, child)

		}
	}
	return children, nil
}

func getNode(path string) (DigestTree, error) {
	log.Printf("getNode: %s", path)

	fileInfo, err := os.Lstat(path)
	if err != nil {
		return DigestTree{}, err
	}

	// Leaf Node
	if !fileInfo.IsDir() {
		digestLeaf, err := getLeaf(path, fileInfo)
		if err != nil {
			return DigestTree{}, err
		}
		return digestLeaf, nil

	} else {
		digestInfo, err := digester.Entry(path, fileInfo)
		if err != nil {
			return DigestTree{}, err
		}

		// childrenDigestInfos, err := digester.Directory(path)
		childrenDigestInfos, err := getChildren(path)
		if err != nil {
			return DigestTree{}, err
		}
		children := make([]DigestTree, 0)
		for idx, childDigestInfo := range childrenDigestInfos {
			fileInfo, err := os.Lstat(path)
			if err != nil {
				return DigestTree{}, err
			}

			if !childDigestInfo.Info.Mode.IsDir() {
				log.Printf("...should getLeaf[%d]: %s", idx, childDigestInfo.Info.Path)

				child, err := getLeaf(childDigestInfo.Info.Path, fileInfo)
				if err != nil {
					return DigestTree{}, err
				}
				children = append(children, child)
			} else {
				log.Printf("...should getDir[%d]: %s", idx, childDigestInfo.Info.Path)
				child, err := getDir(childDigestInfo.Info.Path, fileInfo)
				if err != nil {
					return DigestTree{}, err
				}
				children = append(children, child)
			}
		}

		return DigestTree{
			Info:     digestInfo,
			Children: children,
		}, nil
	}
}

func showNode(node DigestTree, depth int) {
	pad := fmt.Sprintf("%*s", depth, " ")
	log.Printf("%s%s (%d) d:%d", pad, node.Info.Path, len(node.Children), depth)
	for _, child := range node.Children {
		showNode(child, depth+1)
	}
}
func main() {
	logsetup.SetupFormat()

	// Define the directory to walk recursively
	root := "/Users/daniel/Downloads"
	if len(os.Args) > 1 {
		root = os.Args[len(os.Args)-1]
	}
	log.Printf("directory-digester root:%s\n", root) // TODO(daneroo): add version,buildDate

	// Check for any errors while walking the directory tree
	node, err := getNode(root)
	if err != nil {
		log.Fatal(err)
	}
	showNode(node, 0)
}
