package metrics

import (
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/veertuinc/anka-prometheus-exporter/src/types"
)

func histogramSampleCount(t *testing.T, vec *prometheus.HistogramVec, labels prometheus.Labels) uint64 {
	t.Helper()
	metric := &dto.Metric{}
	if err := vec.With(labels).(prometheus.Metric).Write(metric); err != nil {
		t.Fatalf("Write histogram metric: %v", err)
	}
	return metric.GetHistogram().GetSampleCount()
}

func histogramSampleSum(t *testing.T, vec *prometheus.HistogramVec, labels prometheus.Labels) float64 {
	t.Helper()
	metric := &dto.Metric{}
	if err := vec.With(labels).(prometheus.Metric).Write(metric); err != nil {
		t.Fatalf("Write histogram metric: %v", err)
	}
	return metric.GetHistogram().GetSampleSum()
}

func counterValue(t *testing.T, vec *prometheus.CounterVec, labels prometheus.Labels) float64 {
	t.Helper()
	metric := &dto.Metric{}
	if err := vec.With(labels).(prometheus.Metric).Write(metric); err != nil {
		t.Fatalf("Write counter metric: %v", err)
	}
	return metric.GetCounter().GetValue()
}

func TestInstanceStartDuration_EmptyIPNoObserve(t *testing.T) {
	m := newInstanceStartDurationMetric()
	fixedNow := time.Date(2024, 1, 24, 16, 45, 30, 0, time.UTC)
	m.now = func() time.Time { return fixedNow }

	labels := prometheus.Labels{
		"template_uuid": "tmpl-1",
		"template_name": "macos",
		"arch":          "arm64",
		"external_id":   "job-1",
	}

	m.processInstances([]types.Instance{{
		InstanceID: "inst-1",
		ExternalID: "job-1",
		Vm: types.VmData{
			State:        "Started",
			TemplateUUID: "tmpl-1",
			TemplateName: "macos",
			Arch:         "arm64",
			CreationTime: "2024-01-24T16:44:30Z",
			VmInfo:       types.VmInfo{IP: ""},
		},
	}})

	if got := histogramSampleCount(t, m.duration, labels); got != 0 {
		t.Errorf("expected 0 observations, got %d", got)
	}
}

func TestInstanceStartDuration_FirstIPObservesOnce(t *testing.T) {
	m := newInstanceStartDurationMetric()
	fixedNow := time.Date(2024, 1, 24, 16, 45, 30, 0, time.UTC)
	m.now = func() time.Time { return fixedNow }

	labels := prometheus.Labels{
		"template_uuid": "tmpl-1",
		"template_name": "macos",
		"arch":          "arm64",
		"external_id":   "job-1",
	}
	instance := types.Instance{
		InstanceID: "inst-1",
		ExternalID: "job-1",
		Vm: types.VmData{
			State:        "Started",
			TemplateUUID: "tmpl-1",
			TemplateName: "macos",
			Arch:         "arm64",
			CreationTime: "2024-01-24T16:44:30Z",
			VmInfo:       types.VmInfo{IP: "192.168.64.4"},
		},
	}

	m.processInstances([]types.Instance{instance})
	if got := histogramSampleCount(t, m.duration, labels); got != 1 {
		t.Fatalf("expected 1 observation, got %d", got)
	}
	if sum := histogramSampleSum(t, m.duration, labels); sum != 60 {
		t.Errorf("expected sum 60s, got %v", sum)
	}

	m.processInstances([]types.Instance{instance})
	if got := histogramSampleCount(t, m.duration, labels); got != 1 {
		t.Errorf("expected still 1 observation after second poll, got %d", got)
	}
}

func TestInstanceStartDuration_BadCreationTimeSkipped(t *testing.T) {
	m := newInstanceStartDurationMetric()
	m.now = func() time.Time { return time.Date(2024, 1, 24, 16, 45, 30, 0, time.UTC) }

	labels := prometheus.Labels{
		"template_uuid": "tmpl-1",
		"template_name": "macos",
		"arch":          "arm64",
		"external_id":   "",
	}

	m.processInstances([]types.Instance{{
		InstanceID: "inst-bad",
		Vm: types.VmData{
			State:        "Started",
			TemplateUUID: "tmpl-1",
			TemplateName: "macos",
			Arch:         "arm64",
			CreationTime: "not-a-time",
			VmInfo:       types.VmInfo{IP: "10.0.0.1"},
		},
	}})

	if got := histogramSampleCount(t, m.duration, labels); got != 0 {
		t.Errorf("expected 0 observations for bad cr_time, got %d", got)
	}
}

