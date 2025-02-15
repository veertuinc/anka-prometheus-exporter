package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/log"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type InstanceStatePerMetric struct {
	BaseAnkaMetric
	HandleData func([]types.Instance, *prometheus.GaugeVec)
}

func (ispm InstanceStatePerMetric) GetEventHandler() func(interface{}) error {
	return func(instancesData interface{}) error {
		instances, err := ConvertToInstancesData(instancesData)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(ispm.metric)
		if err != nil {
			return err
		}
		ispm.HandleData(
			instances,
			metric,
		)
		return nil
	}
}

var ankaInstanceStatePerMetrics = []InstanceStatePerMetric{
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_instance_state_per_template_count", "Count of Instances in a particular state, per Template (label: state, template_uuid, template_name)", []string{"state", "template_uuid", "template_name"}),
			event:  events.EVENT_VM_DATA_UPDATED,
		},
		HandleData: func(instances []types.Instance, metric *prometheus.GaugeVec) {
			var InstanceStatePerTemplateCountMap = map[string]map[string]int{}
			var instanceTemplates []string
			var instanceTemplatesMap = map[string]string{}
			for _, instance := range instances {
				instanceTemplates = append(instanceTemplates, instance.Vm.TemplateUUID)
				instanceTemplatesMap[instance.Vm.TemplateUUID] = instance.Vm.TemplateName
			}
			instanceTemplates = uniqueThisStringArray(instanceTemplates)
			for _, wantedState := range types.InstanceStates {
				if _, ok := InstanceStatePerTemplateCountMap[wantedState]; !ok {
					InstanceStatePerTemplateCountMap[wantedState] = make(map[string]int)
				}
				for _, wantedInstanceTemplate := range instanceTemplates {
					count := 0
					for _, instance := range instances {
						if instance.Vm.State == wantedState {
							if instance.Vm.TemplateUUID == wantedInstanceTemplate {
								count++
							}
						}
					}
					if _, ok := InstanceStatePerTemplateCountMap[wantedState][wantedInstanceTemplate]; !ok {
						InstanceStatePerTemplateCountMap[wantedState][wantedInstanceTemplate] = count
					}
				}
			}
			checkAndHandleResetOfGaugeVecMetric((len(instances) + len(instanceTemplates)), "anka_instance_state_per_template_count", metric)
			for wantedState, wantedStateMap := range InstanceStatePerTemplateCountMap {
				for wantedTemplateUUID, count := range wantedStateMap {
					metric.With(prometheus.Labels{"state": wantedState, "template_uuid": wantedTemplateUUID, "template_name": instanceTemplatesMap[wantedTemplateUUID]}).Set(float64(count))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_instance_state_per_group_count", "Count of Instances in a particular state, per Group (label: state, group_name)", []string{"state", "group_uuid"}),
			event:  events.EVENT_VM_DATA_UPDATED,
		},
		HandleData: func(instances []types.Instance, metric *prometheus.GaugeVec) {
			var InstanceStatePerGroupCountMap = map[string]map[string]int{}
			var instanceGroups []string
			for _, instance := range instances {
				if instance.Vm.GroupUUID != "" {
					instanceGroups = append(instanceGroups, instance.Vm.GroupUUID)
				}
			}
			instanceGroups = uniqueThisStringArray(instanceGroups)
			for _, wantedState := range types.InstanceStates {
				if _, ok := InstanceStatePerGroupCountMap[wantedState]; !ok {
					InstanceStatePerGroupCountMap[wantedState] = make(map[string]int)
				}
				for _, wantedInstanceGroup := range instanceGroups {
					count := 0
					for _, instance := range instances {
						if instance.Vm.State == wantedState {
							if instance.Vm.GroupUUID == wantedInstanceGroup {
								count++
							}
						}
					}
					if _, ok := InstanceStatePerGroupCountMap[wantedState][wantedInstanceGroup]; !ok {
						InstanceStatePerGroupCountMap[wantedState][wantedInstanceGroup] = count
					}
				}
			}
			checkAndHandleResetOfGaugeVecMetric((len(instances) + len(instanceGroups)), "anka_instance_state_per_group_count", metric)
			for wantedState, wantedStateMap := range InstanceStatePerGroupCountMap {
				for wantedGroupUUID, count := range wantedStateMap {
					metric.With(prometheus.Labels{"state": wantedState, "group_uuid": wantedGroupUUID}).Set(float64(count))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_instance_state_per_node_count", "Count of Instances in a particular state, per Node (label: state, node_uuid)", []string{"state", "node_uuid"}),
			event:  events.EVENT_VM_DATA_UPDATED,
		},
		HandleData: func(instances []types.Instance, metric *prometheus.GaugeVec) {
			var InstanceStatePerNodeCountMap = map[string]map[string]int{}
			var instanceNodes []string
			for _, instance := range instances {
				if instance.Vm.NodeUUID != "" {
					instanceNodes = append(instanceNodes, instance.Vm.NodeUUID)
				}
			}
			instanceNodes = uniqueThisStringArray(instanceNodes)
			for _, wantedState := range types.InstanceStates {
				if _, ok := InstanceStatePerNodeCountMap[wantedState]; !ok {
					InstanceStatePerNodeCountMap[wantedState] = make(map[string]int)
				}
				for _, wantedInstanceNode := range instanceNodes {
					count := 0
					for _, instance := range instances {
						if instance.Vm.State == wantedState {
							if instance.Vm.NodeUUID == wantedInstanceNode {
								count++
							}
						}
					}
					if _, ok := InstanceStatePerNodeCountMap[wantedState][wantedInstanceNode]; !ok {
						InstanceStatePerNodeCountMap[wantedState][wantedInstanceNode] = count
					}
				}
			}
			checkAndHandleResetOfGaugeVecMetric((len(instances) + len(instanceNodes)), "anka_instance_state_per_node_count", metric)
			for wantedState, wantedStateMap := range InstanceStatePerNodeCountMap {
				for wantedNodeUUID, count := range wantedStateMap {
					metric.With(prometheus.Labels{"state": wantedState, "node_uuid": wantedNodeUUID}).Set(float64(count))
				}
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_instance_max_age_per_template_seconds", "Age of oldest Instance in a particular state, per Template. Visible only for templates with at least one instance (label: state, template_uuid, template_name)", []string{"state", "template_uuid", "template_name"}),
			event:  events.EVENT_VM_DATA_UPDATED,
		},
		HandleData: func(instances []types.Instance, metric *prometheus.GaugeVec) {
			var InstanceAgePerTemplateMaximumMap = map[string]map[string]int{}
			var instanceTemplates []string
			var instanceTemplatesMap = map[string]string{}
			now := time.Now()
			for _, instance := range instances {
				instanceTemplates = append(instanceTemplates, instance.Vm.TemplateUUID)
				instanceTemplatesMap[instance.Vm.TemplateUUID] = instance.Vm.TemplateName
			}
			instanceTemplates = uniqueThisStringArray(instanceTemplates)
			for _, wantedState := range types.InstanceStates {
				if _, ok := InstanceAgePerTemplateMaximumMap[wantedState]; !ok {
					InstanceAgePerTemplateMaximumMap[wantedState] = make(map[string]int)
				}
				for _, wantedInstanceTemplate := range instanceTemplates {
					age := 0.0
					for _, instance := range instances {
						if instance.Vm.State == wantedState {
							if instance.Vm.TemplateUUID == wantedInstanceTemplate {
								var instanceTime time.Time
								var err error
								// cr_time only set on an instance creation (in the DB) and never changed
								// ts gets updated from time to time due to different events, like save image, termination etc
								// both ts and cr_time are members of an Instance object, they do not depend on the vm (regardless if the vm has started or not)
								if instance.Vm.State != "Started" && instance.Vm.State != "Scheduling" { // can't use CreationTime because it doesn't change after Scheduling happens
									instanceTime, err = time.Parse(time.RFC3339, instance.Vm.LastUpdateTime)
									if err != nil {
										log.Error(fmt.Sprintf("Error parsing LastUpdateTime %s for template %s: %s", instance.Vm.LastUpdateTime, wantedInstanceTemplate, err.Error()))
									}
								} else {
									instanceTime, err = time.Parse(time.RFC3339, instance.Vm.CreationTime)
									if err != nil {
										log.Error(fmt.Sprintf("Error parsing CreationTime %s for template %s: %s", instance.Vm.CreationTime, wantedInstanceTemplate, err.Error()))
									}
								}
								thisAge := now.Sub(instanceTime).Seconds()
								age = max(age, thisAge)
							}
						}
					}
					if _, ok := InstanceAgePerTemplateMaximumMap[wantedState][wantedInstanceTemplate]; !ok {
						InstanceAgePerTemplateMaximumMap[wantedState][wantedInstanceTemplate] = int(age)
					}
				}
			}
			checkAndHandleResetOfGaugeVecMetric((len(instances) + len(instanceTemplates)), "anka_instance_max_age_per_template_seconds", metric)
			for wantedState, wantedStateMap := range InstanceAgePerTemplateMaximumMap {
				for wantedTemplateUUID, age := range wantedStateMap {
					metric.With(prometheus.Labels{"state": wantedState, "template_uuid": wantedTemplateUUID, "template_name": instanceTemplatesMap[wantedTemplateUUID]}).Set(float64(age))
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
