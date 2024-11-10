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
	"runtime"
	"time"

	"github.com/daneroo/directory-digester/go/logsetup"
)

// export VERSION=$(git describe --dirty --always)
// export COMMIT=$(git rev-parse --short HEAD)
// export BUILDDATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
// go build -ldflags="-X 'main.version=${VERSION}' -X 'main.commit=${COMMIT}' -X 'main.buildDate=${BUILDDATE}'"
var (
	version   string = "v0.0.0-dev"
	commit    string = "feedbac"              // "c0ffee5"
	buildDate string = "1970-01-01T00:00:00Z" // must be static, not time.Now().UTC().Format(time.RFC3339)
)

type DigestTreeNode struct {
	Path     string
	Info     DigestInfo
	Children []DigestTreeNode
}

// This is likely the structure that will be serialized to JSON
// There is a question of whether the full path (from root of our tree) should be used
// for now, we are using Name (basename of Path), which appropriate for recursive display (with indentation)
// TODO(daneroo): rename mod_time to mtime
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
		start := time.Now()

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

		elapsed := time.Since(start).Seconds()
		sizeMB := float64(node.Info.Size) / 1024 / 1024
		rate := sizeMB / elapsed

		if *verboseFlag {
			log.Printf("digestNode(%s) = %s (leaf) - size: %.2fMB elapsed: %.2fs rate: %.2f MB/s\n",
				node.Path, node.Info.Sha256, sizeMB, elapsed, rate)
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
func ignoreName(name string) bool {
	ignorePatterns := []string{".DS_Store", "@eaDir"}
	for _, pattern := range ignorePatterns {
		if match, _ := filepath.Match(pattern, name); match {
			return true
		}
	}
	return false
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

		// ignore patterns
		ignore := ignoreName(file.Name())
		if ignore {
			if *verboseFlag {
				log.Printf("buildTree(%s) ignoring %s\n", parentPath, path)
			}
			continue
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
	// These two lines are printed to stderr even if !verboseFlag
	// TODO(daneroo) add a silent flag to suppress even these
	log.Printf("directory-digester %s - commit:%s - build:%s - %s\n", version, commit, buildDate, runtime.Version())

	log.Printf("directory-digester start root: %s\n", rootDirectory)

	// TODO(daneroo) replace with newDigestInfo()
	start := time.Now()

	rootInfo, err := os.Stat(rootDirectory)
	if err != nil {
		panic(err)
	}
	rootNode, err := buildTree(rootDirectory, rootInfo)
	if err != nil {
		panic(err)
	}

	elapsed := time.Since(start).Seconds()
	totalSizeMB := float64(rootNode.Info.Size) / 1024 / 1024
	rate := totalSizeMB / elapsed

	log.Printf("directory-digester done  root: %s files: %d - size: %.2fMB  elapsed:  %.2fs rate: %.2f MB/s\n",
		rootNode.Info.Name,
		len(rootNode.Children),
		totalSizeMB,
		elapsed,
		rate)
	if *jsonFlag {
		showTreeAsJson(rootNode)
	} else {
		showAsIndented(rootNode, 0, 0)
	}
}
