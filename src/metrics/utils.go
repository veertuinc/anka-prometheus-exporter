package metrics

import (
	"github.com/veertuinc/anka-prometheus/types"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)


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

func ConvertToNodeData(d interface{}) ([]types.NodeInfo, error) {
	data, ok := d.([]types.NodeInfo)
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