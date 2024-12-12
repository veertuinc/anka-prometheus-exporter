package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

type StatusMetric struct {
	BaseAnkaMetric
	HandleData func(*types.Status, *prometheus.GaugeVec)
}

func (sm StatusMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		status, err := ConvertToStatusData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGaugeVec(sm.metric)
		if err != nil {
			return err
		}
		sm.HandleData(
			status,
			metric,
		)
		return nil
	}
}

var ankaStatusMetrics = []StatusMetric{
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_controller_state_count", "Status of the Anka Controller", []string{"state"}),
			event:  events.EVENT_STATUS_UPDATED,
		},
		HandleData: func(status *types.Status, metric *prometheus.GaugeVec) {
			for _, state := range types.ControllerStates {
				counter := 0
				if status.Status == state {
					counter++
				}
				metric.With(prometheus.Labels{"state": state}).Set(float64(counter))
			}
		},
	},
	{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: CreateGaugeMetricVec("anka_registry_state_count", "Status of the Anka Registry", []string{"state"}),
			event:  events.EVENT_STATUS_UPDATED,
		},
		HandleData: func(status *types.Status, metric *prometheus.GaugeVec) {
			for _, state := range types.RegistryStates {
				counter := 0
				if status.RegistryStatus == state {
					counter++
				}
				metric.With(prometheus.Labels{"state": state}).Set(float64(counter))
			}
		},
	},
}

func init() { // runs on exporter init only (updates are made with the above EventHandler; triggered by the Client)
	for _, statusMetric := range ankaStatusMetrics {
		AddMetric(statusMetric)
	}
}
