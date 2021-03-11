package contracts

type MetricType string

const (
	MetricTypeUnknown MetricType = "Unknown"
	MetricTypeCounter MetricType = "Counter"
	MetricTypeGauge   MetricType = "Gauge"
)
