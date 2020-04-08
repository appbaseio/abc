#!/bin/sh
# https://golang.org/doc/install/source#environment
mkdir -p build && cd build

VERSION=1.0.0-alpha.6

export GOARCH=amd64

export GOOS=darwin

go build -o "abc-${VERSION}" -tags 'seabolt !oss' ./../cmd/abc/...
zip -r "abc-${GOOS}-${VERSION}.zip" "abc-${VERSION}"

export GOOS=windows

go build -o "abc-${VERSION}.exe"  -tags 'seabolt !oss' ./../cmd/abc/...
zip -r "abc-${GOOS}-${VERSION}.zip" "abc-${VERSION}.exe"

export GOOS=linux

rm "abc-${VERSION}"
go build -o "abc-${VERSION}" -tags 'seabolt !oss' ./../cmd/abc/...
zip -r "abc-${GOOS}-${VERSION}.zip" "abc-${VERSION}"
