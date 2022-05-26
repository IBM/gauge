#!/bin/bash

VERSION=$(git describe --tags --always)
COMMIT=$(git rev-parse HEAD)
BUILD=$(date +%FT%T%z)
PKG=github.com/tap8stry/gauge/cmd/discover/cli

LDFLAGS="-X $PKG.version=$VERSION -X $PKG.commit=$COMMIT -X $PKG.date=$BUILD"

archs=(amd64 arm64 ppc64)
platforms=(darwin linux)

for os in ${platforms[@]}
do
    for arch in ${archs[@]}
    do
        env GOOS=${os} GOARCH=${arch} go build -o gauge-${os}-${arch} cmd/gauge/main.go
    done    
done