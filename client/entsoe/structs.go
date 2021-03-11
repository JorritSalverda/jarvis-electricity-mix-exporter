package entsoe

import (
	"encoding/xml"
	"fmt"
	"time"
)

type GetAggregatedGenerationPerTypeResponse struct {
	DocumentType DocumentType `xml:"type"`
	ProcessType  ProcessType  `xml:"process.processType"`
	TimePeriod   TimeInterval `xml:"time_Period.timeInterval"`
	TimeSeries   []TimeSerie  `xml:"TimeSeries"`
}

type TimeInterval struct {
	Start time.Time `xml:"start"`
	End   time.Time `xml:"end"`
}

type TimeSerie struct {
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
	Points       []TimeSeriePoint `xml:"Point"`
}

type TimeSeriePoint struct {
	Position int     `xml:"position"`
	Quantity float64 `xml:"quantity"`
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
	AreaUnknown     Area = ""
	AreaNetherlands Area = "10YNL----------L"
)

type ProcessType string

const (
	ProcessTypeUnknown  ProcessType = ""
	ProcessTypeRealised ProcessType = "A16"
)

type DocumentType string

const (
	DocumentTypeUnknown                 DocumentType = ""
	DocumentTypeSystemTotalLoad         DocumentType = "A65"
	DocumentTypeActualGenerationPerType DocumentType = "A75"
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
