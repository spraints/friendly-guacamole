#!/bin/bash

set -e

prereq() {
  echo "==> $*" 1>&2
  "$@" 1>&2
}

prereq go get github.com/golang/protobuf/protoc-gen-go
prereq protoc p/svc.proto --go_out=plugins=grpc:.

go run server.go &
trap "kill -TERM $!" EXIT
sleep 1
go run client.go "$@"
