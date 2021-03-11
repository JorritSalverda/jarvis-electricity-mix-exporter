package contracts

import (
	"time"
)

type Measurement struct {
	ID             string
	Source         string
	Area           string
	Samples        []*Sample
	MeasuredAtTime time.Time
}
