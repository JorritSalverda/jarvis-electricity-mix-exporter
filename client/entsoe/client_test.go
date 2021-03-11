package entsoe

import (
	"os"
	"testing"
	"time"

	"github.com/alecthomas/assert"
)

func TestGetAggregatedGenerationPerType(t *testing.T) {
	t.Run("ReturnsGetAggregatedGenerationPerType", func(t *testing.T) {

		if testing.Short() {
			t.Skip("skipping test in short mode.")
		}

		// set first with 'export ENTSOE_TOKEN=...'
		client, err := NewClient(os.Getenv("ENTSOE_TOKEN"))
		assert.Nil(t, err)

		area := AreaNetherlands

		now := time.Now().UTC()
		nowRoundedToTimeSlotSize := now.Round(time.Duration(TimeSlotsInMinutes * time.Minute))
		startTime := nowRoundedToTimeSlotSize.Add(time.Duration(-3 * time.Hour))
		endTime := nowRoundedToTimeSlotSize

		timeInterval := TimeInterval{
			Start: startTime,
			End:   endTime,
		}

		// act
		response, err := client.GetAggregatedGenerationPerType(area, timeInterval)

		assert.Nil(t, err)
		assert.Equal(t, "", response)
	})
}
