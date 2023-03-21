# Deno implementation

This was ported from tor go reference implementation

## Usage

```bash
deno run --allow-read deno/reference.ts --verbose go/
deno run --allow-read deno/reference.ts --json  | jq '.[]|.sha256'

```
