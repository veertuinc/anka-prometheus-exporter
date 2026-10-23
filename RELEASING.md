# Releasing

Primary branch: **`edge`** (GitHub app repo and Bitbucket wrapper repo).

| Where | What it does |
|---|---|
| GitHub Actions | Unit tests, GoReleaser snapshot on PRs; on `v*` tags: GitHub Release + Docker Hub |
| Jenkins (`prometheus-exporter` / `edge`) | Full private suite (Buzyform, Anka, vault, signing). **Required before you tag.** |

Do not push a release tag until Jenkins is green on that commit.

## One-time setup

1. GitHub: default branch = `edge`.
2. Bitbucket wrapper: branch `edge` (multibranch job discovers it). Keep the `anka-prometheus-exporter` submodule pointed at GitHub `edge`.
3. GitHub secrets: `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`.

## Day to day

1. Open PRs into GitHub **`edge`**. Public CI runs unit tests + GoReleaser snapshot.
2. Merge to GitHub `edge`.
3. On the Bitbucket wrapper `edge`: update the submodule to the GitHub `edge` SHA, commit, push.
4. Run Jenkins job **`prometheus-exporter/edge`** (or the multibranch job for `edge`). Wait for a green private suite.

## Cut a release

1. On GitHub `edge`, set `VERSION` to the version you will ship (example: `4.3.6`). Commit and push.
2. Sync that commit into the Bitbucket wrapper `edge` submodule and push.
3. Run Jenkins **`prometheus-exporter/edge`**. Confirm green (artifacts + signed zips as today).
4. Only then, on the **same** GitHub commit, tag and push:

```bash
git checkout edge
git pull
VERSION="$(tr -d '[:space:]' < VERSION)"
git tag -a "v${VERSION}" -m "v${VERSION}"
git push origin "v${VERSION}"
```

5. The GitHub **Release** workflow on `v*` tags:
   - Checks tag matches `VERSION` (`v` + file contents)
   - Runs GoReleaser (GitHub Release + zip archives + checksums)
   - Pushes `veertu/anka-prometheus-exporter:latest` and `veertu/anka-prometheus-exporter:vX.Y.Z`

6. Confirm the [GitHub Release](https://github.com/veertuinc/anka-prometheus-exporter/releases) and Docker Hub tags.

## Manual Docker republish (optional)

If you need to rebuild images from an existing GitHub Release without re-tagging:

```bash
./generate-dockerhub-tags.bash
```

## Notes

- No `release/vX.X.X` branches for normal releases.
- Tag format: `vMAJOR.MINOR.PATCH`, must match `VERSION`.
- Public GitHub Actions cannot replace Jenkins: the private suite needs internal hosts, licenses, and vault material.
