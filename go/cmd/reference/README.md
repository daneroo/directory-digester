# reference implementation

Let's get this as simple as possible.

Using our own Leaf data structure

```go
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

## Running / Benchmarking

```bash
time go run go/cmd/reference/reference.go go/
go build go/cmd/ref/ref.go; hyperfine './ref go'
```