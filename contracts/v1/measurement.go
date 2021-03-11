package contracts

import (
	"time"
)

type Measurement struct {
	ID             string
	Source         string
	Location       string
	Samples        []*Sample
	MeasuredAtTime time.Time
}
