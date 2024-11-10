# Deno implementation

This was ported from tor go reference implementation

## Usage

```bash
deno run --allow-read --allow-env deno/reference.ts --verbose testDirectories/rootDir01/
# select just the digest or name from json
deno run --allow-read --allow-env deno/reference.ts --json testDirectories/rootDir01/ | jq '.[]|.sha256'
deno run --allow-read --allow-env deno/reference.ts --json testDirectories/rootDir01/ | jq '.[]|.name'
```

### Run from GitHub (jsr later, when published)

```bash
# with deno installed
time deno run --allow-read --allow-env https://raw.githubusercontent.com/daneroo/directory-digester/main/deno/reference.ts --verbose  /Volumes/Space/Home-Movies/Tapes/

# with docker  tag latest, or 2.0.4?
# (on syno)
docker pull denoland/deno:latest;
time docker run --rm -it -v /volume1/Home-Movies/:/Volumes/Space/Home-Movies:ro denoland/deno:latest run --quiet --allow-read --allow-env https://raw.githubusercontent.com/daneroo/directory-digester/main/deno/reference.ts --verbose /Volumes/Space/Home-Movies/Tapes/

# MacOS
docker pull denoland/deno:latest;
time docker run --rm -it -v /Volumes/Space/:/Volumes/Space/:ro denoland/deno:latest run --quiet --allow-read --allow-env https://raw.githubusercontent.com/daneroo/directory-digester/main/deno/reference.ts --verbose /Volumes/Space/Home-Movies/Tapes/
```
