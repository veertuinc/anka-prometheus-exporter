package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type InstanceStateMetric struct {
	BaseAnkaMetric
}

func CountInstanceState(checkForState string, data []types.Instance) int {
	counter := 0
	for _, instanceData := range data {
		if instanceData.Vm.State == checkForState {
			counter++
		}
	}
	return counter
}

func (this InstanceStateMetric) GetEventHandler() func(interface{}) error {
	return func(instancesData interface{}) error {
		instances, err := ConvertToInstancesData(instancesData)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}
		for _, state := range types.InstanceStates {
			metric.With(prometheus.Labels{"state": state}).Set(float64(CountInstanceState(state, instances)))
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(InstanceStateMetric{BaseAnkaMetric{
		metric: CreateGaugeMetricVec("anka_instance_state_count", "Count of Instances in a particular State (label: state)", []string{"state"}),
		event:  events.EVENT_VM_DATA_UPDATED,
	}})

}
