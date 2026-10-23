# Releasing

Primary branch: **`edge`** (GitHub app repo and Bitbucket wrapper repo).

There is no `VERSION` file. The release number comes from the Jenkins **Version** parameter. After a green private suite on `edge`, Jenkins pushes git tag `vX.Y.Z`, which starts the GitHub Release + Docker Hub publish.

| Where | What it does |
|---|---|
| GitHub Actions | Unit tests, GoReleaser snapshot on PRs; on `v*` tags: GitHub Release + Docker Hub |
| Jenkins (`prometheus-exporter` / `edge`) | Full private suite. You set **Version** when you start the job. On success it tags GitHub and triggers publish. |

## One-time setup

1. GitHub: default branch = `edge`.
2. Bitbucket wrapper: branch `edge` (multibranch job discovers it). Keep the `anka-prometheus-exporter` submodule pointed at GitHub `edge`.
3. Jenkins: GitHub token credential with permission to push tags to this repo.
4. GitHub secrets: `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`.

## Day to day

1. Open PRs into GitHub **`edge`**. Public CI runs unit tests + GoReleaser snapshot.
2. Merge to GitHub `edge`.
3. On the Bitbucket wrapper `edge`: update the submodule to the GitHub `edge` SHA, commit, push.

## Cut a release

1. Run Jenkins **`prometheus-exporter/edge`**.
2. Set parameter **Version** to the version you will ship (example: `4.3.6`, no `v` prefix).
3. Wait for green. Jenkins then:
   - Tags the submodule commit as `vX.Y.Z`
   - Pushes that tag to GitHub
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
- Tag push runs only on branch **`edge`** after a successful suite (not on UNSTABLE/FAILURE).
- Public GitHub Actions cannot replace Jenkins: the private suite needs internal hosts, licenses, and vault material.
- Local builds: `make VERSION=4.3.6 go.build` (or `go.releaser`).
