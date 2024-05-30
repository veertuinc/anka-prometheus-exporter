package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

// TODO: can we make prometheus.GaugeVec support also .Gauge?
type NodeMetric struct {
	BaseAnkaMetric
	HandleData func([]types.Node, *prometheus.GaugeVec)
}

func (nm NodeMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		nodes, err := ConvertToNodeData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(nm.metric)
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

var ankaNodeMetrics = []NodeMetric{
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_instance_count", "Count of Instances running on the Node", []string{"id", "name", "arch"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGaugeVecMetric(len(nodes), "anka_node_instance_count", metric)
			for _, node := range nodes {
				if node.NodeName != "" {
					metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "arch": node.HostArch}).Set(float64(node.VMCount))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_instance_capacity", "Total Instance slots (capacity) on the Node", []string{"id", "name", "arch"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGaugeVecMetric(len(nodes), "anka_node_instance_capacity", metric)
			for _, node := range nodes {
				if node.NodeName != "" {
					metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "arch": node.HostArch}).Set(float64(node.Capacity))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_disk_free_space", "Amount of free disk space on the Node in Bytes", []string{"id", "name", "arch"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGaugeVecMetric(len(nodes), "anka_node_disk_free_space", metric)
			for _, node := range nodes {
				if node.NodeName != "" {
					metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "arch": node.HostArch}).Set(float64(node.FreeDiskSpace))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_disk_total_space", "Amount of total available disk space on the Node in Bytes", []string{"id", "name", "arch"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGaugeVecMetric(len(nodes), "anka_node_disk_total_space", metric)
			for _, node := range nodes {
				if node.NodeName != "" {
					metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "arch": node.HostArch}).Set(float64(node.DiskSize))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_disk_anka_used_space", "Amount of disk space used by Anka on the Node in Bytes", []string{"id", "name", "arch"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGaugeVecMetric(len(nodes), "anka_node_disk_anka_used_space", metric)
			for _, node := range nodes {
				if node.NodeName != "" {
					metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "arch": node.HostArch}).Set(float64(node.AnkaDiskUsage))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_cpu_core_count", "Number of CPU Cores in Node", []string{"id", "name", "arch"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGaugeVecMetric(len(nodes), "anka_node_cpu_core_count", metric)
			for _, node := range nodes {
				if node.NodeName != "" {
					metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "arch": node.HostArch}).Set(float64(node.CPU))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_cpu_util", "CPU utilization in node", []string{"id", "name", "arch"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGaugeVecMetric(len(nodes), "anka_node_cpu_util", metric)
			for _, node := range nodes {
				if node.NodeName != "" {
					metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "arch": node.HostArch}).Set(float64(node.CPUUtilization))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_ram_gb", "Total RAM available for the Node in GB", []string{"id", "name", "arch"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGaugeVecMetric(len(nodes), "anka_node_ram_gb", metric)
			for _, node := range nodes {
				if node.NodeName != "" {
					metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "arch": node.HostArch}).Set(float64(node.RAM))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_ram_util", "Total RAM utilized for the Node", []string{"id", "name", "arch"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGaugeVecMetric(len(nodes), "anka_node_ram_util", metric)
			for _, node := range nodes {
				if node.NodeName != "" {
					metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "arch": node.HostArch}).Set(float64(node.RAMUtilization))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_used_virtual_cpu_count", "Total Used Virtual CPU cores for the Node", []string{"id", "name", "arch"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGaugeVecMetric(len(nodes), "anka_node_used_virtual_cpu_count", metric)
			for _, node := range nodes {
				if node.NodeName != "" {
					metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "arch": node.HostArch}).Set(float64(node.UsedVCPUCount))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_used_virtual_ram_mb", "Total Used Virtual RAM for the Node in MB", []string{"id", "name", "arch"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGaugeVecMetric(len(nodes), "anka_node_used_virtual_ram_mb", metric)
			for _, node := range nodes {
				if node.NodeName != "" {
					metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "arch": node.HostArch}).Set(float64(node.UsedVRAM))
				}
			}
		},
	},
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)
	for _, nodeMetric := range ankaNodeMetrics {
		AddMetric(nodeMetric)
	}
}
