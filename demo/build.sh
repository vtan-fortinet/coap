#! /bin/bash

GOPATH=$(pwd)
ln -s ../../ src

go build demo1.go

rm -fr src
