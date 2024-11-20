package middleware

import (
	"errors"
	authReq "github.com/Allen-Career-Institute/common-protos/authentication/v1/request"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	internal "github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/intrnl/clients"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/status"
	"net/http"
	"strings"
)

const (
	BearerSchema          = "Bearer "
	UnauthorizedErrorCode = 401
	StatusForbidden       = 403
	TooManyRequests       = 429
)

func (m *Manager) AuthNMiddleware(conf *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			m.logger.WithContext(c).Debug("AuthN Middleware - Execution ")
			authHeader := c.Request().Header.Get("Authorization")
			// Continue if auth header is empty or not handled
			if authHeader == "" || !strings.HasPrefix(authHeader, BearerSchema) {
				m.logger.WithContext(c).Infof("AuthN Middleware - Execution Skipped %s", authHeader)
				err := next(c)
				if err != nil {
					return err
				}
				return nil
			}
			// Not Required ... Keeping it for now.
			if !strings.HasPrefix(authHeader, BearerSchema) {
				c.Response().Status = UnauthorizedErrorCode
				return &echo.HTTPError{Message: "invalid Authorization header format", Code: UnauthorizedErrorCode}
			}

			tokenString := strings.TrimPrefix(authHeader, BearerSchema)
			token, err := m.parseToken(tokenString)
			// don't refresh the tokens if the access token in the request is valid
			if !(authHeader != "" && err == nil) || authHeader == "" {
				errRefresh := m.updateRefreshToken(c, conf, m.logger)
				if errRefresh != nil {
					m.logger.WithContext(c).Errorf("Error_updateRefreshToken %v", errRefresh)
					return errRefresh
				}
				// if the tokens are refreshed successfully don't throw error and instead use the refreshed tokens for the current request
				if c.Response().Header().Get(utils.XRefreshToken) == "" {
					c.Response().Status = UnauthorizedErrorCode
					return &echo.HTTPError{Message: err.Error(), Code: UnauthorizedErrorCode, Internal: err}
				}
				token, err = m.parseToken(c.Response().Header().Get(utils.XAccessToken))
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				internal.PopulateClaims(&c, claims)
				nextErr := next(c)
				if nextErr != nil {
					m.logger.WithContext(c).Errorf("Error_afterPopulateClaims %v", nextErr)
					return nextErr
				} //con
				return nil
			}
			c.Response().Status = UnauthorizedErrorCode
			return &echo.HTTPError{Message: err.Error(), Code: UnauthorizedErrorCode, Internal: err}
		}
	}
}

func (m *Manager) parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(m.secret), nil
	})
}

func (m *Manager) updateRefreshToken(c echo.Context, conf *config.Config, log logger.Logger) error {
	cm := clients.NewClientManager(conf, log, m.grpc)
	refreshToken := c.Request().Header.Get(utils.RefreshTokenHeader)
	if refreshToken != "" {
		m.logger.WithContext(c).Infof("AuthN Middleware:updateRefreshToken -  Request to refresh tokens - %s", refreshToken)
		tenantID, err := internal.GetTenantID(c)
		if err != nil {
			return err
		}
		refResp, errTokens := cm.RefreshToken(c, conf, &authReq.RefreshTokenRequest{
			TenantId:     tenantID,
			RefreshToken: refreshToken,
		})
		if refResp.GetResponseCode() == "403" {
			m.logger.WithContext(c).Warnf("AuthN Middleware:updateRefreshToken - Invalid RT - Expired - %s", refreshToken)
			c.Response().Status = StatusForbidden
			return &echo.HTTPError{Message: "invalid credentials", Code: StatusForbidden, Internal: err}
		}
		if errTokens != nil {
			e, _ := status.FromError(errTokens)
			respCode := utils.GetHTTPStatusCode(e.Code())
			if respCode == TooManyRequests {
				c.Response().Status = StatusForbidden
				m.logger.WithContext(c).Warnf("AuthN Middleware:updateRefreshToken - Duplicate RT call - %s", refreshToken)
				return &echo.HTTPError{Message: "invalid credentials", Code: StatusForbidden, Internal: err}
			}
			c.Response().Status = UnauthorizedErrorCode
			m.logger.WithContext(c).Warnf("AuthN Middleware:updateRefreshToken - Error refreshing tokens - %s", refreshToken)
			return &echo.HTTPError{Message: "invalid credentials", Code: UnauthorizedErrorCode, Internal: err}
		}
		c.Response().Header().Set(utils.XAccessToken, refResp.GetAccessToken())
		c.Response().Header().Set(utils.XRefreshToken, refResp.GetRefreshToken())
		c.Request().Header.Set("Authorization", "Bearer "+refResp.GetAccessToken())
		cookie := &http.Cookie{
			Name:     "isAuth",
			Value:    refResp.GetAccessToken(),
			HttpOnly: true,
		}
		c.SetCookie(cookie)
	}
	return nil
}
