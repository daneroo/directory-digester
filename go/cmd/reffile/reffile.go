package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/daneroo/directory-digester/go/logsetup"
)

type TreeNode struct {
	Info     os.FileInfo
	Children []*TreeNode
}

func buildTree(parentNode *TreeNode, parentPath string) {

	files, err := ioutil.ReadDir(parentPath)
	if err != nil {
		panic(err)
	}

	// log.Printf("-buildTree: %s  files(%d) nodes(%d)", parentPath, len(files), len(parentNode.Children))
	for _, file := range files {
		// log.Printf("node[%d]: %s", idx, file.Name())
		node := TreeNode{
			Info:     file,
			Children: []*TreeNode{},
		}

		if file.IsDir() {
			buildTree(&node, filepath.Join(parentPath, file.Name()))
		}

		parentNode.Children = append(parentNode.Children, &node)
	}
	// log.Printf("+buildTree: %s  files(%d) nodes(%d)", parentPath, len(files), len(parentNode.Children))
}

func printTree(node TreeNode, depth int) {
	pad := fmt.Sprintf("%*s", depth*2, " ")
	fmt.Printf("%s%s - (%d)\n", pad, node.Info.Name(), len(node.Children))
	for _, child := range node.Children {
		printTree(*child, depth+1)
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

	info, err := os.Stat(root)
	if err != nil {
		panic(err)
	}
	rootNode := TreeNode{
		Info:     info,
		Children: []*TreeNode{},
	}

	buildTree(&rootNode, root)
	log.Printf("-- built tree : %s (%d)\n", rootNode.Info.Name(), len(rootNode.Children))
	printTree(rootNode, 0)
}
