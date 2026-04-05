#!/bin/bash

# Script based on https://gist.github.com/Razikus/86e78534c64529f1f62bf803ebf495c1

VERSION=$1
COMMIT=${COMMIT:-$(git rev-parse HEAD | head -c 8)}
BUILDTIME=$(date +%s)

LDFLAGS="-X 'main.VERSION=${VERSION}' -X 'main.BUILDTIME=${BUILDTIME}' -X 'main.COMMIT=${COMMIT}'"

platforms=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "windows/arm64"
)
OUTPUT=builds/
mkdir -p $OUTPUT

total_platforms=${#platforms[@]}
current=0
for platform in "${platforms[@]}"; do
    current=$((current + 1))
    GOOS=$(echo $platform | cut -d'/' -f1)
    GOARCH=$(echo $platform | cut -d'/' -f2)

    current_platform="$platform"
    
    if [ "$GOOS" = "windows" ]; then
      echo "[$current/$total_platforms] Processing exe platform, static: $current_platform"
      CGO_ENABLED=0 GOARCH=$GOARCH GOOS=$GOOS go build -ldflags "${LDFLAGS}" -o "${OUTPUT}/protonpe.${VERSION}.${GOOS}-${GOARCH}.exe"  -buildmode=exe $TOBUILD
    else
      echo "[$current/$total_platforms] Processing platform, static: $current_platform"
      CGO_ENABLED=0 GOARCH=$GOARCH GOOS=$GOOS go build -ldflags "${LDFLAGS}" -o "${OUTPUT}/protonpe.${VERSION}.${GOOS}.${GOARCH}" $TOBUILD
    fi
done