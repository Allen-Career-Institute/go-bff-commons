package metrics

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"time"
)

const (
	circuitBreakerStateChangeCount     = "circuit_breaker_state_change_count"
	circuitBreakerStateChangeCountDesc = "Circuit breaker state change counter"
	circuitBreakerStateDuration        = "circuit_breaker_state_duration"
	circuitBreakerStateDurationDesc    = "Circuit breaker state duration"
)

type CircuitBreakerMetrics struct {
	StateChangeCount metric.Int64Counter
	StateDuration    metric.Float64Histogram
}

func createCircuitBreakerMetrics(meter metric.Meter, counters, histograms []MetricParams) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	for _, counterParams := range counters {
		counter, err := createCounter(meter, counterParams)
		if err != nil {
			return nil, err
		}
		metrics[counterParams.Name] = counter
	}

	for _, histogramParams := range histograms {
		histogram, err := createHistogram(meter, histogramParams)
		if err != nil {
			return nil, err
		}
		metrics[histogramParams.Name] = histogram
	}

	return metrics, nil
}

func NewCircuitBreakerMetrics(meter metric.Meter) (*CircuitBreakerMetrics, error) {
	counters := []MetricParams{
		{Name: circuitBreakerStateChangeCount, Desc: circuitBreakerStateChangeCountDesc},
	}
	histograms := []MetricParams{
		{Name: circuitBreakerStateDuration, Desc: circuitBreakerStateDurationDesc},
	}
	metrics, err := createCircuitBreakerMetrics(meter, counters, histograms)

	if err != nil {
		return nil, err
	}
	return &CircuitBreakerMetrics{
		StateChangeCount: metrics[circuitBreakerStateChangeCount].(metric.Int64Counter),
		StateDuration:    metrics[circuitBreakerStateDuration].(metric.Float64Histogram),
	}, nil
}

func createAttributes(attributeMap map[string]string, uniqueIdentifier string) []attribute.KeyValue {
	attributes := make([]attribute.KeyValue, 0, len(attributeMap)+1)
	for key, val := range attributeMap {
		attributes = append(attributes, attribute.String(key, val))
	}
	// result added explicitly to force this field
	attributes = append(attributes, attribute.String(UniqueIdentifierKey, uniqueIdentifier))
	return attributes
}

func HistogramRecord(metricHistogram metric.Float64Histogram, ctx context.Context, uniqueIdentifier string,
	attributeMap map[string]string, duration time.Duration) {
	attributes := createAttributes(attributeMap, uniqueIdentifier)
	metricHistogram.Record(ctx, duration.Seconds(), metric.WithAttributes(attributes...))
}

func AddCounter(metricCounter metric.Int64Counter, ctx context.Context, uniqueIdentifier string, attributeMap map[string]string) {
	attributes := createAttributes(attributeMap, uniqueIdentifier)
	metricCounter.Add(ctx, 1, metric.WithAttributes(attributes...))
}
