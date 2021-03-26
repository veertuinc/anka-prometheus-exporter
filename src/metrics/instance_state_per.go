package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type InstanceStatePerMetric struct {
	BaseAnkaMetric
	HandleData func([]types.Instance, *prometheus.GaugeVec)
}

func (this InstanceStatePerMetric) GetEventHandler() func(interface{}) error {
	return func(instancesData interface{}) error {
		instances, err := ConvertToInstancesData(instancesData)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(this.metric)
		if err != nil {
			return err
		}
		this.HandleData(
			instances,
			metric,
		)
		return nil
	}
}

var ankaInstanceStatePerMetrics = []InstanceStatePerMetric{
	InstanceStatePerMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_instance_state_per_template_count", "Count of Instances in a particular state, per Template (label: state, template_name)", []string{"state", "template_uuid"}),
			event:  events.EVENT_VM_DATA_UPDATED,
		},
		HandleData: func(instances []types.Instance, metric *prometheus.GaugeVec) {
			var instanceTemplates []string
			for _, instance := range instances {
				instanceTemplates = append(instanceTemplates, instance.Vm.TemplateUUID)
			}
			instanceTemplates = uniqueThisStringArray(instanceTemplates)
			checkAndHandleResetOfGuageVecMetric((len(instances) + len(instanceTemplates)), "anka_instance_state_per_template_count", metric)
			for _, wantedState := range types.InstanceStates {
				for _, wantedInstanceTemplate := range instanceTemplates {
					count := 0
					for _, instance := range instances {
						if instance.Vm.State == wantedState {
							if instance.Vm.TemplateUUID == wantedInstanceTemplate {
								count++
							}
						}
					}
					metric.With(prometheus.Labels{"state": wantedState, "template_uuid": wantedInstanceTemplate}).Set(float64(count))
				}
			}
		},
	},
	InstanceStatePerMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_instance_state_per_group_count", "Count of Instances in a particular state, per Group (label: state, group_name)", []string{"state", "group_uuid"}),
			event:  events.EVENT_VM_DATA_UPDATED,
		},
		HandleData: func(instances []types.Instance, metric *prometheus.GaugeVec) {
			var instanceGroups []string
			for _, instance := range instances {
				if instance.Vm.GroupUUID != "" {
					instanceGroups = append(instanceGroups, instance.Vm.GroupUUID)
				}
			}
			instanceGroups = uniqueThisStringArray(instanceGroups)
			checkAndHandleResetOfGuageVecMetric((len(instances) + len(instanceGroups)), "anka_instance_state_per_group_count", metric)
			for _, wantedState := range types.InstanceStates {
				for _, wantedInstanceGroup := range instanceGroups {
					count := 0
					for _, instance := range instances {
						if instance.Vm.State == wantedState {
							if instance.Vm.GroupUUID == wantedInstanceGroup {
								count++
							}
						}
					}
					metric.With(prometheus.Labels{"state": wantedState, "group_uuid": wantedInstanceGroup}).Set(float64(count))
				}
			}
		},
	},
	InstanceStatePerMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_instance_state_per_node_count", "Count of Instances in a particular state, per Node (label: state, node_uuid)", []string{"state", "node_uuid"}),
			event:  events.EVENT_VM_DATA_UPDATED,
		},
		HandleData: func(instances []types.Instance, metric *prometheus.GaugeVec) {
			var instanceNodes []string
			for _, instance := range instances {
				if instance.Vm.NodeUUID != "" {
					instanceNodes = append(instanceNodes, instance.Vm.NodeUUID)
				}
			}
			instanceNodes = uniqueThisStringArray(instanceNodes)
			checkAndHandleResetOfGuageVecMetric((len(instances) + len(instanceNodes)), "anka_instance_state_per_node_count", metric)
			for _, wantedState := range types.InstanceStates {
				for _, wantedInstanceNode := range instanceNodes {
					count := 0
					for _, instance := range instances {
						if instance.Vm.State == wantedState {
							if instance.Vm.NodeUUID == wantedInstanceNode {
								count++
							}
						}
					}
					metric.With(prometheus.Labels{"state": wantedState, "node_uuid": wantedInstanceNode}).Set(float64(count))
				}
			}
		},
	},
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)
	for _, instanceStatePerMetric := range ankaInstanceStatePerMetrics {
		AddMetric(instanceStatePerMetric)
	}
}
