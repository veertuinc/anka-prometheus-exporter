#!/bin/bash
set -exo pipefail
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
NAME="anka-prometheus-exporter"
DOCKERFILE_PATH="$SCRIPT_DIR/docker/scratch"
LINUX_BINARY="$SCRIPT_DIR/bin/${NAME}_linux_amd64"
[[ ! -f $LINUX_BINARY ]] && make build-linux
cp -f $LINUX_BINARY $DOCKERFILE_PATH/$NAME
cd $DOCKERFILE_PATH
# docker build --no-cache -t $NAME:latest .
docker buildx build --platform linux/amd64,linux/arm64,linux/386 -t veertu/$NAME:latest -t veertu/$NAME:v$(cat VERSION) --push .
rm -f "$DOCKERFILE_PATH/$NAME"