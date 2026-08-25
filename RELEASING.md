# Releasing

1. Create PR merging into master and ensure CI actions pass.
2. Run `https://INTERNAL_JENKINS_URL/job/prometheus-exporter/` to ensure it builds and all tests pass.
3. Download the artifacts from the job.
4. Do release in github and attach artifacts.
5. `./anka-prometheus-exporter/generate-dockerhub-tags.bash`