package middleware

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	echo "github.com/labstack/echo/v4"
)

// ResponseHeader Response header modify and "Access-Control-Expose-Headers"
func (m *Manager) ResponseHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		hl := utils.GetListOfCustomHeadersAsCommaSeparated()
		ctx.Response().Header().Set(echo.HeaderAccessControlExposeHeaders, hl)
		return next(ctx)
	}
}
