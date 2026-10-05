# Releasing

1. Create release/vX.X.X branch in both repos and update VERSION file.
2. Create PR merging into master and ensure CI actions pass.
3. Run `https://INTERNAL_JENKINS_URL/job/prometheus-exporter/` to ensure it builds and all tests pass.
4. Download the artifacts from the job.
5. Update ~/prometheus-exporter, merging release branch into master. Make sure it targets the proper and updated master branch post-merge for the submodule too.
6. Do release in github and attach artifacts.
7. `./anka-prometheus-exporter/generate-dockerhub-tags.bash`
