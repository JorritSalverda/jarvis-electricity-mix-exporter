package entsoe

import (
	"os"
	"testing"
	"time"

	apiv1 "github.com/JorritSalverda/jarvis-electricity-mix-exporter/api/v1"
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

		area := apiv1.AreaNetherlands

		now := time.Now().UTC()
		nowRoundedToTimeSlotSize := now.Round(time.Duration(15 * time.Minute))
		startTime := nowRoundedToTimeSlotSize.Add(time.Duration(-3 * time.Hour))
		endTime := nowRoundedToTimeSlotSize

		timeInterval := apiv1.TimeInterval{
			Start: startTime,
			End:   endTime,
		}

		// act
		response, err := client.GetAggregatedGenerationPerType(area, timeInterval)

		assert.Nil(t, err)
		assert.Equal(t, "", response)
	})
}
