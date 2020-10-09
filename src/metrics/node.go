package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus/src/events"
	"github.com/veertuinc/anka-prometheus/src/types"
)

type NodeMetric struct {
	BaseAnkaMetric
}

func CountNodeState(checkForState types.NodeState, data []types.NodeInfo) int {
	counter := 0
	for _, nodeData := range data {
		if nodeData.State == checkForState {
			counter++
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
		if this.name == "anka_node_states" {
			metric.With(prometheus.Labels{"state": "Offline"}).Set(float64(0))
			metric.With(prometheus.Labels{"state": "Inactive (Invalid License)"}).Set(float64(0))
			metric.With(prometheus.Labels{"state": "Active"}).Set(float64(0))
			metric.With(prometheus.Labels{"state": "Updating"}).Set(float64(0))
		}
		for _, nodeData := range data { // Loop over each node
			nodeLabels := prometheus.Labels{
				"id":   nodeData.NodeID,
				"name": nodeData.NodeName,
			}
			if this.name == "anka_node_disk_free_space" {
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
			} else if this.name == "anka_node_vm_capacity" {
				metric.With(nodeLabels).Set(float64(nodeData.Capacity))
			} else if this.name == "anka_node_states" {
				metric.With(prometheus.Labels{"state": string(nodeData.State)}).Set(float64(CountNodeState(nodeData.State, data)))
			}
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

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

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_vm_capacity",
		metric: CreateGaugeMetricVec("anka_node_vm_capacity", "Total VM slots (capacity) for the Node", []string{"id", "name"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodeMetric{BaseAnkaMetric{
		name:   "anka_node_states",
		metric: CreateGaugeMetricVec("anka_node_states", "Node States", []string{"state"}),
		event:  events.EVENT_NODE_UPDATED,
	}})

}
