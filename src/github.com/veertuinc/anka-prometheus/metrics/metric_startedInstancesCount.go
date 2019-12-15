package metrics

import (
	"github.com/veertuinc/anka-prometheus/events"
	"github.com/veertuinc/anka-prometheus/types"
)

type StartedInstanceMetric struct {
	BaseAnkaMetric
}

func (this StartedInstanceMetric) GetEventHandler() func(interface{}) error {
	return func(d interface{}) error {
		data, err := ConvertToInstancesData(d)
		if err != nil {
			return err
		}
		metric, err := ConvertMetricToGauge(this.metric)
		if err != nil {
			return err
		}
		counter := 0
		for _, instanceData := range data {
			if instanceData.Vm.State == types.StateStarted {
				counter++
			}
		}
		metric.Set(float64(counter))
		return nil
	}
}

func init() {
	metricName := "anka_started_instances_count"
	metricDesc := "Number of instances in started state"
	ankaMetric := StartedInstanceMetric{}
	ankaMetric.metric = CreateGaugeMetric(metricName, metricDesc)
	ankaMetric.event = events.EVENT_VM_DATA_UPDATED
	
	AddMetric(ankaMetric)
}