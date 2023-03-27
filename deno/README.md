# Deno implementation

This was ported from tor go reference implementation

## Usage

```bash
deno run --allow-read --allow-env deno/reference.ts --verbose testDirectories/rootDir01/
# select just the digest or name from json
deno run --allow-read --allow-env deno/reference.ts --json testDirectories/rootDir01/ | jq '.[]|.sha256'
deno run --allow-read --allow-env deno/reference.ts --json testDirectories/rootDir01/ | jq '.[]|.name'
```

### Remotely (no git checkout)

```bash
deno run --allow-read --allow-env https://raw.githubusercontent.com/daneroo/directory-digester/main/deno/reference.ts --verbose  /Volumes/Space/archive/media/audiobooks/
# with docker (on syno)
docker run --rm -it -v /volume1/Archive/:/Volumes/Space/archive:ro denoland/deno:1.32.1 run --allow-read --allow-env https://raw.githubusercontent.com/daneroo/directory-digester/main/deno/reference.ts --verbose /Volumes/Space/archive/media/MAARIF-IRM/
```
