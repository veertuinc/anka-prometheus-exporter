package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus/src/events"
	"github.com/veertuinc/anka-prometheus/src/types"
)

type InstancesMetric struct {
	BaseAnkaMetric
}

func CountVMState(checkForState types.InstanceState, data []types.InstanceInfo) int {
	counter := 0
	for _, instanceData := range data {
		if instanceData.Vm.State == checkForState {
			counter++
		}
	}
	return counter
}

func (this InstancesMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		data, err := ConvertToInstancesData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}
		metric.With(prometheus.Labels{"state": "Scheduling"}).Set(float64(0))
		metric.With(prometheus.Labels{"state": "Pulling"}).Set(float64(0))
		metric.With(prometheus.Labels{"state": "Started"}).Set(float64(0))
		metric.With(prometheus.Labels{"state": "Stopping"}).Set(float64(0))
		metric.With(prometheus.Labels{"state": "Stopped"}).Set(float64(0))
		metric.With(prometheus.Labels{"state": "Terminating"}).Set(float64(0))
		metric.With(prometheus.Labels{"state": "Terminated"}).Set(float64(0))
		metric.With(prometheus.Labels{"state": "Error"}).Set(float64(0))
		metric.With(prometheus.Labels{"state": "Pushing"}).Set(float64(0))
		for _, instanceData := range data {
			metric.With(prometheus.Labels{"state": string(instanceData.Vm.State)}).Set(float64(CountVMState(instanceData.Vm.State, data)))
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(InstancesMetric{BaseAnkaMetric{
		name:   "anka_instance_states",
		metric: CreateGaugeMetricVec("anka_instance_states", "Instance state counts", []string{"state"}),
		event:  events.EVENT_VM_DATA_UPDATED,
	}})

}
