package api

type Sample struct {
	EnergyType         EnergyType
	OriginalEnergyType string
	IsRenewable        bool
	MetricType         MetricType
	Resolution         string
	SampleDirection    SampleDirection
	SampleUnit         SampleUnit
	Value              float64
}
