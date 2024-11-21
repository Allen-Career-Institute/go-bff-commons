package middleware

import (
	"github.com/Allen-Career-Institute/common-protos/authorization/v1/types"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients"
	"github.com/labstack/echo/v4"
	"net/http"
)

const AccessFailedMessage = "You do not have access to this resource. Drop a mail to internal-console-permissions@allen.in in case of issues."

func (m *Manager) RbacMiddleware(conf *config.Config, resource types.ResourceTypes, action types.Action) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			m.logger.WithContext(c).Debug("Rbac Middleware - Execution ")
			if resource == types.ResourceTypes_RESOURCE_UNSPECIFIED || action == types.Action_ACTION_UNSPECIFIED {
				return next(c)
			}

			cm := clients.NewClientManager(conf, m.logger, m.grpc)
			res, err := cm.EnforceRbac(c, conf, resource, action)
			if err != nil || res.Pass == false {
				c.Response().Status = http.StatusPreconditionFailed
				m.logger.WithContext(c).Errorf("RbacMiddleware: Error enforcing RBAC %v", err)
				return &echo.HTTPError{Message: AccessFailedMessage, Code: http.StatusPreconditionFailed, Internal: err}
			}
			return next(c)
		}
	}
}
