package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type InstanceStateMetric struct {
	BaseAnkaMetric
}

func (ism InstanceStateMetric) GetEventHandler() func(interface{}) error {
	return func(instancesData interface{}) error {
		instances, err := ConvertToInstancesData(instancesData)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(ism.metric)
		if err != nil {
			return err
		}
		var archStateMap = intMapFromTwoStringSlices(types.Architectures, types.InstanceStates)
		for _, instance := range instances {
			if instance.Vm.Arch != "" && instance.Vm.State != "" { // prevent panic: assignment to entry in nil map when no Arch for instance
				archStateMap[instance.Vm.Arch][instance.Vm.State] = archStateMap[instance.Vm.Arch][instance.Vm.State] + 1
			}
		}
		for _, arch := range types.Architectures {
			for _, state := range types.InstanceStates {
				metric.With(prometheus.Labels{"arch": arch, "state": state}).Set(float64(archStateMap[arch][state]))
			}
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(InstanceStateMetric{BaseAnkaMetric{
		metric: CreateGaugeMetricVec("anka_instance_state_count", "Count of Instances in a particular State (label: arch, state)", []string{"arch", "state"}),
		event:  events.EVENT_VM_DATA_UPDATED,
	}})

}
