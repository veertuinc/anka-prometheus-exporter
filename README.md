Created and maintained by Veertu Inc.

License: [MIT](https://choosealicense.com/licenses/mit/)

# General Purpose
This package retrieves data from your [Anka Cloud](https://veertu.com) and creates an endpoint for Prometheus to monitor

# Running with Binary
1. Download the appropriate binary (anka_prometheus_XXXXX) from the releases page
2. Execute: ./anka_prometheus_XXXXX --controller_address ENTER-ADDRESS-HERE

*Check out the flags section for additional options* 

##### Building the binary yourself
Prerequisite: Go >= 1.11

1. Clone this repository
2. Execute build.sh
3. Use the binary generated in bin folder

##### Flags
Flag | Value Type | Default | Mandatory | Comments
---- | ---------- | ------- | --------- | --------
controller_address | string | | Yes | Controller's Address
port | integer | 2112 | No | Port to use
disable_interval_optimizer | Boolean | False | No | Disable automatic interval optimizer
interval | Integer | 15 | No | Interval between Controller requests (Seconds)
tls | Boolean | False | No | Enable TLS
skip_tls_verification | Boolean | False | No | Skips certificate verification (for self signed certs)
ca_cert | String |  | No | Path to CA certificate PEM/x509 file
client_cert | String |  | No | Path to client certificate PEM/x509 file
client_cert_key | String |  | No | Path to client key PEM/x509 file

# Running with Docker Compose
1. Add your controller's address to docker-compose.yml (MANDATORY)
2. Edit docker-compose.yml for other configuration options
2. Run docker-compose up -d

# Integrating Prometheus
Point prometheus to the machine running this package. Default port is 2112 (/metrics endpoint is used)

# Using TLS Options
--tls is not required if controller's certificate is valid and no client authentication is configured
For all other TLS configuration options, --tls must be set

For self signed certificates, you can either use --skip_tls_verification or provide your ca cert with --ca_cert
If client authentication is set on the controller, use --client_cert and --client_cert_key

# Development

- metrics/metric_*.go