func TestInstanceStartDuration_ErrorWithoutIPIncrementsFailureOnce(t *testing.T) {
	m := newInstanceStartDurationMetric()
	m.now = func() time.Time { return time.Date(2024, 1, 24, 16, 45, 30, 0, time.UTC) }

	labels := prometheus.Labels{
		"template_uuid": "tmpl-1",
		"template_name": "macos",
		"arch":          "arm64",
		"external_id":   "job-fail",
	}
	instance := types.Instance{
		InstanceID: "inst-err",
		ExternalID: "job-fail",
		Vm: types.VmData{
			State:        "Error",
			TemplateUUID: "tmpl-1",
			TemplateName: "macos",
			Arch:         "arm64",
			CreationTime: "2024-01-24T16:44:30Z",
			VmInfo:       types.VmInfo{IP: ""},
		},
	}

	m.processInstances([]types.Instance{instance})
	if got := counterValue(t, m.failures, labels); got != 1 {
		t.Fatalf("expected failure count 1, got %v", got)
	}

	m.processInstances([]types.Instance{instance})
	if got := counterValue(t, m.failures, labels); got != 1 {
		t.Errorf("expected failure count still 1, got %v", got)
	}
	if got := histogramSampleCount(t, m.duration, labels); got != 0 {
		t.Errorf("expected no duration observations on failure, got %d", got)
	}
}

func TestInstanceStartDuration_PruneAllowsRecycledInstanceID(t *testing.T) {
	m := newInstanceStartDurationMetric()
	fixedNow := time.Date(2024, 1, 24, 16, 45, 30, 0, time.UTC)
	m.now = func() time.Time { return fixedNow }

	labels := prometheus.Labels{
		"template_uuid": "tmpl-1",
		"template_name": "macos",
		"arch":          "arm64",
		"external_id":   "job-1",
	}
	instance := types.Instance{
		InstanceID: "inst-1",
		ExternalID: "job-1",
		Vm: types.VmData{
			State:        "Started",
			TemplateUUID: "tmpl-1",
			TemplateName: "macos",
			Arch:         "arm64",
			CreationTime: "2024-01-24T16:44:30Z",
			VmInfo:       types.VmInfo{IP: "192.168.64.4"},
		},
	}

	m.processInstances([]types.Instance{instance})
	if got := histogramSampleCount(t, m.duration, labels); got != 1 {
		t.Fatalf("expected 1 observation, got %d", got)
	}

	// Instance disappears from snapshot — prune tracking.
	m.processInstances([]types.Instance{})

	// Same instance_id returns later — should observe again.
	m.processInstances([]types.Instance{instance})
	if got := histogramSampleCount(t, m.duration, labels); got != 2 {
		t.Errorf("expected 2 observations after recycle, got %d", got)
	}
}

func TestInstanceStartDuration_EventHandlerConvertsData(t *testing.T) {
	m := newInstanceStartDurationMetric()
	m.now = func() time.Time { return time.Date(2024, 1, 24, 16, 45, 30, 0, time.UTC) }

	err := m.GetEventHandler()([]types.Instance{{
		InstanceID: "inst-1",
		ExternalID: "ext",
		Vm: types.VmData{
			State:        "Started",
			TemplateUUID: "tmpl-1",
			TemplateName: "macos",
			Arch:         "arm64",
			CreationTime: "2024-01-24T16:44:30Z",
			VmInfo:       types.VmInfo{IP: "192.168.1.1"},
		},
	}})
	if err != nil {
		t.Fatalf("handler error: %v", err)
	}

	labels := prometheus.Labels{
		"template_uuid": "tmpl-1",
		"template_name": "macos",
		"arch":          "arm64",
		"external_id":   "ext",
	}
	if got := histogramSampleCount(t, m.duration, labels); got != 1 {
		t.Errorf("expected 1 observation via handler, got %d", got)
	}
}
