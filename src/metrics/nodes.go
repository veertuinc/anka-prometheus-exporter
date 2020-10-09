package metrics

import (
	"github.com/veertuinc/anka-prometheus/src/events"
)

type NodesMetric struct {
	BaseAnkaMetric
}

func (this NodesMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		data, err := ConvertToNodeData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGauge(this.metric)
		if err != nil {
			return err
		}
		if this.name == "anka_nodes_count" {
			metric.Set(float64(len(data)))
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(NodesMetric{BaseAnkaMetric{
		name:   "anka_nodes_count",
		metric: CreateGaugeMetric("anka_nodes_count", "Total count of Anka Nodes"),
		event:  events.EVENT_NODE_UPDATED,
	}})

}
