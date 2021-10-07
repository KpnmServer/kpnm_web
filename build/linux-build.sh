#!/bin/sh

cd $(dirname $0)

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build\
 -trimpath -gcflags "-trimpath=${GOPATH}" -asmflags "-trimpath=${GOPATH}" -ldflags "-w -s"\
 -o "./output/linux-amd64-binary" "${PWD}/../src/main.go"
