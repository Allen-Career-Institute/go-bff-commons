package pagehandler

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/utils"
	"net/http"
)

func UserFacingMessage(code int) string {
	switch code {
	case http.StatusInternalServerError:
		return utils.GenericError
	case http.StatusUnauthorized:
		return UnauthorizedPageMsg
	case http.StatusBadRequest:
		return BadRequestMsg
	case http.StatusNotFound:
		return EntityNotExistMsg
	}
	return utils.GenericError
}
