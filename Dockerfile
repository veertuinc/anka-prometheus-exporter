#Name of container: docker-anka-controller

# Pull base image
FROM centos:7
MAINTAINER Veertu Inc. "support@veertu.com"

COPY bin/anka_prometheus_linux /usr/bin/anka-prometheus
ENTRYPOINT ["/bin/bash", "-c", "anka-prometheus --controller_address $CONTROLLER_ADDR $DISABLE_INTERVAL_OPTIMIZER $ANKA_PROMETHEUS_PORT $ANKA_PROMETHEUS_INTERVAL"]