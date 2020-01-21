FROM golang:1.13.5-alpine as builder

COPY src /data
WORKDIR /data
RUN GOOS=linux GOARCH=amd64 go build -o /anka_prometheus

FROM alpine
LABEL maintainer="Veertu Inc. support@veertu.com"

COPY --from=builder /anka_prometheus /usr/bin/anka-prometheus
ENTRYPOINT ["/bin/sh", "-c", "anka-prometheus --controller_address $CONTROLLER_ADDR $DISABLE_INTERVAL_OPTIMIZER $ANKA_PROMETHEUS_PORT $ANKA_PROMETHEUS_INTERVAL $USE_TLS $SKIP_TLS $CA_CERT $CLIENT_CERT $CLIENT_KEY"]
