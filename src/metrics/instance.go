package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type InstanceStateMetric struct {
	BaseAnkaMetric
}

type InstanceStatePerTemplateMetric struct {
	BaseAnkaMetric
}

type InstanceStatePerGroupMetric struct {
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

func CountInstanceTemplateState(templateWeWant string, stateWeWant string, data []types.Instance) int {
	counter := 0
	for _, instanceData := range data {
		if instanceData.Vm.State == stateWeWant {
			if instanceData.Vm.TemplateUUID == templateWeWant {
				counter++
			}
		}
	}
	return counter
}

func CountInstanceGroupState(groupWeWant string, stateWeWant string, data []types.Instance) int {
	counter := 0
	for _, instanceData := range data {
		if instanceData.Vm.State == stateWeWant {
			if instanceData.Vm.GroupUUID == groupWeWant {
				counter++
			}
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

func (this InstanceStatePerTemplateMetric) GetEventHandler() func(interface{}) error {
	return func(instancesData interface{}) error {
		instances, err := ConvertToInstancesData(instancesData)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}
		// TODO: Make sure all templates, even if not used show up (API call?)
		var instanceTemplates []string
		for _, instance := range instances {
			instanceTemplates = append(instanceTemplates, instance.Vm.TemplateUUID)
		}
		instanceTemplates = uniqueThisStringArray(instanceTemplates)
		for _, state := range types.InstanceStates {
			for _, instanceTemplate := range instanceTemplates {
				metric.With(prometheus.Labels{"state": state, "template_uuid": instanceTemplate}).Set(float64(CountInstanceTemplateState(instanceTemplate, state, instances)))
			}
		}
		return nil
	}
}

func (this InstanceStatePerGroupMetric) GetEventHandler() func(interface{}) error {
	return func(instancesData interface{}) error {
		instances, err := ConvertToInstancesData(instancesData)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}
		// TODO: Make sure all groups, even if not used show up (API call?)
		var instanceGroups []string
		for _, instance := range instances {
			if instance.Vm.GroupUUID != "" {
				instanceGroups = append(instanceGroups, instance.Vm.GroupUUID)
			}
		}
		instanceGroups = uniqueThisStringArray(instanceGroups)
		for _, state := range types.InstanceStates {
			for _, instanceGroup := range instanceGroups {
				metric.With(prometheus.Labels{"state": state, "group_uuid": instanceGroup}).Set(float64(CountInstanceGroupState(instanceGroup, state, instances)))
			}
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(InstanceStateMetric{BaseAnkaMetric{
		metric: CreateGaugeMetricVec("anka_instance_state_count", "Count of Instances in a particular State (label: state)", []string{"state"}),
		event:  events.EVENT_VM_DATA_UPDATED,
	}})

	AddMetric(InstanceStatePerTemplateMetric{BaseAnkaMetric{
		metric: CreateGaugeMetricVec("anka_instance_state_per_template_count", "Count of Instances in a particular state, per Template (label: state, template_name)", []string{"state", "template_uuid"}),
		event:  events.EVENT_VM_DATA_UPDATED,
	}})

	AddMetric(InstanceStatePerGroupMetric{BaseAnkaMetric{
		metric: CreateGaugeMetricVec("anka_instance_state_per_group_count", "Count of Instances in a particular state, per Group (label: state, group_name)", []string{"state", "group_uuid"}),
		event:  events.EVENT_VM_DATA_UPDATED,
	}})

}
