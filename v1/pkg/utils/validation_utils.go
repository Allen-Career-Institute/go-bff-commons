package utils

import (
	"errors"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

const (
	ValidationFailed = "ValidationFailed "
)

type Util struct {
	validate *validator.Validate
	log      logger.Logger
}

func NewUtils(validate *validator.Validate, log logger.Logger) *Util {
	return &Util{validate: validate, log: log}
}

func ReadValidationError(c echo.Context, err error, l logger.Logger) validator.FieldError {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	l.WithContext(c).Errorf("request validation failed for following fields")

	for _, e := range validationErrors {
		l.WithContext(c).Errorf(e.Error())
	}

	valE := validationErrors[0]

	return valE
}

func (u *Util) ValidateStruct(s interface{}) error {
	err := u.validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors

	errors.As(err, &validationErrors)

	for _, e := range validationErrors {
		u.log.Errorf(e.Error())
	}

	return validationErrors[0]
}
