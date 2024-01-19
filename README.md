# Anka Prometheus Exporter

> Created and maintained by Veertu Inc.

This package retrieves data from your [Anka Build Cloud](https://veertu.com) and exports it for Prometheus at http://localhost/metrics.

## Usage

Ensure you have a functioning Prometheus instance before using this.

### Available ENVs

> ENVs take priority over flags

| ENV | Flag |
| --- | ---  |
| ANKA_PROMETHEUS_EXPORTER_CONTROLLER_ADDRESS (string) | --controller-address (string) |
| ANKA_PROMETHEUS_EXPORTER_INTERVAL (int) | --interval (int) |
| ANKA_PROMETHEUS_EXPORTER_PORT (int) | --port (int) |
| ANKA_PROMETHEUS_EXPORTER_DISABLE_INTERVAL_OPTIMIZER (bool) | --disable-interval-optimizer |
| ANKA_PROMETHEUS_EXPORTER_TLS (bool) | --tls |
| ANKA_PROMETHEUS_EXPORTER_SKIP_TLS_VERIFICATION (bool) | --skip-tls-verification |
| ANKA_PROMETHEUS_EXPORTER_CA_CERT (string) | --ca-cert (string) |
| ANKA_PROMETHEUS_EXPORTER_CLIENT_CERT (string) | --client-cert (string) |
| ANKA_PROMETHEUS_EXPORTER_CLIENT_CERT_KEY (string) | --client-cert-key (string) |
| ANKA_PROMETHEUS_EXPORTER_UAK_ID (string) | --uak-id (string) |
| ANKA_PROMETHEUS_EXPORTER_UAK_PATH (string) | --uak-path (string) |
| ANKA_PROMETHEUS_EXPORTER_UAK_STRING (string) | --uak-string (string) |

```bash
Usage of anka-prometheus-exporter:
  -ca-cert string
        Path to ca PEM/x509 file (cert file path as arg)
  -client-cert string
        Path to client cert PEM/x509 file (cert file path as arg)
  -client-cert-key string
        Path to client key PEM/x509 file (cert file path as arg)
  -controller-address string
        Controller address to monitor (url as arg) (required)
  -controller-username string
        Controller username with basic root token (username as arg)
  -controller-password string
        Controller password with basic root token (password as arg) 
  -disable-interval-optimizer
        Optimize interval according to /metric api requests receieved (no args)
  -interval int
        Seconds to wait between data requests to controller (int as arg) (default 15)
  -port int
        Port to server /metrics endpoint (int as arg) (default 2112)
  -skip-tls-verification
        Skip TLS verification (no args)
  -tls
        Enable TLS (no args)
  -uak-id string
        UAK ID you wish to use for Controller requests (string as arg)
  -uak-string string
        String form (cat myUAK.pem | sed '1,1d' | sed '$d' | tr -d '\n') of the key file contents for Controller requests (string as arg)
  -uak-path string
        Path to the UAK for Controller requests (path as arg) (takes priority over -uak-string if both are specified)
```

> `LOG_LEVEL` can be set using an environment variable

### Using the Binary

1. [Download the appropriate binary (anka-prometheus-exporter) from the releases page](https://github.com/veertuinc/anka-prometheus-exporter/releases)
2. Execute: `./anka-prometheus-exporter --controller_address http://{controller IP or URL}`

### Using Docker

1. Create a `docker-compose.yml`:
```
version: '3.1'
services:
  anka-prometheus-exporter:
    image: veertu/anka-prometheus-exporter:latest
    container_name: anka-prometheus-exporter
    ports:
      - "2112:2112"
    environment:
       - ANKA_PROMETHEUS_EXPORTER_CONTROLLER_ADDRESS # Defaults to using what is under the user executing docker-compose up; you can specify ="http://Your Controller URL and Port Here" if not set in user's env
```
2. `docker-compose pull && docker-compose up --remove-orphans -d`

---

## Adding a Prometheus target

Once running, add the scrape endpoint to your prometheus.yml:

> `host.docker.internal` is only needed if running on the same host as your prometheus container and on mac. For linux, use `172.17.0.1`.

```yaml
scrape_configs:
. . .
  - job_name: 'anka build cloud'
    static_configs:
      - targets: ['host.docker.internal:2112']
```

## Using TLS

The `--tls` flag is not required if your controller certificate is valid and no client authentication is configured.
For all other TLS configuration options, `--tls` must be set.

For self signed certificates, you can either use `--skip-tls-verification` or provide your ca-cert with `--ca-cert`.

### Using Auth

- If client authentication is set on the controller, use `--client-cert` and `--client-cert-key`.
- If RTA is enabled, you can use basic auth through:
    ```bash
      -controller-username string
            Controller username with basic root token (username as arg)
      -controller-password string
            Controller password with basic root token (password as arg) 
    ```

---

## Exposed Metrics

Metric name | Description
---- | ----------
anka_instance_state_count | Count of Instances in a particular State (label: state)
anka_instance_state_per_template_count | Count of Instances in a particular state, per Template (label: state, template_uuid, template_name)
anka_instance_state_per_group_count | Count of Instances in a particular state, per Group (label: state, group_name)
-- | --
anka_node_instance_count | Count of Instances running on the Node
anka_node_instance_capacity | Total Instance slots (capacity) on the Node
anka_node_states | Node state (1 = current state) (label: id, name, state)
anka_node_states_count | Count of Nodes in a particular state (label: state)
anka_node_disk_free_space | Amount of free disk space on the Node in Bytes
anka_node_disk_total_space | Amount of total available disk space on the Node in Bytes
anka_node_disk_anka_used_space | Amount of disk space used by Anka on the Node in Bytes
anka_node_cpu_core_count | Number of CPU Cores in Node
anka_node_cpu_util | CPU utilization in node
anka_node_ram_gb | Total RAM available for the Node in GB
anka_node_ram_util | Total RAM utilized for the Node
anka_node_used_virtual_cpu_count | Total Used Virtual CPU cores for the Node
anka_node_used_virtual_ram_mb | Total Used Virtual RAM for the Node in MB
-- | --
anka_node_group_nodes_count | Count of Nodes in a particular Group
anka_node_group_states_count | Count of Groups in a particular State (labels: group, state)
anka_node_group_instance_count | Count of Instances slots in use for the Group (and Nodes)
anka_node_group_disk_free_space | Amount of free disk space for the Group (and Nodes) in Bytes
anka_node_group_disk_total_space | Amount of total available disk space for the Group (and Nodes) in Bytes
anka_node_group_disk_anka_used_space | Amount of disk space used by Anka for the Group (and Nodes) in Bytes
anka_node_group_cpu_core_count | Number of CPU Cores for the Group (and Nodes)
anka_node_group_cpu_util | CPU utilization for the Group (and Nodes)
anka_node_group_ram_gb | Total RAM available for the Group (and Nodes) in GB
anka_node_group_ram_util | Total RAM utilized for the Group (and Nodes)
anka_node_group_used_virtual_cpu_count | Total Used Virtual CPU cores for the Group (and Nodes)
anka_node_group_used_virtual_ram_mb | Total Used Virtual RAM for the Group (and Nodes) in MB
anka_node_group_instance_capacity | Total Instance slots (capacity) for the Group (and Nodes)
-- | --
anka_nodes_count | Count of total Anka Nodes
anka_nodes_instance_count | Count of Instance slots in use across all Nodes
anka_nodes_instance_capacity | Count of total Instance Capacity across all Nodes
anka_nodes_disk_free_space | Amount of free disk space across all Nodes in Bytes
anka_nodes_disk_total_space | Amount of total available disk space across all Nodes in Bytes
anka_nodes_disk_anka_used_space | Amount of disk space used by Anka across all Nodes in Bytes
anka_nodes_cpu_core_count | Count of CPU Cores across all Nodes
anka_nodes_cpu_util | Total CPU utilization across all Nodes
anka_nodes_ram_gb | Total RAM available across all Nodes in GB
anka_nodes_ram_util | Total RAM utilized across all Nodes
anka_nodes_used_virtual_cpu_count | Total Used Virtual CPU cores across all Nodes
anka_nodes_used_virtual_ram_mb | Total Used Virtual RAM across all Nodes
-- | --
anka_registry_disk_total_space | Anka Build Cloud Registry total disk space
anka_registry_disk_free_space| Anka Build Cloud Registry free disk space
anka_registry_disk_used_space | Anka Build Cloud Registry used disk space
-- | --
anka_registry_template_count | Count of VM Templates in the Registry
anka_registry_template_disk_used | Total disk usage of the Template in the Registry
anka_registry_template_tag_disk_used | Total disk used by the Template's Tag in the Registry
anka_registry_template_tags_count | Count of Tags in the Registry for the Template

---

# Development

```bash
make build-and-run ARGUMENTS="--controller-username root --controller-password 1111111111"
```
