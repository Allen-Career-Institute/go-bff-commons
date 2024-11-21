package httperr

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

var (
	ErrBadRequest          = errors.New("bad Request ")
	ErrNotFound            = errors.New("not Found")
	ErrUnauthorized        = errors.New("user not logged In")
	ErrForbidden           = errors.New("forbidden")
	ErrInternalServerError = errors.New("internal Server Error")
	ErrRequestTimeoutError = errors.New("request Timeout")
)

// RestErr Rest error interface
type RestErr interface {
	Status() int
	Error() string
	Causes() interface{}
}

// RestError Rest error struct
type RestError struct {
	ErrStatus int         `json:"status,omitempty"`
	ErrError  string      `json:"error,omitempty"`
	ErrCauses interface{} `json:"-"`
}

// Error  Error() interface method
func (e RestError) Error() string {
	return e.ErrError
}

// Status Error status
func (e RestError) Status() int {
	return e.ErrStatus
}

// Causes RestError Causes
func (e RestError) Causes() interface{} {
	return e.ErrCauses
}

// NewRestError New Rest Error
func NewRestError(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

// NewRestErrorWithMessage New Rest Error With Message
func NewRestErrorWithMessage(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError:  err,
		ErrCauses: causes,
	}
}

func FromError(err error) (s RestErr) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, ErrRequestTimeoutError.Error(), err)
	case strings.Contains(err.Error(), "Unmarshal"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest.Error(), err)
	}

	var echoError *echo.HTTPError
	if errors.As(err, &echoError) {
		return NewRestError(echoError.Code, echoError.Error(), echoError.Internal)
	}

	var restError RestError
	if errors.As(err, &restError) {
		return NewRestError(restError.Status(), restError.Error(), restError.Causes())
	}

	return NewRestError(http.StatusInternalServerError, ErrInternalServerError.Error(), err)
}

// NewRestErrorFromBytes New Rest Error From Bytes
func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr RestError
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}

	return apiErr, nil
}

// NewBadRequestError New Bad Request Error
func NewBadRequestError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  ErrBadRequest.Error(),
		ErrCauses: causes,
	}
}

// NewNotFoundError New Not Found Error
func NewNotFoundError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusNotFound,
		ErrError:  ErrNotFound.Error(),
		ErrCauses: causes,
	}
}

// NewUnauthorizedError New ErrUnauthorized Error
func NewUnauthorizedError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusUnauthorized,
		ErrError:  ErrUnauthorized.Error(),
		ErrCauses: causes,
	}
}

// New ErrForbidden Error
func NewForbiddenError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusForbidden,
		ErrError:  ErrForbidden.Error(),
		ErrCauses: causes,
	}
}

// New Internal Server Error
func NewInternalServerError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  ErrInternalServerError.Error(),
		ErrCauses: causes,
	}

	return result
}
