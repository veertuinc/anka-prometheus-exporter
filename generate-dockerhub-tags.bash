#!/bin/bash
set -exo pipefail
docker login -u veertuserviceuser
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
DOCKERFILE_PATH="$SCRIPT_DIR/docker/scratch"
NAME="anka-prometheus-exporter"
cleanup() {
  popd
  rm -f "${DOCKERFILE_PATH}/${NAME}"*
  rm -f ./bin/"${NAME}"*
}
ARCH=amd64 make build-linux
ARCH=arm64 make build-linux
cp -f ./bin/"${NAME}"_linux* "${DOCKERFILE_PATH}/"
ls -alht "${DOCKERFILE_PATH}/"
trap cleanup EXIT
pushd "${DOCKERFILE_PATH}"

docker buildx create --name buildx-multi-arch --use || true
# make sure docker login is handled
docker buildx build --builder buildx-multi-arch --no-cache --platform linux/amd64,linux/arm64 -t veertu/"${NAME}":latest -t veertu/"$NAME:v$(cat "${SCRIPT_DIR}"/VERSION)" --push .