package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/daneroo/directory-digester/go/logsetup"
)

type TreeNode struct {
	Path     string
	Info     os.FileInfo
	Children []TreeNode
}

func newLeaf(path string, info os.FileInfo) TreeNode {
	return TreeNode{
		Path: path,
		Info: info,
	}
}

func buildTree(parentPath string, parentInfo fs.FileInfo) (TreeNode, error) {
	log.Printf("buildTree(%s)\n", parentPath)
	parentNode := newLeaf(parentPath, parentInfo)

	// The children of the node we are building : could be empty (dir)
	files, err := ioutil.ReadDir(parentPath)
	if err != nil {
		return TreeNode{}, err
	}

	for _, file := range files {
		path := filepath.Join(parentPath, file.Name())

		var node TreeNode
		if !file.IsDir() { // not a directory, so leaf node
			node = newLeaf(path, file)
		} else { // directory, so recurse
			node, err = buildTree(path, file)
			if err != nil {
				return TreeNode{}, err
			}
		}
		parentNode.Children = append(parentNode.Children, node)
	}
	return parentNode, nil
}

func showTree(node TreeNode, depth int) {
	pad := fmt.Sprintf("%*s", depth*2, "")
	isDirIndicator := "" // "🍁" // leaf
	if node.Info.IsDir() {
		isDirIndicator = fmt.Sprintf("/ (%d)", len(node.Children))
	}
	fmt.Printf("%s%s%s\n", pad, node.Info.Name(), isDirIndicator)
	for _, child := range node.Children {
		showTree(child, depth+1)
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

	rootInfo, err := os.Stat(root)
	if err != nil {
		panic(err)
	}
	rootNode, err := buildTree(root, rootInfo)
	if err != nil {
		panic(err)
	}
	log.Printf("-- built tree : %s (%d)\n\n", rootNode.Info.Name(), len(rootNode.Children))
	showTree(rootNode, 0)
}
