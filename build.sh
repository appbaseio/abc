#!/bin/sh
# https://golang.org/doc/install/source#environment
mkdir -p build && cd build

VERSION=0.6.4

export GOARCH=amd64

export GOOS=darwin

go build -o "abc-${VERSION}" -tags '!oss' ./../cmd/abc/...
zip -r "abc-${GOOS}-${VERSION}.zip" "abc-${VERSION}"

export GOOS=windows

go build -o "abc-${VERSION}.exe"  -tags '!oss' ./../cmd/abc/...
zip -r "abc-${GOOS}-${VERSION}.zip" "abc-${VERSION}.exe"

export GOOS=linux

rm "abc-${VERSION}"
go build -o "abc-${VERSION}" -tags '!oss' ./../cmd/abc/...
zip -r "abc-${GOOS}-${VERSION}.zip" "abc-${VERSION}"
