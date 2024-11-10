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

## Building

To build the multi-architecture Go implementation (`./bin/`):

```bash
# build
./scripts/build-go.sh go/cmd/reference/reference.go
```

## Performance

We will need to balance CPU/digest and IO to optimize speed.

## Operation

```bash
# go verbose, text output
time go run go/cmd/reference/reference.go --verbose testDirectories/rootDir01/
docker run --rm -v "$(pwd)":/app -w /app -e HOSTALIAS=$(hostname -s) golang:latest go run go/cmd/reference/reference.go testDirectories/
# deno verbose, text output
time deno run --allow-read --allow-env deno/reference.ts --verbose testDirectories/rootDir01/
docker run --rm -v "$(pwd)":/app -w /app -e HOSTALIAS=$(hostname -s) denoland/deno:latest deno run --allow-read --allow-env deno/reference.ts testDirectories/rootDir01/
```

### Archive

| Machine | Exec | Time (s) |  Data (MB) | Rate (MB/s) |
|:--------|:-----|---------:|-----------:|------------:|
| galois  | go   |  3830.25 | 1015701.94 |      265.18 |
| syno    | go   |  4009.89 | 1015701.94 |      253.30 |
| davinci | go   |  6695.43 | 1015701.94 |      151.70 |
| dirac   | go   |  8238.48 | 1015701.94 |      123.29 |

### Home-Movies (go/deno/deno-docker)

Careful of caching. Even on Syno...

| Machine | Exec        | Time (s) | Data (MB) | Rate (MB/s) |
|:--------|:------------|---------:|----------:|------------:|
| galois  | deno        |   76.469 |    130075 |     1701.02 |
| galois  | deno-docker |   95.905 |    130075 |     1356.36 |
| davinci | deno        |  147.672 |    130075 |      880.84 |
| syno    | deno-docker |  332.970 |    130075 |      390.65 |
| galois  | go          |   439.72 |    130075 |      295.81 |
| davinci | go          |   963.39 |    130075 |      135.02 |
| syno    | go          |   469.51 |    130075 |      277.04 |

### Reading/audiobooks

| Machine | Exec | Time (s) | Data (MB) | Rate (MB/s) |
|:--------|:-----|---------:|----------:|------------:|
| galois  | go   |  1102.88 |    325593 |      295.22 |
| galois  | deno |  1092.59 |    325593 |      298.00 |
| davinci | go   |  2030.87 |    325593 |      160.32 |
| syno    | go   |  2260.44 |    325593 |      144.04 |
