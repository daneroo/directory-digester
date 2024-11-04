# directory-digester

`directory-digester` is a tool to produce message digests for files and directories.

It has similar goals to [`md5deep`](https://github.com/jessek/hashdeep) (a.k.a `hashdeep`) but,

- it also produces digests for directories themselves as defined below.
- it also produces digests that account for file name, size, permissions modification time.

At the simplest level, it produces a digest for a file's content.

## Stretch goals

- Should have multiple implementations (go, typescript (node,deno), rust)
- Could distinguish between content only digests and content + metadata digests.
- Could include compare/verify functionality.
  - should compare visually like [difftastic](https://github.com/Wilfred/difftastic)
  - could compare on different hosts
  - could compare using nats as a message bus
  - could compare using ipfs as a signing mechanism
- Could include filtering functionality. (exclude file patterns for example)
- Could use different algorithms, including SHA-1, SHA-256, SHA-512, and SHA-3.

## Digests for directories

- The digest for a directory is the digest of the concatenation of the digests of the directory's entries.
  - The digest for a directory's entry is the digest of the concatenation of the entry's name, size, permissions, modification time, and digest. (as JSON)
- traversal order: lexicographic
- if file - print json line
- if directory
  - descend into directory and print json line for each entry
  - then print json line for directory itself

## Reference Implementation

Let's start with a simple, reference implementation in go. That does not optimize for memory or speed.

## Publishing and Releasing

Git Branch Tagging: As per the `goreleaser` conventions we are tagging releases with `v0.1.0` format.

This [article](https://namiops.medium.com/golang-multi-arch-docker-image-with-github-action-b59a62c8d2bd)
and [repo](https://github.com/namiops/go_multiarch/tree/master) have instructions on building a _Go_ binary for multiple architectures and publishing them into a multi-arch docker image.

## Till we get CI/CD and goreleaser working

```bash
# build
./build/build-go.sh go/cmd/reference/reference.go

#  from davinci,dirac,shannon: copy from galois
scp -p galois:Downloads/directory-digester/reference.go-darwin-amd64 .
time ./reference.go-darwin-amd64 --verbose  /Volumes/Space/archive/
#  for syno: copy from galois
scp -p galois:Downloads/directory-digester/reference.go-linux-amd64 .
time ./reference.go-linux-amd64 --verbose  /volume1/Archive/
```

## Performance

We will need to balance CPU/digest and IO to optimize speed.

Initial speed reference for **Go** version

### `/Volumes/Space (estimate)

| Machine | Estimated Time |
| :------ | -------------: |
| galois  |     1.79 hours |
| davinci |     3.85 hours |
| syno    |    25.43 hours |

### Home-Movies

| Machine | Time (s) | Data (MB) | Rate (MB/s) |
| :------ | -------: | --------: | ----------: |
| galois  |  346.751 |    130075 |      375.12 |
| davinci |  799.854 |    130075 |      162.62 |
| syno    | 2402.813 |    130075 |       54.13 |

### Archive

| Machine |  Time (s) | Data (MB) | Rate (MB/s) |
| :------ | --------: | --------: | ----------: |
| galois  |  2866.307 |   1015967 |      354.45 |
| davinci |  6161.966 |   1015967 |      164.88 |
| syno    | 40708.423 |   1015967 |       24.96 |
