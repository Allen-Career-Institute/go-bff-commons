// nolint : gocritic,
package utils

import (
	"github.com/Allen-Career-Institute/go-bff-commons/v1/config"
	"github.com/Allen-Career-Institute/go-bff-commons/v1/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestReadValidationError(t *testing.T) {
	type testCase struct {
		name          string
		expectedError error
		err           validator.ValidationErrors
	}

	testCases := []testCase{
		{
			name:          "No validation errors",
			expectedError: nil,
		},
		{
			name:          "With validation errors",
			expectedError: nil,
			err:           validator.ValidationErrors{},
		},
	}

	c := &config.Config{Logger: config.Logger{Level: "info"}}
	log := logger.NewAPILogger(c)

	log.InitLogger()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			result := ReadValidationError(ctx, tc.expectedError, log)

			assert.Equal(t, tc.expectedError, result)
		})
	}
}
