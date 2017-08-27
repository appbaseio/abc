#!/bin/sh
# https://golang.org/doc/install/source#environment

VERSION=0.3.0

export GOARCH=amd64

export GOOS=darwin

go build -o "build/abc-${GOOS}-${VERSION}" -tags '!oss' ./cmd/abc/...

export GOOS=windows

go build -o "build/abc-${GOOS}-${VERSION}.exe"  -tags '!oss' ./cmd/abc/...

export GOOS=linux

go build -o "build/abc-${GOOS}-${VERSION}" -tags '!oss' ./cmd/abc/...
