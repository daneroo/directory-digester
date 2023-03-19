# reference implementaion

Let's get this as simple as possible.

This implementation only builds a tree of directories and files.

Using our own Leaf data structure

```go
// from
type TreeNode struct {
  Path     string
  Info     os.FileInfo
  Children []TreeNode
}

// to
type DigestTreeNode struct {
  Path     string
  Info     DigestInfo
  Children []DigestTree
}

type DigestInfo struct {
  Name    string      `json:"name"`
  Size    int64       `json:"size"`
  ModTime time.Time   `json:"mod_time"`
  Mode    os.FileMode `json:"mode"`
  Sha256  string      `json:"sha256"`
}
```

## Benchmarking

```bash
go build go/cmd/ref/ref.go; hyperfine './ref go'
```
