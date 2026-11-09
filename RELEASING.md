# Releasing

Primary branch: **`edge`**.

There is no `VERSION` file. The release number is the Jenkins parameter **Version** (`X.Y.Z`, no `v`). After a green private suite, Jenkins pushes git tag `vX.Y.Z`. That tag starts the GitHub Release and the Docker Hub publish.

The Jenkins pipeline and the private suite live in private GitHub repos. Those repos are moving from Bitbucket. Use the GitHub repos.

| Private repo | Role | Branch | Previous location |
|---|---|---|---|
| [`veertuinc/prometheus-exporter`](https://github.com/veertuinc/prometheus-exporter) | Jenkins pipeline | `edge` | `https://bitbucket.org/veertunew/prometheus-exporter` |
| [`veertuinc/prometheus-exporter-tests`](https://github.com/veertuinc/prometheus-exporter-tests) | Private test suite | `master` | `https://bitbucket.org/veertunew/prometheus-exporter-tests` |

This public repo has no Jenkinsfile. The Automation does not clone the private repos to cut a release, and it does not bump a submodule SHA. The Jenkins job checks out the current public `origin/edge` when it starts.

| Where | What it does |
|---|---|
| GitHub Actions | Pull requests: `Check PR / goreleaser` only. Push to `edge`: `Tests / Run Tests`. Tag `v*`: workflow `Release` |
| Jenkins [`prometheus-exporter`](https://jenkins.veertu.com/view/prometheus-exporter/job/prometheus-exporter/) | Runs the private suite. On success, tags the public `edge` commit `vX.Y.Z`. FAILURE or UNSTABLE does not push the tag. |

## Automation

Next security version is the latest GitHub release tag plus one patch, without the `v`. On 2026-11-09 the latest tag is `v4.3.5`, so parameter **Version** is `4.3.6`.

1. `Tests / Run Tests` must be green on the `edge` commit. If that check failed, stop. Do not start Jenkins.
2. If `veertuinc/prometheus-exporter` or `veertuinc/prometheus-exporter-tests` is missing, stop and notify. Do not use Bitbucket. Do not push the tag.
3. Trigger Jenkins job `prometheus-exporter`.
   - URL: `https://jenkins.veertu.com/view/prometheus-exporter/job/prometheus-exporter/`
   - Parameter **Version** = `X.Y.Z` (no `v`). Pass only this parameter.
4. Wait for the build. The job runs the private suite. On success it pushes tag `vX.Y.Z` to this repo. FAILURE or UNSTABLE does not push the tag. Stop and notify. Do not push the tag yourself.
5. Tag `v*` starts GitHub workflow `Release`. The workflow needs repo secrets `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN`. Do not paste them. The workflow publishes the GitHub Release and Docker tags `veertu/anka-prometheus-exporter:latest` and `veertu/anka-prometheus-exporter:vX.Y.Z`.
6. Confirm the [GitHub Release](https://github.com/veertuinc/anka-prometheus-exporter/releases) and both Docker tags, then post the links.
7. Do not run `./generate-dockerhub-tags.bash`.

The private suite needs internal hosts, licenses, and vault material. Public GitHub Actions cannot replace the Jenkins job.

## Manual Docker republish

People use this only to republish images. The Automation does not run it.

```bash
./generate-dockerhub-tags.bash
```

## Notes

- No `release/vX.X.X` branches for normal releases.
- Tag format: `vMAJOR.MINOR.PATCH`.
- The tag push runs after a successful suite. It does not run on UNSTABLE or FAILURE.
- Local force-sync switches inside the private pipeline: `SYNC_EXPORTER_TO_EDGE=true` or `BRANCH_NAME=edge`. These are not parameters of job `prometheus-exporter`.
- Local builds: `make VERSION=4.3.6 go.build` (or `go.releaser`).
