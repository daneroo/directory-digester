#!/usr/bin/env bash

# Show commands for benchmarking
# Go: no done yet, build install, run

# Deno | Docker Deno

REPO_SRC="https://raw.githubusercontent.com/daneroo/directory-digester/main/deno/reference.ts"

# Home is about 23G
BENCH_DIRECTORY="/Volumes/Space/Home-Movies/Tapes/"
DENO_FLAGS="--quiet --allow-read --allow-env"
MACOS_VOLS="-v /Volumes/Space/:/Volumes/Space/:ro"
SYNO_VOLS="-v /volume1/Home-Movies/:/Volumes/Space/Home-Movies:ro"
DOCKER_DENO_TAG="denoland/deno:latest"
cat << EOF
# Deno

## Deno installed (galois, davinci)

time deno run ${DENO_FLAGS} ${REPO_SRC} --verbose ${BENCH_DIRECTORY}

## Docker deno (galois, davinci)

docker pull --quiet ${DOCKER_DENO_TAG}
time docker run --rm -it ${MACOS_VOLS} ${DOCKER_DENO_TAG} run ${DENO_FLAGS} ${REPO_SRC} --verbose ${BENCH_DIRECTORY}

## Docker deno (syno)

docker pull --quiet ${DOCKER_DENO_TAG}
time -f %es docker run --rm -it ${SYNO_VOLS} ${DOCKER_DENO_TAG} run ${DENO_FLAGS} ${REPO_SRC} --verbose ${BENCH_DIRECTORY}

EOF
