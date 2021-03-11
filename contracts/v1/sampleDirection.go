package contracts

type SampleDirection string

const (
	SampleDirectionUnknown                              SampleDirection = ""
	SampleDirection_SAMPLE_TYPE_ELECTRICITY_CONSUMPTION SampleDirection = "SAMPLE_TYPE_ELECTRICITY_CONSUMPTION"
	SampleDirection_SAMPLE_TYPE_ELECTRICITY_PRODUCTION  SampleDirection = "SAMPLE_TYPE_ELECTRICITY_PRODUCTION"
	SampleDirection_SAMPLE_TYPE_ENERGY                  SampleDirection = "SAMPLE_TYPE_ENERGY"
	SampleDirection_SAMPLE_TYPE_GAS                     SampleDirection = "SAMPLE_TYPE_GAS"
	SampleDirection_SAMPLE_TYPE_TEMPERATURE             SampleDirection = "SAMPLE_TYPE_TEMPERATURE"
	SampleDirection_SAMPLE_TYPE_TEMPERATURE_SETPOINT    SampleDirection = "SAMPLE_TYPE_TEMPERATURE_SETPOINT"
	SampleDirection_SAMPLE_TYPE_PRESSURE                SampleDirection = "SAMPLE_TYPE_PRESSURE"
	SampleDirection_SAMPLE_TYPE_FLOW                    SampleDirection = "SAMPLE_TYPE_FLOW"
	SampleDirection_SAMPLE_TYPE_HUMIDITY                SampleDirection = "SAMPLE_TYPE_HUMIDITY"
	SampleDirection_SAMPLE_TYPE_TIME                    SampleDirection = "SAMPLE_TYPE_TIME"
	SampleDirection_SAMPLE_TYPE_HEAT_DEMAND             SampleDirection = "SAMPLE_TYPE_HEAT_DEMAND"
)
