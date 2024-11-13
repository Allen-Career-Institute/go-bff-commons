package models

import "time"

type DSResponse struct {
	Status   int         `json:"status" validate:"required,lte=30"`
	Reason   string      `json:"reason" validate:"required,lte=30"`
	ViewType string      `json:"-"`
	Data     interface{} `json:"data"`
}

type DataSourceConfig struct {
	URI       string
	Method    string
	Timeout   time.Duration
	DsName    string
	Resource  string
	Action    string
	PreloadDS []string
}
