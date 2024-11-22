# Deno implementation

This was ported from tor go reference implementation

## Usage

```bash
deno run --allow-sys --allow-read --allow-env deno/reference.ts --verbose testDirectories/rootDir01/
# select just the digest or name from json
deno run --allow-sys --allow-read --allow-env deno/reference.ts --json testDirectories/rootDir01/ | jq '.[]|.sha256'
deno run --allow-sys--allow-read --allow-env deno/reference.ts --json testDirectories/rootDir01/ | jq '.[]|.name'
```

### Run from GitHub (jsr later, when published)

Use `deno run --reload` to bust the github src url cache (REPO_SRC)

```bash
# with deno installed
time deno run --reload --allow-sys --allow-read --allow-env https://raw.githubusercontent.com/daneroo/directory-digester/main/deno/reference.ts --verbose  /Volumes/Space/Home-Movies/Tapes/

# with docker  tag latest, or 2.0.4?
# (on syno)
docker pull denoland/deno:latest;
time docker run --rm -it -v /volume1/Home-Movies/:/Volumes/Space/Home-Movies:ro denoland/deno:latest run --quiet --reload --allow-sys --allow-read --allow-env https://raw.githubusercontent.com/daneroo/directory-digester/main/deno/reference.ts --verbose /Volumes/Space/Home-Movies/Tapes/

# MacOS
docker pull denoland/deno:latest;
time docker run --rm -it -v /Volumes/Space/:/Volumes/Space/:ro denoland/deno:latest run --quiet --reload --allow-sys --allow-read --allow-env https://raw.githubusercontent.com/daneroo/directory-digester/main/deno/reference.ts --verbose /Volumes/Space/Home-Movies/Tapes/
```

## Progress Bar

Also, in prog-cli, we are not actuall digesting, it's just to test the UI and progress bars.

We want to capture the state for the progress bar as you said:

- bytes / totalBytes
- entries / totalEntries

we will do this in 2 phases.

- phase 1: discover the count/sizes of files - to establish totals (all levels.)
  - indeterminate progress bar (still with levels)
- phase 2: second phase digest each file (although synthetic in this case.), and mark progress

```bash
deno run --allow-sys --allow-read --allow-env deno/prog-cli.ts /Volumes/Space/Reading/audiobooks/
```
