package contracts

type MetricType string

const (
	MetricTypeUnknown MetricType = ""
	MetricTypeCounter MetricType = "METRIC_TYPE_COUNTER"
	MetricTypeGauge   MetricType = "METRIC_TYPE_GAUGE"
)
