#!/bin/bash

set -e

prereq() {
  echo "==> $*" 1>&2
  "$@" 1>&2
}

prereq go get github.com/golang/protobuf/protoc-gen-go
prereq make

./server &
trap "kill -TERM $!" EXIT
sleep 1
./client "$@"
