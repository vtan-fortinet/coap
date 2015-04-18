#! /bin/bash

GOPATH=$(pwd)
ln -s ../../ src

go run demo1.go

rm -fr src
