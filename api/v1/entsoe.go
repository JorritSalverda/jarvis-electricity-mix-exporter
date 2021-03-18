package api

import (
	"encoding/xml"
	"fmt"
	"time"
)

type GetAggregatedGenerationPerTypeResponse struct {
	DocumentType DocumentType                    `xml:"type"`
	ProcessType  ProcessType                     `xml:"process.processType"`
	TimePeriod   TimeInterval                    `xml:"time_Period.timeInterval"`
	TimeSeries   []AggregatedGenerationTimeSerie `xml:"TimeSeries"`
}

type TimeInterval struct {
	Start time.Time `xml:"start"`
	End   time.Time `xml:"end"`
}

type AggregatedGenerationTimeSerie struct {
	ID                     int             `xml:"mRID"`
	InBiddingZone          Area            `xml:"inBiddingZone_Domain.mRID"`
	OutBiddingZone         Area            `xml:"outBiddingZone_Domain.mRID"`
	QuanityMeasurementUnit MeasurementUnit `xml:"quantity_Measure_Unit.name"`
	MktPsrType             struct {
		PsrType PsrType `xml:"psrType"`
	} `xml:"MktPSRType"`
	Period TimeSeriePeriod `xml:"Period"`
}

type TimeSeriePeriod struct {
	TimeInterval TimeInterval     `xml:"timeInterval"`
	Resolution   Resolution       `xml:"resolution"`
	Points       []TimeSeriePoint `xml:"Point"`
}

type TimeSeriePoint struct {
	Position int     `xml:"position"`
	Quantity float64 `xml:"quantity"`
}

type GetPhysicalCrossBorderFlowResponse struct {
	TimePeriod TimeInterval            `xml:"period.timeInterval"`
	TimeSeries []PhysicalFlowTimeSerie `xml:"TimeSeries"`
}

type PhysicalFlowTimeSerie struct {
	ID                     int             `xml:"mRID"`
	InDomain               Area            `xml:"in_Domain.mRID"`
	OutDomain              Area            `xml:"out_Domain.mRID"`
	QuanityMeasurementUnit MeasurementUnit `xml:"quantity_Measure_Unit.name"`
	Period                 TimeSeriePeriod `xml:"Period"`
}

const timeIntervalLayout = "2006-01-02T15:04Z"

func (t *TimeInterval) FormatAsParameter() string {
	return fmt.Sprintf("%v/%v", t.Start.Format(timeIntervalLayout), t.End.Format(timeIntervalLayout))

}

func (t *TimeInterval) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	v := struct {
		Start string `xml:"start"`
		End   string `xml:"end"`
	}{}
	d.DecodeElement(&v, &start)

	startTime, err := time.Parse(timeIntervalLayout, v.Start)
	if err != nil {
		return err
	}
	endTime, err := time.Parse(timeIntervalLayout, v.End)
	if err != nil {
		return err
	}

	*t = TimeInterval{Start: startTime, End: endTime}

	return nil
}

type Area string

const (
	AreaUnknown      Area = ""
	AreaBelgium      Area = "10YBE----------2"
	AreaDenmark      Area = "10YDK-1--------W"
	AreaGermany      Area = "10Y1001A1001A83F"
	AreaGreatBritain Area = "10YGB----------A"
	AreaNetherlands  Area = "10YNL----------L"
	AreaNorway       Area = "10YNO-0--------C"
)

type ProcessType string

const (
	ProcessTypeUnknown  ProcessType = ""
	ProcessTypeRealised ProcessType = "A16"
)

type DocumentType string

const (
	DocumentTypeUnknown                    DocumentType = ""
	DocumentTypeAggregatedEnergyDataReport DocumentType = "A11"
	DocumentTypeSystemTotalLoad            DocumentType = "A65"
	DocumentTypeActualGenerationPerType    DocumentType = "A75"
)

type MeasurementUnit string

const (
	MeasurementUnitUnknown  MeasurementUnit = ""
	MeasurementUnitMegaWatt MeasurementUnit = "MAW"
)

type PsrType string

const (
	PsrTypeUnknown              PsrType = ""
	PsrTypeMixed                PsrType = "A03"
	PsrTypeGeneration           PsrType = "A04"
	PsrTypeLoad                 PsrType = "A05"
	PsrTypeBiomass              PsrType = "B01"
	PsrTypeFossilBrownCoal      PsrType = "B02"
	PsrTypeFossilCoalDerivedGas PsrType = "B03"
	PsrTypeFossilGas            PsrType = "B04"
	PsrTypeFossilHardCoal       PsrType = "B05"
	PsrTypeFossilOil            PsrType = "B06"
	PsrTypeFossilOilShale       PsrType = "B07"
	PsrTypeFossilOilPeat        PsrType = "B08"
	PsrTypeGeothermal           PsrType = "B09"
	PsrTypeHydroPumpedStorage   PsrType = "B10"
	PsrTypeHydroRunOfRiver      PsrType = "B11"
	PsrTypeHydroWaterReservoir  PsrType = "B12"
	PsrTypeMarin                PsrType = "B13"
	PsrTypeNuclear              PsrType = "B14"
	PsrTypeOtherRenewable       PsrType = "B15"
	PsrTypeSolar                PsrType = "B16"
	PsrTypeWaste                PsrType = "B17"
	PsrTypeWindOffshore         PsrType = "B18"
	PsrTypeWindOnshore          PsrType = "B19"
	PsrTypeOther                PsrType = "B20"
	PsrTypeACLink               PsrType = "B21"
	PsrTypeDCLink               PsrType = "B22"
	PsrTypeSubstation           PsrType = "B23"
	PsrTypeTransformer          PsrType = "B24"
)

type Resolution string

const (
	ResolutionUnknown Resolution = ""
	ResolutionPT15M   Resolution = "PT15M"
	ResolutionPT60M   Resolution = "PT60M"
)
