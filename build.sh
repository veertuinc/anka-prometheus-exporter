#!/bin/bash
mkdir bin &> /dev/null
GOPATH=$(pwd) GOOS=darwin GOARCH=amd64 go build -o bin/anka_prometheus_mac src/github.com/veertuinc/anka-prometheus/main.go
GOPATH=$(pwd) GOOS=linux GOARCH=amd64 go build -o bin/anka_prometheus_linux src/github.com/veertuinc/anka-prometheus/main.go