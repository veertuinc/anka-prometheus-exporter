package metrics

import (
	"github.com/veertuinc/anka-prometheus/src/events"
)

type RegistryMetric struct {
	BaseAnkaMetric
}

func (this RegistryMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		data, err := ConvertToRegistryData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGauge(this.metric)
		if err != nil {
			return err
		}
		if this.name == "anka_registry_disk_free_space" {
			metric.Set(float64(data.Free))
		} else if this.name == "anka_registry_disk_used_space" {
			metric.Set(float64(data.Total))
		}
		return nil
	}
}

func init() {
	AddMetric(RegistryMetric{BaseAnkaMetric{
		name:   "anka_registry_disk_free_space",
		metric: CreateGaugeMetric("anka_registry_disk_free_space", "Anka Build Cloud Registry free disk space"),
		event:  events.EVENT_REGISTRY_DATA_UPDATED,
	}})
	AddMetric(RegistryMetric{BaseAnkaMetric{
		name:   "anka_registry_disk_used_space",
		metric: CreateGaugeMetric("anka_registry_disk_used_space", "Anka Build Cloud Registry used disk space"),
		event:  events.EVENT_REGISTRY_DATA_UPDATED,
	}})
}
