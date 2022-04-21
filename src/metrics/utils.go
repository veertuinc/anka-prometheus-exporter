package metrics

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

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

func intMapFromStringSlice(stringSlice []string) map[string]int {
	intMap := map[string]int{}
	for _, item := range stringSlice {
		intMap[item] = 0
	}
	return intMap
}

func uniqueThisStringArray(arr []string) []string {
	occurred := map[string]bool{}
	var result []string
	for e := range arr {
		if !occurred[arr[e]] {
			occurred[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
}

func uniqueNodeGroupsArray(arr []types.NodeGroup) []types.NodeGroup {
	occurred := map[types.NodeGroup]bool{}
	var result []types.NodeGroup
	for e := range arr {
		if !occurred[arr[e]] {
			occurred[arr[e]] = true
			result = append(result, arr[e])
		}
	}
	return result
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
		return nil, fmt.Errorf("could not convert incoming data to required node information. original data: %v", d)
	}
	return data, nil
}

func ConvertToRegistryDiskData(d interface{}) (*types.RegistryDisk, error) {
	data, ok := d.(types.RegistryDisk)
	if !ok {
		return nil, fmt.Errorf("could not convert incoming data to required registry disk information. original data: %v", d)
	}
	return &data, nil
}

func ConvertToRegistryTemplatesData(d interface{}) ([]types.Template, error) {
	data, ok := d.([]types.Template)
	if !ok {
		return nil, fmt.Errorf("could not convert incoming data to required registry template information. original data: %v", d)
	}
	return data, nil
}

func ConvertToInstancesData(d interface{}) ([]types.Instance, error) {
	data, ok := d.([]types.Instance)
	if !ok {
		return nil, fmt.Errorf("could not convert incoming data to required instances information. original data: %v", d)
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

const DefaultRefreshMetricSeconds = 600

var mutex = &sync.Mutex{}
var lastRefreshTime int64
var metricCountMap = make(map[string]int)

func checkAndHandleReset(count int, metricName string) bool {
	mutex.Lock()
	val, ok := metricCountMap[metricName]
	mutex.Unlock()
	if ok { // Check if key exists
		// Check if the count has changed OR if the lastRefreshTime is Greater Than or Eq to 10 minutes
		if val != 0 && (count != val || (time.Now().Unix()-atomic.LoadInt64(&lastRefreshTime)) >= DefaultRefreshMetricSeconds) {
			mutex.Lock()
			delete(metricCountMap, metricName)
			mutex.Unlock()
			atomic.StoreInt64(&lastRefreshTime, time.Now().Unix())
			return true
		}
	}
	mutex.Lock()
	metricCountMap[metricName] = count
	mutex.Unlock()
	return false
}

func checkAndHandleResetOfGuageVecMetric(count int, metricName string, metric *prometheus.GaugeVec) {
	if ok := checkAndHandleReset(count, metricName); ok {
		metric.Reset()
	}
}
