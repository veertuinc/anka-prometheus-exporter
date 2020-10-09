package metrics

import (
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
)

type NodesMetric struct {
	BaseAnkaMetric
}

func (this NodesMetric) GetEventHandler() func(interface{}) error {
	return func(nodesData interface{}) error {
		nodes, err := ConvertToNodeData(nodesData)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGauge(this.metric)
		if err != nil {
			return err
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
		if this.name == "anka_nodes_count" {
			metric.Set(float64(len(nodes)))
		} else if this.name == "anka_nodes_instance_count" {
			metric.Set(float64(instanceCount))
		} else if this.name == "anka_nodes_instance_capacity" {
			metric.Set(float64(instanceCapacity))
		} else if this.name == "anka_nodes_disk_free_space" {
			metric.Set(float64(diskFreeSpace))
		} else if this.name == "anka_nodes_disk_total_space" {
			metric.Set(float64(diskTotalSpace))
		} else if this.name == "anka_nodes_disk_anka_used_space" {
			metric.Set(float64(diskAnkaUsedSpace))
		} else if this.name == "anka_nodes_cpu_core_count" {
			metric.Set(float64(cpuCoreCount))
		} else if this.name == "anka_nodes_cpu_util" {
			metric.Set(float64(cpuUtil))
		} else if this.name == "anka_nodes_ram_gb" {
			metric.Set(float64(ramGB))
		} else if this.name == "anka_nodes_ram_util" {
			metric.Set(float64(ramUtil))
		} else if this.name == "anka_nodes_virtual_cpu_count" {
			metric.Set(float64(virtualCPUCount))
		} else if this.name == "anka_nodes_virtual_ram_gb" {
			metric.Set(float64(virtualRAMGB))
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_count",
		metric: CreateGaugeMetric("anka_nodes_count", "Count of total Anka Nodes"),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_instance_count",
		metric: CreateGaugeMetric("anka_nodes_instance_count", "Count of Instance slots in use across all Nodes"),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_instance_capacity",
		metric: CreateGaugeMetric("anka_nodes_instance_capacity", "Count of total Instance Capacity across all Nodes"),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_disk_free_space",
		metric: CreateGaugeMetric("anka_nodes_disk_free_space", "Amount of free disk space across all Nodes in Bytes"),
		event:  events.EVENT_NODE_UPDATED,
	}})
	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_disk_total_space",
		metric: CreateGaugeMetric("anka_nodes_disk_total_space", "Amount of total available disk space across all Nodes in Bytes"),
		event:  events.EVENT_NODE_UPDATED,
	}})
	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_disk_anka_used_space",
		metric: CreateGaugeMetric("anka_nodes_disk_anka_used_space", "Amount of disk space used by Anka across all Nodes in Bytes"),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_cpu_core_count",
		metric: CreateGaugeMetric("anka_nodes_cpu_core_count", "Count of CPU Cores across all Nodes"),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_cpu_util",
		metric: CreateGaugeMetric("anka_nodes_cpu_util", "Total CPU utilization across all Nodes"),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_ram_gb",
		metric: CreateGaugeMetric("anka_nodes_ram_gb", "Total RAM available across all Nodes in GB"),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_ram_util",
		metric: CreateGaugeMetric("anka_nodes_ram_util", "Total RAM utilized across all Nodes"),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_virtual_cpu_count",
		metric: CreateGaugeMetric("anka_nodes_virtual_cpu_count", "Total Virtual CPU cores across all Nodes"),
		event:  events.EVENT_NODE_UPDATED,
	}})

	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_virtual_ram_gb",
		metric: CreateGaugeMetric("anka_nodes_virtual_ram_gb", "Total Virtual RAM across all Nodes"),
		event:  events.EVENT_NODE_UPDATED,
	}})

}
