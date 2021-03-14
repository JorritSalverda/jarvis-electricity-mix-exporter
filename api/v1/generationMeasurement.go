package api

import (
	"time"
)

type GenerationMeasurement struct {
	ID             string
	Source         string
	Area           string
	Samples        []*Sample
	MeasuredAtTime time.Time
}
