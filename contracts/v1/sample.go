package contracts

type Sample struct {
	EnergyType      EnergyType
	IsRenewable     bool
	MetricType      MetricType
	SampleDirection SampleDirection
	SampleUnit      SampleUnit
	Value           float64
}
