package exporter

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
	"time"

	"github.com/JorritSalverda/jarvis-electricity-mix-exporter/client/entsoe"
	"github.com/alecthomas/assert"
)

func TestCreateMeasurementForTimeSlot(t *testing.T) {
	t.Run("CreatesSamplesForEachTimeSeriesThatHasAPointForFirstTimeSlot", func(t *testing.T) {

		service := service{}
		testResponse, _ := ioutil.ReadFile("../../client/entsoe/A75-response.xml")
		var response entsoe.GetAggregatedGenerationPerTypeResponse
		err := xml.Unmarshal([]byte(testResponse), &response)
		assert.Nil(t, err)

		// act
		measurement := service.createMeasurementForTimeSlot(response, response.TimePeriod.Start, entsoe.AreaNetherlands)

		assert.Equal(t, 19, len(measurement.Samples))
	})

	t.Run("CreatesSamplesForEachTimeSeriesThatHasAPointForLastTimeSlot", func(t *testing.T) {

		service := service{}
		testResponse, _ := ioutil.ReadFile("../../client/entsoe/A75-response.xml")
		var response entsoe.GetAggregatedGenerationPerTypeResponse
		err := xml.Unmarshal([]byte(testResponse), &response)
		assert.Nil(t, err)

		// act
		measurement := service.createMeasurementForTimeSlot(response, response.TimePeriod.End.Add(time.Duration(-1*entsoe.TimeSlotsInMinutes)*time.Minute), entsoe.AreaNetherlands)

		assert.Equal(t, 19, len(measurement.Samples))
	})

	t.Run("CreatesSamplesForEachTimeSeriesForEachTimeSlot", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		service := service{}
		testResponse, _ := ioutil.ReadFile("../../client/entsoe/A75-response.xml")
		var response entsoe.GetAggregatedGenerationPerTypeResponse
		err := xml.Unmarshal([]byte(testResponse), &response)
		assert.Nil(t, err)

		// act
		nrOfSlots := int(response.TimePeriod.End.Sub(response.TimePeriod.Start).Minutes() / entsoe.TimeSlotsInMinutes)
		assert.Equal(t, 96, nrOfSlots)
		for i := 0; i < nrOfSlots; i++ {
			timeSlotStartTime := response.TimePeriod.Start.Add(time.Duration(i*entsoe.TimeSlotsInMinutes) * time.Minute)
			measurement := service.createMeasurementForTimeSlot(response, timeSlotStartTime, entsoe.AreaNetherlands)

			assert.Equal(t, 19, len(measurement.Samples), "Number of samples for time slot %v does not match expectation", timeSlotStartTime)
			assert.Equal(t, timeSlotStartTime, measurement.MeasuredAtTime)
		}
	})

	t.Run("CreatesSamplesForEachTimeSeriesThatHasAPointForLastTimeSlot", func(t *testing.T) {

		service := service{}
		testResponse, _ := ioutil.ReadFile("../../client/entsoe/A75-response.xml")
		var response entsoe.GetAggregatedGenerationPerTypeResponse
		err := xml.Unmarshal([]byte(testResponse), &response)
		assert.Nil(t, err)

		// act
		measurement := service.createMeasurementForTimeSlot(response, response.TimePeriod.End.Add(time.Duration(-1*entsoe.TimeSlotsInMinutes)*time.Minute), entsoe.AreaNetherlands)

		assert.Equal(t, 19, len(measurement.Samples))
	})
}
