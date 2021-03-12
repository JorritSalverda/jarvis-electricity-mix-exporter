package entsoe

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
	"time"

	"github.com/alecthomas/assert"
)

func TestUnmarshal(t *testing.T) {
	t.Run("GetsDocumentTypeUnknownIfNotPresent", func(t *testing.T) {
		xmlDocument := `
<?xml version="1.0" encoding="UTF-8"?>
<GL_MarketDocument xmlns="urn:iec62325.351:tc57wg16:451-6:generationloaddocument:3:0">
</GL_MarketDocument>`

		var response GetAggregatedGenerationPerTypeResponse

		// act
		err := xml.Unmarshal([]byte(xmlDocument), &response)

		assert.Nil(t, err)
		assert.Equal(t, DocumentTypeUnknown, response.DocumentType)
	})

	t.Run("GetsDocumentType", func(t *testing.T) {
		xmlDocument := `
<?xml version="1.0" encoding="UTF-8"?>
<GL_MarketDocument xmlns="urn:iec62325.351:tc57wg16:451-6:generationloaddocument:3:0">
	<type>A75</type>
</GL_MarketDocument>`

		var response GetAggregatedGenerationPerTypeResponse

		// act
		err := xml.Unmarshal([]byte(xmlDocument), &response)

		assert.Nil(t, err)
		assert.Equal(t, DocumentTypeActualGenerationPerType, response.DocumentType)
	})

	t.Run("GetsProcessTypeUnknownIfNotPresent", func(t *testing.T) {
		xmlDocument := `
<?xml version="1.0" encoding="UTF-8"?>
<GL_MarketDocument xmlns="urn:iec62325.351:tc57wg16:451-6:generationloaddocument:3:0">
</GL_MarketDocument>`

		var response GetAggregatedGenerationPerTypeResponse

		// act
		err := xml.Unmarshal([]byte(xmlDocument), &response)

		assert.Nil(t, err)
		assert.Equal(t, ProcessTypeUnknown, response.ProcessType)
	})

	t.Run("GetsProcessType", func(t *testing.T) {
		xmlDocument := `
<?xml version="1.0" encoding="UTF-8"?>
<GL_MarketDocument xmlns="urn:iec62325.351:tc57wg16:451-6:generationloaddocument:3:0">
	<process.processType>A16</process.processType>
</GL_MarketDocument>`

		var response GetAggregatedGenerationPerTypeResponse

		// act
		err := xml.Unmarshal([]byte(xmlDocument), &response)

		assert.Nil(t, err)
		assert.Equal(t, ProcessTypeRealised, response.ProcessType)
	})

	t.Run("GetsTimePeriod", func(t *testing.T) {
		xmlDocument := `
<?xml version="1.0" encoding="UTF-8"?>
<GL_MarketDocument xmlns="urn:iec62325.351:tc57wg16:451-6:generationloaddocument:3:0">
	<time_Period.timeInterval>
		<start>2021-03-11T00:00Z</start>
		<end>2021-03-11T07:30Z</end>
	</time_Period.timeInterval>
</GL_MarketDocument>`

		var response GetAggregatedGenerationPerTypeResponse

		// act
		err := xml.Unmarshal([]byte(xmlDocument), &response)

		assert.Nil(t, err)
		assert.Equal(t, time.Date(2021, 3, 11, 0, 0, 0, 0, time.UTC), response.TimePeriod.Start)
		assert.Equal(t, time.Date(2021, 3, 11, 7, 30, 0, 0, time.UTC), response.TimePeriod.End)
	})

	t.Run("GetsTimeSeries", func(t *testing.T) {

		xmlDocument := `
<?xml version="1.0" encoding="UTF-8"?>
<GL_MarketDocument xmlns="urn:iec62325.351:tc57wg16:451-6:generationloaddocument:3:0">
	<TimeSeries>
		<mRID>1</mRID>
		<businessType>A01</businessType>
		<objectAggregation>A08</objectAggregation>
		<inBiddingZone_Domain.mRID codingScheme="A01">10YNL----------L</inBiddingZone_Domain.mRID>
		<quantity_Measure_Unit.name>MAW</quantity_Measure_Unit.name>
		<curveType>A01</curveType>
		<MktPSRType>
			<psrType>B18</psrType>
		</MktPSRType>
		<Period>
			<timeInterval>
				<start>2021-03-11T07:00Z</start>
				<end>2021-03-11T07:30Z</end>
			</timeInterval>
			<resolution>PT15M</resolution>
				<Point>
					<position>1</position>
                        <quantity>739</quantity>
				</Point>
				<Point>
					<position>2</position>
                        <quantity>750</quantity>
				</Point>
		</Period>
	</TimeSeries>
	<TimeSeries>
		<mRID>2</mRID>
		<businessType>A01</businessType>
		<objectAggregation>A08</objectAggregation>
		<outBiddingZone_Domain.mRID codingScheme="A01">10YNL----------L</outBiddingZone_Domain.mRID>
		<quantity_Measure_Unit.name>MAW</quantity_Measure_Unit.name>
		<curveType>A01</curveType>
		<MktPSRType>
			<psrType>B04</psrType>
		</MktPSRType>
		<Period>
			<timeInterval>
				<start>2021-03-11T07:00Z</start>
				<end>2021-03-11T07:30Z</end>
			</timeInterval>
			<resolution>PT15M</resolution>
				<Point>
					<position>1</position>
                        <quantity>1554</quantity>
				</Point>
				<Point>
					<position>2</position>
                        <quantity>1581</quantity>
				</Point>
		</Period>
	</TimeSeries>	
</GL_MarketDocument>`

		var response GetAggregatedGenerationPerTypeResponse

		// act
		err := xml.Unmarshal([]byte(xmlDocument), &response)

		assert.Nil(t, err)
		assert.Equal(t, 2, len(response.TimeSeries))

		assert.Equal(t, AreaNetherlands, response.TimeSeries[0].InBiddingZone)
		assert.Equal(t, MeasurementUnitMegaWatt, response.TimeSeries[0].QuanityMeasurementUnit)
		assert.Equal(t, PsrTypeWindOffshore, response.TimeSeries[0].MktPsrType.PsrType)
		assert.Equal(t, time.Date(2021, 3, 11, 7, 0, 0, 0, time.UTC), response.TimeSeries[0].Period.TimeInterval.Start)
		assert.Equal(t, time.Date(2021, 3, 11, 7, 30, 0, 0, time.UTC), response.TimeSeries[0].Period.TimeInterval.End)
		assert.Equal(t, 2, len(response.TimeSeries[0].Period.Points))
		assert.Equal(t, 1, response.TimeSeries[0].Period.Points[0].Position)
		assert.Equal(t, 739.0, response.TimeSeries[0].Period.Points[0].Quantity)
		assert.Equal(t, 2, response.TimeSeries[0].Period.Points[1].Position)
		assert.Equal(t, 750.0, response.TimeSeries[0].Period.Points[1].Quantity)

		assert.Equal(t, AreaNetherlands, response.TimeSeries[1].OutBiddingZone)
		assert.Equal(t, MeasurementUnitMegaWatt, response.TimeSeries[1].QuanityMeasurementUnit)
		assert.Equal(t, PsrTypeFossilGas, response.TimeSeries[1].MktPsrType.PsrType)
		assert.Equal(t, time.Date(2021, 3, 11, 7, 0, 0, 0, time.UTC), response.TimeSeries[1].Period.TimeInterval.Start)
		assert.Equal(t, time.Date(2021, 3, 11, 7, 30, 0, 0, time.UTC), response.TimeSeries[1].Period.TimeInterval.End)
		assert.Equal(t, 2, len(response.TimeSeries[1].Period.Points))
		assert.Equal(t, 1, response.TimeSeries[1].Period.Points[0].Position)
		assert.Equal(t, 1554.0, response.TimeSeries[1].Period.Points[0].Quantity)
		assert.Equal(t, 2, response.TimeSeries[1].Period.Points[1].Position)
		assert.Equal(t, 1581.0, response.TimeSeries[1].Period.Points[1].Quantity)
	})

	t.Run("ReadsA75Response", func(t *testing.T) {

		testResponse, _ := ioutil.ReadFile("A75-response.xml")
		var response GetAggregatedGenerationPerTypeResponse

		// act
		err := xml.Unmarshal([]byte(testResponse), &response)

		assert.Nil(t, err)
		assert.Equal(t, 39, len(response.TimeSeries))
		assert.Equal(t, 5, response.TimeSeries[4].ID)
		assert.Equal(t, ResolutionPT15M, response.TimeSeries[4].Period.Resolution)
		assert.Equal(t, 92, len(response.TimeSeries[4].Period.Points))
		assert.Equal(t, 5046.0, response.TimeSeries[4].Period.Points[0].Quantity)
	})

	t.Run("ReadsA11Response", func(t *testing.T) {

		testResponse, _ := ioutil.ReadFile("A11-response.xml")
		var response GetPhysicalCrossBorderFlowResponse

		// act
		err := xml.Unmarshal([]byte(testResponse), &response)

		assert.Nil(t, err)
		assert.Equal(t, 1, len(response.TimeSeries))
		assert.Equal(t, 1, response.TimeSeries[0].ID)
		assert.Equal(t, ResolutionPT60M, response.TimeSeries[0].Period.Resolution)
		assert.Equal(t, 14, len(response.TimeSeries[0].Period.Points))
		assert.Equal(t, 701.0, response.TimeSeries[0].Period.Points[0].Quantity)
	})
}
