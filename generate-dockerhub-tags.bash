#!/bin/bash
# Manual fallback: rebuild/push Docker Hub images from the latest GitHub Release.
# Normal releases do this in .github/workflows/release.yml on v* tags.
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
# Tag is vX.Y.Z; GoReleaser .Version strips the leading v.
VERSION="${RELEASE_TAG#v}"
CHECKSUMS="${NAME}_v${VERSION}_SHA256SUMS"

gh release download "$RELEASE_TAG" --repo "$REPO" --dir "$DOWNLOAD_DIR" --clobber \
  --pattern "${CHECKSUMS}"

for arch in amd64 arm64; do
  asset="${NAME}_v${VERSION}_linux_${arch}"
  gh release download "$RELEASE_TAG" --repo "$REPO" --dir "$DOWNLOAD_DIR" --clobber \
    --pattern "${asset}.zip"
  if command -v shasum >/dev/null 2>&1; then
    ( cd "$DOWNLOAD_DIR" && shasum -a 256 -c <(grep " ${asset}.zip\$" "${CHECKSUMS}") )
  else
    ( cd "$DOWNLOAD_DIR" && sha256sum -c <(grep " ${asset}.zip\$" "${CHECKSUMS}") )
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
docker buildx build --builder buildx-multi-arch --no-cache --platform linux/amd64,linux/arm64 \
  -t veertu/"${NAME}":latest \
  -t veertu/"${NAME}:${RELEASE_TAG}" \
  --push .
