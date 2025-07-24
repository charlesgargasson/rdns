#!/bin/bash
cd -- "$(dirname -- "$0")"
go version
go get main
set -x
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildmode=pie -o "bin/rdns" -a -ldflags '-s -w -extldflags "--static-pie"' -buildvcs=false -tags 'osusergo,netgo,static_build' &
env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 CC=x86_64-w64-mingw32-gcc go build -buildmode=pie -o "bin/rdns.exe" -a -ldflags '-s -w -extldflags "--static-pie"' -buildvcs=false -tags 'osusergo,netgo,static_build' &
wait
