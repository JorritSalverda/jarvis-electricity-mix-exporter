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
		nowRoundedTo15Minutes := now.Round(time.Duration(15 * time.Minute))
		startTime := nowRoundedTo15Minutes.Add(time.Duration(-3 * time.Hour))
		endTime := nowRoundedTo15Minutes

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
