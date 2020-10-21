package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type NodeGroupMetric struct {
	BaseAnkaMetric
}

func uniqueNodeGroupsArray(arr []types.NodeGroup) []types.NodeGroup {
	occured := map[types.NodeGroup]bool{}
	result := []types.NodeGroup{}
	for e := range arr {
		if occured[arr[e]] != true {
			occured[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}

func CountNodeGroupNodes(GroupIdWeWant string, nodesData []types.Node) int {
	counter := 0
	for _, node := range nodesData {
		for _, group := range node.Groups {
			if group.Id == GroupIdWeWant {
				counter++
			}
		}
	}
	return counter
}

func CountNodeGroupState(groupIdWeWant string, stateWeWant string, nodesData []types.Node) int {
	counter := 0
	for _, node := range nodesData {
		if node.State == stateWeWant {
			for _, group := range node.Groups {
				if group.Id == groupIdWeWant {
					counter++
				}
			}
		}
	}
	return counter
}

func (this NodeGroupMetric) GetEventHandler() func(interface{}) error {
	return func(nodesData interface{}) error {

		nodes, err := ConvertToNodeData(nodesData)
		if err != nil {
			return err
		}

		// Collect all groups; TODO: Figure out how to make multiple endpoint calls and aggregate it all into a single data object
		var nodeGroups []types.NodeGroup
		for _, node := range nodes { // EACH GROUP
			for _, group := range node.Groups {
				nodeGroups = append(nodeGroups, group)
			}
		}
		nodeGroups = uniqueNodeGroupsArray(nodeGroups)

		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}

		for _, focusGroup := range nodeGroups { // EACH GROUP
			if this.name == "anka_node_group_states_count" {
				for _, state := range types.NodeStates {
					metric.With(prometheus.Labels{"group_name": focusGroup.Name, "state": state}).Set(float64(0))
					metric.With(prometheus.Labels{"group_name": focusGroup.Name, "state": state}).Set(float64(CountNodeGroupState(focusGroup.Id, state, nodes)))
				}
			}
			if this.name == "anka_node_group_nodes_count" {
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(0))
				metric.With(prometheus.Labels{"group_name": focusGroup.Name}).Set(float64(CountNodeGroupNodes(focusGroup.Id, nodes)))
			}
			nodeGroupLabels := prometheus.Labels{
				"group_name": focusGroup.Name,
			}
			var instanceCount uint = 0
			var instanceCapacity uint = 0
			var diskFreeSpace uint = 0
			var diskTotalSpace uint = 0
			var diskAnkaUsedSpace uint = 0
			var cpuCoreCount uint = 0
			var ramGB uint = 0
			var cpuUtil float32 = 0
			var ramUtil float32 = 0
			var virtualCPUCount uint = 0
			var virtualRAMGB uint = 0
			for _, node := range nodes { // For each node
				nodeMatchesGroup := false
				for index := range node.Groups { // see if it matches the group
					if node.Groups[index].Name == focusGroup.Name {
						nodeMatchesGroup = true
						break
					}
				}
				if nodeMatchesGroup == true { // if it matches the group, collect counts
					instanceCount = instanceCount + node.VMCount
					instanceCapacity = instanceCapacity + node.Capacity
					diskFreeSpace = diskFreeSpace + node.FreeDiskSpace
					diskTotalSpace = diskTotalSpace + node.DiskSize
					diskAnkaUsedSpace = diskAnkaUsedSpace + node.AnkaDiskUsage
					cpuCoreCount = cpuCoreCount + node.CPU
					ramGB = ramGB + node.RAM
					cpuUtil = cpuUtil + node.CPUUtilization
					ramUtil = ramUtil + node.RAMUtilization
					virtualCPUCount = virtualCPUCount + node.VCPUCount
					virtualRAMGB = virtualRAMGB + node.VRAM
				}
			}
			if this.name == "anka_node_group_instance_count" {
				metric.With(nodeGroupLabels).Set(float64(instanceCount))
			} else if this.name == "anka_node_group_instance_capacity" {
				metric.With(nodeGroupLabels).Set(float64(instanceCapacity))
			} else if this.name == "anka_node_group_disk_free_space" {
				metric.With(nodeGroupLabels).Set(float64(diskFreeSpace))
			} else if this.name == "anka_node_group_disk_total_space" {
				metric.With(nodeGroupLabels).Set(float64(diskTotalSpace))
			} else if this.name == "anka_node_group_disk_anka_used_space" {
				metric.With(nodeGroupLabels).Set(float64(diskAnkaUsedSpace))
			} else if this.name == "anka_node_group_cpu_core_count" {
				metric.With(nodeGroupLabels).Set(float64(cpuCoreCount))
			} else if this.name == "anka_node_group_cpu_util" {
				metric.With(nodeGroupLabels).Set(float64(cpuUtil))
			} else if this.name == "anka_node_group_ram_gb" {
				metric.With(nodeGroupLabels).Set(float64(ramGB))
			} else if this.name == "anka_node_group_ram_util" {
				metric.With(nodeGroupLabels).Set(float64(ramUtil))
			} else if this.name == "anka_node_group_virtual_cpu_count" {
				metric.With(nodeGroupLabels).Set(float64(virtualCPUCount))
			} else if this.name == "anka_node_group_virtual_ram_gb" {
				metric.With(nodeGroupLabels).Set(float64(virtualRAMGB))
			}
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_nodes_count",
		metric: CreateGaugeMetricVec("anka_node_group_nodes_count", "Count of Nodes in a particular Group", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_states_count",
		metric: CreateGaugeMetricVec("anka_node_group_states_count", "Count of Groups in a particular state (labels: group, state)", []string{"group_name", "state"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_instance_count",
		metric: CreateGaugeMetricVec("anka_node_group_instance_count", "Count of Instances slots in use for the Group (and Nodes)", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_disk_free_space",
		metric: CreateGaugeMetricVec("anka_node_group_disk_free_space", "Amount of free disk space for the Group (and Nodes) in Bytes", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})
	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_disk_total_space",
		metric: CreateGaugeMetricVec("anka_node_group_disk_total_space", "Amount of total available disk space for the Group (and Nodes) in Bytes", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})
	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_disk_anka_used_space",
		metric: CreateGaugeMetricVec("anka_node_group_disk_anka_used_space", "Amount of disk space used by Anka for the Group (and Nodes) in Bytes", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_cpu_core_count",
		metric: CreateGaugeMetricVec("anka_node_group_cpu_core_count", "Number of CPU Cores for the Group (and Nodes)", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_cpu_util",
		metric: CreateGaugeMetricVec("anka_node_group_cpu_util", "CPU utilization for the Group (and Nodes)", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_ram_gb",
		metric: CreateGaugeMetricVec("anka_node_group_ram_gb", "Total RAM available for the Group (and Nodes) in GB", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_ram_util",
		metric: CreateGaugeMetricVec("anka_node_group_ram_util", "Total RAM utilized for the Group (and Nodes)", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_virtual_cpu_count",
		metric: CreateGaugeMetricVec("anka_node_group_virtual_cpu_count", "Total Virtual CPU cores for the Group (and Nodes)", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_virtual_ram_gb",
		metric: CreateGaugeMetricVec("anka_node_group_virtual_ram_gb", "Total Virtual RAM for the Group (and Nodes) in GB", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeGroupMetric{BaseAnkaMetric{
		name:   "anka_node_group_instance_capacity",
		metric: CreateGaugeMetricVec("anka_node_group_instance_capacity", "Total Instance slots (capacity) for the Group (and Nodes)", []string{"group_name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

}
