package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type InstanceStateMetric struct {
	BaseAnkaMetric
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
		var stateIntMap = intMapFromStringSlice(types.InstanceStates)
		for _, instance := range instances {
			stateIntMap[instance.Vm.State]++
		}
		for _, state := range types.InstanceStates {
			metric.With(prometheus.Labels{"state": state}).Set(float64(stateIntMap[state]))
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(InstanceStateMetric{BaseAnkaMetric{
		metric: CreateGaugeMetricVec("anka_instance_state_count", "Count of Instances in a particular State (label: state)", []string{"state"}),
		event:  events.EventVmDataUpdated,
	}})
}
