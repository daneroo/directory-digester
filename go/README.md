# directory-digester Go implementation

Not sure how to layout the directories yet. Especiall for a multi-language project.
See [this article](https://appliedgo.com/blog/go-project-layout/) for some ideas.

## Usage

```bash
time go run go/cmd/reference/reference.go --verbose go/
time go run go/cmd/reference/reference.go --json go/ | jq '.[]|.name'
```

### Remotely (no git checkout)
This does not work. Cross compile, push exec to destination, and run there. Until we get CI/CD (go-releaser) working, both exec and docker

```bash
time go run https://raw.githubusercontent.com/daneroo/directory-digester/main/go/cmd/reference/reference.go --verbose go/
##  /Volumes/Space/archive/media/audiobooks/
```

## Invoke tests

```bash
go test ./...
# or from the repository root:
go test -v ./go/...

```

## Setup

```bash
# setup in this directory
go mod init github.com/daneroo/directory-digester/go
go mod tidy

# setup in repository root to setup a go workspace, so our editor can find this module
cd ..
go work init
go work use ./go
```
