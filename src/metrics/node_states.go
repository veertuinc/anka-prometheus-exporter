package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

// TODO: can we make prometheus.GaugeVec support also .Gauge?
type NodeStatesMetric struct {
	BaseAnkaMetric
	HandleData func(types.Node, *prometheus.GaugeVec, prometheus.Labels)
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

func (this NodeStatesMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		data, err := ConvertToNodeData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}
		// TODO: don't loop over states, instead create a map with counts and just increment
		for _, state := range types.NodeStates {
			metric.With(prometheus.Labels{"state": state}).Set(float64(0))
			metric.With(prometheus.Labels{"state": state}).Set(float64(CountNodeState(state, data)))
		}
		return nil
	}
}

var ankaNodeStatesMetrics = []NodeStatesMetric{
	NodeStatesMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_states_count", "Count of Nodes in a particular State (label: state)", []string{"state"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodeData types.Node, metric *prometheus.GaugeVec, labels prometheus.Labels) {
			metric.With(labels).Set(float64(nodeData.VMCount))
		},
	},
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)
	for _, nodeStatesMetric := range ankaNodeStatesMetrics {
		AddMetric(nodeStatesMetric)
	}
}
