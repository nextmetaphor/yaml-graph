#!/usr/bin/env bash

printf "### go fmt ###\n"
go fmt ./...

printf "### go vet ###\n"
go vet ./...

printf "\n### golint ###\n"
golint ./...

printf "### go build ###\n"
go build -i -o yaml-graph