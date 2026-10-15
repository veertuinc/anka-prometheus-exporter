#!/bin/bash
set -exo pipefail
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
DOCKERFILE_PATH="$SCRIPT_DIR/docker/scratch"
NAME="anka-prometheus-exporter"
REPO="veertuinc/anka-prometheus-exporter"
DOWNLOAD_DIR="$(mktemp -d)"
DID_PUSHD=0
cleanup() {
  status=$?
  rm -rf "$DOWNLOAD_DIR"
  rm -f "${DOCKERFILE_PATH}/${NAME}_linux_"*
  if [[ "$DID_PUSHD" == 1 ]]; then
    popd
  fi
  exit "$status"
}
trap cleanup EXIT

git -C "$SCRIPT_DIR" fetch origin --tags --force
RELEASE_TAG="$(gh release view --repo "$REPO" --json tagName --jq .tagName)"

for arch in amd64 arm64; do
  asset="${NAME}_${RELEASE_TAG}_linux_${arch}"
  gh release download "$RELEASE_TAG" --repo "$REPO" --dir "$DOWNLOAD_DIR" --clobber \
    --pattern "${asset}.zip" \
    --pattern "${asset}_SHA256SUMS"
  if command -v shasum >/dev/null 2>&1; then
    ( cd "$DOWNLOAD_DIR" && shasum -a 256 -c "${asset}_SHA256SUMS" )
  else
    ( cd "$DOWNLOAD_DIR" && sha256sum -c "${asset}_SHA256SUMS" )
  fi
  unzip -o "$DOWNLOAD_DIR/${asset}.zip" -d "$DOWNLOAD_DIR"
  cp -f "$DOWNLOAD_DIR/${asset}" "${DOCKERFILE_PATH}/${NAME}_linux_${arch}"
  chmod +x "${DOCKERFILE_PATH}/${NAME}_linux_${arch}"
done
ls -alht "${DOCKERFILE_PATH}/"

docker login -u veertuserviceuser
docker buildx create --name buildx-multi-arch --use || true
pushd "${DOCKERFILE_PATH}"
DID_PUSHD=1
docker buildx build --builder buildx-multi-arch --no-cache --platform linux/amd64,linux/arm64 -t veertu/"${NAME}":latest -t veertu/"${NAME}:${RELEASE_TAG}" --push .
