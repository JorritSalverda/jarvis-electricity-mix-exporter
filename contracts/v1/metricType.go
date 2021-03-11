package contracts

type MetricType string

const (
	MetricType_METRIC_TYPE_INVALID MetricType = ""
	MetricType_METRIC_TYPE_COUNTER MetricType = "METRIC_TYPE_COUNTER"
	MetricType_METRIC_TYPE_GAUGE   MetricType = "METRIC_TYPE_GAUGE"
)
