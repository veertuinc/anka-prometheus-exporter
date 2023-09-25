package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

// TODO: can we make prometheus.GaugeVec support also .Gauge?
type RegistryTemplatesMetric struct {
	BaseAnkaMetric
	HandleData func([]types.Template, prometheus.Gauge)
}

func (rtm RegistryTemplatesMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		templates, err := ConvertToRegistryTemplatesData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGauge(rtm.metric)
		if err != nil {
			return err
		}
		rtm.HandleData(
			templates,
			metric,
		)
		return nil
	}
}

var ankaRegistryTemplatesMetrics = []RegistryTemplatesMetric{
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetric("anka_registry_template_count", "Count of VM Templates in the Registry"),
			event:  events.EVENT_REGISTRY_TEMPLATES_UPDATED,
		},
		HandleData: func(templates []types.Template, metric prometheus.Gauge) {
			metric.Set(float64(len(templates)))
		},
	},
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)
	for _, vmRegistryTemplateMetric := range ankaRegistryTemplatesMetrics {
		AddMetric(vmRegistryTemplateMetric)
	}
}
