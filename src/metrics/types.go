package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
)

type AnkaMetric interface {
	GetPrometheusMetric() prometheus.Collector
	GetEvent() events.Event
	GetEventHandler() func(interface{}) error
}

type BaseAnkaMetric struct {
	event  events.Event
	metric prometheus.Collector
}

func (bam BaseAnkaMetric) GetEvent() events.Event {
	return bam.event
}

func (bam BaseAnkaMetric) GetPrometheusMetric() prometheus.Collector {
	return bam.metric
}
