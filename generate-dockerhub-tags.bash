#!/bin/bash
set -exo pipefail
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
DOCKERFILE_PATH="$SCRIPT_DIR/docker/scratch"
NAME="anka-prometheus-exporter"
FULL_BINARY_NAME="${NAME}_linux_amd64"
cleanup() {
  rm -f ${DOCKERFILE_PATH}/${NAME}*
}
LINUX_BINARY="$SCRIPT_DIR/bin/${FULL_BINARY_NAME}"
rm -f $LINUX_BINARY
make build-linux
cp -f $LINUX_BINARY $DOCKERFILE_PATH/
trap cleanup EXIT
cd $DOCKERFILE_PATH
docker buildx build --platform linux/amd64,linux/arm64,linux/386 -t veertu/$NAME:latest -t veertu/$NAME:v$(cat $SCRIPT_DIR/VERSION) --push .