FROM golang:1.15.2-alpine as builder
COPY . /build
WORKDIR /build
RUN apk add --update make
RUN make build-linux
######################
FROM alpine
LABEL maintainer="Veertu Inc. support@veertu.com"
COPY --from=builder /build/bin/anka-prometheus-exporter_linux_amd64 /usr/bin/anka-prometheus-exporter
ENTRYPOINT ["/bin/sh", "-c", "/usr/bin/anka-prometheus-exporter --controller-address $CONTROLLER_ADDR $DISABLE_INTERVAL_OPTIMIZER $ANKA_PROMETHEUS_PORT $ANKA_PROMETHEUS_INTERVAL $USE_TLS $SKIP_TLS $CA_CERT $CLIENT_CERT $CLIENT_KEY"]
