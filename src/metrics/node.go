package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

// TODO: can we make prometheus.GaugeVec support also .Gauge?
type NodeMetric struct {
	BaseAnkaMetric
	HandleData func(types.Node, *prometheus.GaugeVec, prometheus.Labels)
}

func (this NodeMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		data, err := ConvertToNodeData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}
		for _, nodeData := range data { // Loop over each node
			this.HandleData(
				nodeData,
				metric,
				prometheus.Labels{
					"id":   nodeData.NodeID,
					"name": nodeData.NodeName,
				},
			)
		}
		return nil
	}
}

var ankaNodeMetrics = []NodeMetric{
	NodeMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_instance_count", "Count of Instances running on the Node", []string{"id", "name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.VMCount))
		},
	},
	NodeMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_instance_capacity", "Total Instance slots (capacity) on the Node", []string{"id", "name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.Capacity))
		},
	},
	NodeMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_disk_free_space", "Amount of free disk space on the Node in Bytes", []string{"id", "name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.FreeDiskSpace))
		},
	},
	NodeMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_disk_total_space", "Amount of total available disk space on the Node in Bytes", []string{"id", "name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.DiskSize))
		},
	},
	NodeMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_disk_anka_used_space", "Amount of disk space used by Anka on the Node in Bytes", []string{"id", "name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.AnkaDiskUsage))
		},
	},
	NodeMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_cpu_core_count", "Number of CPU Cores in Node", []string{"id", "name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.CPU))
		},
	},
	NodeMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_cpu_util", "CPU utilization in node", []string{"id", "name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.CPUUtilization))
		},
	},
	NodeMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_ram_gb", "Total RAM available for the Node in GB", []string{"id", "name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.RAM))
		},
	},
	NodeMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_ram_util", "Total RAM utilized for the Node", []string{"id", "name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.RAMUtilization))
		},
	},
	NodeMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_virtual_cpu_count", "Total Virtual CPU cores for the Node", []string{"id", "name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.VCPUCount))
		},
	},
	NodeMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_virtual_ram_gb", "Total Virtual RAM for the Node in GB", []string{"id", "name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.VRAM))
		},
	},
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)
	for _, nodeMetric := range ankaNodeMetrics {
		AddMetric(nodeMetric)
	}
}
