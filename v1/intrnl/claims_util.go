package intrnl

import (
	"errors"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func PopulateClaims(c *echo.Context, claims jwt.MapClaims) {
	ec := *c

	uid := claims[utils.AuthUserID]
	if uid != nil {
		userID, ok := uid.(string)
		if ok {
			ec.Set(utils.UserID, userID)
			ec.Set(utils.LoggedIn, true)
		}
	}

	internalUser := claims[utils.InternalStudentUserInToken]
	if internalUser != nil {
		ec.Set(utils.InternalStudentUser, internalUser)
	}

	tid := claims[utils.AuthTenantID]
	if tid != nil {
		tenantID, ok := tid.(string)
		if ok {
			ec.Set(utils.TenantID, tenantID)
		}
	}

	pt := claims[utils.AuthPersonaType]
	if pt != nil {
		personaType, ok := pt.(string)
		if ok {
			ec.Set(utils.PersonaType, personaType)
		}
	}

	sid := claims[utils.AuthSessionID]
	if sid != nil {
		sessionID, ok := sid.(string)
		if ok {
			ec.Set(utils.SessionID, sessionID)
		}
	}
	eid := claims[utils.ExternalID]
	if eid != nil {
		externalID, ok := eid.(string)
		if ok {
			ec.Set(utils.UserExternalID, externalID)
		}
	}
}

func GetUserID(c echo.Context) (string, error) {
	uid := c.Get(utils.UserID)
	if uid != nil {
		return uid.(string), nil
	}
	return "", errors.New("user ID Not Present")
}

func GetUserExternalIDAsStr(c echo.Context) (string, error) {
	euid := c.Get(utils.UserExternalID)
	if euid != nil {
		return euid.(string), nil
	}
	return utils.EmptyString, errors.New("external User ID Not Present")
}

func GetTenantID(_ echo.Context) (string, error) {
	return utils.AuthServTenantID, nil
	//tid := c.Get(utils.TenantID)
	//if tid != nil {
	//	return tid.(string), nil
	//}
	//return "", errors.New("tenant ID Not Present")
}

func GetPersonaType(c echo.Context) (string, error) {
	pt := c.Get(utils.PersonaType)
	if pt != nil {
		return pt.(string), nil
	}
	return "", errors.New("persona Type Not Present")
}

func GetSessionID(c echo.Context) (string, error) {
	sid := c.Get(utils.SessionID)
	if sid != nil {
		return sid.(string), nil
	}
	return "", errors.New("session ID Not Present")
}

func IsUserLoggedIn(c echo.Context) bool {
	lin := c.Get(utils.LoggedIn)
	if lin == nil {
		return false
	}
	loggedIn, ok := lin.(bool)
	if ok {
		return loggedIn
	}
	return false
}

func IsInternalUser(c echo.Context) bool {
	if IsUserLoggedIn(c) {
		pt := c.Get(utils.PersonaType)
		if pt != nil && (pt.(string) == utils.PersonaTypeInternalUser || pt.(string) == utils.PersonaTypeTeacher) {
			return true
		}
	}
	return false
}

func GetUserEnrolledStatus(c echo.Context) string {
	isEnrolled, ok := c.Get(utils.UserEnrolledStatus).(bool)
	if ok && isEnrolled {
		return utils.IsEnrolledTrue
	}
	return utils.IsEnrolledFalse
}
