#!/bin/bash
cd -- "$(dirname -- "$0")"
VCS="-buildvcs=false"
go version
go get main
set -x
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -buildmode=pie -o bin/rdns $VCS &
env GOOS=windows GOARCH=amd64 CGO_ENABLED=0 CC=x86_64-w64-mingw32-gcc go build -buildmode=pie -o "bin/rdns.exe" $VCS &
wait
