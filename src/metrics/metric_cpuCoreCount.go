package metrics

import (
	"github.com/veertuinc/anka-prometheus/events"
	"github.com/prometheus/client_golang/prometheus"
)

type CpuCoreMetric struct {
	BaseAnkaMetric
}

func (this CpuCoreMetric) GetEventHandler() func(interface{}) error {
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
			labls := prometheus.Labels{"id": nodeData.NodeID, "name": nodeData.NodeName}
			value := nodeData.CPU
			metric.With(labls).Set(float64(value))
		}
		return nil
	}
}

func init() {
	metricName := "anka_node_cpu_core_count"
	metricDesc := "Number of CPU Cores in Node"
	labels := []string{"id", "name"}

	ankaMetric := CpuCoreMetric{}
	ankaMetric.metric = CreateGaugeMetricVec(metricName, metricDesc, labels)
	ankaMetric.event = events.EVENT_NODE_UPDATED
	AddMetric(ankaMetric)
}