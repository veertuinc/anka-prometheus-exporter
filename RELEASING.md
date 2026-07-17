# Releasing

1. Run `https://INTERNAL_JENKINS_URL/job/prometheus-exporter/` to ensure it builds and all tests pass.
2. Download the artifacts from the job.
3. Do release in github and attach artifacts.
4. `./anka-prometheus-exporter/generate-dockerhub-tags.bash`