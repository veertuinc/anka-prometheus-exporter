FROM alpine:latest as certs
RUN apk --update add ca-certificates
FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ./anka-prometheus-exporter_linux_amd64 /anka-prometheus-exporter
ENTRYPOINT ["/anka-prometheus-exporter"]