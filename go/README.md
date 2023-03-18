# directory-digester Go implementation

Not sure how to layout the directories yet. Especiall for a multi-language project.
See [this article](https://appliedgo.com/blog/go-project-layout/) for some ideas.

```bash
# setup in this directory
go mod init github.com/daneroo/directory-digester/go
go mod tidy

# setup in repository root to setup a go workspace, so our editor can find this module
cd ..
go work init
go work use ./go
```
