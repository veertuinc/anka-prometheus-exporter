Created and maintained by Veertu Inc.

License: [MIT](https://choosealicense.com/licenses/mit/)

# General Purpose
This package retrieves data from your [Anka Cloud](https://veertu.com) and creates an endpoint for Prometheus to monitor

# Running with Binary
1. Download the appropriate binary (anka_prometheus_XXXXX) from the releases page
2. Execute: ./anka_prometheus_XXXXX --controller_address ENTER-ADDRESS-HERE

*Check out the flags section for additional options* 

##### Building the binary yourself
1. Clone this repository
2. Execute build.sh
3. Use the binary generated in bin folder

##### Flags
Flag | Value Type | Default | Mandatory | Comments
---- | ---------- | ------- | --------- | --------
controller_address | string | | Yes | Controller's Address
port | integer | 2112 | No | Port to use
disable_interval_optimizer | Boolean | False | No | Disavles automatic interval optimizer
interval | Integer | 15 | No | Interval between Controller requests (Seconds)

# Running with Docker Compose
1. Edit and follow instructions on docker-compose.yml
2. Run docker-compose up -d

# Integrating Prometheus
Point prometheus to the machine running this package. Default port is 2112 (/metrics endpoint is used)
