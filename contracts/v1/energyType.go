package contracts

type EnergyType string

const (
	EnergyTypeUnknown        EnergyType = "Unknown"
	EnergyTypeCoal           EnergyType = "Coal"
	EnergyTypeGas            EnergyType = "Gas"
	EnergyTypeOil            EnergyType = "Oil"
	EnergyTypeBiomass        EnergyType = "Biomass"
	EnergyTypeNuclear        EnergyType = "Nuclear"
	EnergyTypeWaste          EnergyType = "Waste"
	EnergyTypeGeothermal     EnergyType = "Geothermal"
	EnergyTypeHydro          EnergyType = "Hydro"
	EnergyTypeSolar          EnergyType = "Solar"
	EnergyTypeWindOffshore   EnergyType = "WindOffshore"
	EnergyTypeWindOnshore    EnergyType = "WindOnshore"
	EnergyTypeOtherRenewable EnergyType = "OtherRenewable"
)

func (e EnergyType) IsRenewable() bool {
	switch e {
	case EnergyTypeGeothermal,
		EnergyTypeHydro,
		EnergyTypeSolar,
		EnergyTypeWindOffshore,
		EnergyTypeWindOnshore,
		EnergyTypeOtherRenewable:
		return true
	}

	return false
}
