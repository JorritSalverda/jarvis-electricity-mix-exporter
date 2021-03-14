package api

type Sample struct {
	EnergyType         EnergyType
	OriginalEnergyType string
	IsRenewable        bool
	MetricType         MetricType
	SampleDirection    SampleDirection
	SampleUnit         SampleUnit
	Value              float64
}
