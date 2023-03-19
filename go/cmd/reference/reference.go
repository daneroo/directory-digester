package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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
		Size:    info.Size(), // not used for directories, will be replaced by sum of children
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

// setSizeOfParent: sets the size of the parent node by summing the size of the children
func setSizeOfParent(node *DigestTreeNode) {
	// sum the size of the children
	var size int64
	for _, child := range node.Children {
		size += child.Info.Size
	}
	node.Info.Size = size
}

// digestNode: calculates the digest of a node
// This can be invoked on a leaf node, or a directory node.
// On the directory it is assumed that the children have been previously digested
func digestNode(node *DigestTreeNode) error {
	digester := sha256.New()
	if !node.Info.Mode.IsDir() {
		// Calculate the sha256 digest of the file
		file, err := os.Open(node.Path)
		if err != nil {
			return err
		}
		defer file.Close()

		if _, err := io.Copy(digester, file); err != nil {
			return err
		}
		node.Info.Sha256 = fmt.Sprintf("%x", digester.Sum(nil))
		if *verboseFlag {
			log.Printf("digestNode(%s) = %s (leaf)\n", node.Path, node.Info.Sha256)
		}
	} else {
		// Calculate the sha256 digest of the children
		for _, child := range node.Children {
			digester.Write([]byte(child.Info.Sha256))
		}
		node.Info.Sha256 = fmt.Sprintf("%x", digester.Sum(nil))
		if *verboseFlag {
			log.Printf("digestNode(%s) = %s (node)\n", node.Path, node.Info.Sha256)
		}
	}
	return nil
}

func buildTree(parentPath string, parentInfo fs.FileInfo) (DigestTreeNode, error) {
	if *verboseFlag {
		log.Printf("buildTree(%s)\n", parentPath)
	}
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
		info, err := file.Info() // fs.DirEntry.Info() may throw an error
		if err != nil {
			return DigestTreeNode{}, err
		}
		var node DigestTreeNode
		if !file.IsDir() { // not a directory, so leaf node
			node = newLeaf(path, info)
			// digest of the leaf node
			err = digestNode(&node)
			if err != nil {
				return DigestTreeNode{}, err
			}
		} else { // directory, so recurse
			node, err = buildTree(path, info)
			if err != nil {
				return DigestTreeNode{}, err
			}
		}
		parentNode.Children = append(parentNode.Children, node)
	}
	// This is where we can aggregate the size and digest of the children
	setSizeOfParent(&parentNode)
	err = digestNode(&parentNode)
	if err != nil {
		return DigestTreeNode{}, err
	}
	return parentNode, nil
}

func shortDigest(digest string, maxLength int) string {
	if len(digest) > maxLength {
		return digest[:(maxLength/2)] + ".." + digest[len(digest)-(maxLength/2):]
	}
	return digest
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxNameLength(node DigestTreeNode, depth int) int {
	max := len(node.Info.Name) + depth*2
	for _, child := range node.Children {
		max = maxInt(max, maxNameLength(child, depth+1))
	}
	return max
}

func showAsIndented(node DigestTreeNode, depth int, maxLength int) {
	if (depth == 0) && (maxLength == 0) {
		maxLength = maxNameLength(node, 0)
	}
	pad := fmt.Sprintf("%*s", depth*2, "")
	isDirIndicator := " " // leaf or directory
	if node.Info.Mode.IsDir() {
		isDirIndicator = "/" //fmt.Sprintf("/ (%d)", len(node.Children))
	}
	fmt.Printf("%s%-*s%s - %10d bytes digest:%s\n", pad, maxLength-depth*2, node.Info.Name, isDirIndicator, node.Info.Size, shortDigest(node.Info.Sha256, 16))
	for _, child := range node.Children {
		showAsIndented(child, depth+1, maxLength)
	}
}

func convertTreeToListWithPath(node DigestTreeNode, list *[]DigestInfo) {
	nameAsPathInfo := node.Info
	nameAsPathInfo.Name = node.Path
	*list = append(*list, nameAsPathInfo)
	for _, child := range node.Children {
		convertTreeToListWithPath(child, list)
	}
}

func showTreeAsJson(node DigestTreeNode) error {
	var list []DigestInfo
	convertTreeToListWithPath(node, &list)
	// jsonBytes, err := json.MarshalIndent(list, "", "  ")
	jsonBytes, err := json.Marshal(list)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonBytes))
	return nil
}

// make this global so we can use it all over the place
var verboseFlag = flag.Bool("verbose", false, "verbose output")

func main() {
	logsetup.SetupFormat()

	// cli flags
	// --verbose is global
	var jsonFlag = flag.Bool("json", false, "json output")
	flag.Parse()

	// Define the directory to walk recursively
	rootDirectory := "/Users/daniel/Downloads"
	if flag.NArg() > 0 {
		rootDirectory = flag.Arg(0)
	}
	log.Printf("directory-digester root:%s\n", rootDirectory) // TODO(daneroo): add version,buildDate

	rootInfo, err := os.Stat(rootDirectory)
	if err != nil {
		panic(err)
	}
	rootNode, err := buildTree(rootDirectory, rootInfo)
	if err != nil {
		panic(err)
	}
	if *verboseFlag {
		log.Printf("-- built tree : %s (%d)\n\n", rootNode.Info.Name, len(rootNode.Children))
	}
	if *jsonFlag {
		showTreeAsJson(rootNode)
	} else {
		showAsIndented(rootNode, 0, 0)
	}
}
