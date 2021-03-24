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

type HandleDataResult struct {
	Count      int
	MetricName string
}

func (this BaseAnkaMetric) GetEvent() events.Event {
	return this.event
}

func (this BaseAnkaMetric) GetPrometheusMetric() prometheus.Collector {
	return this.metric
}
