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

This [article](https://namiops.medium.com/golang-multi-arch-docker-image-with-github-action-b59a62c8d2bd)
ans [repo](https://github.com/namiops/go_multiarch/tree/master) have instructions on building a *Go* binary for multiple architectures and publishing it to a docker image.