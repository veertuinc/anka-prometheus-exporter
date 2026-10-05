# Releasing

1. Create release/vX.X.X branch in both repos and update VERSION file.
2. Create PR merging into master and ensure CI actions pass.
3. Run `https://INTERNAL_JENKINS_URL/job/prometheus-exporter/` to ensure it builds and all tests pass.
4. Download the artifacts from the job.
5. Do release in github and attach artifacts.
6. `./anka-prometheus-exporter/generate-dockerhub-tags.bash`