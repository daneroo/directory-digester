# directory-digester

`directory-digester` is a tool to produce message digests for files and directories.

It has similar goals to [`md5deep`](https://github.com/jessek/hashdeep) (a.k.a `hashdeep`) but,

- it also produces digests for directories themselves as defined below.
- it also produces digests that account for file name, size, permissions modification time.

At the simplest level, it produces a digest for a file's content.

## Stretch goals

- Could use different algorithms, including SHA-1, SHA-256, SHA-512, and SHA-3.
- Could distinguish between content only digests and content + metadata digests.
- Could include compare/verify functionality.
- Could include filtering functionality. (exclude file patterns for example)


## Digests for directories

- traversal order: lexicographic
- if file - print json line
- if directory 
  - descend into directory and print json line for each entry
  - then print json line for directory itself

## Reference Implementation

- see reffile