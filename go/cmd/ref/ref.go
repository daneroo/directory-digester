package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/daneroo/directory-digester/go/logsetup"
)

type DigestTreeNode struct {
	Path     string
	Info     DigestInfo
	Children []DigestTreeNode
}

// This is likely the structure that will be serialized to JSON
// There is a question of whether the full path (from root of our tree) should be used
// for now, we are using Name (basename of Path), which appropriate for recursive display (with indentation)
type DigestInfo struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	ModTime time.Time   `json:"mod_time"`
	Mode    os.FileMode `json:"mode"`
	Sha256  string      `json:"sha256"`
}

func newDigestInfo(info fs.FileInfo) DigestInfo {
	return DigestInfo{
		Name:    info.Name(),
		Size:    info.Size(), // not used for directories, will be replaced by sum of choildren
		ModTime: info.ModTime(),
		Mode:    info.Mode(),
		// Sha256:  "",
	}
}

func newLeaf(path string, info fs.FileInfo) DigestTreeNode {
	return DigestTreeNode{
		Path: path,
		Info: newDigestInfo(info),
	}
}

func buildTree(parentPath string, parentInfo fs.FileInfo) (DigestTreeNode, error) {
	log.Printf("buildTree(%s)\n", parentPath)
	parentNode := newLeaf(parentPath, parentInfo)

	// The children of the node we are building : could be empty (dir)
	// ioutil.ReadDir is deprecated, so we use os.ReadDir instead as suggested
	// However we still a full os.FileInfo
	// I always need info for Mode, ModTime. Size is is not used (or is overwritten) for Directories.
	// unfortunately os.DirEntry.Info() may throw an error, so we need to handle that
	files, err := os.ReadDir(parentPath)
	if err != nil {
		return DigestTreeNode{}, err
	}

	for _, file := range files {
		path := filepath.Join(parentPath, file.Name())
		info, err := file.Info()
		if err != nil {
			return DigestTreeNode{}, err
		}
		var node DigestTreeNode
		if !file.IsDir() { // not a directory, so leaf node
			node = newLeaf(path, info)
		} else { // directory, so recurse
			node, err = buildTree(path, info)
			if err != nil {
				return DigestTreeNode{}, err
			}
		}
		parentNode.Children = append(parentNode.Children, node)
	}
	return parentNode, nil
}

func showTree(node DigestTreeNode, depth int) {
	pad := fmt.Sprintf("%*s", depth*2, "")
	isDirIndicator := "" // "🍁" // leaf
	if node.Info.Mode.IsDir() {
		isDirIndicator = fmt.Sprintf("/ (%d)", len(node.Children))
	}
	fmt.Printf("%s%s%s\n", pad, node.Info.Name, isDirIndicator)
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
	fmt.Printf("-- built tree : %s (%d)\n\n", rootNode.Info.Name, len(rootNode.Children))
	showTree(rootNode, 0)
}
