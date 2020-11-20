#!/usr/bin/env bash

printf "\n### go test ###\n"
go test -coverpkg=./... -coverprofile=profile.cov ./...

printf "\n### go tool cover ###\n"
go tool cover -html=profile.cov