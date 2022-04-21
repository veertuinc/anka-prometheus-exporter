package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type NodesMetric struct {
	BaseAnkaMetric
	HandleData func([]types.Node, prometheus.Gauge)
}

func (m NodesMetric) GetEventHandler() func(interface{}) error {
	return func(nodesData interface{}) error {
		nodes, err := ConvertToNodeData(nodesData)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGauge(m.metric)
		if err != nil {
			return err
		}
		m.HandleData(
			nodes,
			metric,
		)
		return nil
	}
}

var ankaNodesMetrics = []NodesMetric{
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_count", "Count of total Anka Nodes"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			metric.Set(float64(len(nodes)))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_instance_count", "Count of Instance slots in use across all Nodes"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint = 0
			for _, node := range nodes { // For each node
				count += node.VMCount
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_instance_capacity", "Count of total Instance Capacity across all Nodes"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint = 0
			for _, node := range nodes { // For each node
				count += node.Capacity
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_disk_free_space", "Amount of free disk space across all Nodes in Bytes"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint = 0
			for _, node := range nodes { // For each node
				count += node.FreeDiskSpace
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_disk_total_space", "Amount of total available disk space across all Nodes in Bytes"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint = 0
			for _, node := range nodes { // For each node
				count += node.DiskSize
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_disk_anka_used_space", "Amount of disk space used by Anka across all Nodes in Bytes"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint = 0
			for _, node := range nodes { // For each node
				count += node.AnkaDiskUsage
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_cpu_core_count", "Count of CPU Cores across all Nodes"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint = 0
			for _, node := range nodes { // For each node
				count += node.CPU
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_cpu_util", "Total CPU utilization across all Nodes"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count float32 = 0
			for _, node := range nodes { // For each node
				count += node.CPUUtilization
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_ram_gb", "Total RAM available across all Nodes in GB"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint = 0
			for _, node := range nodes { // For each node
				count += node.RAM
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_ram_util", "Total RAM utilized across all Nodes"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count float32 = 0
			for _, node := range nodes { // For each node
				count += node.RAMUtilization
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_virtual_cpu_count", "Total Virtual CPU cores across all Nodes"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint = 0
			for _, node := range nodes { // For each node
				count += node.VCPUCount
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_virtual_ram_gb", "Total Virtual RAM across all Nodes"),
			event:  events.EventNodeUpdated,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint = 0
			for _, node := range nodes { // For each node
				count += node.VRAM
			}
			metric.Set(float64(count))
		},
	},
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)
	for _, nodesMetric := range ankaNodesMetrics {
		AddMetric(nodesMetric)
	}
}
