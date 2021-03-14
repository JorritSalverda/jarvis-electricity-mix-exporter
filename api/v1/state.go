package api

import (
	"time"
)

type State struct {
	LastRetrievedGenerationTime map[Area]time.Time
}
