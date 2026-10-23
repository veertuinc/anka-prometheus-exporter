# Releasing

Primary branch: **`edge`** (GitHub app repo and Bitbucket wrapper repo).

There is no `VERSION` file. The release number comes from the Jenkins **Version** parameter. After a green private suite on `edge`, Jenkins pushes git tag `vX.Y.Z`, which starts the GitHub Release + Docker Hub publish.

| Where | What it does |
|---|---|
| GitHub Actions | Unit tests, GoReleaser snapshot on PRs; on `v*` tags: GitHub Release + Docker Hub |
| Jenkins [`prometheus-exporter`](https://jenkins.veertu.com/view/prometheus-exporter/job/prometheus-exporter/) | Single-branch Pipeline on Bitbucket wrapper `edge`. Syncs exporter to GitHub `edge`, runs the private suite, tags GitHub on success |

You do **not** need to bump the Bitbucket submodule SHA for each release. Jenkins checks out current GitHub `origin/edge` at job start.

## Day to day

1. Open PRs into GitHub **`edge`**. Public CI runs unit tests + GoReleaser snapshot.
2. Merge to GitHub `edge`.

## Cut a release

1. Run Jenkins [`prometheus-exporter`](https://jenkins.veertu.com/view/prometheus-exporter/job/prometheus-exporter/) (**Build with Parameters**).
2. Set parameter **Version** to the version you will ship (example: `4.3.6`, no `v` prefix).
3. Jenkins:
   - Syncs `anka-prometheus-exporter` to latest GitHub `edge`
   - Runs the private suite against that commit
   - **Only if the private suite passes:** tags that commit as `vX.Y.Z` and pushes the tag to GitHub (skipped on UNSTABLE/FAILURE)
4. The GitHub **Release** workflow on `v*` tags:
   - Takes the version from the tag
   - Runs GoReleaser (GitHub Release + zip archives + checksums)
   - Pushes `veertu/anka-prometheus-exporter:latest` and `veertu/anka-prometheus-exporter:vX.Y.Z`
5. Confirm the [GitHub Release](https://github.com/veertuinc/anka-prometheus-exporter/releases) and Docker Hub tags.

## Manual Docker republish (optional)

```bash
./generate-dockerhub-tags.bash
```

## Notes

- No `release/vX.X.X` branches for normal releases.
- Tag format: `vMAJOR.MINOR.PATCH`.
- Tag push runs after a successful suite (not on UNSTABLE/FAILURE).
- Force sync locally: `SYNC_EXPORTER_TO_EDGE=true` or `BRANCH_NAME=edge`.
- Public GitHub Actions cannot replace Jenkins: the private suite needs internal hosts, licenses, and vault material.
- Local builds: `make VERSION=4.3.6 go.build` (or `go.releaser`).
