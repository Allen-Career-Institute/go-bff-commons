package utils

import (
	echo "github.com/labstack/echo/v4"

	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
)

type EchoUtil struct {
	logger    logger.Logger
	claimKeys map[string]bool
}

func NewEchoUtil(log logger.Logger) EchoUtil {
	return EchoUtil{
		logger:    log,
		claimKeys: make(map[string]bool),
	}
}

// CloneContext creates a copy of the echo context copying the req, resp and claims
func (ec *EchoUtil) CloneContext(c echo.Context) echo.Context {
	// Creating new echo context using existing one
	e := echo.New()
	req := c.Request().Clone(c.Request().Context())
	res := echo.NewResponse(c.Response().Writer, e)
	newCtx := e.NewContext(req, res)

	// copy the populated context as well otherwise login filter, tenantID won't be passed
	for claimKey := range ec.getClaimKeys() {
		newCtx.Set(claimKey, c.Get(claimKey))
	}

	return newCtx
}

func (ec *EchoUtil) getClaimKeys() map[string]bool {
	if len(ec.claimKeys) == 0 {
		ec.initClaimKeys()
	}

	return ec.claimKeys
}

func (ec *EchoUtil) setClaimKey(key string) {
	_, found := ec.claimKeys[key]
	if found {
		ec.logger.Warn("claim key already exists")

		return
	}

	ec.claimKeys[key] = true
}

func (ec *EchoUtil) initClaimKeys() {
	ec.setClaimKey(LoggedIn)
	ec.setClaimKey(TenantID)
	ec.setClaimKey(UserID)
	ec.setClaimKey(PersonaType)
	ec.setClaimKey(SessionID)
	ec.setClaimKey(UserExternalID)
	ec.setClaimKey(SharedDataSource)
	ec.setClaimKey(UserEnrolledStatus)
	ec.setClaimKey(PageURL)
	ec.setClaimKey(CenterName)
	// set claim key for url meta passing in context
	ec.setClaimKey(URLMeta)
	ec.setClaimKey(WidgetData)

}
