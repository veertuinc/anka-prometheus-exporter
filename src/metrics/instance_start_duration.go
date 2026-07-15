package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/veertuinc/anka-prometheus-exporter/src/events"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

var instanceStartDurationBuckets = []float64{30, 60, 120, 300, 600, 900, 1800}

var instanceStartMetricLabels = []string{"template_uuid", "template_name", "arch", "external_id"}

type multiCollector struct {
	collectors []prometheus.Collector
}

func (m multiCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, c := range m.collectors {
		c.Describe(ch)
	}
}

func (m multiCollector) Collect(ch chan<- prometheus.Metric) {
	for _, c := range m.collectors {
		c.Collect(ch)
	}
}

type InstanceStartDurationMetric struct {
	BaseAnkaMetric
	duration         *prometheus.HistogramVec
	failures         *prometheus.CounterVec
	mu               sync.Mutex
	observedSuccess  map[string]bool
	observedFailure  map[string]bool
	now              func() time.Time
}

func (m *InstanceStartDurationMetric) GetEventHandler() func(interface{}) error {
	return func(instancesData interface{}) error {
		instances, err := ConvertToInstancesData(instancesData)
		if err != nil {
			return err
		}
		m.processInstances(instances)
		return nil
	}
}

func (m *InstanceStartDurationMetric) processInstances(instances []types.Instance) {
	nowFn := m.now
	if nowFn == nil {
		nowFn = time.Now
	}
	now := nowFn()

	currentIDs := make(map[string]struct{}, len(instances))

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, instance := range instances {
		id := instance.InstanceID
		if id == "" {
			continue
		}
		currentIDs[id] = struct{}{}

		labels := prometheus.Labels{
			"template_uuid": instance.Vm.TemplateUUID,
			"template_name": instance.Vm.TemplateName,
			"arch":          instance.Vm.Arch,
			"external_id":   instance.ExternalID,
		}

		if instance.Vm.VmInfo.IP != "" {
			if m.observedSuccess[id] {
				continue
			}
			createdAt, err := time.Parse(time.RFC3339, instance.Vm.CreationTime)
			if err != nil {
				continue
			}
			m.duration.With(labels).Observe(now.Sub(createdAt).Seconds())
			m.observedSuccess[id] = true
			continue
		}

		if instance.Vm.State == "Error" && !m.observedSuccess[id] && !m.observedFailure[id] {
			m.failures.With(labels).Inc()
			m.observedFailure[id] = true
		}
	}

	for id := range m.observedSuccess {
		if _, ok := currentIDs[id]; !ok {
			delete(m.observedSuccess, id)
		}
	}
	for id := range m.observedFailure {
		if _, ok := currentIDs[id]; !ok {
			delete(m.observedFailure, id)
		}
	}
}

func newInstanceStartDurationMetric() *InstanceStartDurationMetric {
	duration := CreateHistogramMetricVec(
		"anka_instance_start_duration_seconds",
		"Seconds from instance creation (cr_time) to first observed non-empty vminfo.ip (labels: template_uuid, template_name, arch, external_id). Poll interval can add measurement noise.",
		instanceStartMetricLabels,
		instanceStartDurationBuckets,
	)
	failures := CreateCounterMetricVec(
		"anka_instance_start_failures_total",
		"Count of instances that reached Error before a non-empty vminfo.ip was observed (labels: template_uuid, template_name, arch, external_id)",
		instanceStartMetricLabels,
	)
	return &InstanceStartDurationMetric{
		BaseAnkaMetric: BaseAnkaMetric{
			metric: multiCollector{collectors: []prometheus.Collector{duration, failures}},
			event:  events.EVENT_VM_DATA_UPDATED,
		},
		duration:        duration,
		failures:        failures,
		observedSuccess: make(map[string]bool),
		observedFailure: make(map[string]bool),
		now:             time.Now,
	}
}

func init() {
	AddMetric(newInstanceStartDurationMetric())
}
