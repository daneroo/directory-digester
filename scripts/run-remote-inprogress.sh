#!/usr/bin/env bash

# Configuration
REPO_SRC="https://raw.githubusercontent.com/daneroo/directory-digester/main/deno/reference.ts"
TARGET_DIR="Home-Movies/Tapes"
DOCKER_DENO_TAG="denoland/deno:latest"
DENO_FLAGS="--quiet --allow-read --allow-env"

echo "=== Running benchmarks on Synology NAS ==="
ssh -t syno "/usr/bin/env bash -c \"
    echo '=== Setup ==='
    echo Shell: \$SHELL \(bash \$BASH_VERSION\)
    echo OS: \$(uname -s)
    echo Architecture: \$(uname -m)
    
    # Initialize variables
    declare -i HAVE_DOCKER=0
    declare -i HAVE_DENO=0
    declare -i HAVE_GO=0
    
    GO_BINARY='./directory-digester-reference'
    
    # Detect SPACE
    if [ -d /volume1 ]; then
        SPACE=/volume1
        echo Space: \$SPACE \(Synology\)
    elif [ -d /Volumes/Space ]; then
        SPACE=/Volumes/Space
        echo Space: \$SPACE \(MacOS\)
    else
        echo 'Error: Neither /volume1 nor /Volumes/Space exists'
        exit 1
    fi
    
    TARGET_PATH=\${SPACE}/${TARGET_DIR}
    DOCKER_VOLS='-v '\${SPACE}:\${SPACE}':ro,cached'
    
    # Check prerequisites
    # Check Docker
    if command -v docker >/dev/null && docker version >/dev/null 2>&1; then
        echo 'Docker: available and running'
        HAVE_DOCKER=1
    else
        echo 'Docker: not available or not running'
    fi
    
    # Check Deno
    if command -v deno >/dev/null; then
        echo 'Deno: available'
        HAVE_DENO=1
    else
        echo 'Deno: not available'
    fi
    
    # Check Go binary
    if [ -x \${GO_BINARY} ]; then
        echo 'Go binary: available'
        HAVE_GO=1
    else
        echo 'Go binary: not available'
    fi
    
    # Run available benchmarks
    if ((HAVE_GO == 1)); then
        echo '=== Running Go benchmark ==='
        \${GO_BINARY} --verbose \${TARGET_PATH}
    fi
    
    if ((HAVE_DENO == 1)); then
        echo '=== Running native Deno benchmark ==='
        deno run ${DENO_FLAGS} ${REPO_SRC} --verbose \${TARGET_PATH}
    fi
    
    if ((HAVE_DOCKER == 1)); then
        echo '=== Running Docker Deno benchmark ==='
        docker pull --quiet ${DOCKER_DENO_TAG}
        docker run --rm \${DOCKER_VOLS} ${DOCKER_DENO_TAG} run ${DENO_FLAGS} ${REPO_SRC} --verbose \${TARGET_PATH}
    fi
\""
'