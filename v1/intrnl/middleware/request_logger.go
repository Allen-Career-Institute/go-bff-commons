package middleware

import (
	"github.com/labstack/echo/v4"
	"time"
)

// RequestLoggerMiddleware Request logger middleware
func (m *Manager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		t := time.Since(start).String()

		if status >= 200 && status < 300 {
			m.logger.WithContext(ctx).Infow("logging_request", "method", req.Method, "uri", req.URL, "status", status, "size", size, "time", t, "Request Header", req.Header)
		} else {
			m.logger.WithContext(ctx).Errorw("logging_request", "method", req.Method, "uri", req.URL, "status", status, "size", size, "time", t, "Request Header", req.Header, "err", err)
		}

		return err
	}
}
