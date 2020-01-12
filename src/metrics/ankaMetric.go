package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus/events"	
)

type AnkaMetric interface {
	GetPrometheusMetric() prometheus.Collector
	GetEvent() events.Event
	GetEventHandler() func(interface{}) error
}

type BaseAnkaMetric struct {
	event	events.Event
	metric	prometheus.Collector
}

func (this BaseAnkaMetric) GetEvent() events.Event {
	return this.event
}

func (this BaseAnkaMetric) GetPrometheusMetric() prometheus.Collector {
	return this.metric
}