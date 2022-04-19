package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

// TODO: can we make prometheus.GaugeVec support also .Gauge?
type RegistryTemplateMetric struct {
	BaseAnkaMetric
	HandleData func([]types.Template, *prometheus.GaugeVec)
}

func (m RegistryTemplateMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		templates, err := ConvertToRegistryTemplatesData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(m.metric)
		if err != nil {
			return err
		}
		m.HandleData(
			templates,
			metric,
		)
		return nil
	}
}

var ankaRegistryTemplateMetrics = []RegistryTemplateMetric{
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_registry_template_tags_count", "Count of Tags in the Registry for the Template", []string{"template_uuid", "template_name"}),
			event:  events.EventRegistryTemplatesUpdated,
		},
		HandleData: func(templates []types.Template, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGuageVecMetric(len(templates), "anka_registry_template_tags_count", metric)
			for _, template := range templates {
				metric.With(prometheus.Labels{"template_uuid": template.UUID, "template_name": template.Name}).Set(float64(len(template.Tags)))
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_registry_template_disk_used", "Total disk usage of the Template in the Registry", []string{"template_uuid", "template_name"}),
			event:  events.EventRegistryTemplatesUpdated,
		},
		HandleData: func(templates []types.Template, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGuageVecMetric(len(templates), "anka_registry_template_disk_used", metric)
			for _, template := range templates {
				metric.With(prometheus.Labels{"template_uuid": template.UUID, "template_name": template.Name}).Set(float64(template.Size))
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_registry_template_tag_disk_used", "Total disk used by the Template's Tag in the Registry", []string{"template_uuid", "template_name", "tag_name"}),
			event:  events.EventRegistryTemplatesUpdated,
		},
		HandleData: func(templates []types.Template, metric *prometheus.GaugeVec) {
			checkAndHandleResetOfGuageVecMetric(len(templates), "anka_registry_template_tag_disk_used", metric)
			for _, template := range templates {
				for _, tag := range template.Tags {
					metric.With(prometheus.Labels{"template_uuid": template.UUID, "template_name": template.Name, "tag_name": tag.Name}).Set(float64(tag.Size))

				}
			}
		},
	},
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)
	for _, RegistryTemplateMetric := range ankaRegistryTemplateMetrics {
		AddMetric(RegistryTemplateMetric)
	}
}
