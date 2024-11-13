package metrics

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
)

const (
	datasourceRequestCount = "datasource_req_count"
	datasourceRequestDesc  = "Request counts from inside of ds"
)

const servicePrefix = "bff_service_"

type CustomMetrics struct {
	RequestCount metric.Int64Counter
	log          logger.Logger
}

func createMetrics(meter metric.Meter, counters []MetricParams) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	for _, counterParams := range counters {
		counter, err := createCounter(meter, counterParams)
		if err != nil {
			return nil, err
		}

		metrics[counterParams.Name] = counter
	}

	return metrics, nil
}

func NewCustomMetrics(l logger.Logger) (*CustomMetrics, error) {
	meter := otel.GetMeterProvider().Meter("bff-service")

	counters := []MetricParams{
		{Name: datasourceRequestCount, Desc: datasourceRequestDesc},
	}

	metrics, err := createMetrics(meter, counters)
	if err != nil {
		l.Errorf("Failed to create metric: %v", err)

		return nil, err
	}

	return &CustomMetrics{
		RequestCount: metrics[datasourceRequestCount].(metric.Int64Counter),
		log:          l,
	}, nil
}

func (v *CustomMetrics) AddCounter(ctx context.Context, uniqueIdentifier string, attributeMap map[string]string) {
	attributes := make([]attribute.KeyValue, 0, len(attributeMap)+1)

	for key, val := range attributeMap {
		attributes = append(attributes, attribute.String(key, val))
	}
	// result added explicitly to force this field
	attributes = append(attributes, attribute.String(UniqueIdentifierKey, uniqueIdentifier))
	v.RequestCount.Add(ctx, 1, metric.WithAttributes(attributes...))
}
