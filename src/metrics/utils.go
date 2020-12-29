package metrics

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func uniqueThisStringArray(arr []string) []string {
	occured := map[string]bool{}
	result := []string{}
	for e := range arr {
		if occured[arr[e]] != true {
			occured[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}

func uniqueNodeGroupsArray(arr []types.NodeGroup) []types.NodeGroup {
	occured := map[types.NodeGroup]bool{}
	result := []types.NodeGroup{}
	for e := range arr {
		if occured[arr[e]] != true {
			occured[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}

func CountNodeGroupNodes(GroupIdWeWant string, nodesData []types.Node) int {
	counter := 0
	for _, node := range nodesData {
		for _, group := range node.Groups {
			if group.Id == GroupIdWeWant {
				counter++
			}
		}
	}
	return counter
}

func CountNodeGroupState(groupIdWeWant string, stateWeWant string, nodesData []types.Node) int {
	counter := 0
	for _, node := range nodesData {
		if node.State == stateWeWant {
			for _, group := range node.Groups {
				if group.Id == groupIdWeWant {
					counter++
				}
			}
		}
	}
	return counter
}

func CountVMState(checkForState string, data []types.InstanceInfo) int {
	counter := 0
	for _, instanceData := range data {
		if instanceData.Vm.State == checkForState {
			counter++
		}
	}
	return counter
}

func CountInstanceTemplateState(templateWeWant string, stateWeWant string, data []types.InstanceInfo) int {
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

func CountInstanceGroupState(groupWeWant string, stateWeWant string, data []types.InstanceInfo) int {
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

func CreateGaugeMetric(name string, help string) prometheus.Gauge {
	m := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		})
	return m
}

func CreateGaugeMetricVec(name string, help string, labels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}, labels)
}

func ConvertToNodeData(d interface{}) ([]types.Node, error) {
	data, ok := d.([]types.Node)
	if !ok {
		return nil, fmt.Errorf("could not convert incoming data to required node information. original data: ", d)
	}
	return data, nil
}

func ConvertToRegistryData(d interface{}) (*types.RegistryInfo, error) {
	data, ok := d.(types.RegistryInfo)
	if !ok {
		return nil, fmt.Errorf("could not convert incoming data to required registry information. original data: ", d)
	}
	return &data, nil
}

func ConvertToInstancesData(d interface{}) ([]types.InstanceInfo, error) {
	data, ok := d.([]types.InstanceInfo)
	if !ok {
		return nil, fmt.Errorf("could not convert incoming data to required instances information. original data: ", d)
	}
	return data, nil
}

func ConvertMetricToGauge(m prometheus.Collector) (prometheus.Gauge, error) {
	data, ok := m.(prometheus.Gauge)
	if !ok {
		return nil, fmt.Errorf("could not convert metric to gauge type")
	}
	return data, nil
}

func ConvertMetricToGaugeVec(m prometheus.Collector) (*prometheus.GaugeVec, error) {
	data, ok := m.(*prometheus.GaugeVec)
	if !ok {
		return nil, fmt.Errorf("could not convert metric to gauge vector type")
	}
	return data, nil
}
