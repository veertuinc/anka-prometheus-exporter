package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type NodeGroupMetric struct {
	BaseAnkaMetric
	HandleData func([]types.Node, []types.NodeGroup, *prometheus.GaugeVec)
}

func (this NodeGroupMetric) GetEventHandler() func(interface{}) error {
	return func(nodesData interface{}) error {
		nodes, err := ConvertToNodeData(nodesData)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}
		// TODO: make an API call for all groups
		var nodeGroups []types.NodeGroup
		for _, node := range nodes { // EACH GROUP
			for _, group := range node.Groups {
				nodeGroups = append(nodeGroups, group)
			}
		}
		nodeGroups = uniqueNodeGroupsArray(nodeGroups)
		this.HandleData(
			nodes,
			nodeGroups,
			metric,
		)
		return nil
	}
}

var ankaNodeGroupMetrics = []NodeGroupMetric{
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_nodes_count", "Count of Nodes in a particular Group", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(0))
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(CountNodeGroupNodes(focusGroup.Id, nodes)))
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_states_count", "Count of Groups in a particular state (labels: group, state)", []string{"group_name", "state"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, state := range types.NodeStates {
				for _, focusGroup := range nodeGroups { // EACH GROUP
					metric.With(prometheus.Labels{"group_name": focusGroup.Name, "state": state}).Set(float64(0))
					metric.With(prometheus.Labels{"group_name": focusGroup.Name, "state": state}).Set(float64(CountNodeGroupState(focusGroup.Id, state, nodes)))
				}
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_instance_capacity", "Total Instance slots (capacity) for the Group and its Nodes", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				var count uint = 0
				for _, node := range nodes {
					count = count + node.Capacity
				}
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(count))
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_instance_count", "Count of Instances slots in use for the Group (and Nodes)", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				var count uint = 0
				for _, node := range nodes {
					count = count + node.VMCount
				}
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(count))
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_disk_free_space", "Amount of free disk space for the Group (and Nodes) in Bytes", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				var count uint = 0
				for _, node := range nodes {
					count = count + node.FreeDiskSpace
				}
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(count))
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_disk_total_space", "Amount of total available disk space for the Group (and Nodes) in Bytes", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				var count uint = 0
				for _, node := range nodes {
					count = count + node.DiskSize
				}
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(count))
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_disk_anka_used_space", "Amount of disk space used by Anka for the Group (and Nodes) in Bytes", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				var count uint = 0
				for _, node := range nodes {
					count = count + node.AnkaDiskUsage
				}
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(count))
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_cpu_core_count", "Number of CPU Cores for the Group (and Nodes)", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				var count uint = 0
				for _, node := range nodes {
					count = count + node.CPU
				}
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(count))
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_cpu_util", "CPU utilization for the Group (and Nodes)", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				var count float32 = 0
				for _, node := range nodes {
					count = count + node.CPUUtilization
				}
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(count))
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_ram_gb", "Total RAM available for the Group (and Nodes) in GB", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				var count uint = 0
				for _, node := range nodes {
					count = count + node.RAM
				}
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(count))
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_ram_util", "Total RAM utilized for the Group (and Nodes)", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				var count float32 = 0
				for _, node := range nodes {
					count = count + node.RAMUtilization
				}
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(count))
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_virtual_cpu_count", "Total Virtual CPU cores for the Group (and Nodes)", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				var count uint = 0
				for _, node := range nodes {
					count = count + node.VCPUCount
				}
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(count))
			}
		},
	},
	NodeGroupMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_group_virtual_ram_gb", "Total Virtual RAM for the Group (and Nodes) in GB", []string{"group_name"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, nodeGroups []types.NodeGroup, metric *prometheus.GaugeVec) {
			for _, focusGroup := range nodeGroups { // EACH GROUP
				var count uint = 0
				for _, node := range nodes {
					count = count + node.VRAM
				}
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(count))
			}
		},
	},
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)
	for _, nodeGroupMetric := range ankaNodeGroupMetrics {
		AddMetric(nodeGroupMetric)
	}
}
