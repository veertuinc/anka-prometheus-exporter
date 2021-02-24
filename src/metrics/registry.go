package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type RegistryMetric struct {
	BaseAnkaMetric
	HandleData func(*types.Registry, prometheus.Gauge)
}

func (this RegistryMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		registryData, err := ConvertToRegistryData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGauge(this.metric)
		if err != nil {
			return err
		}
		this.HandleData(
			registryData,
			metric,
		)
		return nil
	}
}

var ankaRegistryMetrics = []RegistryMetric{
	RegistryMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_registry_disk_free_space", "Anka Build Cloud Registry free disk space"),
			event:  events.EVENT_REGISTRY_DATA_UPDATED,
		},
		HandleData: func(registry *types.Registry, metric prometheus.Gauge) {
			metric.Set(float64(registry.Free))
		},
	},
	RegistryMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_registry_disk_total_space", "Anka Build Cloud Registry total disk size"),
			event:  events.EVENT_REGISTRY_DATA_UPDATED,
		},
		HandleData: func(registry *types.Registry, metric prometheus.Gauge) {
			metric.Set(float64(registry.Total))
		},
	},
	RegistryMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_registry_disk_used_space", "Anka Build Cloud Registry used disk space"),
			event:  events.EVENT_REGISTRY_DATA_UPDATED,
		},
		HandleData: func(registry *types.Registry, metric prometheus.Gauge) {
			var used uint64 = 0
			used = registry.Total - registry.Free
			metric.Set(float64(used))
		},
	},
}

func init() {
	for _, registryMetric := range ankaRegistryMetrics {
		AddMetric(registryMetric)
	}
}
