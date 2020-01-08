package metrics

import (
	"github.com/veertuinc/anka-prometheus/events"
	"math"
)

type RegTotalSpaceMetric struct {
	BaseAnkaMetric
}

func (this RegTotalSpaceMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		data, err := ConvertToRegistryData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGauge(this.metric)
		if err != nil {
			return err
		}
		value := float64(data.Total) / 1024 / 1024 / 1024
		metric.Set(math.Round(value*100)/100)
		return nil
	}
}

func init() {
	metricName := "anka_registry_total_space_gb"
	metricDesc := "Total space in Registry"
	ankaMetric := RegTotalSpaceMetric{}
	ankaMetric.metric = CreateGaugeMetric(metricName, metricDesc)
	ankaMetric.event = events.EVENT_REGISTRY_DATA_UPDATED
	
	AddMetric(ankaMetric)
}