#!/usr/bin/env bash

# Configuration
HOSTS="syno davinci"

for host in $HOSTS; do
    echo
    echo "=== Running directory-digester on $host ==="
    ssh $host /bin/bash << 'EOF'
        echo '=== Init ==='
        echo "Shell: $0 ($BASH_VERSION)"
        echo "OS: $(uname -s)"
        echo "Architecture: $(uname -m)"

        # Configuration
        # Volume - SPACE
        # Detect SPACE
        if [ -d /volume1 ]; then
            SPACE=/volume1
        elif [ -d /Volumes/Space ]; then
            SPACE=/Volumes/Space
        else
            echo 'Error: Neither /volume1 nor /Volumes/Space exists'
            exit 1
        fi
        # Go
        GO_BINARY='./directory-digester-reference'
        # Deno
        REPO_SRC="https://raw.githubusercontent.com/daneroo/directory-digester/main/deno/reference.ts"
        TARGET_DIR="Home-Movies/Tapes"
        DENO_PERMS="--allow-read --allow-env"
        # Docker
        DOCKER_DENO_TAG="denoland/deno:latest"

        echo '=== Setup ==='
        echo "SPACE: $SPACE"
        echo "TARGET_DIR: $TARGET_DIR"
        echo "DOCKER_DENO_TAG: $DOCKER_DENO_TAG"
        echo "Deno Version $(deno --version|head -1)"
EOF
done

