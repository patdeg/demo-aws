#!/usr/bin/env bash
# Stops the process if something fails
set -xe

# All of the dependencies needed/fetched for the project.
go get -d ./...

# create the application binary that eb uses
GOOS=linux GOARCH=amd64 go build -o bin/application -ldflags="-s -w"

