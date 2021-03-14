package exporter

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
	"time"

	apiv1 "github.com/JorritSalverda/jarvis-electricity-mix-exporter/api/v1"
	"github.com/alecthomas/assert"
)

func TestcreateGenerationMeasurementForTimeSlot(t *testing.T) {
	t.Run("CreatesSamplesForEachTimeSeriesThatHasAPointForFirstTimeSlot", func(t *testing.T) {

		service := service{}
		testResponse, _ := ioutil.ReadFile("../../client/entsoe/A75-response.xml")
		var response apiv1.GetAggregatedGenerationPerTypeResponse
		err := xml.Unmarshal([]byte(testResponse), &response)
		assert.Nil(t, err)

		// act
		measurement := service.createGenerationMeasurementForTimeSlot(response, response.TimePeriod.Start, apiv1.AreaConfig{Area: apiv1.AreaNetherlands})

		assert.Equal(t, 19, len(measurement.Samples))
	})

	t.Run("CreatesSamplesForEachTimeSeriesThatHasAPointForLastTimeSlot", func(t *testing.T) {

		service := service{}
		testResponse, _ := ioutil.ReadFile("../../client/entsoe/A75-response.xml")
		var response apiv1.GetAggregatedGenerationPerTypeResponse
		err := xml.Unmarshal([]byte(testResponse), &response)
		assert.Nil(t, err)

		// act
		measurement := service.createGenerationMeasurementForTimeSlot(response, response.TimePeriod.End.Add(time.Duration(-1*15)*time.Minute), apiv1.AreaConfig{Area: apiv1.AreaNetherlands})

		assert.Equal(t, 19, len(measurement.Samples))
	})

	t.Run("CreatesSamplesForEachTimeSeriesForEachTimeSlot", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		service := service{}
		testResponse, _ := ioutil.ReadFile("../../client/entsoe/A75-response.xml")
		var response apiv1.GetAggregatedGenerationPerTypeResponse
		err := xml.Unmarshal([]byte(testResponse), &response)
		assert.Nil(t, err)

		// act
		nrOfSlots := int(response.TimePeriod.End.Sub(response.TimePeriod.Start).Minutes() / 15)
		assert.Equal(t, 96, nrOfSlots)
		for i := 0; i < nrOfSlots; i++ {
			timeSlotStartTime := response.TimePeriod.Start.Add(time.Duration(i*15) * time.Minute)
			measurement := service.createGenerationMeasurementForTimeSlot(response, timeSlotStartTime, apiv1.AreaConfig{Area: apiv1.AreaNetherlands})

			assert.Equal(t, 19, len(measurement.Samples), "Number of samples for time slot %v does not match expectation", timeSlotStartTime)
			assert.Equal(t, timeSlotStartTime, measurement.MeasuredAtTime)
		}
	})

	t.Run("CreatesSamplesForEachTimeSeriesThatHasAPointForLastTimeSlot", func(t *testing.T) {

		service := service{}
		testResponse, _ := ioutil.ReadFile("../../client/entsoe/A75-response.xml")
		var response apiv1.GetAggregatedGenerationPerTypeResponse
		err := xml.Unmarshal([]byte(testResponse), &response)
		assert.Nil(t, err)

		// act
		measurement := service.createGenerationMeasurementForTimeSlot(response, response.TimePeriod.End.Add(time.Duration(-1*15)*time.Minute), apiv1.AreaConfig{Area: apiv1.AreaNetherlands})

		assert.Equal(t, 19, len(measurement.Samples))
	})
}
