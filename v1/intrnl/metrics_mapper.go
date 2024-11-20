package intrnl

import (
	"go.opentelemetry.io/otel/metric"
)

type Mapper struct {
	meter metric.Meter
}

func NewMapper(meter metric.Meter) *Mapper {
	return &Mapper{meter: meter}
}

func (m *Mapper) GetDuration(message string) (metric.Int64Histogram, error) {
	return m.meter.Int64Histogram(
		message,
		metric.WithUnit("ms"),
		metric.WithDescription("Request Latency in ms"))
}

func (m *Mapper) GetCount(message string) (metric.Int64Counter, error) {
	return m.meter.Int64Counter(
		message,
		metric.WithDescription("Number of requests"),
	)
}
