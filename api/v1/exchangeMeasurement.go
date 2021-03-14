package api

import (
	"time"
)

type ExchangeMeasurement struct {
	ID               string
	Source           string
	Area             string
	ExchangeWithArea string
	Samples          []*Sample
	MeasuredAtTime   time.Time
}
