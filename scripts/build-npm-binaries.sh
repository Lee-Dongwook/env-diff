#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BIN_DIR="$ROOT_DIR/npm/vendor/bin"

mkdir -p "$BIN_DIR"
cd "$ROOT_DIR"

build() {
    local goos="$1"
    local goarch="$2"
    local filename="$3"

    echo "Building $filename"
    CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
        go build -trimpath -ldflags="-s -w" \
        -o "$BIN_DIR/$filename" ./cmd/envdiff
}

build darwin arm64 envdiff-darwin-arm64
build darwin amd64 envdiff-darwin-x64
build linux arm64 envdiff-linux-arm64
build linux amd64 envdiff-linux-x64
build windows arm64 envdiff-win32-arm64.exe
build windows amd64 envdiff-win32-x64.exe

