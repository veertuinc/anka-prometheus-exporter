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

func (nm NodesMetric) GetEventHandler() func(interface{}) error {
	return func(nodesData interface{}) error {
		nodes, err := ConvertToNodeData(nodesData)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGauge(nm.metric)
		if err != nil {
			return err
		}
		nm.HandleData(
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
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			metric.Set(float64(len(nodes)))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_instance_count", "Count of Instance slots in use across all Nodes"),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint
			for _, node := range nodes { // For each node
				count = count + node.VMCount
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_instance_capacity", "Count of total Instance Capacity across all Nodes"),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint
			for _, node := range nodes { // For each node
				count = count + node.Capacity
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_disk_free_space", "Amount of free disk space across all Nodes in Bytes"),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint
			for _, node := range nodes { // For each node
				count = count + node.FreeDiskSpace
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_disk_total_space", "Amount of total available disk space across all Nodes in Bytes"),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint
			for _, node := range nodes { // For each node
				count = count + node.DiskSize
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_disk_anka_used_space", "Amount of disk space used by Anka across all Nodes in Bytes"),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint
			for _, node := range nodes { // For each node
				count = count + node.AnkaDiskUsage
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_cpu_core_count", "Count of CPU Cores across all Nodes"),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint
			for _, node := range nodes { // For each node
				count = count + node.CPU
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_cpu_util", "Total CPU utilization across all Nodes"),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count float64
			for _, node := range nodes { // For each node
				count = count + node.CPUUtilization
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_ram_gb", "Total RAM available across all Nodes in GB"),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint
			for _, node := range nodes { // For each node
				count = count + node.RAM
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_ram_util", "Total RAM utilized across all Nodes"),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count float64
			for _, node := range nodes { // For each node
				count = count + node.RAMUtilization
			}
			metric.Set(count)
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_used_virtual_cpu_count", "Total Used Virtual CPU cores across all Nodes"),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint
			for _, node := range nodes { // For each node
				count = count + node.UsedVCPUCount
			}
			metric.Set(float64(count))
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_nodes_used_virtual_ram_mb", "Total Used Virtual RAM across all Nodes in MB"),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric prometheus.Gauge) {
			var count uint
			for _, node := range nodes { // For each node
				count = count + node.UsedVRAM
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
