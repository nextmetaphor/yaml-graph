#!/usr/bin/env bash

if [ "$#" -gt 0 ]; then
    export GOOS=$1
fi
if [ "$#" -gt 1 ]; then
    export GOARCH=$2
fi

# install utilities to build a static binary
apt-get update -y && apt-get upgrade -y && apt-get install musl-tools -y

printf "### go fmt ###\n"
go fmt ./...

printf "\n### go vet ###\n"
go vet ./...

printf "\n### go build ###\n"
CGO_ENABLED=1 CC=musl-gcc go build --ldflags '-linkmode=external -extldflags=-static' -o yaml-graph