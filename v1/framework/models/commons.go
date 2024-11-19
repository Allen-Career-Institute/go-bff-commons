package models

import "time"

type DataSourceConfig struct {
	URI       string
	Method    string
	Timeout   time.Duration
	DsName    string
	Resource  string
	Action    string
	PreloadDS []string
}
