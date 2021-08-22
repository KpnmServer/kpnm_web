#!/bin/sh

cd $(dirname $0)

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build\
 -trimpath -gcflags "-trimpath=${GOPATH}" -asmflags "-trimpath=${GOPATH}" -ldflags "-w -s"\
 -o "./output/windows-amd64-binary.exe" "${PWD}/../src/main.go"
