package api

import (
	"time"
)

type GenerationMeasurement struct {
	ID             string
	Source         string
	Area           string
	Country        string
	Resolution     string
	Samples        []*Sample
	MeasuredAtTime time.Time
}
