#!/bin/bash
set -e

mkdir bin &> /dev/null
ver=$(cat VERSION)
cd src/
GOOS=darwin GOARCH=amd64 go build -o ../bin/anka_prometheus_mac-$ver main.go
GOOS=linux GOARCH=amd64 go build -o ../bin/anka_prometheus_linux-$ver main.go
