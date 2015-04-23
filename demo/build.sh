#! /bin/bash

GOPATH=$(pwd)
ln -s ../../ src

go run demo1.go -a '"-b"'

rm -fr src
