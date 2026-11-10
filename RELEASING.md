# Releasing

Primary branch: **`edge`**.

There is no `VERSION` file. The release number comes from the Jenkins **Version** parameter. After a green private suite on `edge`, Jenkins pushes a git tag, which starts the GitHub Release + Docker Hub publish.

The Jenkins pipeline and the private suite live in private GitHub repos. Those repos are moving from Bitbucket. Use the GitHub repos.

| Private repo | Role | Branch | Previous location |
|---|---|---|---|
| [`veertuinc/prometheus-exporter`](https://github.com/veertuinc/prometheus-exporter) | Jenkins pipeline | `edge` | `https://bitbucket.org/veertunew/prometheus-exporter` |
| [`veertuinc/prometheus-exporter-tests`](https://github.com/veertuinc/prometheus-exporter-tests) | Private test suite | `master` | `https://bitbucket.org/veertunew/prometheus-exporter-tests` |

This public repo has no Jenkinsfile. You do not bump a submodule SHA for each release. The Jenkins job checks out the current public GitHub `origin/edge` at job start. If either private repo is missing, stop and notify. Do not use Bitbucket.

| Where | What it does |
|---|---|
| GitHub Actions | Pull requests: `Check PR / goreleaser` only. Push to `edge`: `Tests / Run Tests`. Tag `v*`: workflow `Release` |
| Jenkins [`prometheus-exporter`](https://jenkins.veertu.com/job/prometheus-exporter/) | Runs the private suite. On success, tags the public `edge` commit. FAILURE or UNSTABLE does not push the tag. |

## Automation

A security patch is the latest GitHub shipping tag plus one patch, without the `v`. Do not use `Version=edge` for that release. On 2026-11-09 the latest shipping tag is `v4.3.5`, so **Version** is `4.3.6`.

1. `Tests / Run Tests` must be green on the `edge` commit. If that check failed, stop. Do not start Jenkins.
2. If `veertuinc/prometheus-exporter` or `veertuinc/prometheus-exporter-tests` is missing, stop and notify. Do not use Bitbucket. Do not push the tag.
3. POST Jenkins job `prometheus-exporter` at `https://jenkins.veertu.com/job/prometheus-exporter/`. The controller is anonymous. A POST needs the `Jenkins-Crumb` header and the session cookie from `GET https://jenkins.veertu.com/crumbIssuer/api/json`.
   - Parameter **Version** = `X.Y.Z` (no `v`).
   - Parameter **OverwriteExistingRelease** = `false`.
   - Leave **AnkaLicense** and the credential parameters at their job defaults.
4. Wait for the build. The job runs the private suite. On success it pushes tag `vX.Y.Z`. FAILURE or UNSTABLE does not push the tag. Stop and notify. Do not push the tag yourself.
5. Tag `v*` starts GitHub workflow `Release`. The workflow needs repo secrets `DOCKERHUB_USERNAME` and `DOCKERHUB_TOKEN`. Do not paste them. For a shipping tag it publishes the GitHub Release and Docker tags `veertu/anka-prometheus-exporter:latest` and `veertu/anka-prometheus-exporter:vX.Y.Z`.
6. Confirm the [GitHub Release](https://github.com/veertuinc/anka-prometheus-exporter/releases) and both Docker tags, then post the links.
7. Do not run `./generate-dockerhub-tags.bash`.
8. Do not refresh the edge channel as part of a security patch.

The private suite needs internal hosts, licenses, and vault material. Public GitHub Actions cannot replace the Jenkins job.

## Cut a shipping release

1. Run Jenkins [`prometheus-exporter`](https://jenkins.veertu.com/job/prometheus-exporter/) (**Build with Parameters**).
2. Set parameter **Version** to the version you will ship (example: `4.3.6`, no `v` prefix).
3. Leave **OverwriteExistingRelease** off unless you intend to replace that tag.
4. Jenkins:
   - Syncs `anka-prometheus-exporter` to latest GitHub `edge`
   - Runs the private suite against that commit
   - **Only if the private suite passes:** tags that commit as `vX.Y.Z` and pushes the tag to GitHub (skipped on UNSTABLE/FAILURE)
5. The GitHub **Release** workflow on `v*` tags:
   - Takes the version from the tag
   - Runs GoReleaser in replace mode (GitHub Release + zip archives + checksums)
   - Pushes `veertu/anka-prometheus-exporter:latest` and `veertu/anka-prometheus-exporter:vX.Y.Z`
6. Confirm the [GitHub Release](https://github.com/veertuinc/anka-prometheus-exporter/releases) and Docker Hub tags.

## Cut / refresh the edge release

The floating edge channel uses semver tag **`v0.0.0-edge`** (prerelease) so GoReleaser and GitHub stay happy.

1. Run the same Jenkins job.
2. Set **Version** to `edge` (or `0.0.0-edge`).
3. Overwrite is always on for this Version: Jenkins deletes the existing `v0.0.0-edge` GitHub Release and tag, then recreates them on the tested commit.
4. Actions publishes the GitHub Release (prerelease) and Docker tags:
   - `veertu/anka-prometheus-exporter:edge`
   - `veertu/anka-prometheus-exporter:v0.0.0-edge`
   - Does **not** move `:latest`

## Manual Docker republish (optional)

The Automation does not run this.

```bash
./generate-dockerhub-tags.bash
```

## Notes

- No `release/vX.X.X` branches for normal releases.
- Shipping tag format: `vMAJOR.MINOR.PATCH`.
- Edge channel tag: `v0.0.0-edge` (Jenkins `Version=edge`).
- Tag push runs after a successful suite (not on UNSTABLE/FAILURE).
- **OverwriteExistingRelease** deletes the existing GitHub Release and tag for that Version, then recreates the tag on the tested commit so Actions republishes artifacts. Always enabled for `edge` / `0.0.0-edge`. Leave it off (`false`) for a security patch.
- Force sync locally: `SYNC_EXPORTER_TO_EDGE=true` or `BRANCH_NAME=edge`. These are not parameters of job `prometheus-exporter`.
- Public GitHub Actions cannot replace Jenkins: the private suite needs internal hosts, licenses, and vault material.
- Local builds: `make VERSION=4.3.6 go.build` (or `go.releaser`). For edge: `make VERSION=0.0.0-edge go.releaser`.
