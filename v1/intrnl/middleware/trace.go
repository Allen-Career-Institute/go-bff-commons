package middleware

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

func (m *Manager) Trace(tracer trace.Tracer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.HasPrefix(c.Request().URL.Path, "/health") {
				return next(c)
			}

			spanName := c.Request().URL.Path
			ctx, span := tracer.Start(utils.GetRequestCtx(c), spanName, trace.WithAttributes(attribute.KeyValue{
				Key:   "http.method",
				Value: attribute.StringValue(c.Request().Method),
			}, attribute.KeyValue{
				Key:   "http.status_code",
				Value: attribute.IntValue(c.Response().Status),
			}, attribute.KeyValue{
				Key:   "resource.name",
				Value: attribute.StringValue(spanName),
			}))
			defer span.End()
			c.SetRequest(c.Request().WithContext(ctx))
			err := next(c)
			if err != nil {
				m.logger.WithContext(c).Errorf("Error_traceMiddleware %v", err)
			}
			return err
		}
	}
}
