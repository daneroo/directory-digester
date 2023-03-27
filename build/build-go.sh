#!/usr/bin/env bash

# See <https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04>
# 
# This script builds a Go package for multiple platforms.
# It is intended to be run from the root of the repository.
# Usage: ./build/build-go.sh <package-name>
#   e.g. ./build/build-go.sh go/cmd/reference/reference.go
package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi
package_split=(${package//\// })
package_name=${package_split[-1]}
	
platforms=("linux/amd64" "darwin/amd64" "darwin/arm64")

export VERSION=$(git describe --dirty --always)
export COMMIT=$(git rev-parse --short HEAD)
export BUILDDATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })
	export GOOS=${platform_split[0]}
	export GOARCH=${platform_split[1]}
	output_name=$package_name'-'$GOOS'-'$GOARCH
	if [ $GOOS = "windows" ]; then
		output_name+='.exe'
	fi	

	# here is the actual build!
	go build -ldflags="-X 'main.version=${VERSION}' -X 'main.commit=${COMMIT}' -X 'main.buildDate=${BUILDDATE}'" -o $output_name $package
	if [ $? -ne 0 ]; then
   		echo 'An error has occurred! Aborting the script execution...'
		exit 1
	fi
  echo Built: $output_name GOOS=$GOOS GOARCH=$GOARCH
done