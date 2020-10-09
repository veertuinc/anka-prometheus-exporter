package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type InstanceMetric struct {
	BaseAnkaMetric
}

func (this InstanceMetric) GetEventHandler() func(interface{}) error {
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
		var instanceTemplates []string
		var instanceGroups []string
		if this.name == "anka_instance_state_per_template_count" || this.name == "anka_instance_state_per_group_count" {
			for _, instance := range instances { // EACH GROUP
				instanceTemplates = append(instanceTemplates, instance.Vm.TemplateUUID)
				if instance.Vm.GroupUUID != "" {
					instanceGroups = append(instanceGroups, instance.Vm.GroupUUID)
				}
			}
			instanceTemplates = uniqueThisStringArray(instanceTemplates)
			instanceGroups = uniqueThisStringArray(instanceGroups)
		}
		// Populate
		for _, state := range types.InstanceStates {
			if this.name == "anka_instance_state_count" {
				metric.With(prometheus.Labels{"state": state}).Set(float64(0))
				metric.With(prometheus.Labels{"state": state}).Set(float64(CountVMState(state, instances)))
			} else if this.name == "anka_instance_state_per_template_count" {
				for _, instanceTemplate := range instanceTemplates {
					metric.With(prometheus.Labels{"state": state, "template_uuid": instanceTemplate}).Set(float64(0))
					metric.With(prometheus.Labels{"state": state, "template_uuid": instanceTemplate}).Set(float64(CountInstanceTemplateState(instanceTemplate, state, instances)))
				}
			} else if this.name == "anka_instance_state_per_group_count" {
				for _, instanceGroup := range instanceGroups {
					metric.With(prometheus.Labels{"state": state, "group_uuid": instanceGroup}).Set(float64(0))
					metric.With(prometheus.Labels{"state": state, "group_uuid": instanceGroup}).Set(float64(CountInstanceGroupState(instanceGroup, state, instances)))
				}
			}
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(InstanceMetric{BaseAnkaMetric{
		name:   "anka_instance_state_count",
		metric: CreateGaugeMetricVec("anka_instance_state_count", "Count of Instances in a particular State (label: state)", []string{"state"}),
		event:  events.EVENT_VM_DATA_UPDATED,
	}})

	AddMetric(InstanceMetric{BaseAnkaMetric{
		name:   "anka_instance_state_per_template_count",
		metric: CreateGaugeMetricVec("anka_instance_state_per_template_count", "Count of Instances in a particular state, per Template (label: state, template_name)", []string{"state", "template_uuid"}),
		event:  events.EVENT_VM_DATA_UPDATED,
	}})

	AddMetric(InstanceMetric{BaseAnkaMetric{
		name:   "anka_instance_state_per_group_count",
		metric: CreateGaugeMetricVec("anka_instance_state_per_group_count", "Count of Instances in a particular state, per Group (label: state, group_name)", []string{"state", "group_uuid"}),
		event:  events.EVENT_VM_DATA_UPDATED,
	}})

}
