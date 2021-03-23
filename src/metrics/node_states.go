package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

// TODO: can we make prometheus.GaugeVec support also .Gauge?
type NodeStatesMetric struct {
	BaseAnkaMetric
	HandleData func([]types.Node, *prometheus.GaugeVec)
}

func (this NodeStatesMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		nodesData, err := ConvertToNodeData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}
		this.HandleData(
			nodesData,
			metric,
		)
		return nil
	}
}

var ankaNodeStatesMetrics = []NodeStatesMetric{
	NodeStatesMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_states_count", "Count of Nodes in a particular State (label: state)", []string{"state"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			var stateIntMap = intMapFromStringSlice(types.NodeStates)
			for _, node := range nodes {
				stateIntMap[node.State] = stateIntMap[node.State] + 1
			}
			for _, state := range types.NodeStates {
				metric.With(prometheus.Labels{"state": state}).Set(float64(stateIntMap[state]))
			}
		},
	},
	NodeStatesMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_node_states", "Node state (1 = current state) (label: id, name, state)", []string{"id", "name", "state"}),
			event:  events.EVENT_NODE_UPDATED,
		},
		HandleData: func(nodes []types.Node, metric *prometheus.GaugeVec) {
			for _, node := range nodes {
				for _, state := range types.NodeStates {
					if state == node.State {
						metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "state": node.State}).Set(float64(1))
					} else {
						metric.With(prometheus.Labels{"id": node.NodeID, "name": node.NodeName, "state": state}).Set(float64(0))
					}
				}
			}
			checkAndHandleResetOfMetric(len(nodes), "anka_node_states", metric)
		},
	},
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)
	for _, nodeStatesMetric := range ankaNodeStatesMetrics {
		AddMetric(nodeStatesMetric)
	}
}
