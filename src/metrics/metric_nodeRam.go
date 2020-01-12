package metrics

import (
	"github.com/veertuinc/anka-prometheus/events"
	"github.com/prometheus/client_golang/prometheus"
)

type NodeRamMetric struct {
	BaseAnkaMetric
}

func (this NodeRamMetric) GetEventHandler() func(interface{}) error {
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
			value := nodeData.RAM
			metric.With(labels).Set(float64(value))
		}
		return nil
	}
}

func init() {
	metricName := "anka_node_ram_gb"
	metricDesc := "RAM in Node"
	labels := []string{"id", "name"}

	ankaMetric := NodeRamMetric{}
	ankaMetric.metric = CreateGaugeMetricVec(metricName, metricDesc, labels)
	ankaMetric.event = events.EVENT_NODE_UPDATED
	AddMetric(ankaMetric)
}