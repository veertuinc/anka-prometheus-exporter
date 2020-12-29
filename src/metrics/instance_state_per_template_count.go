package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
)

type InstanceStatePerTemplateCountMetric BaseAnkaMetric

func (this InstanceStatePerTemplateCountMetric) GetEventHandler() func(interface{}) error {
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
		var instanceTemplates []string
		var instanceGroups []string
		for _, instance := range instances { // EACH GROUP
			instanceTemplates = append(instanceTemplates, instance.Vm.TemplateUUID)
			if instance.Vm.GroupUUID != "" {
				instanceGroups = append(instanceGroups, instance.Vm.GroupUUID)
			}
		}
		instanceTemplates = uniqueThisStringArray(instanceTemplates)
		instanceGroups = uniqueThisStringArray(instanceGroups)

		// Populate
		for _, state := range InstanceStates {
			for _, instanceTemplate := range instanceTemplates {
				metric.With(prometheus.Labels{"state": state, "template_uuid": instanceTemplate}).Set(float64(CountInstanceTemplateState(instanceTemplate, state, instances)))
			}
		}
		return nil
	}
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)

	AddMetric(BaseAnkaMetric{
		metric: CreateGaugeMetricVec("anka_instance_state_per_template_count", "Count of Instances in a particular state, per Template (label: state, template_name)", []string{"state", "template_uuid"}),
		event:  events.EVENT_VM_DATA_UPDATED,
	})

}
