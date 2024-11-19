package middleware

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// MetricsMiddleware Prometheus metrics middleware
func (m *Manager) MetricsMiddleware() echo.MiddlewareFunc {
	reqCount, mapperErr := m.mapper.GetCount(utils.MetricPrefix + utils.Count)
	if mapperErr != nil {
		m.logger.Fatalf("Error instrumenting request for count, %v", mapperErr)
	}

	reqDuration, mapperErr := m.mapper.GetDuration(utils.MetricPrefix + utils.Duration)
	if mapperErr != nil {
		m.logger.Fatalf("Error instrumenting request for duration, %v", mapperErr)
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()

			ctx := c.Request().Context()
			err := next(c)
			if err != nil {
				m.logger.WithContext(c).Errorf("Error_MetricsMiddleware %v", err)
			}

			if c.Request().Method != http.MethodOptions {
				reqCount.Add(ctx, 1, getAddMetricTags(c)...)
				reqDuration.Record(ctx, time.Since(startTime).Milliseconds(), getRecordMetricTags(c)...)
			}
			return err
		}
	}
}

func getRecordMetricTags(c echo.Context) []metric.RecordOption {
	return []metric.RecordOption{
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.ServiceName,
			Value: attribute.StringValue("bff-service"),
		}),
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.URI,
			Value: attribute.StringValue(getURI(c)),
		}),
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.StatusCode,
			Value: attribute.IntValue(c.Response().Status),
		}),
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.ServiceEnv,
			Value: attribute.StringValue(os.Getenv("ENV")),
		}),
	}
}

func getAddMetricTags(c echo.Context) []metric.AddOption {
	return []metric.AddOption{
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.ServiceName,
			Value: attribute.StringValue("bff-service"),
		}),
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.URI,
			Value: attribute.StringValue(getURI(c)),
		}),
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.StatusCode,
			Value: attribute.IntValue(c.Response().Status),
		}),
		metric.WithAttributes(attribute.KeyValue{
			Key:   utils.ServiceEnv,
			Value: attribute.StringValue(os.Getenv("ENV")),
		}),
	}
}

func getURI(c echo.Context) string {
	pathParams := c.ParamNames()
	if pathParams != nil && len(pathParams) > 0 {
		return c.Path()
	}
	return c.Request().URL.Path
}
