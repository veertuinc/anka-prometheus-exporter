package metrics

import (
	"github.com/veertuinc/anka-prometheus/events"
)

type AnkaNodesNumMetric struct {
	BaseAnkaMetric
}

func (this AnkaNodesNumMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		data, err := ConvertToNodeData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGauge(this.metric)
		if err != nil {
			return err
		}
		value := len(data)
		metric.Set(float64(value))
		return nil
	}
}

func init() {
	metricName := "anka_nodes_num"
	metricDesc := "Current number of nodes"
	ankaMetric := AnkaNodesNumMetric{}
	ankaMetric.metric = CreateGaugeMetric(metricName, metricDesc)
	ankaMetric.event = events.EVENT_NODE_UPDATED
	
	AddMetric(ankaMetric)
}