package utils

import (
	"time"
)

type EnableResultInsightsV2Struct struct {
	EnableResultInsightsFlag bool
	ConfiguredDate           time.Time
	TestDate                 time.Time
}
