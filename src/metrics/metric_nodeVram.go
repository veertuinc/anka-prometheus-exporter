package metrics

import (
	"github.com/veertuinc/anka-prometheus/events"
	"github.com/prometheus/client_golang/prometheus"
)

type VramMetric struct {
	BaseAnkaMetric
}

func (this VramMetric) GetEventHandler() func(interface{}) error {
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
			value := nodeData.VRAM / 1024
			metric.With(labels).Set(float64(value))
		}
		return nil
	}
}

func init() {
	metricName := "anka_node_virtual_ram_gb"
	metricDesc := "Virtual RAM in Node"
	labels := []string{"id", "name"}

	ankaMetric := VramMetric{}
	ankaMetric.metric = CreateGaugeMetricVec(metricName, metricDesc, labels)
	ankaMetric.event = events.EVENT_NODE_UPDATED
	AddMetric(ankaMetric)
}