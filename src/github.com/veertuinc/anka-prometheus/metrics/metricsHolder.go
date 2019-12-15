package metrics

var MetricsHolder []AnkaMetric

func AddMetric(m AnkaMetric) {
	if MetricsHolder == nil {
		MetricsHolder = make([]AnkaMetric, 0)
	}

	MetricsHolder = append(MetricsHolder, m)
}