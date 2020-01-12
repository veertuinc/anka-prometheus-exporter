#!/bin/bash
mkdir bin &> /dev/null
cd src/

GOOS=darwin GOARCH=amd64 go build -o ../bin/anka_prometheus_mac main.go
GOOS=linux GOARCH=amd64 go build -o ../bin/anka_prometheus_linux main.go
