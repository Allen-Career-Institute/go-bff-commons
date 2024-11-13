package metrics

import "go.opentelemetry.io/otel/metric"

const UniqueIdentifierKey = "unique_identifier"

type MetricParams struct {
	Name string
	Desc string
}

func createCounter(meter metric.Meter, params MetricParams) (metric.Int64Counter, error) {
	counter, err := meter.Int64Counter(
		servicePrefix+params.Name,
		metric.WithDescription(params.Desc),
	)
	if err != nil {
		return nil, err
	}

	return counter, nil
}

func createHistogram(meter metric.Meter, params MetricParams) (metric.Float64Histogram, error) {
	histogram, err := meter.Float64Histogram(
		servicePrefix+params.Name,
		metric.WithDescription(params.Desc),
	)
	if err != nil {
		return nil, err
	}
	return histogram, nil
}
