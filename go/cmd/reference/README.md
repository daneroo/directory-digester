# reference implementation

Let's get this as simple as possible.

Using our own Leaf data structure

BUG: character encoding for filenames: fixed example: "/Volumes/Archive/media/audiobooks/Paul Halpern - Einstein's Dice and Schrodinger's Cat/.."

## Data Structures

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
time go run go/cmd/reference/reference.go --verbose testDirectories/rootDir01/
# select just the digest from json
time go run go/cmd/reference/reference.go --json testDirectories/rootDir01/ | jq '.[] | .sha256'
```

```bash
time go run go/cmd/reference/reference.go testDirectories/rootDir01/
go build go/cmd/reference/reference.go; hyperfine './reference testDirectories/rootDir01/'
```
