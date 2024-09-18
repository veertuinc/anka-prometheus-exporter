FROM golang:1.15.2-alpine as builder
COPY . /build
WORKDIR /build
RUN apk add --update make
RUN make build-linux
######################
FROM alpine
LABEL maintainer="Veertu Inc. support@veertu.com"
ARG TARGETARCH
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ./anka-prometheus-exporter_linux_${TARGETARCH} /anka-prometheus-exporter
ENTRYPOINT ["/anka-prometheus-exporter"]