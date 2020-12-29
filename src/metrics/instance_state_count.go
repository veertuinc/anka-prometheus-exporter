package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
)

type InstanceStateCountMetric BaseAnkaMetric

func (this InstanceStateCountMetric) GetEventHandler() func(interface{}) error {
	return func(instancesData interface{}) error {
		instances, err := ConvertToInstancesData(instancesData)
		if err != nil {
			return err
		}

		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}

		// Collect templateUUIDs and GroupUUIDs
		// TODO: Get template name and group name and include it as a label in the metrics
		// TODO: Make sure all groups, even if not used show up (API call?)
		// TODO: Make sure all templates, even if not used show up (API call?)

		// Populate
		for _, state := range InstanceStates {
			metric.With(prometheus.Labels{"state": state}).Set(float64(CountVMState(state, instances)))
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(BaseAnkaMetric{
		metric: CreateGaugeMetricVec("anka_instance_state_count", "Count of Instances in a particular State (label: state)", []string{"state"}),
		event:  events.EVENT_VM_DATA_UPDATED,
	})

}
