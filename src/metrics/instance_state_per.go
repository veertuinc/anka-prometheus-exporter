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
			metric: CreateGaugeMetricVec("anka_instance_state_per_template_count", "Count of Instances in a particular state, per Template (label: state, template_uuid, template_name)", []string{"state", "template_uuid", "template_name"}),
			event:  events.EVENT_VM_DATA_UPDATED,
		},
		HandleData: func(instances []types.Instance, metric *prometheus.GaugeVec) {
			var InstanceStatePerTemplateCountMap = map[string]map[string]int{}
			var instanceTemplates []string
			var instanceTemplatesMap = map[string]string{}
			for _, instance := range instances {
				instanceTemplates = append(instanceTemplates, instance.Vm.TemplateUUID)
				instanceTemplatesMap[instance.Vm.TemplateUUID] = instance.Vm.TemplateNAME
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
			checkAndHandleResetOfGuageVecMetric((len(instances) + len(instanceTemplates)), "anka_instance_state_per_template_count", metric)
			for wantedState, wantedStateMap := range InstanceStatePerTemplateCountMap {
				for wantedTemplateUUID, count := range wantedStateMap {
					metric.With(prometheus.Labels{"state": wantedState, "template_uuid": wantedTemplateUUID, "template_name": instanceTemplatesMap[wantedTemplateUUID]}).Set(float64(count))
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
			checkAndHandleResetOfGuageVecMetric((len(instances) + len(instanceGroups)), "anka_instance_state_per_group_count", metric)
			for wantedState, wantedStateMap := range InstanceStatePerGroupCountMap {
				for wantedGroupUUID, count := range wantedStateMap {
					metric.With(prometheus.Labels{"state": wantedState, "group_uuid": wantedGroupUUID}).Set(float64(count))
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
			checkAndHandleResetOfGuageVecMetric((len(instances) + len(instanceNodes)), "anka_instance_state_per_node_count", metric)
			for wantedState, wantedStateMap := range InstanceStatePerNodeCountMap {
				for wantedNodeUUID, count := range wantedStateMap {
					metric.With(prometheus.Labels{"state": wantedState, "node_uuid": wantedNodeUUID}).Set(float64(count))
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
