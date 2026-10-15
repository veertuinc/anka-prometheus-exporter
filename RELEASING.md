# Releasing

1. Create release/vX.X.X branch in both repos and update VERSION file.
2. Make your changes in the release branch.
3. Create PR merging into master and ensure CI actions pass.
4. Run `https://INTERNAL_JENKINS_URL/job/prometheus-exporter/` to ensure it builds and all tests pass.
5. Download the artifacts from the job.
6. Update ~/prometheus-exporter, merging release branch into master. Make sure it targets the proper and updated master branch post-merge for the submodule too.
7. Do release in github and attach artifacts.
8. `./anka-prometheus-exporter/generate-dockerhub-tags.bash`
