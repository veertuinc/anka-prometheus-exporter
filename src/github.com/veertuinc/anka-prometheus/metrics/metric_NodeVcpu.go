package metrics

import (
	"github.com/veertuinc/anka-prometheus/events"
	"github.com/prometheus/client_golang/prometheus"
)

type VcpuMetric struct {
	BaseAnkaMetric
}

func (this VcpuMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {

		data, err := ConvertToNodeData(d)
		if err != nil {
			return err
		}

		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}
		
		for _, nodeData := range data {
			labels := prometheus.Labels{"id": nodeData.NodeID, "name": nodeData.NodeName}
			value := nodeData.VCPUCount
			metric.With(labels).Set(float64(value))
		}
		return nil
	}
}

func init() {
	metricName := "anka_node_virtual_cpu_count"
	metricDesc := "Virtual CPU cores in Node"
	labels := []string{"id", "name"}

	ankaMetric := VcpuMetric{}
	ankaMetric.metric = CreateGaugeMetricVec(metricName, metricDesc, labels)
	ankaMetric.event = events.EVENT_NODE_UPDATED
	AddMetric(ankaMetric)
}