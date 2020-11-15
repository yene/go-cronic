#!/bin/bash
env GOOS=linux GOARCH=amd64 go build -o cronic-linux
upx cronic-linux
env GOOS=darwin GOARCH=amd64 go build -o cronic-macos
upx cronic-macos
