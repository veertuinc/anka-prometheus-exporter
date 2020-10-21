package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type NodeMetric struct {
	BaseAnkaMetric
}

func CountNodeState(checkForState string, data []types.Node) int {
	counter := 0
	for _, nodeData := range data {
		if nodeData.State == checkForState {
			counter++
		}
	}
	return counter
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
		if this.name == "anka_node_states_count" {
			for _, state := range types.NodeStates {
				metric.With(prometheus.Labels{"state": state}).Set(float64(0))
				metric.With(prometheus.Labels{"state": state}).Set(float64(CountNodeState(state, data)))
			}
		}
		for _, nodeData := range data { // Loop over each node
			nodeLabels := prometheus.Labels{
				"id":   nodeData.NodeID,
				"name": nodeData.NodeName,
			}
			if this.name == "anka_node_instance_count" {
				metric.With(nodeLabels).Set(float64(nodeData.VMCount))
			} else if this.name == "anka_node_instance_capacity" {
				metric.With(nodeLabels).Set(float64(nodeData.Capacity))
			} else if this.name == "anka_node_disk_free_space" {
				metric.With(nodeLabels).Set(float64(nodeData.FreeDiskSpace))
			} else if this.name == "anka_node_disk_total_space" {
				metric.With(nodeLabels).Set(float64(nodeData.DiskSize))
			} else if this.name == "anka_node_disk_anka_used_space" {
				metric.With(nodeLabels).Set(float64(nodeData.AnkaDiskUsage))
			} else if this.name == "anka_node_cpu_core_count" {
				metric.With(nodeLabels).Set(float64(nodeData.CPU))
			} else if this.name == "anka_node_cpu_util" {
				metric.With(nodeLabels).Set(float64(nodeData.CPUUtilization))
			} else if this.name == "anka_node_ram_gb" {
				metric.With(nodeLabels).Set(float64(nodeData.RAM))
			} else if this.name == "anka_node_ram_util" {
				metric.With(nodeLabels).Set(float64(nodeData.RAMUtilization))
			} else if this.name == "anka_node_virtual_cpu_count" {
				metric.With(nodeLabels).Set(float64(nodeData.VCPUCount))
			} else if this.name == "anka_node_virtual_ram_gb" {
				metric.With(nodeLabels).Set(float64(nodeData.VRAM))
			}
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_instance_count",
		metric: CreateGaugeMetricVec("anka_node_instance_count", "Count of Instances running on the Node", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_instance_capacity",
		metric: CreateGaugeMetricVec("anka_node_instance_capacity", "Total Instance slots (capacity) on the Node", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_states_count",
		metric: CreateGaugeMetricVec("anka_node_states_count", "Count of Nodes in a particular State (label: state)", []string{"state"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_disk_free_space",
		metric: CreateGaugeMetricVec("anka_node_disk_free_space", "Amount of free disk space on the Node in Bytes", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})
	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_disk_total_space",
		metric: CreateGaugeMetricVec("anka_node_disk_total_space", "Amount of total available disk space on the Node in Bytes", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})
	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_disk_anka_used_space",
		metric: CreateGaugeMetricVec("anka_node_disk_anka_used_space", "Amount of disk space used by Anka on the Node in Bytes", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_cpu_core_count",
		metric: CreateGaugeMetricVec("anka_node_cpu_core_count", "Number of CPU Cores in Node", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_cpu_util",
		metric: CreateGaugeMetricVec("anka_node_cpu_util", "CPU utilization in node", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_ram_gb",
		metric: CreateGaugeMetricVec("anka_node_ram_gb", "Total RAM available for the Node in GB", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_ram_util",
		metric: CreateGaugeMetricVec("anka_node_ram_util", "Total RAM utilized for the Node", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_virtual_cpu_count",
		metric: CreateGaugeMetricVec("anka_node_virtual_cpu_count", "Total Virtual CPU cores for the Node", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_virtual_ram_gb",
		metric: CreateGaugeMetricVec("anka_node_virtual_ram_gb", "Total Virtual RAM for the Node in GB", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

}
